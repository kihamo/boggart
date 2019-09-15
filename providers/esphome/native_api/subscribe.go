package native_api

import (
	"context"
	"github.com/golang/protobuf/proto"
	"time"
)

type Subscribe struct {
	client         *Client
	ctx            context.Context
	tickerInterval time.Duration
	startMessage   proto.Message
	messages       chan proto.Message
	errors         chan error
}

func NewSubscribe(client *Client, ctx context.Context, startMessage proto.Message, interval time.Duration) *Subscribe {
	s := &Subscribe{
		client:         client,
		ctx:            ctx,
		tickerInterval: interval,
		startMessage:   startMessage,
		messages:       make(chan proto.Message),
		errors:         make(chan error),
	}
	go s.start()

	return s
}

func (s *Subscribe) start() {
	conn, err := s.client.connect()
	if err != nil {
		return
	}

	conn.Lock()
	defer conn.Unlock()

	ticker := time.NewTicker(s.tickerInterval)
	defer func() {
		ticker.Stop()
		close(s.messages)
		close(s.errors)
	}()

	var init bool

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			if !init {
				if err = conn.Write(s.ctx, s.startMessage); err != nil {
					s.errors <- err
					continue
				}

				init = true
			}

			message, err := conn.Read(s.ctx)

			if err != nil {
				s.errors <- err
			}

			if message != nil {
				s.messages <- message
			}
		}
	}
}

func (s *Subscribe) NextMessage() <-chan proto.Message {
	return s.messages
}

func (s *Subscribe) NextError() <-chan error {
	return s.errors
}
