package natsevent_test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TexHik620953/natsevent-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func TestEvent(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
		return
	}

	js, err := jetstream.New(nc)
	if err != nil {
		t.Error(err)
		return
	}

	ne := natsevent.New(js, natsevent.WithSubjectRoot("TestEvent"))

	ne.AddEventHandler("g1", "event1", func(c natsevent.NatsEventContext) {
		fmt.Println("event1")
	})
	ne.AddEventHandler("g2", "event2", func(c natsevent.NatsEventContext) {
		fmt.Println("event2")
	})

	err = ne.StartWithContext(t.Context())
	if err != nil {
		t.Error(err)
		return
	}

	err = ne.FireEvent(t.Context(), "TestEvent.event1", "sraka")
	if err != nil {
		t.Error(err)
		return
	}
	<-time.After(time.Second)
}

func TestEventsDiffGroups(t *testing.T) {
	var dat []natsevent.NatsEventService

	var i atomic.Int32
	h := func(c natsevent.NatsEventContext) {
		i.Add(1)
	}
	for range 5 {
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			t.Error(err)
			return
		}
		js, err := jetstream.New(nc)
		if err != nil {
			t.Error(err)
			return
		}

		ne := natsevent.New(js, natsevent.WithSubjectRoot("TestEventsDiffGroups"))

		ne.AddEventHandler("g1", "event", h)

		err = ne.StartWithContext(t.Context())
		if err != nil {
			t.Error(err)
			return
		}
		dat = append(dat, ne)
	}

	err := dat[0].FireEvent(t.Context(), "TestEventsDiffGroups.event", "123")
	if err != nil {
		t.Error(err)
		return
	}

	<-time.After(time.Second)
	if i.Load() != 1 {
		t.Error("invalid received events cnt")
	}

}

func TestEventNoConsumer(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
		return
	}

	js, err := jetstream.New(nc)
	if err != nil {
		t.Error(err)
		return
	}

	ne := natsevent.New(js, natsevent.WithSubjectRoot("TestEventNoConsumer"))
	err = ne.StartWithContext(t.Context())
	if err != nil {
		t.Error(err)
		return
	}

	err = ne.FireEvent(t.Context(), "TestEventNoConsumer.event", "sraka")
	if err != nil {
		t.Error(err)
		return
	}
}

type testEvent struct {
	Text   string
	Number int
}

func TestEventSerializer(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
		return
	}

	js, err := jetstream.New(nc)
	if err != nil {
		t.Error(err)
		return
	}

	ne := natsevent.New(js, natsevent.WithSubjectRoot("TestEventNoConsumer"))

	ne.AddEventHandler("g1", "event", func(c natsevent.NatsEventContext) {
		var event testEvent
		err := c.EventJSON(&event)
		fmt.Println(err)

	})

	err = ne.StartWithContext(t.Context())
	if err != nil {
		t.Error(err)
		return
	}

	err = ne.FireEvent(t.Context(), "TestEventNoConsumer.event", &testEvent{
		Text:   "Hi ",
		Number: 13,
	})
	if err != nil {
		t.Error(err)
		return
	}
	<-time.After(time.Second * 5)
}
