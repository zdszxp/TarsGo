package sync

import (
	"github.com/TarsCloud/TarsGo/tars/data/store"
	"github.com/TarsCloud/TarsGo/tars/sync/lock"
)

// WithLock sets the locking implementation option
func WithLock(l lock.Lock) Option {
	return func(o *Options) {
		o.Lock = l
	}
}

// WithStore sets the store implementation option
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithTime sets the time implementation option
func WithTime(t time.Time) Option {
	return func(o *Options) {
		o.Time = t
	}
}
