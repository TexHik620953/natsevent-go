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

	ne.AddEventHandler("g1", "event1", func(msg jetstream.Msg) {
		fmt.Println("event1")
	})
	ne.AddEventHandler("g2", "event2", func(msg jetstream.Msg) {
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
	h := func(msg jetstream.Msg) {
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
