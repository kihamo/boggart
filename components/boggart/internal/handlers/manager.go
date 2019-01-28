package handlers

import (
	"bytes"
	"time"

	"github.com/kihamo/boggart/components/boggart/internal/manager"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	listeners "github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

// easyjson:json
type managerHandlerDevice struct {
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
	HasWidget       bool     `json:"has_widget"`
}

// easyjson:json
type managerListener struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Events     map[string]string `json:"events"`
	Fires      int64             `json:"fires"`
	FiredFirst *time.Time        `json:"fire_first"`
	FiredLast  *time.Time        `json:"fire_last"`
}

// TODO: rename to ManagerHandler
type ManagerHandler struct {
	dashboard.Handler

	manager          *manager.Manager
	listenersManager *listeners.ListenersManager
}

func NewManagerHandler(manager *manager.Manager, listenersManager *listeners.ListenersManager) *ManagerHandler {
	return &ManagerHandler{
		manager:          manager,
		listenersManager: listenersManager,
	}
}

func (h *ManagerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), "manager", map[string]interface{}{
			"device_types": boggart.GetBindTypes(),
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
		list := make([]managerHandlerDevice, 0, 0)
		buf := bytes.NewBuffer(nil)
		enc := yaml.NewEncoder(buf)
		defer enc.Close()

		for _, bindItem := range h.manager.BindItems() {
			buf.Reset()
			if err := enc.Encode(bindItem); err != nil {
				panic(err.Error())
			}

			item := managerHandlerDevice{
				Id:              bindItem.ID(),
				Type:            bindItem.Type(),
				Description:     bindItem.Description(),
				SerialNumber:    bindItem.Bind().SerialNumber(),
				Status:          bindItem.Status().String(),
				Tags:            bindItem.Tags(),
				Tasks:           make([]string, 0, len(bindItem.Tasks())),
				MQTTPublishes:   make([]string, 0, len(bindItem.MQTTPublishes())),
				MQTTSubscribers: make([]string, 0, len(bindItem.MQTTSubscribers())),
				Config:          buf.String(),
			}

			if _, ok := bindItem.BindType().(boggart.BindTypeHasWidget); ok {
				item.HasWidget = ok
			}

			for _, task := range bindItem.Tasks() {
				item.Tasks = append(item.Tasks, task.Name())
			}

			for _, topic := range bindItem.MQTTPublishes() {
				item.MQTTPublishes = append(item.MQTTPublishes, topic.String())
			}

			for _, topic := range bindItem.MQTTSubscribers() {
				item.MQTTSubscribers = append(item.MQTTSubscribers, topic.Topic())
			}

			list = append(list, item)
		}

		entities.Data = list
		entities.Total = len(list)

	case "listeners":
		list := make([]managerListener, 0, 0)

		for _, l := range h.listenersManager.Listeners() {
			item := managerListener{
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
