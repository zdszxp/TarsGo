package tars

import (
	"context"
	"strings"

	"github.com/TarsCloud/TarsGo/tars/sync"
	"github.com/TarsCloud/TarsGo/tars/sync/lock"
	redisLock "github.com/TarsCloud/TarsGo/tars/sync/lock/redis"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/broker/redis"

	"github.com/TarsCloud/TarsGo/tars/data/store"
	"github.com/TarsCloud/TarsGo/tars/data/store/memory"
	redisStore "github.com/TarsCloud/TarsGo/tars/data/store/redis"

	"github.com/TarsCloud/TarsGo/tars/api/session"
	_ "github.com/TarsCloud/TarsGo/tars/api/session/redis" //must import
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
			rawurl := strings.TrimPrefix(storeConfig, "redis"+"@")
			opts = append(opts, Store(redisStore.NewStore(store.Nodes(rawurl))))
		} else if strings.HasPrefix(storeConfig, "memory") {
			opts = append(opts, Store(memory.NewStore()))
		}
	}

	brokerConfig := configs["broker"]
	if len(brokerConfig) > 0 {
		var b broker.Broker
		if strings.HasPrefix(brokerConfig, "redis") {
			rawurl := strings.TrimPrefix(brokerConfig, "redis"+"@")
			b = redis.NewBroker(broker.Addrs(rawurl))
		}

		if b != nil {
			opts = append(opts, Broker(b))
			opts = append(opts, BeforeStart(func() error{
				err := b.Init()
				if err != nil {
					TLOG.Errorf("init broker error: %v", err)
					return err
				} else {
					TLOG.Debugf("Broker [%s] init successfully", b.String())
				}

				if err = b.Connect(); err != nil {
					TLOG.Errorf("connect broker[%v] error: %v", err)
					return err
				} else {
					TLOG.Debugf("Broker Listening on [%s]", brokerConfig)
				}

				return nil
			}))

			opts = append(opts, AfterStop(func() error{
				err := b.Disconnect()
				if err != nil {
					TLOG.Errorf("unload broker error: %v", err)
					return err
				} else {
					TLOG.Debug("unload broker successfully")
				}
				
				return nil
			}))
		}
	}

	var syncOpts sync.Options
	lockConfig := configs["lock"]
	if len(lockConfig) > 0 {
		if strings.HasPrefix(lockConfig, "redis") {
			rawurl := strings.TrimPrefix(lockConfig, "redis"+"@")

			syncOpts.Lock = redisLock.NewLock(lock.Nodes(rawurl))
		} //else if strings.HasPrefix(lockConfig, "") {}
	}

	opts = append(opts, SyncOptions(&syncOpts))

	sessionConfig := configs["session"]
	if len(sessionConfig) > 0 {
		if strings.HasPrefix(sessionConfig, "redis") {
			rawurl := strings.TrimPrefix(sessionConfig, "redis"+"@")

			sessionManager, err := session.NewManager("redis", rawurl)
			if err != nil {
				TLOG.Errorf("create session manager error: %v", err)
				return
			}

			opts = append(opts, SessionManager(sessionManager))
		}
	}

	configure(opts...)
}

func configure(opts ...Option) {
	tarsOpts = newOptions(opts...)
}

type Option func(*Options)

type Options struct {
	sync.Options
	broker broker.Broker
	store  store.Store
	sessionManager session.Manager

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error //not support//TO DO
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

func AppendOptions(opts ...Option) {
	opt := getOptions()
	for _, o := range opts {
		o(opt)
	}

	return
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

func (o *Options) SessionManager() session.Manager {
	return o.sessionManager
}

func SessionManager(s session.Manager) Option {
	return func(o *Options) {
		o.sessionManager = s
	}
}

func (o *Options) Lock() lock.Lock {
	return o.Options.Lock
}

func SyncOptions(opts *sync.Options) Option {
	return func(o *Options) {
		o.Options = *opts
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

func GetLock() lock.Lock {
	opts := getOptions()
	if opts == nil {
		return nil
	}

	return opts.Lock()
}

func GetSessionManager() session.Manager {
	opts := getOptions()
	if opts == nil {
		return nil
	}

	return opts.SessionManager()
}