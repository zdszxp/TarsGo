package tars

import (
	"context"
	"sync"
)

var brokerFHelperSingleton *brokerFHelper //Singleton
var brokerFHelperSingletonInitOnce sync.Once

func BrokerHelper() *brokerFHelper {
	brokerFHelperSingletonInitOnce.Do(func() {
		brokerFHelperSingleton = newBrokerFHelper()
	})

	return brokerFHelperSingleton
}

type BrokerOptions struct {
	//HdlrWrappers []HandlerWrapper
	SubWrappers []SubscriberWrapper

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type BrokerOption func(*BrokerOptions)

func newBrokerOptions(opt ...BrokerOption) BrokerOptions {
	opts := BrokerOptions{}
	
	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// Adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w SubscriberWrapper) BrokerOption {
	return func(o *BrokerOptions) {
		o.SubWrappers = append(o.SubWrappers, w)
	}
}

//brokerFHelper is helper struct for the broker module.
type brokerFHelper struct {
	subscriberHelper

	opts     BrokerOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newBrokerFHelper(opts ...BrokerOption) *brokerFHelper {
	options := newBrokerOptions(opts...)

	srv := &brokerFHelper{
		opts: options,
		subscriberHelper: subscriberHelper{
			subscribers: make(map[string]*subscriber),
			wg:          wait(options.Context),
			opts:        &options,
		},
	}

	return srv
}