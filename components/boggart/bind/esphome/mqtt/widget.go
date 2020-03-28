package mqtt

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	q := r.URL().Query()

	switch q.Get("action") {
	case "component":
		t.handleComponent(w, r, bind)

	case "command":
		t.handleCommand(w, r, bind)

	default:
		t.handleIndex(w, r, bind)
	}
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (t Type) handleComponent(w *dashboard.Response, r *dashboard.Request, bind *Bind) {
	q := r.URL().Query()

	componentID := q.Get("id")

	if componentID == "" {
		t.NotFound(w, r)
		return
	}

	component := bind.Component(componentID)
	if component == nil {
		t.NotFound(w, r)
		return
	}

	if component.GetType() == ComponentTypeLight {
		if light, ok := component.(*ComponentLight); ok {
			t.handleComponentLight(w, r, bind, light)
			return
		}
	}

	t.NotFound(w, r)
}

func (t Type) handleComponentLight(w *dashboard.Response, r *dashboard.Request, bind *Bind, component *ComponentLight) {
	ctx := r.Context()

	var err error

	if r.IsPost() {
		err = r.Original().ParseForm()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Parse form failed with error %s", "", err.Error()))
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
					r.Session().FlashBag().Error(t.Translate(ctx, "Parse param %s failed with error %s", "", key, err.Error()))
					err = nil
				}
			}

			ctx := r.Context()
			err := bind.MQTT().PublishRawWithoutCache(ctx, component.GetCommandTopic(), 1, false, component.CommandToPayload(command))

			if err != nil {
				r.Session().FlashBag().Error(err.Error())
			} else {
				r.Session().FlashBag().Success(t.Translate(ctx, "Success toggle", ""))

				t.Redirect(r.URL().Path+"?action=component&id="+component.GetID(), http.StatusFound, w, r)
				return
			}
		}
	}

	t.Render(ctx, "light", map[string]interface{}{
		"component": component,
		"error":     err,
	})
}

func (t Type) handleCommand(w *dashboard.Response, r *dashboard.Request, bind *Bind) {
	q := r.URL().Query()

	componentID := q.Get("id")
	command := q.Get("cmd")

	if componentID == "" || command == "" {
		t.NotFound(w, r)
		return
	}

	component := bind.Component(componentID)
	if component == nil || component.GetCommandTopic() == "" {
		t.NotFound(w, r)
		return
	}

	ctx := r.Context()
	err := bind.MQTT().PublishRawWithoutCache(ctx, component.GetCommandTopic(), 1, false, component.CommandToPayload(command))

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	} else {
		r.Session().FlashBag().Success(t.Translate(ctx, "Success toggle", ""))
	}

	redirectURL := &url.URL{}
	*redirectURL = *r.Original().URL
	redirectURL.RawQuery = ""

	t.Redirect(redirectURL.String(), http.StatusFound, w, r)
}

func (t Type) handleIndex(_ *dashboard.Response, r *dashboard.Request, bind *Bind) {
	ctx := r.Context()

	t.Render(ctx, "index", map[string]interface{}{
		"components": bind.Components(),
	})
}
