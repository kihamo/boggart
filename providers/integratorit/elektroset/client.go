package elektroset

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/internal"
	"go.uber.org/multierr"
)

const (
	dateFormatLayout = "2006-01-02"
)

type Client struct {
	base *internal.Client
}

func New(login, password string) *Client {
	c := &Client{}
	c.base = internal.New("https://lkk.oao-elektroset.ru/gate_lkk_myt", login, password, c.Auth)

	return c
}

func (c *Client) Auth(ctx context.Context) error {
	data := make([]struct {
		KDResult  uint64 `json:"kd_result"`
		NMResult  string `json:"nm_result"`
		IDProfile uint64 `json:"id_profile"`
		AuthCount uint64 `json:"cnt_auth"`
		Session   string `json:"session"`
	}, 0)

	hash := md5.New()
	hash.Write([]byte(c.base.Password()))

	err := c.base.DoRequest(ctx, "auth", nil, map[string]string{
		"nm_email":         c.base.Login(),
		"nm_psw":           hex.EncodeToString(hash.Sum(nil)),
		"plugin":           "captchaChecker",
		"captcha_response": "",
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
	data := make([]struct {
		Accounts []Account `json:"houses"`
	}, 0)

	err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "lka_get_houses"}, nil, &data)
	if err != nil {
		return nil, err
	}

	accounts := make([]Account, 0)
	var e error

	for _, row := range data {
		for _, account := range row.Accounts {
			for i, service := range account.Services {
				account.Services[i].Balance, e = strconv.ParseFloat(service.BalanceRAW, 10)
				if e != nil {
					err = multierr.Append(err, e)
				}

				account.Services[i].NNAccount, e = strconv.ParseUint(service.NNAccountRAW, 10, 64)
				if e != nil {
					err = multierr.Append(err, e)
				}
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, err
}

func (c *Client) BalanceDetails(ctx context.Context, number string, provider Provider, dateStart, dateEnd time.Time) ([]BalanceDetail, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}

	data := make([]BalanceDetail, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_balance_detail"}, map[string]string{
		"dt_start":    dateStart.Format(dateFormatLayout),
		"dt_end":      dateEnd.Format(dateFormatLayout),
		"nn_ls":       number,
		"nm_email":    c.base.Login(),
		"vl_provider": string(vlProvider),
		"page":        "1",
		"start":       "0",
		"limit":       "25",
	}, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
