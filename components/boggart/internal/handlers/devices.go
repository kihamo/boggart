package handlers

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/manager"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type deviceHandlerDevice struct {
	RegisterId      string                 `json:"register_id"`
	Id              string                 `json:"id"`
	Description     string                 `json:"description"`
	SerialNumber    string                 `json:"serial_number"`
	Status          string                 `json:"status"`
	Tasks           []string               `json:"tasks"`
	MQTTTopics      []string               `json:"mqtt_topics"`
	MQTTSubscribers []string               `json:"mqtt_subscribers"`
	Tags            []string               `json:"tags"`
	Config          map[string]interface{} `json:"config"`
}

// easyjson:json
type deviceHandlerListener struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Events     map[string]string `json:"events"`
	Fires      int64             `json:"fires"`
	FiredFirst *time.Time        `json:"fire_first"`
	FiredLast  *time.Time        `json:"fire_last"`
}

type DevicesHandler struct {
	dashboard.Handler

	devicesManager   boggart.DevicesManager
	listenersManager *manager.ListenersManager
}

func NewDevicesHandler(devicesManager boggart.DevicesManager, listenersManager *manager.ListenersManager) *DevicesHandler {
	return &DevicesHandler{
		devicesManager:   devicesManager,
		listenersManager: listenersManager,
	}
}

func (h *DevicesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if !r.IsAjax() {
		h.Render(r.Context(), "devices", nil)
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
		list := make([]deviceHandlerDevice, 0, 0)

		for registerId, d := range h.devicesManager.Devices() {
			bind := d.Bind()

			item := deviceHandlerDevice{
				RegisterId:      registerId,
				Id:              d.Id(),
				Description:     d.Description(),
				Status:          bind.Status().String(),
				Tags:            make([]string, 0, len(d.Tags())),
				Tasks:           make([]string, 0),
				MQTTTopics:      make([]string, 0),
				MQTTSubscribers: make([]string, 0),
				Config:          d.Config(),
			}

			if sn, ok := bind.(boggart.DeviceHasSerialNumber); ok {
				item.SerialNumber = sn.SerialNumber()
			}

			if tasks, ok := bind.(boggart.DeviceHasTasks); ok {
				for _, task := range tasks.Tasks() {
					item.Tasks = append(item.Tasks, task.Name())
				}
			}

			item.Tags = append(item.Tags, d.Tags()...)

			if topics, ok := bind.(boggart.DeviceHasMQTTTopics); ok {
				for _, topic := range topics.MQTTTopics() {
					item.MQTTTopics = append(item.MQTTTopics, topic.String())
				}
			}

			if subscribers, ok := bind.(boggart.DeviceHasMQTTSubscribers); ok {
				for _, topic := range subscribers.MQTTSubscribers() {
					item.MQTTSubscribers = append(item.MQTTSubscribers, topic.Topic())
				}
			}

			list = append(list, item)
		}

		entities.Data = list
		entities.Total = len(list)

	case "listeners":
		list := make([]deviceHandlerListener, 0, 0)

		for _, l := range h.listenersManager.Listeners() {
			item := deviceHandlerListener{
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
