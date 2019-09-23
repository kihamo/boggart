package esphome

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/esphome/native_api"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicState,
		b.config.TopicStateColor,
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
		mqtt.NewSubscriber(b.config.TopicStateSet, 0, boggart.WrapMQTTSubscribeDeviceIsOnline(b.Status, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 4) {
				return nil
			}

			parts := message.Topic().Split()

			entity, err := b.EntityByObjectID(ctx, parts[len(parts)-3])
			if err != nil {
				return err
			}

			switch native_api.EntityType(entity) {
			case native_api.EntityTypeFan:
				err = b.provider.FanCommand(ctx, &native_api.FanCommandRequest{
					Key:      entity.(native_api.MessageEntity).GetKey(),
					HasState: true,
					State:    message.Bool(),
				})

			case native_api.EntityTypeLight:
				var (
					state   *native_api.LightStateResponse
					light   = entity.(native_api.MessageEntity)
					cmd     = &native_api.LightCommandRequest{}
					request struct {
						State            *bool    `json:"state,omitempty"`
						Brightness       *float32 `json:"brightness,omitempty"`
						Red              *float32 `json:"red,omitempty"`
						Green            *float32 `json:"green,omitempty"`
						Blue             *float32 `json:"blue,omitempty"`
						White            *float32 `json:"white,omitempty"`
						ColorTemperature *float32 `json:"color_temperature,omitempty"`
						TransitionLength *uint32  `json:"transition_length,omitempty"`
						FlashLength      *uint32  `json:"flash_length,omitempty"`
						Effect           *string  `json:"effect,omitempty"`
					}
				)

				states, err := b.States(ctx, entity)
				if err != nil {
					return err
				}

				if s, ok := states[light.GetKey()]; !ok {
					return errors.New("failed get entity state")
				} else {
					state = s.(*native_api.LightStateResponse)
				}

				if err := message.UnmarshalJSON(&request); err == nil {
					if request.State != nil {
						cmd.HasState = true
						cmd.State = *request.State
					}

					if request.Brightness != nil {
						val := *request.Brightness

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
						val := *request.Red

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
						val := *request.Green

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
						val := *request.Blue

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
						val := *request.White

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
						val := *request.ColorTemperature

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

			case native_api.EntityTypeSwitch:
				err = b.provider.SwitchCommand(ctx, &native_api.SwitchCommandRequest{
					Key:   entity.(native_api.MessageEntity).GetKey(),
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
