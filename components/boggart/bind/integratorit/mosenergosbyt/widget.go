package mosenergosbyt

import (
	"errors"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()

	account, err := b.Account(ctx)
	if err != nil {
		widget.NotFound(w, r)
		return
	}

	query := r.URL().Query()

	if query.Get("action") != "download" {
		widget.NotFound(w, r)
		return
	}

	bills, err := b.client.Bills(ctx, account)
	if err != nil {
		widget.InternalError(w, r, err)
		return
	}

	var bill *mosenergosbyt.Bill
	if id := query.Get("bill"); id == "" && len(bills) > 0 {
		bill = &bills[0]
	} else {
		for _, b := range bills {
			if id == b.ID {
				bill = &b
				break
			}
		}
	}

	if bill == nil {
		widget.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\"mosenergosbyt_bill_"+bill.Period.Format("20060102")+".pdf\"")
	w.Header().Set("Content-Type", "application/pdf")

	if err = b.client.BillDownload(ctx, account, *bill, w); err != nil {
		if errors.Is(err, mosenergosbyt.ErrProviderMethodNotSupported) {
			widget.NotFound(w, r)
			return
		}

		widget.InternalError(w, r, err)
	}
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return nil
}
