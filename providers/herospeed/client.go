package herospeed

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
)

type Client struct {
	connection *connection.Client
	username   string
	password   string
	address    string

	snapshotLock sync.Mutex
	snapshotURL  string
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

	separator := "<br>"
	separatorLen := len(separator)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), separator); i >= 0 {
			return i + separatorLen, data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return
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

func (c *Client) Snapshot(ctx context.Context, writer io.Writer) error {
	// камера криво поддерживает множество потоков, поэтому организуем мьютекс
	c.snapshotLock.Lock()
	defer c.snapshotLock.Unlock()

	u, err := c.lazyLoadClickSnapStorage(ctx)
	if err != nil {
		return err
	}

	response, err := c.Do(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, response.Body)
	response.Body.Close()

	return err
}

func (c *Client) lazyLoadClickSnapStorage(ctx context.Context) (string, error) {
	if c.snapshotURL == "" {
		cfg, err := c.Configuration(ctx)
		if err != nil {
			return "", err
		}

		clickSnapStorage, ok := cfg["clicksnapstorage"]
		if !ok {
			return "", errors.New("clicksnapstorage option isn't exists")
		}

		result, err := url.Parse(clickSnapStorage)
		if err != nil {
			return "", err
		}

		api, err := url.Parse(c.address)
		if err != nil {
			return "", err
		}

		result.Host = api.Host

		c.snapshotURL = result.String()
	}

	return c.snapshotURL, nil
}
