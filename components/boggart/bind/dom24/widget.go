package dom24

import (
	"sort"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/dom24/client/accounting"
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

	accounts := make(map[string]*widgetAccountView, 0)

	userResponse, err := b.provider.User.UserInfo(user.NewUserInfoParamsWithContext(ctx))
	if err != nil {
		widget.FlashError(r, "Get user info failed %v", "", err)
	} else {
		for _, item := range userResponse.GetPayload().Accounts {
			accounts[item.Ident] = &widgetAccountView{
				ID:           item.Ident,
				Address:      item.Address,
				FIO:          item.FIO,
				Company:      item.Company,
				Transactions: make([]*widgetAccountTransactionView, 0),
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
				account = &widgetAccountView{
					ID:           item.Ident,
					Address:      item.Address,
					Transactions: make([]*widgetAccountTransactionView, 0),
				}
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

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"accounts": accounts,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
