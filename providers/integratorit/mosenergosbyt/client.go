package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/kihamo/boggart/providers/integratorit/internal"
)

const (
	BaseURL    = "https://my.mosenergosbyt.ru/gate_lkcomu"
	DeviceInfo = `{"appver":"1.17.1","type":"browser","userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"}`
)

type Client struct {
	base *internal.Client
}

func New(login, password string) *Client {
	c := &Client{}
	c.base = internal.New(BaseURL, login, password, c.Auth)

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

	err := c.base.DoRequest(ctx, "auth", map[string]string{
		"query": "login",
	}, map[string]string{
		"login":          c.base.Login(),
		"psw":            c.base.Password(),
		"remember":       "true",
		"vl_device_info": DeviceInfo,
	}, &data)

	if err != nil {
		return err
	}

	// тут просит аттрибуты, пропускаем этот шаг
	// {"success":true,"total":1,"data":[{"kd_result":192,"nm_result":"Необходимо дозаполнить атрибуты профиля","id_profile":"***","cnt_auth":null,"new_token":"***"}],"metaData":{"responseTime":0.033}}

	err = c.base.DoRequest(ctx, "auth", map[string]string{
		"query": "login",
	}, map[string]string{
		"login":          c.base.Login(),
		"psw_token":      data[0].Token,
		"remember":       "true",
		"vl_device_info": DeviceInfo,
	}, &data)

	if data[0].Session == "" {
		return errors.New(data[0].NMResult)
	}

	c.base.SetSession(data[0].Session)

	return err
}

func (c *Client) Accounts(ctx context.Context) ([]Account, error) {
	data := make([]Account, 0)

	if err := c.base.DoRequest(ctx, "sql", map[string]string{
		"query": "LSList",
	}, nil, &data); err != nil {
		return nil, err
	}

	for i, account := range data {
		if account.ProviderRAW != "" {
			if err := json.Unmarshal([]byte(account.ProviderRAW), &account.Provider); err != nil {
				return nil, multierror.Append(err, errors.New("bad accounts content "+account.ProviderRAW))
			}

			data[i].Provider = account.Provider
		}
	}

	return data, nil
}

func (c *Client) CurrentBalance(ctx context.Context, account *Account) (float64, error) {
	provider, err := NewProvider(c.base, account)
	if err != nil {
		return 0, err
	}

	support, ok := provider.(HasSupportCurrentBalance)
	if !ok {
		return 0, ErrProviderMethodNotSupported
	}

	return support.CurrentBalance(ctx)
}

func (c *Client) Bills(ctx context.Context, account *Account) ([]Bill, error) {
	provider, err := NewProvider(c.base, account)
	if err != nil {
		return nil, err
	}

	support, ok := provider.(HasSupportBills)
	if !ok {
		return nil, ErrProviderMethodNotSupported
	}

	return support.Bills(ctx, time.Now().Add(-time.Hour*24*31*6), time.Now())
}

func (c *Client) BillDownload(ctx context.Context, account *Account, bill Bill, writer io.Writer) error {
	provider, err := NewProvider(c.base, account)
	if err != nil {
		return err
	}

	support, ok := provider.(HasSupportBills)
	if !ok {
		return ErrProviderMethodNotSupported
	}

	return support.BillDownload(ctx, bill, writer)
}

func (c *Client) IsSupportCurrentBalance(account *Account) bool {
	provider, err := NewProvider(c.base, account)
	if err != nil {
		return false
	}

	_, ok := provider.(HasSupportCurrentBalance)
	return ok
}

func (c *Client) IsSupportBills(account *Account) bool {
	provider, err := NewProvider(c.base, account)
	if err != nil {
		return false
	}

	_, ok := provider.(HasSupportBills)
	return ok
}
