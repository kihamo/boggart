package herospeed

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
)

type Client struct {
	connection *connection.Client
	username   string
	password   string
	address    string
}

func New(host string, port int64, username, password string) *Client {
	return &Client{
		connection: connection.NewClient(),
		username:   username,
		password:   password,
		address:    "http://" + net.JoinHostPort(host, strconv.FormatInt(port, 10)) + "/",
	}
}

func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return c.DoRequest(ctx, request)
}

func (c *Client) DoRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	if _, _, ok := request.BasicAuth(); !ok {
		request.SetBasicAuth(c.username, c.password)
	}

	return c.connection.Do(request.WithContext(ctx))
}

func (c *Client) Configuration(ctx context.Context) (map[string]string, error) {
	response, err := c.Do(ctx, http.MethodGet, c.address+"ini.htm", nil)
	if err != nil {
		return nil, err
	}

	separator := []byte("<br>")
	separatorLen := len(separator)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		i := bytes.IndexByte(data, separator[0])
		if i < 0 {
			return len(data), data, nil
		}

		if len(data[i:]) >= separatorLen {
			if bytes.Equal(data[i:i+separatorLen], separator) {
				return i + separatorLen, data[:i], nil
			}
		}

		return len(data), data, nil
	}

	scanner := bufio.NewScanner(response.Body)
	scanner.Split(split)

	options := make(map[string]string)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")

		if len(parts) > 1 {
			options[parts[0]] = strings.Join(parts[1:], "=")
		} else {
			options[parts[0]] = ""
		}
	}

	return options, scanner.Err()
}

// FIXME: скорее всего однопоточная реализация, нужен мьютекс
func (c *Client) Snapshot(ctx context.Context, writer io.Writer) error {
	response, err := c.Do(ctx, http.MethodGet, c.address+"snap.jpg", nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, response.Body)
	response.Body.Close()

	return err
}
