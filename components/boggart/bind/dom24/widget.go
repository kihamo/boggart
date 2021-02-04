package dom24

import (
	"bytes"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/dom24/client/accounting"
	"github.com/kihamo/boggart/providers/dom24/client/bill"
	"github.com/kihamo/boggart/providers/dom24/client/user"
	"github.com/kihamo/boggart/providers/dom24/static/models"
	"github.com/kihamo/shadow/components/dashboard"
)

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

	exists := make(map[string]bool)
	accounts := make(map[string]*widgetAccountView)
	vars := make(map[string]interface{}, 2)

	query := r.URL().Query()
	action := query.Get("action")

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

	default:
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

			widget.Render(r.Context(), "widget", vars)
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
		if b.config.AutoRegisterIfNotExists {
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

				for _, bill := range item.Bills {
					account.Transactions = append(account.Transactions, &widgetAccountTransactionView{
						IsPayment: false,
						Date:      bill.Date,
						Period:    strings.TrimRight(bill.Period, " г."),
						Amount:    bill.Total,
						FileURL:   bill.FileLink,
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
		if b.config.AutoRegisterIfNotExists {
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

	widget.Render(r.Context(), "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
