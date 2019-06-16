package hikvision2

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client"
)

type UploadRoundTripper struct {
	Proxied http.RoundTripper
}

/*
 * Протокол не поддерживает загрузку прошивки через multipart, поэтому парсим реквест,
 * сформированный swagger, получаем нужную часть и преобразуем в ревест понятный устройству
 */
func (rt UploadRoundTripper) modify(req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		return
	}

	d, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return
	}

	if d != "application/x-www-form-urlencoded" {
		return
	}

	var buf bytes.Buffer

	boundary, ok := params["boundary"]
	if !ok {
		return
	}

	mr := multipart.NewReader(req.Body, boundary)

	p, err := mr.NextPart()
	if err != nil {
		return
	}

	if _, err = buf.ReadFrom(p); err != nil {
		return
	}

	if err = req.Body.Close(); err != nil {
		return
	}

	req.ContentLength = int64(buf.Len())
	req.Header.Set("Content-Type", d)
	req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
}

func (rt UploadRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.modify(req)
	return rt.Proxied.RoundTrip(req)
}

func New(host string, port int64, user, password string, debug bool, logger logger.Logger) *client.HikVision {
	cfg := client.DefaultTransportConfig().WithHost(net.JoinHostPort(host, strconv.FormatInt(port, 10)))
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(user, password)
		rt.Transport = UploadRoundTripper{rt.Transport}

		// что бы скачивались файлы с изображениями
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/gif"] = runtime.ByteStreamConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
