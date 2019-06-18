package tars

import (
	"context"
	//"time"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/broker/redis"
)

var tarsOpts *Options

func getOptions() *Options{
	return tarsOpts
}

func Configure(opts ...Option) {
	tarsOpts = newOptions(opts...)
}

type Option func(*Options)

type Options struct {
	broker    broker.Broker

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
		broker:    redis.NewBroker(),
		Context:   context.Background(),
	}

	for _, o := range opts {
		o(opt)
	}

	return opt
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
