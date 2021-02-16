package fcm

import (
	"context"
	"errors"

	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/kihamo/boggart/components/boggart/di"
	"go.uber.org/multierr"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	messaging *messaging.Client
	projectID string
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	ctx := context.Background()
	cfg := b.config()

	credentialsJSON := []byte(cfg.Credentials)
	credentials, err := google.CredentialsFromJSON(ctx, credentialsJSON)
	if err != nil {
		return err
	}

	b.projectID = credentials.ProjectID

	opts := option.WithCredentialsJSON(credentialsJSON)

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return err
	}

	b.messaging, err = app.Messaging(ctx)

	return err
}

func (b *Bind) Send(ctx context.Context, text string) (err error) {
	if len(text) == 0 {
		return errors.New("text is empty")
	}

	for _, token := range b.config().Tokens {
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
