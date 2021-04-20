package elektroset

import (
	"encoding/json"
	"io"
	"sort"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/protocols/http"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
	"github.com/kihamo/shadow/components/dashboard"
)

type widgetAccountView struct {
	MeterNumber          string
	MeterCalibrationDate time.Time
	MeterPhaseName       string
	MeterState           string
	MeterModel           string
	MeterRate            string
	MeterQuantity        float64
	AccountID            string
	HouseName            string
	HouseAddress         string
	ProviderName         string
	ServiceName          string
	Values               []*widgetMeterView
	Zones                [3]string // 0 - T1, 1 - T2, 2 - T3
}

type widgetMeterView struct {
	Date    time.Time
	Current [3]float64
	Deltas  [3]float64
	Trends  [3]float64
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	vars := map[string]interface{}{}
	ctx := r.Context()
	query := r.URL().Query()
	widget := b.Widget()

	action := query.Get("action")
	vars["action"] = action

	switch action {
	case "bill":
		accountID := query.Get("account")

		period, err := time.Parse(layoutPeriod, query.Get("period"))
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		var provider elektroset.Provider
		err = json.Unmarshal([]byte(query.Get("provider")), &provider)

		if err != nil {
			widget.NotFound(w, r)
			return
		}

		billLink, err := b.client.BillFile(ctx, accountID, provider, period)
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		response, err := http.NewClient().Get(ctx, billLink)
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=\"elektroset_bill_"+period.Format("20060102")+".pdf\"")
		w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", response.Header.Get("Content-Length"))

		_, err = io.Copy(w, response.Body)
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		response.Body.Close()

		return

	default:
		houses, err := b.Houses(ctx)
		if err != nil {
			widget.FlashError(r, "Get houses failed with error %v", "")
			break
		}

		dateStart := time.Now().Add(-time.Hour * 24 * 365)
		dateEnd := time.Now()
		accountsView := make([]*widgetAccountView, 0, len(houses))

		for _, house := range houses {
			for _, service := range house.Services {
				balances, err := b.client.BalanceDetails(ctx, service.AccountID, service.Provider, dateStart, dateEnd)
				if err != nil {
					widget.FlashError(r, "Get balance details failed with error %v", "", err)
					continue
				}

				meter, err := b.client.Meter(ctx, service.AccountID, service.Provider)
				if err != nil {
					widget.FlashError(r, "Get meter info failed with error %v", "", err)
					continue
				}

				accountView := &widgetAccountView{
					AccountID:            service.AccountID,
					HouseName:            house.Name,
					HouseAddress:         house.Address,
					ProviderName:         service.ProviderName,
					ServiceName:          service.ServiceTypeName,
					MeterNumber:          meter.Number,
					MeterCalibrationDate: meter.CalibrationDate.Time,
					MeterPhaseName:       meter.PhaseName,
					MeterState:           meter.State,
					MeterModel:           meter.Model,
					MeterRate:            meter.Rate,
					MeterQuantity:        meter.Quantity,
				}

				metersByTime := make(map[int]*widgetMeterView)
				metersKeys := make([]int, 0)

				for _, balance := range balances {
					if balance.TariffPlanEntityID != elektroset.TariffPlanEntityValue {
						continue
					}

					unix := int(balance.DatetimeEntity.Time.Unix())
					r, ok := metersByTime[unix]
					if !ok && (balance.T1Value != nil || balance.T2Value != nil || balance.T3Value != nil) {
						r = &widgetMeterView{
							Date: balance.DatetimeEntity.Time,
						}

						metersByTime[unix] = r
						metersKeys = append(metersKeys, unix)
					}

					switch {
					case balance.T1Value != nil:
						r.Current[0] = *balance.T1Value
						accountView.Zones[0] = *balance.Zone1Name
					case balance.T2Value != nil:
						r.Current[1] = *balance.T2Value
						accountView.Zones[1] = *balance.Zone2Name
					case balance.T3Value != nil:
						r.Current[2] = *balance.T3Value
						accountView.Zones[2] = *balance.Zone3Name
					}
				}

				sort.Sort(sort.Reverse(sort.IntSlice(metersKeys)))
				//sort.Ints(metersKeys)

				for _, k := range metersKeys {
					accountView.Values = append(accountView.Values, metersByTime[k])
				}

				// deltas 1
				for i, meter := range accountView.Values {
					for tariff, value := range meter.Current {
						meter.Deltas[tariff] = value - accountView.Values[i+1].Current[tariff]
					}

					if i == len(accountView.Values)-2 {
						break
					}
				}

				// trends 2
				for i, meter := range accountView.Values {
					for tariff, value := range meter.Deltas {
						meter.Trends[tariff] = value - accountView.Values[i+1].Deltas[tariff]
					}

					if i == len(accountView.Values)-2 {
						break
					}
				}

				accountsView = append(accountsView, accountView)
			}
		}

		vars["accounts"] = accountsView
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
