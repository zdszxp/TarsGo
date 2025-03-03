package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/TarsCloud/TarsGo/tars/transport"
)

//MyServer testing tars udp server
type MyServer struct{}

//Invoke recv package and make response.
func (s *MyServer) Invoke(req []byte) (rsp []byte) {
	rsp = make([]byte, 4)
	rsp = append(rsp, []byte("Hello ")...)
	rsp = append(rsp, req...)
	binary.BigEndian.PutUint32(rsp[:4], uint32(len(rsp)))
	return
}

//ParsePackage parse full tars package.
func (s *MyServer) ParsePackage(buff []byte) (pkgLen, status int) {
	if len(buff) < 4 {
		return 0, transport.PACKAGE_LESS
	}
	length := binary.BigEndian.Uint32(buff[:4])

	if length > 1048576000 || len(buff) > 1048576000 { // 1000MB
		return 0, transport.PACKAGE_ERROR
	}
	if len(buff) < int(length) {
		return 0, transport.PACKAGE_LESS
	}
	return int(length), transport.PACKAGE_FULL
}

//InvokeTimeout show how to deal with timeout response.
func (s *MyServer) InvokeTimeout(ctx context.Context, pkg []byte) []byte {
	payload := []byte("timeout")
	ret := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(pkg[:4], uint32(len(ret)))
	copy(pkg[4:], payload)
	return ret
}

func main() {
	conf := &transport.TarsServerConf{
		Proto:   "udp",
		Address: "127.0.0.1:3333",
		//MaxAccept:     500,
		MaxInvoke:     20,
		AcceptTimeout: time.Millisecond * 500,
		ReadTimeout:   time.Millisecond * 100,
		WriteTimeout:  time.Millisecond * 100,
		HandleTimeout: time.Millisecond * 6000,
		IdleTimeout:   time.Millisecond * 600000,
	}
	s := MyServer{}
	svr := transport.NewTarsServer(&s, conf)
	err := svr.Serve()
	if err != nil {
		fmt.Println(err)
	}
}
