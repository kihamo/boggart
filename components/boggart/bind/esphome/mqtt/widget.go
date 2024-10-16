package mqtt

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	action := r.URL().Query().Get("action")

	switch action {
	case "component":
		b.handleComponent(w, r)

	case "command":
		b.handleCommand(w, r)

	case "delete":
		b.handleDelete(w, r)

	case "config", "state":
		b.handleDump(w, r, action)

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

	if component.Type() == ComponentTypeLight {
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
			command := component.State().(*ComponentLightState)
			command.SetState(false)

			for key, value := range r.Original().PostForm {
				if len(value) == 0 {
					continue
				}

				if IsAllowLightStateField(component.ColorModes(), key) {
					switch key {
					case "state":
						command.SetState(value[0] == "on")
					case "color-mode":
						command.SetColorMode(value[0])
					case "brightness":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetBrightness(val)
						} else {
							err = e
						}
					case "red":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetRed(val)
						} else {
							err = e
						}
					case "green":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetGreen(val)
						} else {
							err = e
						}
					case "blue":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetBlue(val)
						} else {
							err = e
						}
					case "white":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetWhite(val)
						} else {
							err = e
						}
					case "cold-white":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetColdWhite(val)
						} else {
							err = e
						}
					case "color-temperature":
						if val, e := strconv.ParseUint(value[0], 10, 64); e == nil {
							command.SetColorTemperature(val)
						} else {
							err = e
						}
					case "effect":
						if len(component.EffectList()) > 0 {
							command.SetEffect(value[0])
						} else {
							err = errors.New("effect " + value[0] + " not exists")
						}
					case "flash-length":
						if val, e := time.ParseDuration(value[0]); e == nil {
							command.SetFlash(val)
						} else {
							err = e
						}
					case "transition-length":
						if val, e := time.ParseDuration(value[0]); e == nil {
							command.SetTransition(val)
						} else {
							err = e
						}
					}
				} else {
					err = errors.New("field \"" + key + "\" not allowed")
				}

				if err != nil {
					widget.FlashError(r, "Parse param %s failed with error %v", "", key, err)
					err = nil
				}
			}

			err = b.MQTT().PublishRawWithoutCache(ctx, component.CommandTopic(), 1, false, component.CommandToPayload(command))

			if err != nil {
				widget.FlashError(r, err, "")
			} else {
				widget.FlashSuccess(r, "Success toggle", "")

				widget.Redirect(r.URL().Path+"?action=component&id="+component.ID(), http.StatusFound, w, r)
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
	widget := b.Widget()
	if !b.Meta().Status().IsStatusOnline() {
		widget.NotFound(w, r)
		return
	}

	q := r.URL().Query()
	componentID := q.Get("id")
	command := q.Get("cmd")

	if componentID == "" || command == "" {
		widget.NotFound(w, r)
		return
	}

	component := b.Component(componentID)
	if component == nil || component.CommandTopic() == "" {
		widget.NotFound(w, r)
		return
	}

	ctx := r.Context()
	err := b.MQTT().PublishRawWithoutCache(ctx, component.CommandTopic(), 1, false, component.CommandToPayload(command))

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

func (b *Bind) handleDelete(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	componentID := r.URL().Query().Get("id")
	if componentID == "" {
		widget.NotFound(w, r)
		return
	}

	component := b.Component(componentID)
	if component == nil {
		widget.NotFound(w, r)
		return
	}

	ctx := r.Context()
	var isFailed bool

	if topic := component.ConfigMessage().Topic(); topic != "" {
		if err := b.MQTT().Delete(ctx, topic); err != nil {
			widget.FlashError(r, "Delete topic %v failed with error %v", "", topic, err)
			isFailed = true
		}
	}

	if topic := component.StateTopic(); topic != "" {
		if err := b.MQTT().Delete(ctx, topic); err != nil {
			widget.FlashError(r, "Delete topic %v failed with error %v", "", topic, err)
			isFailed = true
		}
	}

	if topic := component.CommandTopic(); topic != "" {
		if err := b.MQTT().Delete(ctx, topic); err != nil {
			widget.FlashError(r, "Delete topic %v failed with error %v", "", topic, err)
			isFailed = true
		}
	}

	if !isFailed {
		b.delete(component.UniqueID())
		widget.FlashSuccess(r, "Remove component %s success", "", componentID)
	}

	redirectURL := &url.URL{}
	*redirectURL = *r.Original().URL
	redirectURL.RawQuery = ""

	widget.Redirect(redirectURL.String(), http.StatusFound, w, r)
}

func (b *Bind) handleDump(w *dashboard.Response, r *dashboard.Request, entity string) {
	widget := b.Widget()

	componentID := r.URL().Query().Get("id")
	if componentID == "" {
		widget.NotFound(w, r)
		return
	}

	component := b.Component(componentID)
	if component == nil {
		widget.NotFound(w, r)
		return
	}

	var dump string

	switch entity {
	case "config":
		var v interface{}

		if err := component.ConfigMessage().JSONUnmarshal(&v); err != nil {
			widget.InternalError(w, r, err)
			return
		}

		output, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		dump = string(output)

	case "state":
		if cmp, ok := component.(HasStateRaw); ok {
			dump = cmp.StateRaw()
		} else {
			widget.NotFound(w, r)
			return
		}

	default:
		widget.NotFound(w, r)
		return
	}

	b.Widget().RenderLayout(r.Context(), "dump", "simple", map[string]interface{}{
		"dump":  dump,
		"modal": true,
	})
}

func (b *Bind) handleIndex(_ *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()

	b.Widget().Render(ctx, "index", map[string]interface{}{
		"components": b.Components(),
	})
}
