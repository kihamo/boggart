package dom24

import (
	"bytes"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/dom24/client/accounting"
	"github.com/kihamo/boggart/providers/dom24/client/bill"
	"github.com/kihamo/boggart/providers/dom24/client/meters"
	"github.com/kihamo/boggart/providers/dom24/client/user"
	"github.com/kihamo/boggart/providers/dom24/static/models"
	"github.com/kihamo/shadow/components/dashboard"
)

type widgetMeterView struct {
	AccountID       string
	Name            string
	Address         string
	Resource        string
	FactoryNumber   string
	LastCheckupDate *time.Time
	NextCheckupDate *time.Time
	RecheckInterval uint64
	StartDate       *time.Time
	StartValue      float64
	Units           string
	Values          []*widgetMeterValueView
}

type widgetMeterValueView struct {
	Period time.Time
	Value  float64
	Delta  float64
	Kind   string
}

type widgetAccountView struct {
	ID             string
	Type           string
	Address        string
	FIO            string
	Company        string
	CompanyINN     string
	DebtActualDate string
	TotalAmount    float64
	Transactions   []*widgetAccountTransactionView
}

type widgetAccountTransactionView struct {
	IsPayment bool
	Date      swagger.Date
	Period    string
	Amount    float64
	FileURL   string
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	widget := b.Widget()

	vars := make(map[string]interface{}, 3)

	query := r.URL().Query()
	action := query.Get("action")
	vars["action"] = action

	switch action {
	case "bill":
		billValue := query.Get("bill")
		if billValue == "" {
			widget.NotFound(w, r)
			return
		}

		billID, err := strconv.ParseUint(billValue, 10, 64)
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		params := bill.NewDownloadParamsWithContext(ctx).
			WithID(billID)
		buf := bytes.NewBuffer(nil)

		_, err = b.provider.Bill.Download(params, buf)
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename=\"dom24_bill_"+billValue+".pdf\"")
		w.Header().Set("Content-Type", "application/pdf ")
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

		_, err = io.Copy(w, buf)
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		return

	case "meters":
		response, err := b.provider.Meters.List(meters.NewListParamsWithContext(ctx))
		if err != nil {
			widget.FlashError(r, "Get meters list failed %v", "", err)

			widget.Render(ctx, "widget", vars)
			return
		}

		list := make([]*widgetMeterView, 0, len(response.GetPayload().Data))

		for _, meter := range response.GetPayload().Data {
			if meter.IsDisabled {
				continue
			}

			item := &widgetMeterView{
				AccountID:       meter.Ident,
				Name:            meter.Name,
				Address:         meter.Address,
				Resource:        meter.Resource,
				FactoryNumber:   meter.FactoryNumber,
				RecheckInterval: meter.RecheckInterval,
				StartValue:      meter.StartValue,
				Units:           meter.Units,
				Values:          make([]*widgetMeterValueView, 0, len(meter.Values)),
			}

			if !meter.LastCheckupDate.Time().IsZero() {
				t := meter.LastCheckupDate.Time()
				item.LastCheckupDate = &t
			}

			if !meter.NextCheckupDate.Time().IsZero() {
				t := meter.NextCheckupDate.Time()
				item.NextCheckupDate = &t
			} else if !meter.LastCheckupDate.Time().IsZero() && meter.RecheckInterval > 0 {
				t := meter.LastCheckupDate.Time().AddDate(int(meter.RecheckInterval), 0, 0)
				item.NextCheckupDate = &t
			}

			if !meter.StartDate.Time().IsZero() {
				t := meter.StartDate.Time()
				item.StartDate = &t
			}

			list = append(list, item)

			if len(meter.Values) == 0 {
				continue
			}

			var currentPeriod time.Time
			prevPeriod := meter.Values[0].Period.Time()

			now := time.Now()
			nowYear, nowMonth, nowDay := now.Date()

			// если текущая дата раньше, чем установленный период фиксации показаний
			// то переключаем на предыдущий месяц, чтобы не захватывать пока еще
			// не отчетный месяц
			if meter.ValuesEndDay > 0 && nowDay < int(meter.ValuesEndDay) {
				now = time.Date(nowYear, nowMonth, 1, 0, 0, 0, 0, now.Location()).Add(-time.Nanosecond)
				nowYear, nowMonth, nowDay = now.Date()
			}

			if now.After(prevPeriod) {
				prevPeriod = now
			}

			for _, value := range meter.Values {
				currentPeriod = value.Period.Time()

				for year := prevPeriod.Year(); year >= currentPeriod.Year(); year-- {
					for month := prevPeriod.Month(); month > 0 && !currentPeriod.After(prevPeriod); month-- {
						item.Values = append(item.Values, &widgetMeterValueView{
							Period: prevPeriod,
							Value:  value.Value,
							Kind:   "Unknown",
						})

						prevPeriod = time.Date(year, month, 1, 0, 0, 0, 0, prevPeriod.Location()).Add(-time.Nanosecond)
					}
				}

				last := len(item.Values) - 1
				item.Values[last].Period = value.Period.Time()
				item.Values[last].Kind = value.Kind
			}

			// deltas
			for i, current := range item.Values {
				if i == len(item.Values)-1 {
					continue
				}

				current.Delta = current.Value - item.Values[i+1].Value
			}
		}

		vars["meters"] = list

	default:
		exists := make(map[string]bool)
		accounts := make(map[string]*widgetAccountView)

		if v := query.Get("account"); v != "" {
			if parts := strings.Split(v, ","); len(parts) > 0 {
				for _, part := range parts {
					exists[part] = false
				}

				vars["filter"] = parts
			}
		}

		userResponse, err := b.provider.User.UserInfo(user.NewUserInfoParamsWithContext(ctx))
		if err != nil {
			widget.FlashError(r, "Get user info failed %v", "", err)

			widget.Render(ctx, "widget", vars)
			return
		}

		for _, item := range userResponse.GetPayload().Accounts {
			if _, ok := exists[item.Ident]; len(exists) > 0 && !ok {
				continue
			}

			if len(exists) > 0 {
				exists[item.Ident] = true
			}

			accounts[item.Ident] = &widgetAccountView{
				ID:           item.Ident,
				Address:      item.Address,
				FIO:          item.FIO,
				Company:      item.Company,
				Transactions: make([]*widgetAccountTransactionView, 0),
			}
		}

		// проверяем все ли учетки найдены, если нет то пытаемся их добавить
		if b.config().AutoRegisterIfNotExists {
			for ident, exist := range exists {
				if exist {
					continue
				}

				params := user.NewAddByIdentParamsWithContext(ctx)
				params.Request.Ident = &[]string{ident}[0]

				addedResponse, err := b.provider.User.AddByIdent(params)
				if err != nil {
					widget.FlashError(r, "Add account #%s failed %v", "", ident, err)
					continue
				}

				for _, item := range addedResponse.GetPayload().Accounts {
					if item.Ident != ident {
						continue
					}

					accounts[item.Ident] = &widgetAccountView{
						ID:           item.Ident,
						Address:      item.Address,
						FIO:          item.Fio,
						Company:      item.Company,
						Transactions: make([]*widgetAccountTransactionView, 0),
					}

					break
				}
			}
		}

		accountingResponse, err := b.provider.Accounting.AccountingInfo(accounting.NewAccountingInfoParamsWithContext(ctx))
		if err != nil {
			widget.FlashError(r, "Get accounting info failed %v", "", err)
		} else {
			for _, item := range accountingResponse.GetPayload().Data {
				account, ok := accounts[item.Ident]
				if !ok {
					continue
				}

				account.Type = item.AccountType
				account.CompanyINN = item.INN
				account.DebtActualDate = item.DebtActualDate
				account.TotalAmount = item.TotalSum

				for _, b := range item.Bills {
					account.Transactions = append(account.Transactions, &widgetAccountTransactionView{
						IsPayment: false,
						Date:      b.Date,
						Period:    strings.TrimRight(b.Period, " г."),
						Amount:    b.Total,
						FileURL:   b.FileLink,
					})
				}

				for _, payment := range item.Payments {
					account.Transactions = append(account.Transactions, &widgetAccountTransactionView{
						IsPayment: true,
						Date:      payment.Date,
						Period:    payment.Period,
						Amount:    payment.Sum,
					})
				}

				sort.SliceStable(account.Transactions, func(i, j int) bool {
					return account.Transactions[i].Date.Time().After(account.Transactions[j].Date.Time())
				})

				accounts[item.Ident] = account
			}
		}

		vars["accounts"] = accounts

		// удаляем учетки, которые автоматически добавляли
		if b.config().AutoRegisterIfNotExists {
			for ident, exist := range exists {
				if exist {
					continue
				}

				params := user.NewDeleteByIdentParamsWithContext(ctx)
				params.Request.Ident = ident

				if _, err := b.provider.User.DeleteByIdent(params); err != nil {
					widget.FlashError(r, "Delete account #%s failed %v", "", ident, err)
				}
			}
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
