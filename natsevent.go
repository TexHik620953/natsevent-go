package natsevent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

// Routing
type EventHandler func(msg jetstream.Msg)
type Route struct {
	Subject   string
	Handler   EventHandler
	CCtx      jetstream.ConsumeContext
	GroupName string // Добавляем поле для имени группы
}

// NatsEvent

type NatsEventService interface {
	AddEventHandler(groupName, eventName string, handler EventHandler)
	FireEvent(ctx context.Context, subject string, object any) error
	StartWithContext(ctx context.Context) error
}

type natsEventServiceImpl struct {
	js     jetstream.JetStream
	opts   NatsEventOptions
	routes map[string]*Route
	ctx    context.Context
}

func New(js jetstream.JetStream, options ...NatsEventOptionsFunc) NatsEventService {
	e := &natsEventServiceImpl{
		js:     js,
		opts:   GetDefaultOptions(),
		routes: make(map[string]*Route),
	}
	for _, v := range options {
		v(&e.opts)
	}
	return e
}

func (e *natsEventServiceImpl) AddEventHandler(groupName, eventName string, handler EventHandler) {
	if len(groupName) == 0 {
		panic("groupName cannot be empty")
	}
	if len(eventName) == 0 {
		panic("eventName cannot be empty")
	}
	if e.ctx != nil {
		panic("cannot add event handler after StartWithContext called")
	}

	route := &Route{
		Subject:   fmt.Sprintf("%s.%s", e.opts.SubjectRoot, eventName),
		Handler:   handler,
		GroupName: groupName,
	}

	if _, ex := e.routes[route.Subject]; ex {
		panic(fmt.Errorf("route %s already exists", route.Subject))
	}
	e.routes[route.Subject] = route
}

func (e *natsEventServiceImpl) FireEvent(ctx context.Context, subject string, object any) error {
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = e.js.Publish(ctx, subject, data)
	return err
}

func (e *natsEventServiceImpl) StartWithContext(ctx context.Context) error {
	if e.ctx != nil {
		panic("startWithContext called more than once")
	}
	e.ctx = ctx

	go func() {
		<-ctx.Done()
		for _, route := range e.routes {
			if route.CCtx != nil {
				route.CCtx.Stop()
			}
		}
	}()

	stream, err := e.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      e.opts.StreamName,
		Subjects:  []string{fmt.Sprintf("%s.*", e.opts.SubjectRoot)},
		Retention: jetstream.InterestPolicy,
		MaxAge:    e.opts.StreamMaxAge,
	})
	if err != nil {
		return err
	}

	for _, route := range e.routes {
		consumerConfig := jetstream.ConsumerConfig{
			FilterSubject: route.Subject,
			AckPolicy:     jetstream.AckExplicitPolicy,
			Durable:       route.GroupName,
		}

		consumer, err := stream.CreateOrUpdateConsumer(ctx, consumerConfig)
		if err != nil {
			return err
		}

		consumeContext, err := consumer.Consume(
			func(msg jetstream.Msg) {
				route.Handler(msg)
				msg.Ack() // Явное подтверждение обработки
			},
		)
		if err != nil {
			return err
		}
		route.CCtx = consumeContext
	}

	return nil
}
