package elektroset

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/internal"
	"go.uber.org/multierr"
)

const (
	BaseURL = "https://lkk.oao-elektroset.ru/gate_lkk_myt"

	dateFormatLayout = "2006-01-02"
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

func (c *Client) Houses(ctx context.Context) ([]House, error) {
	data := make([]struct {
		Houses []House `json:"houses"`
	}, 0)

	err := c.base.DoRequest(ctx, "sql", map[string]string{"query": "lka_get_houses"}, nil, &data)
	if err != nil {
		return nil, err
	}

	houses := make([]House, 0)

	var e error

	for _, row := range data {
		for _, house := range row.Houses {
			for i, service := range house.Services {
				house.Services[i].Balance, e = strconv.ParseFloat(service.BalanceRAW, 10)
				if e != nil {
					err = multierr.Append(err, e)
				}
			}

			houses = append(houses, house)
		}
	}

	return houses, err
}

func (c *Client) BalanceDetails(ctx context.Context, accountID string, provider Provider, dateStart, dateEnd time.Time) ([]BalanceDetail, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}

	data := make([]BalanceDetail, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_balance_detail"}, map[string]string{
		"dt_start":    dateStart.Format(dateFormatLayout),
		"dt_end":      dateEnd.Format(dateFormatLayout),
		"nn_ls":       accountID,
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

func (c *Client) IndicatorInfo(ctx context.Context, accountID string, provider Provider) (*IndicationInfo, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}

	data := make([]IndicationInfo, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_info_ind"}, map[string]string{
		"nn_ls":       accountID,
		"nm_email":    c.base.Login(),
		"vl_provider": string(vlProvider),
	}, &data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("bad response")
	}

	return &data[0], nil
}

func (c *Client) CurrentRate(ctx context.Context, houseID uint64, provider Provider) (*Rate, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}

	data := make([]Rate, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_current_rate"}, map[string]string{
		"nn_ls":       strconv.FormatUint(houseID, 10),
		"nm_email":    c.base.Login(),
		"vl_provider": string(vlProvider),
	}, &data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("bad response")
	}

	return &data[0], nil
}

func (c *Client) Meter(ctx context.Context, accountID string, provider Provider) (*Meter, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}

	data := make([]Meter, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_meter"}, map[string]string{
		"nn_ls":       accountID,
		"nm_email":    c.base.Login(),
		"vl_provider": string(vlProvider),
	}, &data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("bad response")
	}

	return &data[0], nil
}

func (c *Client) BillFile(ctx context.Context, accountID string, provider Provider, period time.Time) (string, error) {
	vlProvider, err := json.Marshal(provider)
	if err != nil {
		return "", err
	}

	data := make([]BillFile, 0)

	err = c.base.DoRequest(ctx, "sql", map[string]string{"query": "proxy", "plugin": "bytTmbProxy", "proxyquery": "tmb_file_bill_period"}, map[string]string{
		"nn_ls":       accountID,
		"v_email":     c.base.Login(),
		"vl_provider": string(vlProvider),
		"dt_period":   period.Format(dateFormatLayout),
	}, &data)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New("bad response")
	}

	link, err := url.Parse(data[0].Link)
	if err != nil {
		return "", err
	}

	values := link.Query()
	values.Add("session", c.base.Session())
	link.RawQuery = values.Encode()

	return link.String(), err
}
