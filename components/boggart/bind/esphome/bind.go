package esphome

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/boggart/providers/wifiled"
	"go.uber.org/multierr"
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

	objectIDReplace := mqtt.NameReplace(objectID)

	for _, message := range messages {
		if e, ok := message.(native_api.MessageEntity); ok && mqtt.NameReplace(e.GetObjectId()) == objectIDReplace {
			return message, nil
		}
	}

	return nil, errors.New("entity with object ID " + objectID + " not found")
}

func (b *Bind) States(ctx context.Context, messages ...proto.Message) (map[uint32]proto.Message, error) {
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

func (b *Bind) syncState(ctx context.Context, messages ...proto.Message) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	states, err := b.States(ctx, messages...)
	if err != nil {
		return err
	}

	entities := make(map[uint32]native_api.MessageEntity)
	for _, message := range messages {
		if e, ok := message.(native_api.MessageEntity); ok {
			entities[e.GetKey()] = e
		}
	}

	for _, state := range states {
		messageState, ok := state.(native_api.MessageState)
		if !ok {
			continue
		}

		var entity native_api.MessageEntity
		entity, ok = entities[messageState.GetKey()]
		if !ok {
			continue
		}

		objectID := entity.GetObjectId()

		if stateName, e := native_api.State(entity.(proto.Message), state); e == nil {
			if e = b.MQTTPublishAsync(ctx, b.config.TopicState.Format(sn, objectID), stateName); e != nil {
				err = multierr.Append(err, e)
			}

			switch st := state.(type) {
			case *native_api.LightStateResponse:
				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateBrightness.Format(sn, objectID), st.Brightness); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateRed.Format(sn, objectID), st.Red); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateGreen.Format(sn, objectID), st.Green); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateBlue.Format(sn, objectID), st.Blue); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateWhite.Format(sn, objectID), st.White); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateColorTemperature.Format(sn, objectID), st.ColorTemperature); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateEffect.Format(sn, objectID), st.Effect); e != nil {
					err = multierr.Append(err, e)
				}

				color := wifiled.Color{
					Red:       uint8(st.Red),
					Green:     uint8(st.Green),
					Blue:      uint8(st.Blue),
					WarmWhite: uint8(st.White),
				}

				// in HEX
				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateColor.Format(sn, objectID), color.String()); e != nil {
					err = multierr.Append(err, e)
				}

				// in HSV
				h, s, v := color.HSV()
				if e = b.MQTTPublishAsync(ctx, b.config.TopicStateColorHSV.Format(sn, objectID), fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
					err = multierr.Append(err, e)
				}
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return err
}
