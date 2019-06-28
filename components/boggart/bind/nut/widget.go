package nut

import (
	"errors"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	if r.IsPost() {
		var (
			err        error
			successMsg string
		)

		q := r.URL().Query()

		switch q.Get("action") {
		case "cmd":
			cmd := strings.TrimSpace(q.Get("cmd"))
			if cmd == "" {
				t.NotFound(w, r)
				return
			}

			ok, err := bind.SendCommand(cmd)
			if err != nil {
				// skip
			} else if !ok {
				err = errors.New("Execute command " + cmd + " return false")
			} else {
				successMsg = "Execute command " + cmd + " success"
			}

		case "variable":
			err = r.Original().ParseForm()
			if err == nil {
				keys := make([]string, 0)

				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					ok, err := bind.SetVariable(key, value[0])
					if err != nil {
						break
					} else if !ok {
						err = errors.New("Set variable " + key + " return false")
						break
					}

					keys = append(keys, key)
				}

				if err == nil {
					successMsg = "Change variables (" + strings.Join(keys, ",") + ") value success"
				}
			}

		default:
			t.NotFound(w, r)
			return
		}

		if err != nil {
			_ = w.SendJSON(response{
				Result:  "failed",
				Message: err.Error(),
			})

		} else {
			_ = w.SendJSON(response{
				Result:  "success",
				Message: successMsg,
			})
		}

		return
	}

	ups, err := bind.ups()
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get UPS failed with error %s", "", err.Error()))
	}

	variables := make(map[string]interface{}, len(ups.Variables))

	var (
		charged bool
		runtime time.Duration
	)

	for _, v := range ups.Variables {
		variables[v.Name] = v

		switch v.Name {
		case "ups.status":
			charged = strings.HasSuffix(v.Value.(string), " CHRG")

		case "battery.runtime":
			runtime = time.Second * time.Duration(v.Value.(int64))
		}
	}

	vars := map[string]interface{}{
		"ups":       ups,
		"variables": variables,
		"charged":   charged,
		"runtime":   runtime,
		"error":     err,
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
