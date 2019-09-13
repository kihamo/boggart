package fcm

import (
	"context"

	"firebase.google.com/go"
	"github.com/kihamo/boggart/components/boggart"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

/*
~/openhab-cloud$ mongo 127.0.0.1:27017/openhab --eval "db.userdevices.find()"
*/

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)
	credentialsJSON := []byte(config.Credentials)
	ctx := context.Background()

	credentials, err := google.CredentialsFromJSON(ctx, credentialsJSON)
	if err != nil {
		return nil, err
	}

	opts := option.WithCredentialsJSON(credentialsJSON)

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	messaging, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	config.TopicSend = config.TopicSend.Format(credentials.ProjectID)

	bind := &Bind{
		config:    config,
		messaging: messaging,
	}

	return bind, nil
}
