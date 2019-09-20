module tars

go 1.12

replace github.com/TarsCloud/TarsGo/tars => ./

replace github.com/openzipkin/zipkin-go-opentracing => github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.3

replace github.com/golang/protobuf/protoc-gen-go/tarsrpc => ./tools/pb2tarsgo/protoc-gen-go/tarsrpc/tarsrpc

require (
	github.com/TarsCloud/TarsGo/tars v0.0.0-00010101000000-000000000000
	github.com/beevik/ntp v0.2.0
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/go-redsync/redsync v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/hashicorp/consul/api v1.2.0
	github.com/nats-io/nats.go v1.8.1
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
)
