package elektroset

import (
	"io"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/http"
	"github.com/kihamo/boggart/providers/integratorit/elektroset"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{}
	ctx := r.Context()
	query := r.URL().Query()

	action := query.Get("action")
	vars["action"] = action

	switch action {
	case "bill":
		period, err := time.Parse(layoutPeriod, query.Get("period"))
		if err != nil {
			t.NotFound(w, r)
			return
		}

		accounts, err := bind.client.Accounts(ctx)
		if err != nil {
			t.InternalError(w, r, err)
			return
		}

		for _, account := range accounts {
			for _, service := range account.Services {
				billLink, err := bind.client.BillFile(ctx, account.Number, service.Provider, period)
				if err != nil {
					t.InternalError(w, r, err)
					return
				}

				response, err := http.NewClient().Get(ctx, billLink)
				if err != nil {
					t.InternalError(w, r, err)
					return
				}

				for key, values := range response.Header {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}

				_, err = io.Copy(w, response.Body)
				if err != nil {
					t.InternalError(w, r, err)
					return
				}

				response.Body.Close()
				return
			}
		}

		t.NotFound(w, r)
		return

	default:
		type row struct {
			Date   time.Time
			Values map[string][3]float64
		}

		rows := make([]*row, 0)
		rowsByTime := make(map[time.Time]*row)
		tariffs := make(map[int64]bool, 3)

		accounts, err := bind.client.Accounts(ctx)
		if err == nil {
			dateStart := time.Now().Add(-time.Hour * 24 * 365)
			dateEnd := time.Now()

			for _, account := range accounts {
				for _, service := range account.Services {
					if balances, e := bind.client.BalanceDetails(ctx, account.Number, service.Provider, dateStart, dateEnd); e == nil {
						for _, balance := range balances {
							switch balance.KDTariffPlanEntity {
							// показания
							case elektroset.TariffPlanEntityValue:
								r, ok := rowsByTime[balance.DatetimeEntity.Time]
								if !ok && (balance.ValueT1 != nil || balance.ValueT2 != nil || balance.ValueT3 != nil) {
									r = &row{
										Date:   balance.DatetimeEntity.Time,
										Values: make(map[string][3]float64, 3),
									}
									rowsByTime[balance.DatetimeEntity.Time] = r
									rows = append(rows, r)
								}

								if balance.ValueT1 != nil {
									r.Values["tariff1"] = [3]float64{*balance.ValueT1 * 1000}
									tariffs[1] = true
								}

								if balance.ValueT2 != nil {
									r.Values["tariff2"] = [3]float64{*balance.ValueT2 * 1000}
									tariffs[2] = true
								}

								if balance.ValueT3 != nil {
									r.Values["tariff3"] = [3]float64{*balance.ValueT3 * 1000}
									tariffs[3] = true
								}
							}
						}
					} else {
						r.Session().FlashBag().Error(t.Translate(ctx, "Get balance details failed with error %s", "", err.Error()))
					}
				}

				// deltas 1
				for i, row := range rows {
					if i > 0 {
						for tariff, current := range row.Values {
							if prev, ok := rows[i-1].Values[tariff]; ok {
								current[1] = current[0] - prev[0]
								row.Values[tariff] = current
							}
						}
					}
				}

				// trends 2
				for i, row := range rows {
					if i > 0 {
						for tariff, current := range row.Values {
							if prev, ok := rows[i-1].Values[tariff]; ok {
								current[2] = current[1] - prev[1]
								row.Values[tariff] = current
							}
						}
					}
				}
			}
		} else {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get accounts failed with error %s", "", err.Error()))
		}

		vars["rows"] = rows
		vars["tariffs"] = tariffs
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
