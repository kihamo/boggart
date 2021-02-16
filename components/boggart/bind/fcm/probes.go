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

	if tokens := b.config().Tokens; len(tokens) > 0 {
		message.Token = tokens[0]
	}

	_, err := b.messaging.SendDryRun(ctx, message)

	return err
}
