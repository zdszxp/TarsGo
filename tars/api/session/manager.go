package session

import (
	"github.com/pkg/errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")	
)

const Maxlifetime = 60*60*24

type Manager interface {
	GetSession(sid string) (s Session, err error)
	SessionStart(sid string) (s Session, err error)
	SessionDestroy(sid string) error
}

func NewManager(provideName string, providerConfig string) (Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}

	err := provider.SessionInit(Maxlifetime, providerConfig)
	if err != nil {
		return nil, err
	}

	return &manager{
		provider: provider,
	}, nil
}

//session manager for net connections
type manager struct {
	provider Provider
}

func (manager *manager) GetProvider() Provider {
	return manager.provider
}

func (m *manager) GetSession(sid string) (s Session, err error) {
	if m.provider.SessionExist(sid) == false {
		return nil, ErrNotFound
	}

	return m.provider.SessionRead(sid)
}

func (m *manager) SessionStart(sid string) (s Session, err error) {
	if sid != "" && m.provider.SessionExist(sid) {
		return m.provider.SessionRead(sid)
	}

	s, err = m.provider.SessionRead(sid)
	if err != nil {
		return nil, err
	}

	return
}

func (m *manager) SessionDestroy(sid string) error {
	if m.provider == nil {
		return errors.New("SessionDestroy: provider == nil")
	}

	return m.provider.SessionDestroy(sid)
}