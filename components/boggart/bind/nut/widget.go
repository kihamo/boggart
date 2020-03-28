package nut

import (
	"fmt"
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

			if e := bind.SendCommand(cmd); e != nil {
				err = fmt.Errorf("Execute command %s return error: %w", cmd, e)
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

					if e := bind.SetVariable(key, value[0]); e != nil {
						err = fmt.Errorf("Set variable %s return error: %w", key, e)
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
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get List UPS failed with error %s", "", err.Error()))
	}

	variables, err := bind.Variables()
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get variables failed with error %s", "", err.Error()))
	}

	commands, err := bind.Commands()
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get commands failed with error %s", "", err.Error()))
	}

	variablesView := make(map[string]interface{}, len(variables))

	var (
		charged bool
		runtime time.Duration
	)

	for _, v := range variables {
		variablesView[v.Name] = v

		switch v.Name {
		case "ups.status":
			charged = strings.HasSuffix(v.Value.(string), " CHRG")

		case "battery.runtime":
			runtime = time.Second * time.Duration(v.Value.(int))
		}
	}

	vars := map[string]interface{}{
		"ups":       ups,
		"variables": variablesView,
		"commands":  commands,
		"charged":   charged,
		"runtime":   runtime,
		"error":     err,
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
