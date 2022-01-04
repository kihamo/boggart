package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/kihamo/boggart/performance"
	"github.com/kihamo/boggart/providers/integratorit/internal"
)

type ProvideMosOblERC struct {
	Provider

	base    *internal.Client
	account *Account
}

func NewProvideMosOblERC(base *internal.Client, account *Account) *ProvideMosOblERC {
	return &ProvideMosOblERC{
		base:    base,
		account: account,
	}
}

func (p *ProvideMosOblERC) CurrentBalance(ctx context.Context) (float64, error) {
	var response []struct {
		Value float64 `json:"sm_balance"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "smorodinaTransProxy",
	}, map[string]string{
		"plugin":      "smorodinaTransProxy",
		"proxyquery":  "AbonentCurrentBalance",
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

func (p *ProvideMosOblERC) Bills(ctx context.Context, dateStart, dateEnd time.Time) ([]Bill, error) {
	var response []struct {
		Period internal.Date `json:"dt_period"`
		Child  []struct {
			UUID   string  `json:"vl_report_uuid"`
			Amount float64 `json:"sm_total"`
		} `json:"child"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "smorodinaTransProxy",
	}, map[string]string{
		"plugin":          "smorodinaTransProxy",
		"proxyquery":      "AbonentChargeDetail",
		"vl_provider":     p.account.ProviderRAW,
		"dt_period_start": dateStart.Format(time.RFC3339),
		"dt_period_end":   dateEnd.Format(time.RFC3339),
		"kd_tp_mode":      "1",
	}, &response)

	if err != nil {
		return nil, err
	}

	bills := make([]Bill, 0, len(response))

	for _, item := range response {
		if len(item.Child) < 1 {
			continue
		}

		bills = append(bills, Bill{
			ID:     item.Child[0].UUID,
			Period: item.Period.Time,
			Amount: item.Child[0].Amount,
		})
	}

	sort.SliceStable(bills, func(i, j int) bool {
		return bills[i].Period.After(bills[j].Period)
	})

	return bills, nil
}

func (p *ProvideMosOblERC) BillDownload(ctx context.Context, bill Bill, writer io.Writer) error {
	var params struct {
		AbonentID uint64 `json:"id_abonent"`
	}

	if err := json.Unmarshal(performance.UnsafeString2Bytes(p.account.ProviderRAW), &params); err != nil {
		return err
	}

	return p.base.DoRequest(ctx, "", nil, map[string]string{
		"action":  "sql",
		"plugin":  "smorodinaReportProxy",
		"query":   "smorodinaProxy",
		"rParams": "ID_ABONENT=" + strconv.FormatUint(params.AbonentID, 10) + "&DT_PERIOD=" + bill.Period.Format("02.01.2006"),
		"rUuid":   bill.ID,
	}, writer)
}
