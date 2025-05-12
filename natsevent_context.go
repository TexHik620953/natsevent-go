package natsevent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsEventContext interface {
	context.Context
	EventJSON(value interface{}) error
}

type natsRPCContextImpl struct {
	ctx context.Context
	msg jetstream.Msg
}

func NewNatsRPCContext(ctx context.Context, msg jetstream.Msg) NatsEventContext {
	return &natsRPCContextImpl{
		ctx: ctx,
		msg: msg,
	}
}

// Default context methods

func (c *natsRPCContextImpl) Done() <-chan struct{} {
	return c.ctx.Done()
}
func (c *natsRPCContextImpl) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}
func (c *natsRPCContextImpl) Err() error {
	return c.ctx.Err()
}
func (c *natsRPCContextImpl) Value(key any) any {
	return c.ctx.Value(key)
}

// Custom methods
func (c *natsRPCContextImpl) EventJSON(value interface{}) error {
	return json.Unmarshal(c.msg.Data(), value)
}
