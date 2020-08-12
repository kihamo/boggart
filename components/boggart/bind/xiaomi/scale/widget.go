package scale

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	widget := b.Widget()

	if r.IsPost() {
		var err error

		profile := b.Profile(r.Original().FormValue("profile"))
		if profile == nil {
			err = errors.New("Profile not found")
		} else if profile.Editable {
			if val := r.Original().FormValue("height"); val != "" {
				profile.Height, err = strconv.ParseUint(val, 10, 64)
			} else {
				profile.Height = 0
			}

			if err == nil {
				if val := r.Original().FormValue("age"); val != "" {
					profile.Age, err = strconv.ParseUint(val, 10, 64)
				} else {
					profile.Age = 0
				}
			}

			if err == nil {
				if val := r.Original().FormValue("birthday"); val != "" {
					profile.Birthday, err = time.Parse("2006.01.02", val)
				} else {
					profile.Birthday = time.Time{}
				}
			}

			if err == nil {
				if val := r.Original().FormValue("sex"); val != "" {
					profile.Sex = val == "1"
				}
			}
		}

		if err == nil {
			b.SetProfile(profile.Name)
			err = b.notifyCurrentProfile(ctx)
		}

		if err != nil {
			r.Session().FlashBag().Error(err.Error())
		} else {
			r.Session().FlashBag().Info(widget.Translate(ctx, "Profile set success", ""))
			widget.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"current_profile": b.CurrentProfile(),
		"profiles":        b.Profiles(),
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
