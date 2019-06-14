package hikvision

import (
	"net"
	"strconv"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	apiclient "github.com/kihamo/boggart/components/boggart/providers/hikvision2/client"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	port, _ := strconv.ParseInt(config.Address.Port(), 10, 64)

	clientConfig := apiclient.DefaultTransportConfig().WithHost(net.JoinHostPort(config.Address.Hostname(), strconv.FormatInt(port, 10)))
	client := apiclient.NewHTTPClientWithConfig(nil, clientConfig)

	// auth
	password, _ := config.Address.User.Password()
	if rt, ok := client.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(config.Address.User.Username(), password)

		// что бы скачивались файлы с изображениями
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/gif"] = runtime.ByteStreamConsumer()
	}

	bind := &Bind{
		client:                client,
		isapi:                 hikvision.NewISAPI(config.Address.Hostname(), port, config.Address.User.Username(), password),
		alertStreamingHistory: make(map[string]time.Time),
		address:               config.Address.URL,
		config:                config,
	}

	return bind, nil
}
