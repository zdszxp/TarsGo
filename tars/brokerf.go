package tars

import (
	"context"
	"sync"

	"github.com/TarsCloud/TarsGo/tars/broker"
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
			subscribers: make(map[*subscriber][]broker.Subscriber),
			wg:          wait(options.Context),
			opts:        &options,
		},
	}

	return srv
}

//The connection address may be a fully qualified IANA address such
// as: redis://user:secret@localhost:6379/0?foo=bar&qux=baz
func (bh *brokerFHelper) LoadBroker(opts ...broker.Option) (err error) {
	err = getOptions().Broker().Init(opts...)
	if err != nil {
		TLOG.Errorf("Broker Init error: %v", err)
	} else {
		TLOG.Debug("Broker Init successfully")
	}

	if err = getOptions().Broker().Connect(); err != nil {
		TLOG.Errorf("Broker Connect error: %v", err)
	} else {
		TLOG.Debug("Broker Connect successfully")
	}

	return err
}
