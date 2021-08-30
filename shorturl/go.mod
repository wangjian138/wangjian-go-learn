module shorturl

go 1.16

require (
	cloud.google.com/go v0.81.0 // indirect
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/alicebob/miniredis/v2 v2.14.1
	github.com/antlr/antlr4 v0.0.0-20210105212045-464bcbc32de2
	github.com/bgentry/speakeasy v0.1.0
	github.com/cncf/udpa/go v0.0.0-20210322005330-6414d713912e
	github.com/cockroachdb/datadriven v0.0.0-20190809214429-80d97fb3cbaa
	github.com/coreos/go-semver v0.2.0
	github.com/coreos/go-systemd/v22 v22.0.0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/creack/pty v1.1.7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v0.0.0-20171111073723-bb3d318650d4
	github.com/emicklei/proto v1.9.0
	github.com/envoyproxy/go-control-plane v0.9.9-0.20201210154907-fd9021fe5dad
	github.com/fatih/structtag v1.2.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/builder v0.3.4
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.1
	github.com/google/btree v1.0.0
	github.com/google/go-cmp v0.5.5
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/iancoleman/strcase v0.1.2
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.1.0
	github.com/json-iterator/go v1.1.9
	github.com/justinas/alice v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.3.0
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/modern-go/reflect2 v1.0.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pierrec/lz4 v2.5.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/client_model v0.2.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/soheilhy/cmux v0.1.4
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.1
	github.com/stretchr/testify v1.6.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20171017195756-830351dc03c6
	github.com/urfave/cli v1.22.5
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2
	go.etcd.io/bbolt v1.3.4
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/zap v1.14.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1
	golang.org/x/oauth2 v0.0.0-20210313182246-cd4f82c27b84
	golang.org/x/sys v0.0.0-20210414055047-fe65e336abe0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/genproto v0.0.0-20210402141018-6c239bbf2bb1
	gopkg.in/cheggaaa/pb.v1 v1.0.28
	gopkg.in/h2non/gock.v1 v1.0.15
	gopkg.in/yaml.v2 v2.4.0
	sigs.k8s.io/yaml v1.2.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
