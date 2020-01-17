package native_api

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
	api "github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/boggart/providers/wifiled"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicState,
		b.config.TopicStateColorRGB,
		b.config.TopicStateColorHSV,
		b.config.TopicStateBrightness,
		b.config.TopicStateRed,
		b.config.TopicStateGreen,
		b.config.TopicStateBlue,
		b.config.TopicStateWhite,
		b.config.TopicStateColorTemperature,
		b.config.TopicStateEffect,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicPower, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
				return nil
			}

			parts := message.Topic().Split()

			entity, err := b.EntityByObjectID(ctx, parts[len(parts)-2])
			if err != nil {
				return err
			}

			switch api.EntityType(entity) {
			case api.EntityTypeFan:
				err = b.provider.FanCommand(ctx, &api.FanCommandRequest{
					Key:      entity.(api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case api.EntityTypeLight:
				err = b.provider.LightCommand(ctx, &api.LightCommandRequest{
					Key:      entity.(api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case api.EntityTypeSwitch:
				err = b.provider.SwitchCommand(ctx, &api.SwitchCommandRequest{
					Key:   entity.(api.MessageEntity).GetKey(),
					State: message.Bool(),
				})
			}

			if err == nil {
				err = b.syncState(ctx, entity)
			}

			return err
		})),
		mqtt.NewSubscriber(b.config.TopicColor, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 3) {
				return nil
			}

			parts := message.Topic().Split()

			entity, err := b.EntityByObjectID(ctx, parts[len(parts)-2])
			if err != nil {
				return err
			}

			if api.EntityType(entity) != api.EntityTypeLight {
				return nil
			}

			states, err := b.States(ctx, entity)
			if err != nil {
				return err
			}

			var (
				state *api.LightStateResponse
				key   = entity.(api.MessageEntity).GetKey()
			)

			if s, ok := states[key]; !ok {
				return errors.New("failed get entity state")
			} else {
				state = s.(*api.LightStateResponse)
			}

			cmd := &api.LightCommandRequest{
				Key:      key,
				HasRgb:   true,
				HasWhite: true,
			}

			color, err := wifiled.ColorFromString(message.String())
			if err == nil {
				if color.UseRGB {
					cmd.Red = float32(color.Red) / 255
					cmd.Green = float32(color.Green) / 255
					cmd.Blue = float32(color.Blue) / 255
					cmd.White = state.GetWhite()
				} else {
					cmd.White = float32(color.WarmWhite) / 100
					cmd.Red = state.GetRed()
					cmd.Green = state.GetGreen()
					cmd.Blue = state.GetBlue()

					if cmd.Red == 0 && cmd.Green == 0 && cmd.Blue == 0 {
						cmd.Red = 1
						cmd.Green = 1
						cmd.Blue = 1
					}
				}
			}

			if err != nil {
				return err
			}

			err = b.provider.LightCommand(ctx, cmd)

			if err == nil {
				err = b.syncState(ctx, entity)
			}

			return err
		})),
		mqtt.NewSubscriber(b.config.TopicStateSet, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), 4) {
				return nil
			}

			parts := message.Topic().Split()

			entity, err := b.EntityByObjectID(ctx, parts[len(parts)-3])
			if err != nil {
				return err
			}

			switch api.EntityType(entity) {
			case api.EntityTypeFan:
				err = b.provider.FanCommand(ctx, &api.FanCommandRequest{
					Key:      entity.(api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case api.EntityTypeLight:
				var (
					state   *api.LightStateResponse
					light   = entity.(api.MessageEntity)
					cmd     = &api.LightCommandRequest{}
					request struct {
						State            *bool   `json:"state,omitempty"`
						Brightness       *uint8  `json:"brightness,omitempty"`
						Red              *uint8  `json:"red,omitempty"`
						Green            *uint8  `json:"green,omitempty"`
						Blue             *uint8  `json:"blue,omitempty"`
						White            *uint8  `json:"white,omitempty"`
						ColorTemperature *uint8  `json:"color_temperature,omitempty"`
						TransitionLength *uint32 `json:"transition_length,omitempty"`
						FlashLength      *uint32 `json:"flash_length,omitempty"`
						Effect           *string `json:"effect,omitempty"`
					}
				)

				states, err := b.States(ctx, entity)
				if err != nil {
					return err
				}

				if s, ok := states[light.GetKey()]; !ok {
					return errors.New("failed get entity state")
				} else {
					state = s.(*api.LightStateResponse)
				}

				if err := message.UnmarshalJSON(&request); err == nil {
					if request.State != nil {
						cmd.HasState = true
						cmd.State = *request.State
					}

					if request.Brightness != nil {
						val := float32(*request.Brightness) / 100

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetBrightness() {
							cmd.HasBrightness = true
							cmd.Brightness = val
						}
					}

					if request.Red != nil {
						val := float32(*request.Red) / 255

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetRed() {
							cmd.HasRgb = true
							cmd.Red = val
						}
					}

					if request.Green != nil {
						val := float32(*request.Green) / 255

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetGreen() {
							cmd.HasRgb = true
							cmd.Green = val
						}
					}

					if request.Blue != nil {
						val := float32(*request.Blue) / 255

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetBlue() {
							cmd.HasRgb = true
							cmd.Blue = val
						}
					}

					if request.White != nil {
						val := float32(*request.White) / 255

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetWhite() {
							cmd.HasWhite = true
							cmd.White = val

							// хак, без установки RGB значение white по факту будет 1
							if cmd.GetRed() == 0 && cmd.GetGreen() == 0 && cmd.GetBlue() == 0 {
								cmd.HasRgb = true
								cmd.Red = 1
								cmd.Green = 1
								cmd.Blue = 1
							}
						}
					}

					if request.ColorTemperature != nil {
						val := float32(*request.ColorTemperature) / 100

						if val > 1 {
							val = 1
						} else if val < 0 {
							val = 0
						}

						if val != state.GetColorTemperature() {
							cmd.HasColorTemperature = true
							cmd.ColorTemperature = val
						}
					}

					if request.TransitionLength != nil {
						val := *request.TransitionLength

						if val > 0 {
							cmd.HasTransitionLength = true
							cmd.TransitionLength = val
						}
					}

					if request.FlashLength != nil {
						val := *request.FlashLength

						if val > 0 {
							cmd.HasFlashLength = true
							cmd.FlashLength = val
						}
					}

					if request.Effect != nil {
						val := *request.Effect
						if val == "" {
							val = "None"
						}

						if val != state.GetEffect() {
							cmd.HasEffect = true
							cmd.Effect = val
						}
					}

				} else {
					cmd.State = message.Bool()
				}

				cmd.HasState = true
				cmd.Key = light.GetKey()

				err = b.provider.LightCommand(ctx, cmd)

			case api.EntityTypeSwitch:
				err = b.provider.SwitchCommand(ctx, &api.SwitchCommandRequest{
					Key:   entity.(api.MessageEntity).GetKey(),
					State: message.Bool(),
				})
			}

			if err == nil {
				err = b.syncState(ctx, entity)
			}

			return err
		})),
	}
}
