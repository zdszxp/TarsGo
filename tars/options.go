package tars

import (
	"context"
	"strings"

	//"time"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/broker/redis"
	"github.com/TarsCloud/TarsGo/tars/data/store"
	"github.com/TarsCloud/TarsGo/tars/data/store/memory"
	redisStore "github.com/TarsCloud/TarsGo/tars/data/store/redis"
)

var tarsOpts *Options

func getOptions() *Options {
	return tarsOpts
}

func ConfigureWithConfigs(configs map[string]string) {
	var opts []Option

	storeConfig := configs["store"]
	if len(storeConfig) > 0 {
		if strings.HasPrefix(storeConfig, "redis") {
			var rawurl string
			if strings.HasPrefix(storeConfig, "redis"+"@") {
				rawurl = strings.TrimPrefix(storeConfig, "redis"+"@")
			}

			opts = append(opts, Store(redisStore.NewStore(store.Nodes(rawurl))))
		} else if strings.HasPrefix(storeConfig, "memory") {
			opts = append(opts, Store(memory.NewStore()))
		}
	}

	brokerConfig := configs["broker"]
	if len(brokerConfig) > 0 {
		if strings.HasPrefix(brokerConfig, "redis") {
			var rawurl string
			if strings.HasPrefix(storeConfig, "redis"+"@") {
				rawurl = strings.TrimPrefix(storeConfig, "redis"+"@")
			}

			opts = append(opts, Broker(redis.NewBroker(broker.Addrs(rawurl))))
		}
	}

	configure(opts...)
}

func configure(opts ...Option) {
	tarsOpts = newOptions(opts...)
}

type Option func(*Options)

type Options struct {
	broker broker.Broker
	store  store.Store

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		//broker:  redis.NewBroker(),
		Context: context.Background(),
		//store:   memory.NewStore(),
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

func (o *Options) Store() store.Store {
	return o.store
}

func Store(s store.Store) Option {
	return func(o *Options) {
		o.store = s
	}
}

func (o *Options) Broker() broker.Broker {
	return o.broker
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.broker = b
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// // WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
// func WrapSubscriber(w ...server.SubscriberWrapper) Option {
// 	return func(o *Options) {
// 		var wrappers []server.Option

// 		for _, wrap := range w {
// 			wrappers = append(wrappers, server.WrapSubscriber(wrap))
// 		}

// 		// Init once
// 		o.Server.Init(wrappers...)
// 	}
// }

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

func GetStore() store.Store {
	opts := getOptions()
	if opts == nil {
		return nil
	}

	return opts.Store()
}
