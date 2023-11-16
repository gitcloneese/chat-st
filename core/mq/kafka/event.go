package kafka

import "context"

type Event interface {
	Key() string
	Value() []byte
	MsgKey() string
}

type Handler func(context.Context, Event) error

type Sender interface {
	Send(msg Event) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context, handler Handler) error
	Close() error
}
