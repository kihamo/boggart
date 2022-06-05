package softvideo

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/kihamo/boggart/protocols/http"
)

const (
	AccountURL = "https://user.softvideo.ru/"
)

var (
	balanceRegexp = regexp.MustCompile(`(?:badge-success|badge-important|badge-warning)">([\d.-]+)(?:[\s(]+([+\-\d]+)[\s)]+)*[^<]*</span>`)
)

type Client struct {
	connection *http.Client
	login      string
	password   string
}

func New(login, password string, debug bool) *Client {
	return &Client{
		connection: http.NewClient().WithDebug(debug),
		login:      login,
		password:   password,
	}
}

func (c *Client) AccountID() string {
	return c.login
}

func (c *Client) Balance(ctx context.Context) (balance, promise float64, err error) {
	_, err = c.connection.Get(ctx, AccountURL)
	if err != nil {
		return -1, -1, err
	}

	response, err := c.connection.Post(ctx, AccountURL, map[string]string{
		"bootstrap[username]": c.login,
		"bootstrap[password]": c.password,
		"bootstrap[send]:":    "<i class=\"icon-ok\"></i>",
	})

	if err != nil {
		return -1, -1, err
	}

	submatch := balanceRegexp.FindStringSubmatch(http.BodyFromResponse(response))
	if len(submatch) != 3 {
		return -1, -1, errors.New("balance string not found in page")
	}

	balance, err = strconv.ParseFloat(submatch[1], 64)
	if err != nil {
		return -1, -1, err
	}

	if submatch[2] != "" {
		promise, err = strconv.ParseFloat(submatch[2], 64)
		if err != nil {
			return -1, -1, err
		}
	}

	return balance, promise, nil
}
