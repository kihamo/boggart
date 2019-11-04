package led_wifi

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()

	vars := map[string]interface{}{}

	state, err := bind.bulb.State(ctx)
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(ctx,
			"Get state failed with error %s",
			"",
			err.Error(),
		))
	} else {
		vars["state"] = state
	}

	if r.IsPost() {
		err := r.Original().ParseForm()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Parse form failed with error %s", "", err.Error()))
		} else {
			var power bool

			for key, value := range r.Original().PostForm {
				if len(value) == 0 || key != "state" {
					continue
				}

				power = value[0] == "on"
				break
			}

			if power {
				err = bind.bulb.PowerOn(ctx)
			} else {
				err = bind.bulb.PowerOff(ctx)
			}

			if err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Change state failed with error %s", "", err.Error()))
			}

			t.Redirect(r.URL().Path, http.StatusFound, w, r)
		}
	}

	t.Render(ctx, "index", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
