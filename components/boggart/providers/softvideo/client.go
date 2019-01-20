package softvideo

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/protocols/http"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

const (
	AccountURL    = "https://user.softvideo.ru/"
	ComponentName = "softvideo"
)

var (
	balanceRegexp = regexp.MustCompile(`(?:badge-success|badge-important)">([^\s]+)[^<]*</span>`)
)

type Client struct {
	connection *http.Client
	login      string
	password   string
}

func NewClient(login, password string, debug bool) *Client {
	return &Client{
		connection: http.NewClient().WithDebug(debug),
		login:      login,
		password:   password,
	}
}

func (c *Client) AccountID() string {
	return c.login
}

func (c *Client) Balance(ctx context.Context) (float64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "balance")
	defer span.Finish()

	_, err := c.connection.Get(ctx, AccountURL)
	if err != nil {
		tracing.SpanError(span, err)
		return -1, err
	}

	response, err := c.connection.Post(ctx, AccountURL, map[string]string{
		"bootstrap[username]": c.login,
		"bootstrap[password]": c.password,
		"bootstrap[send]:":    "<i class=\"icon-ok\"></i>",
	})

	if err != nil {
		tracing.SpanError(span, err)
		return -1, err
	}

	submatch := balanceRegexp.FindStringSubmatch(http.BodyFromResponse(response))
	if len(submatch) != 2 {
		err := errors.New("balance string not found in page")

		tracing.SpanError(span, err)
		return -1, err
	}

	balance, err := strconv.ParseFloat(submatch[1], 10)
	if err != nil {
		tracing.SpanError(span, err)
		return -1, err
	}

	span.LogFields(log.Float64("balance", balance))

	return balance, nil
}
