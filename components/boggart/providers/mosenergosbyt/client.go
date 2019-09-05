package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"sync"
	"time"

	connection "github.com/kihamo/boggart/components/boggart/protocols/http"
)

const (
	BaseURL = "https://my.mosenergosbyt.ru/gate_mlkcomu"
)

type Client struct {
	connection *connection.Client

	login       string
	password    string
	session     string
	sessionLock sync.RWMutex
}

func New(login, password string) *Client {
	return &Client{
		connection: connection.NewClient(), //.WithDebug(true),
		login:      login,
		password:   password,
	}
}

func (c *Client) doRequest(ctx context.Context, action, query string, body map[string]string, data interface{}) error {
	values := url.Values{}
	values.Add("action", action)

	if query != "" {
		values.Add("query", query)
	}

	// auto login
	if action != "auth" {
		c.sessionLock.RLock()
		session := c.session
		c.sessionLock.RUnlock()

		if session == "" {
			if err := c.Auth(ctx); err != nil {
				return err
			}

			c.sessionLock.RLock()
			session = c.session
			c.sessionLock.RUnlock()
		}

		values.Add("session", session)
	}

	resp, err := c.connection.Post(ctx, BaseURL+"?"+values.Encode(), body)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := &response{}
	if data != nil {
		r.Data = data
	}

	err = json.Unmarshal(b, r)

	if err == nil && r.ErrorMessage != "" {
		err = errors.New(r.ErrorMessage)
	}

	return err
}

func (c *Client) Auth(ctx context.Context) error {
	data := make([]struct {
		KDResult  uint64 `json:"kd_result"`
		NMResult  string `json:"nm_result"`
		IDProfile string `json:"id_profile"`
		AuthCount uint64 `json:"cnt_auth"`
		Token     string `json:"new_token"`
		Session   string `json:"session"`
	}, 0)

	err := c.doRequest(ctx, "auth", "", map[string]string{
		"login":          c.login,
		"psw":            c.password,
		"remember":       "true",
		"vl_device_info": "",
	}, &data)

	if data[0].Session == "" {
		return errors.New(data[0].NMResult)
	}

	c.sessionLock.Lock()
	c.session = data[0].Session
	c.sessionLock.Unlock()

	return err
}

func (c *Client) Accounts(ctx context.Context) ([]Account, error) {
	data := make([]Account, 0)

	if err := c.doRequest(ctx, "sql", "LSList", nil, &data); err != nil {
		return nil, err
	}

	for i, account := range data {
		if account.VLProviderRAW != "" {
			if err := json.Unmarshal([]byte(account.VLProviderRAW), &account.VLProvider); err != nil {
				return nil, err
			}

			data[i].VLProvider = account.VLProvider
		}
	}

	return data, nil
}

func (c *Client) CurrentBalance(ctx context.Context, IDAbonent uint64) (*Balance, error) {
	var data []Balance

	err := c.doRequest(ctx, "sql", "smorodinaTransProxy", map[string]string{
		"plugin":      "smorodinaTransProxy",
		"proxyquery":  "AbonentCurrentBalance",
		"vl_provider": `{"id_abonent": ` + strconv.FormatUint(IDAbonent, 10) + `}`,
	}, &data)

	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func (c *Client) Payments(ctx context.Context, IDAbonent uint64, dateStart, dateEnd time.Time) ([]Payment, error) {
	var data []Payment

	err := c.doRequest(ctx, "sql", "smorodinaTransProxy", map[string]string{
		"plugin":      "smorodinaTransProxy",
		"proxyquery":  "AbonentPays",
		"vl_provider": `{"id_abonent": ` + strconv.FormatUint(IDAbonent, 10) + `}`,
		"dt_st":       dateStart.Format(time.RFC3339),
		"dt_en":       dateEnd.Format(time.RFC3339),
	}, &data)

	return data, err
}

func (c *Client) ChargeDetail(ctx context.Context, IDAbonent uint64, dateStart, dateEnd time.Time) ([]Charge, error) {
	var data []Charge

	err := c.doRequest(ctx, "sql", "smorodinaTransProxy", map[string]string{
		"plugin":          "smorodinaTransProxy",
		"proxyquery":      "AbonentChargeDetail",
		"vl_provider":     `{"id_abonent": ` + strconv.FormatUint(IDAbonent, 10) + `}`,
		"dt_period_start": dateStart.Format(time.RFC3339),
		"dt_period_end":   dateEnd.Format(time.RFC3339),
		"kd_tp_mode":      "1",
	}, &data)

	return data, err
}
