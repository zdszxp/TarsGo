module tars

go 1.12

replace github.com/TarsCloud/TarsGo/tars => ./

replace github.com/openzipkin/zipkin-go-opentracing => github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.3

replace github.com/golang/protobuf/protoc-gen-go/tarsrpc => ./tools/pb2tarsgo/protoc-gen-go/tarsrpc/tarsrpc

require (
	github.com/TarsCloud/TarsGo/tars v0.0.0-00010101000000-000000000000
	github.com/beevik/ntp v0.2.0
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-redsync/redsync v1.3.1
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.1 // indirect
	github.com/hashicorp/consul/api v1.2.0
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/nats-io/nats-server/v2 v2.0.4 // indirect
	github.com/nats-io/nats.go v1.8.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)
