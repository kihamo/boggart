package hikvision

import (
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
)

type Connection interface {
	Do(request *http.Request) (*http.Response, error)
}

type overrideFloat64 struct {
	value float64
}

/**
 *Проблемы с \n в ответе от видео регистратора
 */
func (f *overrideFloat64) Float64() float64 {
	return f.value
}

func (f *overrideFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	d.DecodeElement(&raw, &start)

	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return err
	}

	*f = overrideFloat64{value}

	return nil
}

type ISAPI struct {
	connection Connection
	username   string
	password   string
	address    string
}

func NewISAPI(host string, port int64, username, password string) *ISAPI {
	return &ISAPI{
		connection: connection.NewClient(),
		username:   username,
		password:   password,
		address:    "http://" + net.JoinHostPort(host, strconv.FormatInt(port, 10)) + "/ISAPI",
	}
}

func (a *ISAPI) Do(request *http.Request) (*http.Response, error) {
	if _, _, ok := request.BasicAuth(); !ok {
		request.SetBasicAuth(a.username, a.password)
	}

	return a.connection.Do(request)
}

func (a *ISAPI) DoAndParse(request *http.Request, v interface{}) error {
	response, err := a.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return xml.Unmarshal(content, &v)
}
