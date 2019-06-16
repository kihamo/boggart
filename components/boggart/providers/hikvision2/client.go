package hikvision2

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"io"
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
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/event"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/models"
)

type Client struct {
	*client.HikVision
}

type UploadRoundTripper struct {
	Proxied http.RoundTripper
}

type AlertStreaming struct {
	ctx context.Context

	client *Client
	buffer *bytes.Buffer

	alerts chan *models.EventNotificationAlert
	errors chan error
	done   chan struct{}
}

var (
	eventTagAlertStart = []byte("<EventNotificationAlert")
	eventTagAlertEnd   = []byte("</EventNotificationAlert>")
)

func New(host string, port int64, user, password string, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().WithHost(net.JoinHostPort(host, strconv.FormatInt(port, 10)))
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(user, password)
		rt.Transport = UploadRoundTripper{rt.Transport}

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
		HikVision: cl,
	}
}

func (c *Client) EventNotificationAlertStream(ctx context.Context) *AlertStreaming {
	s := &AlertStreaming{
		ctx:    ctx,
		client: c,
		buffer: bytes.NewBuffer(nil),
	}
	s.start()

	return s
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

func (s *AlertStreaming) start() {
	s.done = make(chan struct{}, 1)
	s.alerts = make(chan *models.EventNotificationAlert)
	s.errors = make(chan error)

	params := event.NewGetNotificationAlertStreamParamsWithContext(s.ctx)
	params.SetTimeout(0)

	go func() {
		_, err := s.client.Event.GetNotificationAlertStream(params, nil, s.buffer)
		if err != nil {
			s.errors <- err
			s.done <- struct{}{}
		}
	}()

	go s.loop()
}

func (s *AlertStreaming) loop() {
	reader := bufio.NewReader(s.buffer)
	parseBuffer := bytes.NewBuffer(nil)

	defer func() {
		s.buffer.Reset()
		parseBuffer.Reset()

		//close(s.alerts)
		//close(s.errors)
		close(s.done)
	}()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-s.done:
			return

		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF || len(line) == 0 {
				continue
			}

			if err != nil {
				s.errors <- err
				continue
			}

			if bytes.HasPrefix(line, eventTagAlertStart) {
				parseBuffer.Write(line)
				continue
			} else if parseBuffer.Len() == 0 { // если сообщение не началось игнорируем весь контент
				continue
			}

			parseBuffer.Write(line)

			// если сообщение заканчивается запускаем алгоритм
			if !bytes.HasPrefix(line, eventTagAlertEnd) {
				continue
			}

			event := &models.EventNotificationAlert{}
			err = xml.Unmarshal(parseBuffer.Bytes(), event)

			parseBuffer.Reset()

			go func() {
				if err != nil {
					s.errors <- err
				} else {
					s.alerts <- event
				}
			}()
		}
	}
}

func (s *AlertStreaming) NextAlert() <-chan *models.EventNotificationAlert {
	return s.alerts
}

func (s *AlertStreaming) NextError() <-chan error {
	return s.errors
}
