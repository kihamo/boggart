package ledwifi

import (
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/wifiled"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	widget := b.Widget()

	vars := map[string]interface{}{
		"effects": map[wifiled.Mode]string{
			wifiled.ModePreset1:  wifiled.ModePreset1.String(),
			wifiled.ModePreset2:  wifiled.ModePreset2.String(),
			wifiled.ModePreset3:  wifiled.ModePreset3.String(),
			wifiled.ModePreset4:  wifiled.ModePreset4.String(),
			wifiled.ModePreset5:  wifiled.ModePreset5.String(),
			wifiled.ModePreset6:  wifiled.ModePreset6.String(),
			wifiled.ModePreset7:  wifiled.ModePreset7.String(),
			wifiled.ModePreset8:  wifiled.ModePreset8.String(),
			wifiled.ModePreset9:  wifiled.ModePreset9.String(),
			wifiled.ModePreset10: wifiled.ModePreset10.String(),
			wifiled.ModePreset11: wifiled.ModePreset11.String(),
			wifiled.ModePreset12: wifiled.ModePreset12.String(),
			wifiled.ModePreset13: wifiled.ModePreset13.String(),
			wifiled.ModePreset14: wifiled.ModePreset14.String(),
			wifiled.ModePreset15: wifiled.ModePreset15.String(),
			wifiled.ModePreset16: wifiled.ModePreset16.String(),
			wifiled.ModePreset17: wifiled.ModePreset17.String(),
			wifiled.ModePreset18: wifiled.ModePreset18.String(),
			wifiled.ModePreset19: wifiled.ModePreset19.String(),
			wifiled.ModePreset20: wifiled.ModePreset20.String(),
			wifiled.ModePreset21: wifiled.ModePreset21.String(),
			wifiled.ModeCustom:   wifiled.ModeCustom.String(),
			wifiled.ModeStatic:   wifiled.ModeStatic.String(),
			wifiled.ModeMusic:    wifiled.ModeMusic.String(),
			wifiled.ModeTesting:  wifiled.ModeTesting.String(),
		},
	}

	state, err := b.bulb.State(ctx)
	if err != nil {
		widget.FlashError(r, "Get state failed with error %v", "", err)
	} else {
		vars["state"] = state
	}

	if r.IsPost() {
		err := r.Original().ParseForm()
		if err != nil {
			widget.FlashError(r, "Parse form failed with error %v", "", err)
		} else {
			var power bool

			for key, value := range r.Original().PostForm {
				if len(value) == 0 || key != "state" {
					continue
				}

				power = value[0] == "on"
				break
			}

			if power != state.Power {
				if power {
					err = b.bulb.PowerOn(ctx)
				} else {
					err = b.bulb.PowerOff(ctx)
				}

				if err != nil {
					widget.FlashError(r, "Change state failed with error %v", "", err)
				}
			}

			if power {
				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					if key == "effect" {
						mode, err := wifiled.ModeFromString(strings.TrimSpace(value[0]))
						if err == nil {
							err = b.bulb.SetMode(ctx, *mode, state.Speed)
						}

						if err != nil {
							widget.FlashError(r, "Change mode failed with error %v", "", err)
						}
					}
				}
			}

			widget.Redirect(r.URL().Path, http.StatusFound, w, r)
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
