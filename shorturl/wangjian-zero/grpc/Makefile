all: vet test testrace

build: deps
	go build shorturl/wangjian-zero/grpc/...

clean:
	go clean -i shorturl/wangjian-zero/grpc/...

deps:
	go get -d -v shorturl/wangjian-zero/grpc/...

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	go generate shorturl/wangjian-zero/grpc/...

test: testdeps
	go test -cpu 1,4 -timeout 7m shorturl/wangjian-zero/grpc/...

testsubmodule: testdeps
	cd security/advancedtls && go test -cpu 1,4 -timeout 7m shorturl/wangjian-zero/grpc/security/advancedtls/...

testappengine: testappenginedeps
	goapp test -cpu 1,4 -timeout 7m shorturl/wangjian-zero/grpc/...

testappenginedeps:
	goapp get -d -v -t -tags 'appengine appenginevm' shorturl/wangjian-zero/grpc/...

testdeps:
	go get -d -v -t shorturl/wangjian-zero/grpc/...

testrace: testdeps
	go test -race -cpu 1,4 -timeout 7m shorturl/wangjian-zero/grpc/...

updatedeps:
	go get -d -v -u -f shorturl/wangjian-zero/grpc/...

updatetestdeps:
	go get -d -v -t -u -f shorturl/wangjian-zero/grpc/...

vet: vetdeps
	./vet.sh

vetdeps:
	./vet.sh -install

.PHONY: \
	all \
	build \
	clean \
	deps \
	proto \
	test \
	testappengine \
	testappenginedeps \
	testdeps \
	testrace \
	updatedeps \
	updatetestdeps \
	vet \
	vetdeps
