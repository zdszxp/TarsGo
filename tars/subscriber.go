package tars

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
	"github.com/pkg/errors"
	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/codec"
	"github.com/TarsCloud/TarsGo/tars/metadata"
)

var (
	ErrBrokerTopicExists = errors.New("topic exists")
)

const (
	subSig = "func(context.Context, interface{}) error"
)

type subscriberHelper struct {
	sync.RWMutex
	wg *sync.WaitGroup

	// handlers    map[string]server.Handler
	subscribers map[string]*subscriber
	opts        *BrokerOptions
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func (s *subscriberHelper) RegisterSubscriber(topic string, h interface{}, opts ...SubscriberOption) error {
	s.Lock()
	_, ok := s.subscribers[topic]
	if ok {
		return ErrBrokerTopicExists
	}
	s.Unlock()

	//TLOG.Debugf("Subscrib topic %v", topic)
	subscriber := newSubscriber(topic, h, opts...)
	return s.Subscribe(subscriber)
}

func (s *subscriberHelper) NewSubscriber(topic string, handler interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, handler, opts...)
}

func (s *subscriberHelper) Unsubscribe(topic string) error {
	s.Lock()
	defer s.Unlock()

	sub, ok := s.subscribers[topic]
	if !ok {
		return fmt.Errorf("subscriber %v not exists", topic)
	}

	for _, v := range sub.subscribers {
		v.Unsubscribe()
	}

	delete(s.subscribers, topic)
	return nil
}

func (s *subscriberHelper) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *subscriber")
	}
	
	if len(sub.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	_, ok = s.subscribers[sub.topic]
	if ok {
		return ErrBrokerTopicExists
	}

	handler := s.createSubHandler(sub, *s.opts)
	var opts []broker.SubscribeOption
	if queue := sb.Options().Queue; len(queue) > 0 {
		opts = append(opts, broker.Queue(queue))
	}

	if !sb.Options().AutoAck {
		opts = append(opts, broker.DisableAutoAck())
	}

	bsub, err := getOptions().Broker().Subscribe(sub.Topic(), handler, opts...)
	if err != nil {
		return err
	}
	sub.subscribers = append(sub.subscribers, bsub)
	s.subscribers[sub.topic] = sub

	return nil
}

// Subscriber interface represents a subscription to a given topic using
// a specific subscriber function or object with endpoints.
type Subscriber interface {
	Topic() string
	Subscriber() interface{}
	//Endpoints() []*registry.Endpoint
	Options() SubscriberOptions
}

type handler struct {
	method  reflect.Value
	reqType reflect.Type
	ctxType reflect.Type
}

type subscriber struct {
	topic      string
	rcvr       reflect.Value
	typ        reflect.Type
	subscriber interface{}
	handlers   []*handler
	//endpoints  []*registry.Endpoint
	opts SubscriberOptions

	subscribers []broker.Subscriber
}

func newSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	options := SubscriberOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&options)
	}

	//var endpoints []*registry.Endpoint
	var handlers []*handler

	if typ := reflect.TypeOf(sub); typ.Kind() == reflect.Func {
		h := &handler{
			method: reflect.ValueOf(sub),
		}

		switch typ.NumIn() {
		case 1:
			h.reqType = typ.In(0)
		case 2:
			h.ctxType = typ.In(0)
			h.reqType = typ.In(1)
		}

		handlers = append(handlers, h)

		// endpoints = append(endpoints, &registry.Endpoint{
		// 	Name:    "Func",
		// 	Request: extractSubValue(typ),
		// 	Metadata: map[string]string{
		// 		"topic":      topic,
		// 		"subscriber": "true",
		// 	},
		// })
	} else {
		//hdlr := reflect.ValueOf(sub)
		//name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}

			switch method.Type.NumIn() {
			case 2:
				h.reqType = method.Type.In(1)
			case 3:
				h.ctxType = method.Type.In(1)
				h.reqType = method.Type.In(2)
			}

			handlers = append(handlers, h)

			// endpoints = append(endpoints, &registry.Endpoint{
			// 	Name:    name + "." + method.Name,
			// 	Request: extractSubValue(method.Type),
			// 	Metadata: map[string]string{
			// 		"topic":      topic,
			// 		"subscriber": "true",
			// 	},
			// })
		}
	}

	return &subscriber{
		rcvr:       reflect.ValueOf(sub),
		typ:        reflect.TypeOf(sub),
		topic:      topic,
		subscriber: sub,
		handlers:   handlers,
		//endpoints:  endpoints,
		opts: options,
	}
}

//check subscriber signature
func validateSubscriber(sub Subscriber) error {
	typ := reflect.TypeOf(sub.Subscriber())
	var argType reflect.Type

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 2:
			argType = typ.In(1)
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("subscriber %v returns %v not error", name, returnType.String())
		}
	} else {
		hdlr := reflect.ValueOf(sub.Subscriber())
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)

			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), subSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of outs: %v require signature %s",
					name, method.Name, method.Type.NumOut(), subSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}

func (s *subscriberHelper) createSubHandler(sb *subscriber, opts BrokerOptions) broker.Handler {
	return func(p broker.Event) error {
		msg := p.Message()

		// get codec
		ct := msg.Header["Content-Type"]

		// default content type
		if len(ct) == 0 {
			msg.Header["Content-Type"] = DefaultContentType
			ct = DefaultContentType
		}

		// get codec
		cf, err := NewCodec(ct)
		if err != nil {
			return err
		}

		// copy headers
		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}

		// create context
		ctx := metadata.NewContext(context.Background(), hdr)

		results := make(chan error, len(sb.handlers))

		for i := 0; i < len(sb.handlers); i++ {
			handler := sb.handlers[i]

			var isVal bool
			var req reflect.Value

			if handler.reqType.Kind() == reflect.Ptr {
				req = reflect.New(handler.reqType.Elem())
			} else {
				req = reflect.New(handler.reqType)
				isVal = true
			}
			if isVal {
				req = req.Elem()
			}

			b := &buffer{bytes.NewBuffer(msg.Body)}
			co := cf(b)
			defer co.Close()

			if err := co.ReadHeader(&codec.Message{}, codec.Event); err != nil {
				return err
			}

			if err := co.ReadBody(req.Interface()); err != nil {
				return err
			}

			fn := func(ctx context.Context, payload interface{}) error {
				var vals []reflect.Value
				if sb.typ.Kind() != reflect.Func {
					vals = append(vals, sb.rcvr)
				}
				if handler.ctxType != nil {
					vals = append(vals, reflect.ValueOf(ctx))
				}

				vals = append(vals, reflect.ValueOf(payload))

				returnValues := handler.method.Call(vals)
				if err := returnValues[0].Interface(); err != nil {
					return err.(error)
				}
				return nil
			}

			// for i := len(opts.SubWrappers); i > 0; i-- {
			// 	fn = opts.SubWrappers[i-1](fn)
			// }

			if s.wg != nil {
				s.wg.Add(1)
			}

			go func() {
				if s.wg != nil {
					defer s.wg.Done()
				}

				results <- fn(ctx, req.Interface())
			}()
		}

		var errors []string

		for i := 0; i < len(sb.handlers); i++ {
			if err := <-results; err != nil {
				errors = append(errors, err.Error())
			}
		}

		if len(errors) > 0 {
			return fmt.Errorf("subscriber error: %s", strings.Join(errors, "\n"))
		}

		return nil
	}
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Subscriber() interface{} {
	return s.subscriber
}

// func (s *subscriber) Endpoints() []*registry.Endpoint {
// 	return s.endpoints
// }

func (s *subscriber) Options() SubscriberOptions {
	return s.opts
}

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func wait(ctx context.Context) *sync.WaitGroup {
	if ctx == nil {
		return nil
	}
	wg, ok := ctx.Value("wait").(*sync.WaitGroup)
	if !ok {
		return nil
	}
	return wg
}
