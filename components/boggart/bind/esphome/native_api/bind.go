package native_api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/esphome"
	api "github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/boggart/providers/wifiled"
	"go.uber.org/multierr"
)

const (
	subscribeStateTimeoutByEntity = time.Millisecond * 200
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind

	config   *Config
	provider *api.Client
	ota      *esphome.OTA
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
		if e, ok := message.(api.MessageEntity); ok && mqtt.NameReplace(e.GetObjectId()) == objectIDReplace {
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
		e, ok := message.(api.MessageState)
		if !ok {
			return nil, errors.New("input message " + proto.MessageName(message) + " not implement MessageState interface")
		}

		entitiesKeys[e.GetKey()] = struct{}{}
	}

	states := make(map[uint32]proto.Message, length)

	for {
		select {
		case message := <-chMessage:
			e, ok := message.(api.MessageState)
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
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	states, err := b.States(ctx, messages...)
	if err != nil {
		return err
	}

	entities := make(map[uint32]api.MessageEntity)
	for _, message := range messages {
		if e, ok := message.(api.MessageEntity); ok {
			entities[e.GetKey()] = e
		}
	}

	for _, state := range states {
		messageState, ok := state.(api.MessageState)
		if !ok {
			continue
		}

		var entity api.MessageEntity
		entity, ok = entities[messageState.GetKey()]
		if !ok {
			continue
		}

		objectID := entity.GetObjectId()

		if stateName, e := api.State(entity.(proto.Message), state, false); e == nil {
			if e = b.MQTT().PublishAsync(ctx, b.config.TopicState.Format(sn, objectID), stateName); e != nil {
				err = multierr.Append(err, e)
			}

			switch st := state.(type) {
			case *api.LightStateResponse:
				color := wifiled.Color{
					Red:       uint8(st.Red * 255),
					Green:     uint8(st.Green * 255),
					Blue:      uint8(st.Blue * 255),
					WarmWhite: uint8(st.White * 100),
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateBrightness.Format(sn, objectID), st.Brightness); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateRed.Format(sn, objectID), color.Red); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateGreen.Format(sn, objectID), color.Green); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateBlue.Format(sn, objectID), color.Blue); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateWhite.Format(sn, objectID), color.WarmWhite); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateColorTemperature.Format(sn, objectID), uint8(st.ColorTemperature*100)); e != nil {
					err = multierr.Append(err, e)
				}

				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateEffect.Format(sn, objectID), st.Effect); e != nil {
					err = multierr.Append(err, e)
				}

				// in HEX
				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateColorRGB.Format(sn, objectID), color.String()); e != nil {
					err = multierr.Append(err, e)
				}

				// in HSV
				h, s, v := color.HSV()
				if e = b.MQTT().PublishAsync(ctx, b.config.TopicStateColorHSV.Format(sn, objectID), fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
					err = multierr.Append(err, e)
				}
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return err
}
