package tars

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/codec"
	raw "github.com/TarsCloud/TarsGo/tars/codec/bytes"
	"github.com/TarsCloud/TarsGo/tars/codec/grpc"
	"github.com/TarsCloud/TarsGo/tars/codec/json"
	"github.com/TarsCloud/TarsGo/tars/codec/jsonrpc"
	"github.com/TarsCloud/TarsGo/tars/codec/proto"
	"github.com/TarsCloud/TarsGo/tars/codec/protorpc"

)

var (
	DefaultContentType = "application/octet-stream" 

	DefaultCodecs = map[string]codec.NewCodec{
		"application/grpc":         grpc.NewCodec,
		"application/grpc+json":    grpc.NewCodec,
		"application/grpc+proto":   grpc.NewCodec,
		"application/protobuf":     proto.NewCodec,
		"application/json":         json.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": raw.NewCodec,
	}

	// // TODO: remove legacy codec list
	// defaultCodecs = map[string]codec.NewCodec{
	// 	"application/json":         jsonrpc.NewCodec,
	// 	"application/json-rpc":     jsonrpc.NewCodec,
	// 	"application/protobuf":     protorpc.NewCodec,
	// 	"application/proto-rpc":    protorpc.NewCodec,
	// 	"application/octet-stream": protorpc.NewCodec,
	// }
)

func NewCodec(contentType string) (codec.NewCodec, error) {
	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}
