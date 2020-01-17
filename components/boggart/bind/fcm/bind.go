package fcm

import (
	"context"
	"errors"

	"firebase.google.com/go/messaging"
	"github.com/kihamo/boggart/components/boggart/di"
	"go.uber.org/multierr"
)

type Bind struct {
	di.MQTTBind
	di.ProbesBind

	config    *Config
	messaging *messaging.Client
}

func (b *Bind) Send(ctx context.Context, text string) (err error) {
	if len(text) == 0 {
		return errors.New("text is empty")
	}

	for _, token := range b.config.Tokens {
		if e := b.sendByToken(ctx, token, text); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}

func (b *Bind) sendByToken(ctx context.Context, token, text string) error {
	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Body: text,
			},
		},
		Token: token,
	}

	_, err := b.messaging.Send(ctx, message)
	return err
}
