package handlers

import (
	"bytes"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

// easyjson:json
type managerIndexHandlerDevice struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	Description     string   `json:"description"`
	SerialNumber    string   `json:"serial_number"`
	Status          string   `json:"status"`
	Tasks           []string `json:"tasks"`
	MQTTPublishes   []string `json:"mqtt_publishes"`
	MQTTSubscribers []string `json:"mqtt_subscribers"`
	Tags            []string `json:"tags"`
	Config          string   `json:"config"`
}

// easyjson:json
type managerIndexListener struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Events     map[string]string `json:"events"`
	Fires      int64             `json:"fires"`
	FiredFirst *time.Time        `json:"fire_first"`
	FiredLast  *time.Time        `json:"fire_last"`
}

// TODO: rename to ManagerHandler
type ManagerIndexHandler struct {
	dashboard.Handler

	devicesManager   boggart.DevicesManager
	listenersManager *manager.ListenersManager
}

func NewManagerIndexHandler(devicesManager boggart.DevicesManager, listenersManager *manager.ListenersManager) *ManagerIndexHandler {
	return &ManagerIndexHandler{
		devicesManager:   devicesManager,
		listenersManager: listenersManager,
	}
}

func (h *ManagerIndexHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), "manager_index", map[string]interface{}{
			"device_types": boggart.GetDeviceTypes(),
		})
		return
	}

	entities := struct {
		Draw     int         `json:"draw"`
		Total    int         `json:"recordsTotal"`
		Filtered int         `json:"recordsFiltered"`
		Data     interface{} `json:"data"`
		Error    string      `json:"error,omitempty"`
	}{}
	entities.Draw = 1

	switch r.URL().Query().Get("entity") {
	case "devices":
		list := make([]managerIndexHandlerDevice, 0, 0)
		buf := bytes.NewBuffer(nil)
		enc := yaml.NewEncoder(buf)
		defer enc.Close()

		for _, d := range h.devicesManager.Devices() {
			buf.Reset()
			if err := enc.Encode(d.Config()); err != nil {
				panic(err.Error())
			}

			item := managerIndexHandlerDevice{
				Id:              d.ID(),
				Type:            d.Type(),
				Description:     d.Description(),
				SerialNumber:    d.Bind().SerialNumber(),
				Status:          d.Bind().Status().String(),
				Tags:            d.Tags(),
				Tasks:           make([]string, 0, len(d.Tasks())),
				MQTTPublishes:   make([]string, 0, len(d.MQTTPublishes())),
				MQTTSubscribers: make([]string, 0, len(d.MQTTSubscribers())),
				Config:          buf.String(),
			}

			for _, task := range d.Tasks() {
				item.Tasks = append(item.Tasks, task.Name())
			}

			for _, topic := range d.MQTTPublishes() {
				item.MQTTPublishes = append(item.MQTTPublishes, topic.String())
			}

			for _, topic := range d.MQTTSubscribers() {
				item.MQTTSubscribers = append(item.MQTTSubscribers, topic.Topic())
			}

			list = append(list, item)
		}

		entities.Data = list
		entities.Total = len(list)

	case "listeners":
		list := make([]managerIndexListener, 0, 0)

		for _, l := range h.listenersManager.Listeners() {
			item := managerIndexListener{
				Id:     l.Id(),
				Name:   l.Name(),
				Events: make(map[string]string, 0),
			}

			listener := h.listenersManager.GetById(l.Id())
			if listener == nil {
				continue
			}

			md := listener.Metadata()
			item.Fires = md[workers.ListenerMetadataFires].(int64)
			item.FiredFirst = md[workers.ListenerMetadataFirstFiredAt].(*time.Time)
			item.FiredLast = md[workers.ListenerMetadataLastFireAt].(*time.Time)

			events := md[workers.ListenerMetadataEvents].([]workers.Event)
			for _, event := range events {
				item.Events[event.Id()] = event.Name()
			}

			list = append(list, item)
		}

		entities.Data = list
		entities.Total = len(list)

	default:
		h.NotFound(w, r)
		return
	}

	entities.Filtered = entities.Total
	w.SendJSON(entities)
}
