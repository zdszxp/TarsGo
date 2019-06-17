package tars

import (
	"context"
	"log"
	"sync"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/broker/redis"
)

var brokerFHelperSingleton *brokerFHelper //Singleton
var brokerFHelperSingletonInitOnce sync.Once

func BrokerFHelper() *brokerFHelper {
	brokerFHelperSingletonInitOnce.Do(func() {
		brokerFHelperSingleton = newBrokerFHelper()
		brokerFHelperSingleton.LoadBroker("user:foobared@192.168.10.158:6379/0")
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

	initOnce sync.Once
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
func (bh *brokerFHelper) LoadBroker(address string) (err error) {
	//broker.DefaultBroker = redis.NewBroker(broker.Addrs("192.168.10.158", "6387"))
	broker.DefaultBroker = redis.NewBroker(broker.Addrs(address))

	err = broker.Init()
	if err != nil {
		log.Fatal("Broker Init error: %v", err)
	} else {
		log.Print("Broker Init successfully")
	}

	if err = broker.Connect(); err != nil {
		log.Fatal("Broker Connect error: %v", err)
	} else {
		log.Print("Broker Connect successfully")
	}

	return err
}
