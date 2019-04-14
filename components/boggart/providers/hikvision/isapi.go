package hikvision

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
	tracing "github.com/kihamo/shadow/components/tracing/http"
)

const (
	ComponentName = "hikvision"

	SubStatusCodeRebootRequired = "rebootRequired"
)

type ResponseStatus struct {
	StatusCode    uint64 `xml:"statusCode"`
	StatusString  string `xml:"statusString"`
	SubStatusCode string `xml:"subStatusCode"`
}

type ISAPI struct {
	connection *connection.Client
	username   string
	password   string
	address    string
}

func NewISAPI(host string, port int64, username, password string) *ISAPI {
	return &ISAPI{
		connection: connection.NewClient().WithTimeout(0),
		username:   username,
		password:   password,
		address:    "http://" + net.JoinHostPort(host, strconv.FormatInt(port, 10)) + "/ISAPI",
	}
}

func (a *ISAPI) Do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return a.DoRequest(ctx, request)
}

func (a *ISAPI) DoRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	ctx = tracing.ComponentNameToContext(ctx, ComponentName)

	if _, _, ok := request.BasicAuth(); !ok {
		request.SetBasicAuth(a.username, a.password)
	}

	return a.connection.Do(request.WithContext(ctx))
}

func (a *ISAPI) DoXML(ctx context.Context, method, url string, body, v interface{}) error {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	return a.DoRequestXML(ctx, request, body, v)
}

func (a *ISAPI) DoRequestXML(ctx context.Context, request *http.Request, body, v interface{}) error {
	if body != nil {
		data, err := xml.Marshal(body)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)

		request.Header.Set("Content-Type", `application/xml; charset="UTF-8"`)
		request.ContentLength = int64(buf.Len())
		request.Body = ioutil.NopCloser(buf)
	}

	response, err := a.DoRequest(ctx, request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if v == nil {
		return nil
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(content, &v)
}
