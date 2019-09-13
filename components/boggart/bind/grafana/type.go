package grafana

import (
	"github.com/kihamo/boggart/components/boggart"
	g "github.com/kihamo/go-grafana-api"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicAnnotation = config.TopicAnnotation.Format(config.Name)

	client := g.New(config.Address.String())

	if config.ApiKey != "" {
		client = client.WithApiKey(config.ApiKey)
	} else {
		client = client.WithBasicAuth(config.Username, config.Password)
	}

	bind := &Bind{
		config: config,
		client: client,
	}

	return bind, nil
}
