package proverkacheka

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/components/boggart/providers/nalog.ru/proverkacheka/client"
)

type Client struct {
	*client.NalogRU
}

func New(user, password string, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().WithHost("proverkacheka.nalog.ru:9999")
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(user, password)

		// что бы скачивались файлы с изображениями
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/gif"] = runtime.ByteStreamConsumer()

		// для alert stream
		rt.Consumers["multipart/mixed"] = runtime.ByteStreamConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		NalogRU: cl,
	}
}
