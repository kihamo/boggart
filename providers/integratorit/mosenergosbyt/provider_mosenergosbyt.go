package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/internal"
)

type ProvideMosEnergoSbyt struct {
	Provider

	base    *internal.Client
	account *Account
}

func NewProvideMosEnergoSbyt(base *internal.Client, account *Account) *ProvideMosEnergoSbyt {
	return &ProvideMosEnergoSbyt{
		base:    base,
		account: account,
	}
}

func (p *ProvideMosEnergoSbyt) CurrentBalance(ctx context.Context) (float64, error) {
	var response []struct {
		Value float64 `json:"vl_balance"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "bytProxy",
	}, map[string]string{
		"plugin":      "bytProxy",
		"proxyquery":  "CurrentBalance",
		"vl_provider": p.account.ProviderRAW,
	}, &response)

	if err != nil {
		return 0, err
	}

	if len(response) == 0 {
		return 0, errors.New("balance information not found in response")
	}

	return response[0].Value, nil
}

func (p *ProvideMosEnergoSbyt) Bills(ctx context.Context, dateStart, dateEnd time.Time) ([]Bill, error) {
	var response []struct {
		ID     uint64        `json:"id_korr"`
		Period internal.Time `json:"dt_period"`
		Amount float64       `json:"sm_total"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "bytProxy",
	}, map[string]string{
		"plugin":      "bytProxy",
		"proxyquery":  "Invoice",
		"dt_st":       dateStart.Format(time.RFC3339),
		"dt_en":       dateEnd.Format(time.RFC3339),
		"vl_provider": p.account.ProviderRAW,
	}, &response)

	if err != nil {
		return nil, err
	}

	bills := make([]Bill, 0, len(response))

	for _, item := range response {
		bills = append(bills, Bill{
			ID:     strconv.FormatUint(item.ID, 10),
			Period: item.Period.Time,
			Amount: item.Amount,
		})
	}

	sort.SliceStable(bills, func(i, j int) bool {
		return bills[i].Period.After(bills[j].Period)
	})

	return bills, nil
}

func (p *ProvideMosEnergoSbyt) BillDownload(ctx context.Context, bill Bill, writer io.Writer) error {
	var response []struct {
		Link   string `json:"nm_link"`
		Params string `json:"vl_params"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "bytProxy",
	}, map[string]string{
		"plugin":      "bytProxy",
		"proxyquery":  "GetPrintBillLink",
		"dt_period":   bill.Period.Format("2006-01-02 15:04:05"),
		"kd_provider": "1",
		"vl_provider": p.account.ProviderRAW,
	}, &response)

	if len(response) == 0 {
		return errors.New("last bill not found")
	}

	u, err := url.Parse(response[0].Link)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	if err := json.Unmarshal([]byte(response[0].Params), &params); err != nil {
		return err
	}

	body := url.Values{}
	for key, val := range params {
		body[key] = []string{val}
	}

	return p.base.DoRequestRaw(ctx, u, strings.NewReader(body.Encode()), writer)
}
