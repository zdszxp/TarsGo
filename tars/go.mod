module tars

go 1.12

replace github.com/TarsCloud/TarsGo/tars => ./

require (
	github.com/TarsCloud/TarsGo/tars v0.0.0-00010101000000-000000000000
	github.com/go-redsync/redsync v1.3.0 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/consul/api v1.2.0 // indirect
	github.com/pkg/errors v0.8.1
)
