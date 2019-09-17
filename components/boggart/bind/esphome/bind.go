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

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config   *Config
	provider *native_api.Client
}

func (b *Bind) Close() error {
	return b.provider.Close()
}

func (b *Bind) EntityByObjectID(ctx context.Context, objectID string) (proto.Message, error) {
	messages, err := b.provider.ListEntities(ctx)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		if e, ok := message.(native_api.MessageEntity); ok && e.GetObjectId() == objectID {
			return message, nil
		}
	}

	return nil, errors.New("entity with object ID " + objectID + " not found")
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
		e, ok := message.(native_api.MessageState)
		if !ok {
			return nil, errors.New("input message " + proto.MessageName(message) + " not implement MessageState interface")
		}

		entitiesKeys[e.GetKey()] = struct{}{}
	}

	states := make(map[uint32]proto.Message, length)

	for {
		select {
		case message := <-chMessage:
			e, ok := message.(native_api.MessageState)
			if !ok {
				return states, errors.New("output message " + proto.MessageName(message) + " not implement MessageState interface")
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
