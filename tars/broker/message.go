package broker

import (
)

// type event struct {
// 	topic       string
// 	contentType string
// 	payload     interface{}
// }

// func newEvent(topic string, payload interface{}, contentType string, opts ...client.MessageOption) client.Message {
// 	var options client.MessageOptions
// 	for _, o := range opts {
// 		o(&options)
// 	}

// 	if len(options.ContentType) > 0 {
// 		contentType = options.ContentType
// 	}

// 	return &event{
// 		payload:     payload,
// 		topic:       topic,
// 		contentType: contentType,
// 	}
// }

// func (g *event) ContentType() string {
// 	return g.contentType
// }

// func (g *event) Topic() string {
// 	return g.topic
// }

// func (g *event) Payload() interface{} {
// 	return g.payload
// }
