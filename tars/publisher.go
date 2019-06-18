package tars

import (
	"bytes"
	"context"
	"os"

	"github.com/google/uuid"

	"github.com/TarsCloud/TarsGo/tars/broker"
	"github.com/TarsCloud/TarsGo/tars/codec"
	"github.com/TarsCloud/TarsGo/tars/errors"
	"github.com/TarsCloud/TarsGo/tars/metadata"
)

// Publisher is syntactic sugar for publishing
// type Publisher interface {
// 	Publish(ctx context.Context, msg interface{}, opts ...PublishOption) error
// }

// //Publisher is binded with endpoint
// type publisher struct {
// 	topic string
// }

// func (p *publisher) Publish(ctx context.Context, contentType string, msg interface{}, opts ...PublishOption) error {
// 	return publish(ctx, contentType, p.topic, msg, opts...)
// }

// NewPublisher returns a new Publisher
// func NewPublisher(contentType string, topic string, opts ...PublishOption) Publisher {
// 	options := PublishOptions{
// 		Context: context.Background(),
// 	}
// 	for _, o := range opts {
// 		o(&options)
// 	}

// 	return &publisher{contentType, topic}
// }

type PublishOptions struct {
	// Exchange is the routing exchange for the message
	Exchange string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

// WithExchange sets the exchange to route a message through
func WithExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

func (bh *brokerFHelper) Publish(ctx context.Context, contentType string, topic string, msg interface{}, opts ...PublishOption) error {
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

	// get proxy
	if prx := os.Getenv("TARS_PROXY"); len(prx) > 0 {
		options.Exchange = prx
	}

	// get the exchange
	if len(options.Exchange) > 0 {
		topic = options.Exchange
	}

	// encode message body
	cf, err := NewCodec(contentType)
	if err != nil {
		return errors.InternalServerError("tars.publish", err.Error())
	}
	
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

	return getOptions().Broker().Publish(topic, &broker.Message{
		Header: md,
		Body:   b.Bytes(),
	})
}
