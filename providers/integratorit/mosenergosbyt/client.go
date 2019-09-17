package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/internal"
)

type Client struct {
	base *internal.Client
}

func New(login, password string) *Client {
	c := &Client{}
	c.base = internal.New("https://my.mosenergosbyt.ru/gate_mlkcomu", login, password, c.Auth)

	return c
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

	err := c.base.DoRequest(ctx, "auth", nil, map[string]string{
		"login":          c.base.Login(),
		"psw":            c.base.Password(),
		"remember":       "true",
		"vl_device_info": "",
	}, &data)

	if err != nil {
		return err
	}

	if data[0].Session == "" {
		return errors.New(data[0].NMResult)
	}

	c.base.SetSession(data[0].Session)

	return err
}

func (c *Client) Accounts(ctx context.Context) ([]Account, error) {
	data := make([]Account, 0)

	if err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "LSList"}, nil, &data); err != nil {
		return nil, err
	}

	for i, account := range data {
		if account.ProviderRAW != "" {
			if err := json.Unmarshal([]byte(account.ProviderRAW), &account.Provider); err != nil {
				return nil, err
			}

			data[i].Provider = account.Provider
		}
	}

	return data, nil
}

func (c *Client) CurrentBalance(ctx context.Context, IDAbonent uint64) (*Balance, error) {
	var data []Balance

	err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "smorodinaTransProxy"}, map[string]string{
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

	err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "smorodinaTransProxy"}, map[string]string{
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

	err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "smorodinaTransProxy"}, map[string]string{
		"plugin":          "smorodinaTransProxy",
		"proxyquery":      "AbonentChargeDetail",
		"vl_provider":     `{"id_abonent": ` + strconv.FormatUint(IDAbonent, 10) + `}`,
		"dt_period_start": dateStart.Format(time.RFC3339),
		"dt_period_end":   dateEnd.Format(time.RFC3339),
		"kd_tp_mode":      "1",
	}, &data)

	return data, err
}
