// Package redis is a redis implementation of kv
package redis

import (
	"time"
	"strings"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/TarsCloud/TarsGo/tars/config/options"
	"github.com/TarsCloud/TarsGo/tars/data/store"
)

var (
	DefaultMaxActive      = 0
	DefaultMaxIdle        = 5
	DefaultIdleTimeout    = 2 * time.Minute
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 5 * time.Second
	DefaultWriteTimeout   = 5 * time.Second
)

type rediskv struct {
	options.Options
	pool *redis.Pool
}

func (rs *rediskv) Read(key string) (*store.Record, error) {
	c := rs.pool.Get()
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	
	if value == "" {
		return nil, store.ErrNotFound
	}

	return &store.Record{
		Key:   key,
		Value: []byte(value),
	}, nil
}

func (rs *rediskv) Delete(key string) error {
	c := rs.pool.Get()
	defer c.Close()

	c.Do("DEL", key)
	return nil
}

func (rs *rediskv) Write(record *store.Record) error {
	c := rs.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", record.Key, record.Value)

	return err
}

func (rs *rediskv) Exist(key string) bool {
	c := rs.pool.Get()
	defer c.Close()

	if existed, err := redis.Int(c.Do("EXISTS", key)); err != nil || existed == 0 {
		return false
	}
	return true
}

func (rs *rediskv) Dump() ([]*store.Record, error) {
	return nil, store.ErrNotSupport

	// var vals []*store.Record
	// for _, keyv := range keyval {
	// 	vals = append(vals, &store.Record{
	// 		Key:   keyv.Key,
	// 		Value: keyv.Value,
	// 	})
	// }
	// return vals, nil
}

func (rs *rediskv) String() string {
	return "redis"
}

func (rs *rediskv) Disconnect() error {
	err := rs.pool.Close()
	rs.pool = nil
	return err
}

func NewStore(opts ...options.Option) store.Store {
	options := options.NewOptions(opts...)

	var nodes []string

	if n, ok := options.Values().Get("store.nodes"); ok {
		nodes = n.([]string)
	}

	var rawurl string 
	if len(nodes) > 0 {
		rawurl = nodes[0]
	}

	if !strings.HasPrefix(rawurl, "redis://") {
		rawurl = "redis://" + rawurl
	}

	fmt.Printf("Store redis rawurl: %v\n", rawurl)

	pool := &redis.Pool{
		MaxIdle:     DefaultMaxIdle,
		MaxActive:   DefaultMaxActive,
		IdleTimeout: DefaultIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(
				rawurl, //Example: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
				redis.DialConnectTimeout(DefaultConnectTimeout),
				redis.DialReadTimeout(DefaultReadTimeout),
				redis.DialWriteTimeout(DefaultWriteTimeout),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &rediskv{
		Options: options,
		pool:  pool,
	}
}
