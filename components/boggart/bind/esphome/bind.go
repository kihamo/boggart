package esphome

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

const (
	subscribeStateTimeoutByEntity = time.Millisecond * 200
)

type hasKey interface {
	GetKey() uint32
}

type Bind struct {
	boggart.BindBase

	config   *Config
	provider *native_api.Client
}

func (b *Bind) Close() error {
	return b.provider.Close()
}

func (b *Bind) States(ctx context.Context, messages []proto.Message) (map[uint32]proto.Message, error) {
	length := len(messages)

	timeout := subscribeStateTimeoutByEntity*time.Duration(length) + time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	chMessage, err := b.provider.SubscribeStates(ctx)
	if err != nil {
		return nil, err
	}

	entitiesKeys := make(map[uint32]struct{}, length)
	for _, message := range messages {
		e, ok := message.(hasKey)
		if !ok {
			return nil, errors.New("input message " + proto.MessageName(message) + " has not method GetKey() uint32")
		}

		entitiesKeys[e.GetKey()] = struct{}{}
	}

	states := make(map[uint32]proto.Message, length)

	for {
		select {
		case message := <-chMessage:
			e, ok := message.(hasKey)
			if !ok {
				return states, errors.New("output message " + proto.MessageName(message) + " has not method GetKey() uint32")
			}

			if _, ok = entitiesKeys[e.GetKey()]; ok {
				states[e.GetKey()] = message

				length--
				if length == 0 {
					return states, nil
				}
			}

		case <-ctx.Done():
			return states, ctx.Err()
		}
	}
}
