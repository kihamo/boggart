package openweathermap

import (
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/openweathermap/client"
)

type Client struct {
	*client.OpenWeather
}

func New(apiKey string, price int, debug bool, logger logger.Logger) *Client {
	cl := &Client{
		OpenWeather: client.Default,
	}

	if rt, ok := cl.OpenWeather.Transport.(*httptransport.Runtime); ok {
		rl := NewLimiter(apiKey, price)

		// если установлено ограничение, то таймауты в реквестах надо отменить по-умолчанию,
		// так как не известно когда запрос будет доступен для выполнения. Если оставить дефолтный
		// в 30 секунд, то лимитатор сразу отключит такие запросы, так как поймет что дедлайн
		// сработает быстрее и не смысла держать такой запрос. При этом явно установленные таймуты
		// в параметрах запроса продолжат работать как ожидается
		if interval := rl.LimitInterval(); interval > 0 {
			rt = httptransport.NewWithClient(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes, &http.Client{
				Transport: rt.Transport,
				Jar:       rt.Jar,
				Timeout:   0,
			})

			cl.OpenWeather.SetTransport(rt)
		}

		rt.DefaultAuthentication = httptransport.APIKeyAuth("appid", "query", apiKey)

		rt.Transport = RoundTripper{
			original: rt.Transport,
			limiter:  rl,
		}

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func Icon(id string) *url.URL {
	u := &url.URL{}
	u.Host = "openweathermap.org"
	u.Path = "/img/w/" + id + ".png"

	return u
}
