package redis

import (
	"time"

	"github.com/TarsCloud/TarsGo/tars/broker"
)

var (
	DefaultMaxActive      = 0
	DefaultMaxIdle        = 5
	DefaultIdleTimeout    = 2 * time.Minute
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = time.Duration(0)
	DefaultWriteTimeout   = 5 * time.Second

	optionsKey = optionsKeyType{}
)

// options contain additional options for the broker.
type brokerOptions struct {
	maxIdle        int
	maxActive      int
	idleTimeout    time.Duration
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
}

type optionsKeyType struct{}

func ConnectTimeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.connectTimeout = d
	}
}

func ReadTimeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.readTimeout = d
	}
}

func WriteTimeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.writeTimeout = d
	}
}

func MaxIdle(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.maxIdle = n
	}
}

func MaxActive(n int) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.maxActive = n
	}
}

func IdleTimeout(d time.Duration) broker.Option {
	return func(o *broker.Options) {
		bo := o.Context.Value(optionsKey).(*brokerOptions)
		bo.idleTimeout = d
	}
}
