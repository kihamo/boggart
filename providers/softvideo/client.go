package softvideo

import (
	"context"
	"errors"
	"regexp"
	"strconv"

	"github.com/kihamo/boggart/protocols/http"
)

const (
	BaseURL  = "https://user.softvideo.ru/"
	LoginURL = BaseURL + "login"
)

var (
	loginFormTokenRegexp = regexp.MustCompile(`"_token".*?value="(\S+?)"`)
	balanceRegexp        = regexp.MustCompile(`panel-right__balance[^\d]+<span[^\d]+([^\s]+)`)
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
	responseLoginForm, err := c.connection.Get(ctx, BaseURL)
	if err != nil {
		return -1, -1, err
	}

	pageBodyAsString := http.BodyFromResponse(responseLoginForm)

	// если обнаружили форму авторизации, то проводим ее, форма прекращает появляться на
	// втором запросе так как срабатывают куки и сразу появляется форма личного кабинета
	s := loginFormTokenRegexp.FindStringSubmatch(pageBodyAsString)
	if len(s) == 2 {
		responseBalancePage, err := c.connection.Post(ctx, LoginURL, map[string]string{
			"type":     "contract",
			"_token":   s[1],
			"contract": c.login,
			"password": c.password,
		})

		if err != nil {
			return -1, -1, err
		}

		pageBodyAsString = http.BodyFromResponse(responseBalancePage)
	}

	s = balanceRegexp.FindStringSubmatch(pageBodyAsString)
	if len(s) != 2 {
		return -1, -1, errors.New("balance string not found")
	}

	balance, err = strconv.ParseFloat(s[1], 64)
	if err != nil {
		return -1, -1, err
	}

	// TODO: проверить доверительный платеж
	//if submatch[2] != "" {
	//	promise, err = strconv.ParseFloat(submatch[2], 64)
	//	if err != nil {
	//		return -1, -1, err
	//	}
	//}

	return balance, promise, nil
}
