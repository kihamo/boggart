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

type configYAMLDevice struct {
	Type        string
	ID          string
	Description string
	Tags        []string
	Config      interface{}
}

type configYAML struct {
	Devices []configYAMLDevice
}

type ConfigHandler struct {
	dashboard.Handler

	devicesManager boggart.DevicesManager
}

func NewConfigHandler(devicesManager boggart.DevicesManager) *ConfigHandler {
	return &ConfigHandler{
		devicesManager: devicesManager,
	}
}

func (h *ConfigHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	action := r.URL().Query().Get(":action")
	if action != "download" && action != "modal" && action != "view" {
		h.NotFound(w, r)
	}

	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)
	devices := h.devicesManager.Devices()

	config := &configYAML{
		Devices: make([]configYAMLDevice, 0, len(devices)),
	}

	for _, d := range devices {
		config.Devices = append(config.Devices, configYAMLDevice{
			Type:        d.Type(),
			ID:          d.ID(),
			Description: d.Description(),
			Tags:        d.Tags(),
			Config:      d.Config(),
		})
	}

	if err := enc.Encode(config); err != nil {
		panic(err.Error())
	}
	enc.Close()

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
