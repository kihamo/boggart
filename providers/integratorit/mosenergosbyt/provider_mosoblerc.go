package mosenergosbyt

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/performance"
	"github.com/kihamo/boggart/providers/integratorit/internal"
)

var servicesMapping = map[string]string{
	"КАПИТАЛЬНЫЙ РЕМОНТ":        "ВЗНОС НА КАП. РЕМОНТ",
	"ГАЗОСНАБЖЕНИЕ (ОТОПЛЕНИЕ)": "ГАЗОСНАБЖЕНИЕ",
}

type ProvideMosOblERC struct {
	Provider

	base    *internal.Client
	account *Account

	servicesOnce *atomic.Once
	servicesList map[string]uint64
}

func NewProvideMosOblERC(base *internal.Client, account *Account) *ProvideMosOblERC {
	return &ProvideMosOblERC{
		base:         base,
		account:      account,
		servicesOnce: &atomic.Once{},
		servicesList: make(map[string]uint64, 0),
	}
}

func (p *ProvideMosOblERC) Services(ctx context.Context) (_ map[string]uint64, err error) {
	p.servicesOnce.Do(func() {
		var response []struct {
			ID   uint64 `json:"id_service"`
			Name string `json:"nm_service"`
		}

		err = p.base.DoRequest(ctx, "sql", map[string]string{
			"query": "smorodinaTransProxy",
		}, map[string]string{
			"plugin":      "smorodinaTransProxy",
			"proxyquery":  "AbonentContractData",
			"vl_provider": p.account.ProviderRAW,
		}, &response)

		if err != nil {
			return
		}

		for _, item := range response {
			// Dirty hack
			item.Name = strings.ReplaceAll(item.Name, "(", " (")
			item.Name = strings.ReplaceAll(item.Name, "  ", " ")
			item.Name = strings.TrimSpace(item.Name)

			if alias, ok := servicesMapping[item.Name]; ok {
				item.Name = alias
			}

			p.servicesList[item.Name] = item.ID
		}
	})

	if err != nil {
		p.servicesOnce.Reset()
	}

	return p.servicesList, err
}

func (p *ProvideMosOblERC) CurrentBalance(ctx context.Context) (*Balance, error) {
	var response []struct {
		Total float64 `json:"sm_balance"`
		Child []struct {
			Name  string  `json:"nm_service"`
			Total float64 `json:"sm_balance"`
		} `json:"child"`
	}

	err := p.base.DoRequest(ctx, "sql", map[string]string{
		"query": "smorodinaTransProxy",
	}, map[string]string{
		"plugin":      "smorodinaTransProxy",
		"proxyquery":  "AbonentCurrentBalance",
		"vl_provider": p.account.ProviderRAW,
	}, &response)

	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, errors.New("balance information not found in response")
	}

	balance := &Balance{
		Total:    response[0].Total,
		Services: make([]ServiceBalance, 0, len(response[0].Child)),
	}

	services, err := p.Services(ctx)
	if err != nil {
		return nil, errors.New("get information about services failed")
	}

	var (
		serviceID   uint64
		serviceName string
	)

	for _, child := range response[0].Child {
		serviceID = 0
		serviceName = strings.TrimSpace(child.Name)

		for name, id := range services {
			if name == serviceName {
				serviceID = id
				break
			}
		}

		if serviceID == 0 {
			return nil, errors.New("id for service " + child.Name + " not found")
		}

		balance.Services = append(balance.Services, ServiceBalance{
			ID:    strconv.FormatUint(serviceID, 10),
			Name:  child.Name,
			Total: child.Total,
		})
	}

	return balance, nil
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
