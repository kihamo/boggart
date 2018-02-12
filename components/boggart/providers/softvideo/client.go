package softvideo

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/protocols/http"
)

const (
	AccountURL = "https://user.softvideo.ru/"
)

var (
	balanceRegexp = regexp.MustCompile(`badge-success">([^\s]+)[^<]*</span>`)
)

type Client struct {
	connection *http.Client
	login      string
	password   string
}

func NewClient(login, password string) *Client {
	return &Client{
		connection: http.NewClient(),
		login:      login,
		password:   password,
	}
}

func (c *Client) AccountID() string {
	return c.login
}

func (c *Client) Balance(ctx context.Context) (float64, error) {
	_, err := c.connection.Get(ctx, AccountURL)
	if err != nil {
		return -1, err
	}

	response, err := c.connection.Post(ctx, AccountURL, map[string]string{
		"bootstrap[username]": c.login,
		"bootstrap[password]": c.password,
		"bootstrap[send]:":    "<i class=\"icon-ok\"></i>",
	})

	if err != nil {
		return -1, err
	}

	submatch := balanceRegexp.FindStringSubmatch(http.BodyFromResponse(response))
	if len(submatch) != 2 {
		return -1, errors.New("Balance string not found in page")
	}

	return strconv.ParseFloat(submatch[1], 10)
}
