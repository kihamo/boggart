package native_api

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	defaultTimeout = time.Second * 5
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
	ticker := time.NewTicker(s.tickerInterval)
	defer func() {
		fmt.Println("Close")

		ticker.Stop()
		close(s.messages)
		close(s.errors)
	}()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			message, err := s.client.invoke(s.ctx, s.startMessage)

			fmt.Println("Subscribe read:", message, err)

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
