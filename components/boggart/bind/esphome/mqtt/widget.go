package mqtt

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	switch q.Get("action") {
	case "component":
		b.handleComponent(w, r)

	case "command":
		b.handleCommand(w, r)

	default:
		b.handleIndex(w, r)
	}
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) handleComponent(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()

	componentID := q.Get("id")

	if componentID == "" {
		widget.NotFound(w, r)
		return
	}

	component := b.Component(componentID)
	if component == nil {
		widget.NotFound(w, r)
		return
	}

	if component.GetType() == ComponentTypeLight {
		if light, ok := component.(*ComponentLight); ok {
			b.handleComponentLight(w, r, light)
			return
		}
	}

	widget.NotFound(w, r)
}

func (b *Bind) handleComponentLight(w http.ResponseWriter, r *dashboard.Request, component *ComponentLight) {
	ctx := r.Context()
	widget := b.Widget()

	var err error

	if r.IsPost() {
		err = r.Original().ParseForm()
		if err != nil {
			widget.FlashError(r, "Parse form failed with error %v", "", err)
		} else {
			command := component.GetState().(*ComponentLightState)
			command.SetState(false)

			for key, value := range r.Original().PostForm {
				if len(value) == 0 {
					continue
				}

				switch key {
				case "state":
					command.SetState(value[0] == "on")
				case "brightness":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.Brightness = val
					} else {
						err = e
					}
				case "red":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.Color.Red = val
					} else {
						err = e
					}
				case "green":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.Color.Green = val
					} else {
						err = e
					}
				case "blue":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.Color.Blue = val
					} else {
						err = e
					}
				case "white":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.White = val
					} else {
						err = e
					}
				case "color-temperature":
					if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
						command.ColorTemperature = val
					} else {
						err = e
					}
				case "effect":
					command.Effect = value[0]
				case "flash-length":
					if val, e := time.ParseDuration(value[0]); e == nil {
						if val > 0 {
							command.Flash = uint64(val.Seconds())
						}
					} else {
						err = e
					}
				case "transition-length":
					if val, e := time.ParseDuration(value[0]); e == nil {
						command.Transition = uint64(val.Seconds())
					} else {
						err = e
					}
				}

				if err != nil {
					widget.FlashError(r, "Parse param %s failed with error %v", "", key, err)
					err = nil
				}
			}

			ctx := r.Context()
			err := b.MQTT().PublishRawWithoutCache(ctx, component.GetCommandTopic(), 1, false, component.CommandToPayload(command))

			if err != nil {
				widget.FlashError(r, err, "")
			} else {
				widget.FlashSuccess(r, "Success toggle", "")

				widget.Redirect(r.URL().Path+"?action=component&id="+component.GetID(), http.StatusFound, w, r)
				return
			}
		}
	}

	widget.Render(ctx, "light", map[string]interface{}{
		"component": component,
		"error":     err,
	})
}

func (b *Bind) handleCommand(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()

	componentID := q.Get("id")
	command := q.Get("cmd")

	if componentID == "" || command == "" {
		widget.NotFound(w, r)
		return
	}

	component := b.Component(componentID)
	if component == nil || component.GetCommandTopic() == "" {
		widget.NotFound(w, r)
		return
	}

	ctx := r.Context()
	err := b.MQTT().PublishRawWithoutCache(ctx, component.GetCommandTopic(), 1, false, component.CommandToPayload(command))

	if err != nil {
		widget.FlashError(r, err, "")
	} else {
		widget.FlashSuccess(r, "Success toggle", "")
	}

	redirectURL := &url.URL{}
	*redirectURL = *r.Original().URL
	redirectURL.RawQuery = ""

	widget.Redirect(redirectURL.String(), http.StatusFound, w, r)
}

func (b *Bind) handleIndex(_ *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()

	b.Widget().Render(ctx, "index", map[string]interface{}{
		"components": b.Components(),
	})
}
