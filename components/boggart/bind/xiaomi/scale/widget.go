package scale

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()

	if r.IsPost() {
		var err error

		profile := bind.Profile(r.Original().FormValue("profile"))
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
			bind.SetProfile(profile.Name)
			err = bind.notifyCurrentProfile(ctx)
		}

		if err != nil {
			r.Session().FlashBag().Error(err.Error())
		} else {
			r.Session().FlashBag().Info(t.Translate(ctx, "Profile set success", ""))
			t.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	}

	t.Render(r.Context(), "widget", map[string]interface{}{
		"current_profile": bind.CurrentProfile(),
		"profiles":        bind.Profiles(),
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
