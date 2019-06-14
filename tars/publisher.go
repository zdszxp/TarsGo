package tars

import (
	"context"
	"bytes"
	"sync"
	"log"
	// "reflect"

	"github.com/google/uuid"

	"github.com/TarsCloud/TarsGo/tars/errors"
	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/metadata"
	"github.com/TarsCloud/TarsGo/tars/codec"
)

var (
	once sync.Once
)

type PublishOptions struct {
	// Exchange is the routing exchange for the message
	Exchange string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

type Publisher struct {
	ContentType string
}

func (p *Publisher) Publish(ctx context.Context, topic string, msg interface{}, opts ...PublishOption) error {
	return publish(ctx, p.ContentType, topic, msg, opts...)
}

func publish(ctx context.Context, contentType string, topic string, msg interface{}, opts ...PublishOption) error {
	options := PublishOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	id := uuid.New().String()
	md["Content-Type"] = contentType
	md["Tars-Topic"] = topic
	md["Tars-Id"] = id

	// // get proxy
	// if prx := os.Getenv("MICRO_PROXY"); len(prx) > 0 {
	// 	options.Exchange = prx
	// }

	// // get the exchange
	// if len(options.Exchange) > 0 {
	// 	topic = options.Exchange
	// }

	// encode message body
	cf, err := NewCodec(contentType)
	if err != nil {
		return errors.InternalServerError("tars.publish", err.Error())
	}

	//log.Printf("[pub] msg: %v \n %v", msg, string(reflect.ValueOf(msg).Interface().([]uint8)))

	b := &buffer{bytes.NewBuffer(nil)}
	if err := cf(b).Write(&codec.Message{
		Target: topic,
		Type:   codec.Publication,
		Header: map[string]string{
			"Tars-Id":    id,
			"Tars-Topic": topic,
		},
	}, msg); err != nil {
		return errors.InternalServerError("tars.publish", err.Error())
	}

	// var v []byte = make([]byte, 256)
	// if err := cf(b).ReadBody(&v); err != nil {
	// 	return err
	// }

	log.Println("[pub] encode: ", b.String())

	//log.Printf("[pub] msg: %v \n %v", msg, string(reflect.ValueOf(msg).Interface().([]uint8)))

	// once.Do(func() {
	// 	initBroker()
	// })

	return broker.Publish(topic, &broker.Message{
		Header: md,
		Body:   b.Bytes(),
	})
}
