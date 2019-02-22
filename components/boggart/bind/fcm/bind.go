package fcm

import (
	"context"

	"firebase.google.com/go/messaging"
	"github.com/kihamo/boggart/components/boggart"
	"go.uber.org/multierr"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	messaging *messaging.Client
	projectID string
	tokens    []string
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
}

func (b *Bind) Send(ctx context.Context, text string) (err error) {
	for _, token := range b.tokens {
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
