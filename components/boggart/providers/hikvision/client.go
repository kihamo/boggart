package hikvision

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

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/client"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/client/event"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision/models"
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

	alerts chan *models.EventNotificationAlert
	errors chan error
	done   chan struct{}
}

var (
	eventTagAlertStart = []byte("<EventNotificationAlert")
	eventTagAlertEnd   = []byte("</EventNotificationAlert>")
)

func New(address string, user, password string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
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

	// принудительно отключаем дебаг, иначе его вывод заблочит чтение стрима
	// хотя это не потокобезопасно
	if t, ok := s.client.Transport.(*httptransport.Runtime); ok {
		if t.Debug {
			t.SetDebug(false)
		}
	}

	params := event.NewGetNotificationAlertStreamParams().
		WithContext(s.ctx).
		WithTimeout(0)

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		_, err := s.client.Event.GetNotificationAlertStream(params, nil, pw)
		if err != nil {
			s.errors <- err
			s.done <- struct{}{}
		}
	}()

	go s.loop(pr)
}

func (s *AlertStreaming) loop(pr *io.PipeReader) {
	reader := bufio.NewReader(pr)
	parseBuffer := bytes.NewBuffer(nil)

	defer func() {
		parseBuffer.Reset()
		pr.Close()

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

func (s *AlertStreaming) Close() error {
	s.done <- struct{}{}
	return nil
}
