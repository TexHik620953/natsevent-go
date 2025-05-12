package natsevent

import (
	"time"
)

type NatsEventOptionsFunc func(*NatsEventOptions)

// Defaults

func GetDefaultOptions() NatsEventOptions {
	return NatsEventOptions{
		StreamName:   "EVENTS",
		SubjectRoot:  "events",
		StreamMaxAge: time.Hour,
	}
}

// Options

type NatsEventOptions struct {
	StreamName   string
	SubjectRoot  string
	StreamMaxAge time.Duration
}

func WithStreamName(name string) NatsEventOptionsFunc {
	return func(ho *NatsEventOptions) {
		ho.StreamName = name
	}
}
func WithSubjectRoot(name string) NatsEventOptionsFunc {
	return func(ho *NatsEventOptions) {
		ho.SubjectRoot = name
	}
}
func WithStreamMaxAge(maxAge time.Duration) NatsEventOptionsFunc {
	return func(ho *NatsEventOptions) {
		ho.StreamMaxAge = maxAge
	}
}
