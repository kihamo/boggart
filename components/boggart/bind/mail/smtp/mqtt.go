package smtp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	index := -1

	for i, value := range b.config.TopicSend.Split() {
		if value == "#" || value == "+" {
			index = i
			break
		}
	}

	if index == -1 {
		return nil
	}

	type sendRequest struct {
		Subject string          `json:"subject"`
		Body    json.RawMessage `json:"body"`
		To      []string        `json:"to"`
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicSend, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
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
	}
}
