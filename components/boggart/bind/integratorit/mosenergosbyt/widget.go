package mosenergosbyt

import (
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	query := r.URL().Query()

	action := query.Get("action")
	uuid := query.Get("uuid")

	if action != "bill" || len(uuid) == 0 {
		widget.NotFound(w, r)
		return
	}

	period, err := time.Parse(layoutPeriod, query.Get("period"))
	if err != nil {
		widget.NotFound(w, r)
		return
	}

	abonent, err := strconv.ParseUint(query.Get("abonent"), 10, 64)
	if err != nil {
		widget.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\"mosenergosbyt_bill_"+period.Format("20060102")+".pdf\"")
	w.Header().Set("Content-Type", "application/pdf")

	if err := b.client.Bill(r.Context(), abonent, uuid, period, w); err != nil {
		widget.InternalError(w, r, err)
	}
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return nil
}
