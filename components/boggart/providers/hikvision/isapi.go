package hikvision

import (
	"net"
	"net/http"
	"strconv"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
)

type ISAPI struct {
	connection *connection.Client
	username   string
	password   string
	address    string
}

func NewISAPI(host string, port int64, username, password string) *ISAPI {
	return &ISAPI{
		connection: connection.NewClient(),
		username:   username,
		password:   password,
		address:    "http://" + net.JoinHostPort(host, strconv.FormatInt(port, 10)) + "/ISAPI/",
	}
}

func (a *ISAPI) do(request *http.Request) (*http.Response, error) {
	if _, _, ok := request.BasicAuth(); !ok {
		request.SetBasicAuth(a.username, a.password)
	}

	return a.connection.Do(request)
}
