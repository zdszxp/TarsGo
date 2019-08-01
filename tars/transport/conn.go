package transport

import (
	"net"
)

type Conn interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Close()
	Destroy()

	Write([]byte) error

	Filter(PacketFilter)
	FilterRead([]byte) ([]byte, error)
	FilterWrite([]byte) ([]byte, error)

	WithValue(key, value interface{})
	Value(key interface{}) interface{}
}
