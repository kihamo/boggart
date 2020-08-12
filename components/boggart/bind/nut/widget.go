package nut

import (
	"fmt"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

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
				widget.NotFound(w, r)
				return
			}

			if e := b.SendCommand(cmd); e != nil {
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

					if e := b.SetVariable(key, value[0]); e != nil {
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
			widget.NotFound(w, r)
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

	ups, err := b.ups()
	if err != nil {
		widget.FlashError(r, "Get List UPS failed with error %v", "", err)
	}

	variables, err := b.Variables()
	if err != nil {
		widget.FlashError(r, "Get variables failed with error %v", "", err)
	}

	commands, err := b.Commands()
	if err != nil {
		widget.FlashError(r, "Get commands failed with error %v", "", err)
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

	widget.Render(r.Context(), "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
