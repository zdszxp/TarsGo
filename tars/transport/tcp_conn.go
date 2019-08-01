package transport

import (
	"net"
	"sync"
	"context"
	"errors"
)

type ConnSet map[net.Conn]Conn

type TCPConn struct {
	sync.Mutex
	conn      net.Conn
	writeChan chan []byte
	closeFlag bool
	filter PacketFilter
	authed bool
	ctx context.Context
}

func NewTCPConn(conn net.Conn, pendingWriteNum int, filter PacketFilter) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, pendingWriteNum)
	tcpConn.filter = filter
	tcpConn.ctx = context.Background()

	go func() {
		for b := range tcpConn.writeChan {
			if b == nil {
				break
			}

			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}

		conn.Close()
		tcpConn.Lock()
		tcpConn.closeFlag = true
		tcpConn.Unlock()
	}()

	return tcpConn
}

func (tcpConn *TCPConn) Filter(pf PacketFilter) {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	
	tcpConn.filter = pf
}

func (tcpConn *TCPConn) doDestroy() {
	tcpConn.conn.(*net.TCPConn).SetLinger(0)
	tcpConn.conn.Close()

	if !tcpConn.closeFlag {
		close(tcpConn.writeChan)
		tcpConn.closeFlag = true
	}
}

func (tcpConn *TCPConn) Destroy() {
	tcpConn.Lock()
	defer tcpConn.Unlock()

	tcpConn.doDestroy()
}

func (tcpConn *TCPConn) Close() {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag {
		return
	}

	tcpConn.doWrite(nil)
	tcpConn.closeFlag = true
}

func (tcpConn *TCPConn) doWrite(b []byte) error {
	if len(tcpConn.writeChan) == cap(tcpConn.writeChan) {
		tcpConn.doDestroy()
		return errors.New("close conn: channel full")
	}

	tcpConn.writeChan <- b
	return nil
}

// b must not be modified by the others goroutines
func (tcpConn *TCPConn) Write(b []byte) error {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag || b == nil {
		return errors.New("conn has closed")
	}

	return tcpConn.doWrite(b)
}

func (tcpConn *TCPConn) Read(b []byte) (int, error) {
	return tcpConn.conn.Read(b)
}

func (tcpConn *TCPConn) LocalAddr() net.Addr {
	return tcpConn.conn.LocalAddr()
}

func (tcpConn *TCPConn) RemoteAddr() net.Addr {
	return tcpConn.conn.RemoteAddr()
}

func (tcpConn *TCPConn) FilterRead(in []byte) ([]byte, error) {
	tcpConn.Lock()
	defer tcpConn.Unlock()

	if tcpConn.filter != nil {
		out, err := tcpConn.filter.Read(in)
		if err != nil {
			return nil, err
		} else {
			return out, nil
		}
	}

	return in, nil
}

func (tcpConn *TCPConn) FilterWrite(in []byte) ([]byte, error) {
	tcpConn.Lock()
	defer tcpConn.Unlock()

	if tcpConn.filter != nil {
		out, err := tcpConn.filter.Write(in)
		if err != nil {
			return nil, err
		} else {
			return out, nil
		}
	}

	return in, nil
}

func (tcpConn *TCPConn) WithValue(key, value interface{}){
	tcpConn.ctx = context.WithValue(tcpConn.ctx, key, value)
}

func (tcpConn *TCPConn) Value(key interface{}) interface{} {
	return tcpConn.ctx.Value(key)
}