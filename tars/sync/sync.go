// Package sync is a distributed synchronization framework
package sync

import (
	"github.com/TarsCloud/TarsGo/tars/data/store"
	"github.com/TarsCloud/TarsGo/tars/sync/lock"
	"github.com/TarsCloud/TarsGo/tars/sync/time"
)

// Map provides synchronized access to key-value storage.
// It uses the store interface and lock interface to
// provide a consistent storage mechanism.
type Map interface {
	// Read value with given key
	Read(key, val interface{}) error
	// Write value with given key
	Write(key, val interface{}) error
	// Delete value with given key
	Delete(key interface{}) error
	// Iterate over all key/vals. Value changes are saved
	Iterate(func(key, val interface{}) error) error
}

type Options struct {
	Lock   lock.Lock
	Store  store.Store
	Time   time.Time
}

type Option func(o *Options)
