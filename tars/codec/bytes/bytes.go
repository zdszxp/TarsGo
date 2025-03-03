// Package bytes provides a bytes codec which does not encode or decode anything
package bytes

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/TarsCloud/TarsGo/tars/codec"
)

type Codec struct {
	Conn io.ReadWriteCloser
}

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *Codec) ReadBody(b interface{}) error {
	// read bytes
	buf, err := ioutil.ReadAll(c.Conn)
	if err != nil {
		return err
	}

	switch b.(type) {
	case *[]byte:
		v := b.(*[]byte)
		*v = buf
	case *Frame:
		v := b.(*Frame)
		v.Data = buf
	default:
		return fmt.Errorf("failed to read body: %v is not type of *[]byte", b)
	}

	return nil
}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	var v []byte
	switch b.(type) {
	case *Frame:
		v = b.(*Frame).Data
	case *[]byte:
		ve := b.(*[]byte)
		v = *ve
	case []byte:
		v = b.([]byte)
	default:
		return fmt.Errorf("failed to write: %v is not type of *[]byte or []byte", b)
	}
	_, err := c.Conn.Write(v)
	return err
}

func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "bytes"
}

func NewCodec(c io.ReadWriteCloser) codec.Codec {
	return &Codec{
		Conn: c,
	}
}
