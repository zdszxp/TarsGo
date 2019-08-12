package session

import (
	"github.com/pkg/errors"
	"fmt"
)

const Maxlifetime = 60*60*24

type Manager interface {
	GetSession(sid string) (s Session, err error)
	AddSession(sid string) (s Session, err error)
	DelSession(sid string) error
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
		return nil, errors.Errorf("%v not exist", sid)
	}

	return m.provider.SessionRead(sid)
}

func (m *manager) AddSession(sid string) (s Session, err error) {
	if m.provider.SessionExist(sid) {
		return nil, errors.Errorf("%v exists", sid)
	}

	return m.provider.SessionRegenerate(sid, sid)//reset expire
}

func (m *manager) DelSession(sid string) error {
	if m.provider == nil {
		return errors.New("DelSession: provider == nil")
	}

	return m.provider.SessionDestroy(sid)
}