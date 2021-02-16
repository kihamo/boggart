package handlers

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

type MarshalShort interface {
	MarshalShortYAML() (interface{}, error)
}

type MarshallerContainer struct {
	short bool
	value interface{}
}

func (m MarshallerContainer) MarshalYAML() (interface{}, error) {
	if m.short {
		if short, ok := m.value.(MarshalShort); ok {
			return short.MarshalShortYAML()
		}
	}

	if long, ok := m.value.(yaml.Marshaler); ok {
		return long.MarshalYAML()
	}

	return m.value, nil
}

type BindItemsList []interface{}

func (l BindItemsList) MarshalYAML() (interface{}, error) {
	return struct {
		Devices []interface{}
	}{
		Devices: l,
	}, nil
}

type ConfigHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewConfigHandler(component boggart.Component) *ConfigHandler {
	return &ConfigHandler{
		component: component,
	}
}

func (h *ConfigHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	action := q.Get(":action")
	id := q.Get(":id")

	if action == "reload" && r.IsPost() {
		// reload by ID
		if id != "" {
			bind := h.component.Bind(id)
			if bind == nil {
				_ = w.SendJSON(boggart.NewResponseJSON().Failed("Bind " + id + " not found"))
			} else if bind.Type() == boggart.ComponentName {
				_ = w.SendJSON(boggart.NewResponseJSON().Failed("Can't reload config from special bind " + boggart.ComponentName))
			} else if err := h.component.ReloadConfigByID(id); err != nil {
				_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
			} else {
				_ = w.SendJSON(boggart.NewResponseJSON().Success("Bind " + id + " reloaded from file"))
			}

			return
		}

		// reload all
		if loaded, err := h.component.ReloadConfig(); err != nil {
			_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
		} else {
			_ = w.SendJSON(boggart.NewResponseJSON().Success("Loaded " + strconv.FormatInt(int64(loaded), 10) + " binds"))
		}

		return
	}

	if action != "download" && action != "modal" && action != "view" {
		h.NotFound(w, r)
		return
	}

	buf := bytes.NewBuffer(nil)

	enc := yaml.NewEncoder(buf)
	defer enc.Close()

	var (
		value interface{}
		short bool
	)

	if val := q.Get("short"); val != "" {
		short = true
	}

	if id != "" {
		bind := h.component.Bind(id)

		if bind == nil {
			h.NotFound(w, r)
			return
		}

		value = MarshallerContainer{
			short: short,
			value: bind,
		}
	} else {
		binds := h.component.BindItems()
		list := make([]interface{}, 0, len(binds))

		for _, bind := range binds {
			list = append(list, MarshallerContainer{
				short: short,
				value: bind,
			})
		}

		value = BindItemsList(list)
	}

	if err := enc.Encode(value); err != nil {
		h.InternalError(w, r, err)
		return
	}

	switch action {
	case "download":
		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
		w.Header().Set("Content-Disposition", "attachment; filename=\""+time.Now().Format("20060102150405_config.yaml")+"\"")
		w.Header().Set("Content-Type", "text/x-yaml")

		if _, err := io.Copy(w, buf); err != nil {
			panic(err.Error())
		}

	case "modal":
		h.RenderLayout(r.Context(), "config", "simple", map[string]interface{}{
			"yaml":  buf.String(),
			"modal": true,
		})

	case "view":
		h.Render(r.Context(), "config", map[string]interface{}{
			"yaml":  buf.String(),
			"modal": false,
		})
	}
}
