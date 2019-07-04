package session

import (
	"net"
	"encoding/binary"
	"sync"
)

const (
	StatusAuthed = 1
)

// Session is a net session
type Session interface {
	Address() string
	SendToClient(data []byte) error
	IsAuthed() bool
	Authed()
}

func NewSession(conn *net.TCPConn) Session{
	return &localSession{
		Conn:conn,
		status:0,
	}
}

type localSession struct {
	net.Conn
	status int
}

func (ls *localSession) IsAuthed() bool {
	return ls.status == StatusAuthed
}

func (ls *localSession) Authed() {
	ls.status = StatusAuthed
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

type Manager interface{
	GetSession(key string) (s Session, ok bool)
	AddSession(key string, s Session) bool
	DelSession(key string) bool
	BroadcastRaw(buff []byte) bool
	Broadcast(buff []byte) bool
}

func NewSessionManager() Manager{
	return &manager{
	}
}

type manager struct{
	sessions sync.Map
}

func (m *manager) GetSession(key string) (s Session, ok bool) {
	value, ok := m.sessions.Load(key)
	if ok {
		s, ok = value.(Session)
		return
	}

	return nil, false
}

func (m *manager) AddSession(key string, s Session) bool{
	_, ok := m.sessions.LoadOrStore(key, s)
	return ok
}

func (m *manager) DelSession(key string) bool{
	m.sessions.Delete(key)
	return true
}

func (m *manager) BroadcastRaw(buff []byte) bool {
	data := make([]byte, 4+len(buff))
	binary.BigEndian.PutUint32(data[:4], uint32(len(data)))
	copy(data[4:], buff)

	return m.Broadcast(data)
}

func (m *manager) Broadcast(buff []byte) bool {
	return m.RangeSessions(func(k, v interface{}) bool {
		v.(Session).SendToClient(buff)
		return true
	})
}

func (m *manager) RangeSessions(f func(key, value interface{}) bool) bool {
	m.sessions.Range(f)
	return true
}
