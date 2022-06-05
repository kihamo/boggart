package premiergc

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/protocols/http"
)

const (
	AccountURL = "http://stat.premier-gc.ru/client.php"
)

var (
	contractRegexp = regexp.MustCompile(`Договор:</td>\s*<td><b>([^\s]+[^<]*)</b></td>`)
	balanceRegexp  = regexp.MustCompile(`Баланс:</td>\s*<td><b>([^\s]+[^<]*)</b></td>`)
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

func (c *Client) Balance(ctx context.Context) (contract string, balance float64, err error) {
	_, err = c.connection.Get(ctx, AccountURL)
	if err != nil {
		return contract, balance, err
	}

	response, err := c.connection.Post(ctx, AccountURL, map[string]string{
		"ContractNumber": c.login,
		"Password":       c.password,
	})

	if err != nil {
		return contract, balance, err
	}

	body := http.BodyFromResponse(response)

	submatch := contractRegexp.FindStringSubmatch(body)
	if len(submatch) != 2 {
		return contract, balance, errors.New("balance string not found in page")
	}

	contract = submatch[1]

	submatch = balanceRegexp.FindStringSubmatch(body)
	if len(submatch) != 2 {
		return contract, balance, errors.New("balance string not found in page")
	}

	balance, err = strconv.ParseFloat(strings.ReplaceAll(submatch[1], ",", "."), 64)
	if err != nil {
		return contract, balance, err
	}

	return contract, balance, err
}
