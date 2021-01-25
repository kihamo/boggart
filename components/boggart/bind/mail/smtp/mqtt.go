package smtp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	type sendRequest struct {
		Subject string          `json:"subject"`
		Body    json.RawMessage `json:"body"`
		To      []string        `json:"to"`
	}

	subscribers := make([]mqtt.Subscriber, 0, 2)
	subscribers = append(subscribers,
		mqtt.NewSubscriber(b.config.TopicSend, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !b.MQTT().CheckBindIDInTopic(message.Topic(), 3) {
				return nil
			}

			var request sendRequest

			if err := message.JSONUnmarshal(&request); err != nil {
				return err
			}

			return b.Send(request.To, request.Subject, request.Body)
		}),
	)

	index := -1

	for i, value := range b.config.TopicSendMulti.Format(b.Meta().ID()).Split() {
		if value == "#" || value == "+" {
			index = i
			break
		}
	}

	if index != -1 {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(b.config.TopicSendMulti, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				if !b.MQTT().CheckBindIDInTopic(message.Topic(), 3) {
					return nil
				}

				parts := message.Topic().Split()
				if len(parts) < index+1 {
					return errors.New("emails is empty")
				}

				var (
					request sendRequest
					subject string
					body    []byte
					to      []string
				)

				if err := message.JSONUnmarshal(&request); err == nil {
					subject = request.Subject
					body = bytes.Trim(request.Body, `"'`)
					to = request.To
				} else {
					subject = "New email"
					body = message.Payload()
				}

				if len(to) == 0 {
					to = parts[index:]
				}

				return b.Send(to, subject, body)
			}),
		)
	}

	return subscribers
}
