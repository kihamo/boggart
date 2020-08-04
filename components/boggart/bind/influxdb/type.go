package influxdb

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	authToken := config.AuthToken
	if authToken == "" && config.DSN.User != nil {
		authToken = config.DSN.User.String()
	}

	scheme := config.DSN.Scheme
	if scheme == "" {
		scheme = "http"
	}

	bind := &Bind{
		config: config,
		client: influxdb2.NewClient(scheme+"://"+config.DSN.Host, authToken),
	}

	return bind, nil
}
