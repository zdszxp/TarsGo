// Package redis is a redis implemenation of lock
package redis

import (
	"errors"
	"sync"
	"time"
	"strings"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/TarsCloud/TarsGo/tars/sync/lock"
)

var (
	DefaultMaxActive      = 0
	DefaultMaxIdle        = 5
	DefaultIdleTimeout    = 2 * time.Minute
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = 5 * time.Second
)

type redisLock struct {
	sync.Mutex

	locks map[string]*redsync.Mutex
	opts  lock.Options
	c     *redsync.Redsync
}

func (r *redisLock) Acquire(id string, opts ...lock.AcquireOption) error {
	var options lock.AcquireOptions
	for _, o := range opts {
		o(&options)
	}

	var ropts []redsync.Option

	if options.Wait > time.Duration(0) {
		ropts = append(ropts, redsync.SetRetryDelay(options.Wait))
		ropts = append(ropts, redsync.SetTries(1))
	}

	if options.TTL > time.Duration(0) {
		ropts = append(ropts, redsync.SetExpiry(options.TTL))
	}

	m := r.c.NewMutex(r.opts.Prefix+id, ropts...)
	err := m.Lock()
	if err != nil {
		return err
	}

	r.Lock()
	r.locks[id] = m
	r.Unlock()

	return nil
}

func (r *redisLock) Release(id string) error {
	r.Lock()
	defer r.Unlock()
	m, ok := r.locks[id]
	if !ok {
		return errors.New("lock not found")
	}

	unlocked := m.Unlock()
	delete(r.locks, id)

	if !unlocked {
		return errors.New("lock not unlocked")
	}

	return nil
}

func (r *redisLock) String() string {
	return "redis"
}

func NewLock(opts ...lock.Option) lock.Lock {
	var options lock.Options
	for _, o := range opts {
		o(&options)
	}

	nodes := options.Nodes

	if len(nodes) == 0 {
		nodes = []string{"127.0.0.1:6379"}
	}

	var pools []redsync.Pool
	for _, addr := range nodes {
		if !strings.HasPrefix(addr, "redis://") {
			addr = "redis://" + addr
		}

		pools = append(pools, &redis.Pool{
			MaxIdle:     DefaultMaxIdle,
			MaxActive:   DefaultMaxActive,
			IdleTimeout: DefaultIdleTimeout,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(
					addr,
					redis.DialConnectTimeout(DefaultConnectTimeout),
					redis.DialReadTimeout(DefaultReadTimeout),
					redis.DialWriteTimeout(DefaultWriteTimeout),
				)
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		})
	}

	rpool := redsync.New(pools)

	return &redisLock{
		locks: make(map[string]*redsync.Mutex),
		opts:  options,
		c:     rpool,
	}
}
