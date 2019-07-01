package session

import (
	"net"
	"encoding/binary"
	"sync"
)

// Session is a net session
type Session interface {
	Address() string
	SendToClient(data []byte) error
}

func NewSession(conn *net.TCPConn) Session{
	return &localSession{
		conn,
	}
}

type localSession struct {
	net.Conn
}

func (ls *localSession) Address() string {
	return ls.RemoteAddr().String()
}

func (ls *localSession) SendToClient(data []byte) (err error) {
	if _, err = ls.Write(data); err != nil {
		//TLOG.Errorf("send pkg to %v failed %v", ls.RemoteAddr().String(), err)
		return
	}
	return
}

type SessionManager interface{
	GetSession(key string) (s Session, ok bool)
	AddSession(key string, s Session) bool
	DelSession(key string) bool
	BroadcastRaw(buff []byte) bool
	Broadcast(buff []byte) bool
}

func NewSessionManager() SessionManager{
	return &sessionManager{
	}
}

type sessionManager struct{
	sessions sync.Map
}

func (sm *sessionManager) GetSession(key string) (s Session, ok bool) {
	value, ok := sm.sessions.Load(key)
	if ok {
		s, ok = value.(Session)
		return
	}

	return nil, false
}

func (sm *sessionManager) AddSession(key string, s Session) bool{
	_, ok := sm.sessions.LoadOrStore(key, s)
	return ok
}

func (sm *sessionManager) DelSession(key string) bool{
	sm.sessions.Delete(key)
	return true
}

func (sm *sessionManager) BroadcastRaw(buff []byte) bool {
	data := make([]byte, 4+len(buff))
	binary.BigEndian.PutUint32(data[:4], uint32(len(data)))
	copy(data[4:], buff)

	return sm.Broadcast(data)
}

func (sm *sessionManager) Broadcast(buff []byte) bool {
	return sm.RangeSessions(func(k, v interface{}) bool {
		v.(Session).SendToClient(buff)
		return true
	})
}

func (sm *sessionManager) RangeSessions(f func(key, value interface{}) bool) bool {
	sm.sessions.Range(f)
	return true
}
