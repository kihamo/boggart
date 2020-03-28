package fcm

import (
	"context"

	"firebase.google.com/go/messaging"
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Body: "readiness probe",
			},
		},
	}

	if len(b.config.Tokens) > 0 {
		message.Token = b.config.Tokens[0]
	}

	_, err := b.messaging.SendDryRun(ctx, message)

	return err
}
