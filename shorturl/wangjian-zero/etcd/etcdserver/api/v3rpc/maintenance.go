// Copyright 2016 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v3rpc

import (
	"context"
	"crypto/sha256"
	"io"

	"shorturl/wangjian-zero/etcd/auth"
	"shorturl/wangjian-zero/etcd/etcdserver"
	"shorturl/wangjian-zero/etcd/etcdserver/api/v3rpc/rpctypes"
	pb "shorturl/wangjian-zero/etcd/etcdserver/etcdserverpb"
	"shorturl/wangjian-zero/etcd/mvcc"
	"shorturl/wangjian-zero/etcd/mvcc/backend"
	"shorturl/wangjian-zero/etcd/raft"
	"shorturl/wangjian-zero/etcd/version"

	"go.uber.org/zap"
)

type KVGetter interface {
	KV() mvcc.ConsistentWatchableKV
}

type BackendGetter interface {
	Backend() backend.Backend
}

type Alarmer interface {
	// Alarms is implemented in Server interface located in etcdserver/server.go
	// It returns a list of alarms present in the AlarmStore
	Alarms() []*pb.AlarmMember
	Alarm(ctx context.Context, ar *pb.AlarmRequest) (*pb.AlarmResponse, error)
}

type Downgrader interface {
	Downgrade(ctx context.Context, dr *pb.DowngradeRequest) (*pb.DowngradeResponse, error)
}

type LeaderTransferrer interface {
	MoveLeader(ctx context.Context, lead, target uint64) error
}

type AuthGetter interface {
	AuthInfoFromCtx(ctx context.Context) (*auth.AuthInfo, error)
	AuthStore() auth.AuthStore
}

type ClusterStatusGetter interface {
	IsLearner() bool
}

type maintenanceServer struct {
	lg  *zap.Logger
	rg  etcdserver.RaftStatusGetter
	kg  KVGetter
	bg  BackendGetter
	a   Alarmer
	lt  LeaderTransferrer
	hdr header
	cs  ClusterStatusGetter
	d   Downgrader
}

func NewMaintenanceServer(s *etcdserver.EtcdServer) pb.MaintenanceServer {
	srv := &maintenanceServer{lg: s.Cfg.Logger, rg: s, kg: s, bg: s, a: s, lt: s, hdr: newHeader(s), cs: s, d: s}
	if srv.lg == nil {
		srv.lg = zap.NewNop()
	}
	return &authMaintenanceServer{srv, s}
}

func (ms *maintenanceServer) Defragment(ctx context.Context, sr *pb.DefragmentRequest) (*pb.DefragmentResponse, error) {
	ms.lg.Info("starting defragment")
	err := ms.bg.Backend().Defrag()
	if err != nil {
		ms.lg.Warn("failed to defragment", zap.Error(err))
		return nil, err
	}
	ms.lg.Info("finished defragment")
	return &pb.DefragmentResponse{}, nil
}

func (ms *maintenanceServer) Snapshot(sr *pb.SnapshotRequest, srv pb.Maintenance_SnapshotServer) error {
	snap := ms.bg.Backend().Snapshot()
	pr, pw := io.Pipe()

	defer pr.Close()

	go func() {
		snap.WriteTo(pw)
		if err := snap.Close(); err != nil {
			ms.lg.Warn("failed to close snapshot", zap.Error(err))
		}
		pw.Close()
	}()

	// send file data
	h := sha256.New()
	br := int64(0)
	buf := make([]byte, 32*1024)
	sz := snap.Size()
	for br < sz {
		n, err := io.ReadFull(pr, buf)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return togRPCError(err)
		}
		br += int64(n)
		resp := &pb.SnapshotResponse{
			RemainingBytes: uint64(sz - br),
			Blob:           buf[:n],
		}
		if err = srv.Send(resp); err != nil {
			return togRPCError(err)
		}
		h.Write(buf[:n])
	}

	// send sha
	sha := h.Sum(nil)
	hresp := &pb.SnapshotResponse{RemainingBytes: 0, Blob: sha}
	if err := srv.Send(hresp); err != nil {
		return togRPCError(err)
	}

	return nil
}

func (ms *maintenanceServer) Hash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	h, rev, err := ms.kg.KV().Hash()
	if err != nil {
		return nil, togRPCError(err)
	}
	resp := &pb.HashResponse{Header: &pb.ResponseHeader{Revision: rev}, Hash: h}
	ms.hdr.fill(resp.Header)
	return resp, nil
}

func (ms *maintenanceServer) HashKV(ctx context.Context, r *pb.HashKVRequest) (*pb.HashKVResponse, error) {
	h, rev, compactRev, err := ms.kg.KV().HashByRev(r.Revision)
	if err != nil {
		return nil, togRPCError(err)
	}

	resp := &pb.HashKVResponse{Header: &pb.ResponseHeader{Revision: rev}, Hash: h, CompactRevision: compactRev}
	ms.hdr.fill(resp.Header)
	return resp, nil
}

func (ms *maintenanceServer) Alarm(ctx context.Context, ar *pb.AlarmRequest) (*pb.AlarmResponse, error) {
	resp, err := ms.a.Alarm(ctx, ar)
	if err != nil {
		return nil, togRPCError(err)
	}
	if resp.Header == nil {
		resp.Header = &pb.ResponseHeader{}
	}
	ms.hdr.fill(resp.Header)
	return resp, nil
}

func (ms *maintenanceServer) Status(ctx context.Context, ar *pb.StatusRequest) (*pb.StatusResponse, error) {
	hdr := &pb.ResponseHeader{}
	ms.hdr.fill(hdr)
	resp := &pb.StatusResponse{
		Header:           hdr,
		Version:          version.Version,
		Leader:           uint64(ms.rg.Leader()),
		RaftIndex:        ms.rg.CommittedIndex(),
		RaftAppliedIndex: ms.rg.AppliedIndex(),
		RaftTerm:         ms.rg.Term(),
		DbSize:           ms.bg.Backend().Size(),
		DbSizeInUse:      ms.bg.Backend().SizeInUse(),
		IsLearner:        ms.cs.IsLearner(),
	}
	if resp.Leader == raft.None {
		resp.Errors = append(resp.Errors, etcdserver.ErrNoLeader.Error())
	}
	for _, a := range ms.a.Alarms() {
		resp.Errors = append(resp.Errors, a.String())
	}
	return resp, nil
}

func (ms *maintenanceServer) MoveLeader(ctx context.Context, tr *pb.MoveLeaderRequest) (*pb.MoveLeaderResponse, error) {
	if ms.rg.ID() != ms.rg.Leader() {
		return nil, rpctypes.ErrGRPCNotLeader
	}

	if err := ms.lt.MoveLeader(ctx, uint64(ms.rg.Leader()), tr.TargetID); err != nil {
		return nil, togRPCError(err)
	}
	return &pb.MoveLeaderResponse{}, nil
}

func (ms *maintenanceServer) Downgrade(ctx context.Context, r *pb.DowngradeRequest) (*pb.DowngradeResponse, error) {
	resp, err := ms.d.Downgrade(ctx, r)
	if err != nil {
		return nil, togRPCError(err)
	}
	resp.Header = &pb.ResponseHeader{}
	ms.hdr.fill(resp.Header)
	return resp, nil
}

type authMaintenanceServer struct {
	*maintenanceServer
	ag AuthGetter
}

func (ams *authMaintenanceServer) isAuthenticated(ctx context.Context) error {
	authInfo, err := ams.ag.AuthInfoFromCtx(ctx)
	if err != nil {
		return err
	}

	return ams.ag.AuthStore().IsAdminPermitted(authInfo)
}

func (ams *authMaintenanceServer) Defragment(ctx context.Context, sr *pb.DefragmentRequest) (*pb.DefragmentResponse, error) {
	if err := ams.isAuthenticated(ctx); err != nil {
		return nil, err
	}

	return ams.maintenanceServer.Defragment(ctx, sr)
}

func (ams *authMaintenanceServer) Snapshot(sr *pb.SnapshotRequest, srv pb.Maintenance_SnapshotServer) error {
	if err := ams.isAuthenticated(srv.Context()); err != nil {
		return err
	}

	return ams.maintenanceServer.Snapshot(sr, srv)
}

func (ams *authMaintenanceServer) Hash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	if err := ams.isAuthenticated(ctx); err != nil {
		return nil, err
	}

	return ams.maintenanceServer.Hash(ctx, r)
}

func (ams *authMaintenanceServer) HashKV(ctx context.Context, r *pb.HashKVRequest) (*pb.HashKVResponse, error) {
	if err := ams.isAuthenticated(ctx); err != nil {
		return nil, err
	}
	return ams.maintenanceServer.HashKV(ctx, r)
}

func (ams *authMaintenanceServer) Status(ctx context.Context, ar *pb.StatusRequest) (*pb.StatusResponse, error) {
	return ams.maintenanceServer.Status(ctx, ar)
}

func (ams *authMaintenanceServer) MoveLeader(ctx context.Context, tr *pb.MoveLeaderRequest) (*pb.MoveLeaderResponse, error) {
	return ams.maintenanceServer.MoveLeader(ctx, tr)
}

func (ams *authMaintenanceServer) Downgrade(ctx context.Context, r *pb.DowngradeRequest) (*pb.DowngradeResponse, error) {
	return ams.maintenanceServer.Downgrade(ctx, r)
}
