package handlers

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/shadow/components/dashboard"
)

type deviceViewHandlerDevice struct {
	TasksCount  int
	Id          string
	Description string
	Types       []string
	Enabled     bool
}

type listenerViewHandlerDevice struct {
	Id           string
	Name         string
	Events       map[string]string
	Fires        int64
	FirstFiredAt *time.Time
	LastFiredAt  *time.Time
}

type DevicesHandler struct {
	dashboard.Handler
}

func (h *DevicesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	dm := r.Component().(boggart.Component).DevicesManager()

	// devices
	devicesList := dm.Devices()
	devicesListView := make([]*deviceViewHandlerDevice, 0, len(devicesList))

	for _, d := range devicesList {
		viewDevice := &deviceViewHandlerDevice{
			Id:          d.Id(),
			Description: d.Description(),
			Types:       make([]string, 0, len(d.Types())),
			TasksCount:  len(d.Tasks()),
			Enabled:     d.IsEnabled(),
		}

		for _, t := range d.Types() {
			viewDevice.Types = append(viewDevice.Types, t.String())
		}

		devicesListView = append(devicesListView, viewDevice)
	}

	// listeners
	listenersListView := make([]listenerViewHandlerDevice, 0, 0)

	for _, item := range dm.Listeners() {
		listener := listenerViewHandlerDevice{
			Id:     item.Id(),
			Name:   item.Name(),
			Events: make(map[string]string, 0),
		}

		md := dm.GetListenerMetadata(item.Id())
		if md == nil {
			continue
		}

		listener.Fires = md[workers.ListenerMetadataFires].(int64)
		listener.FirstFiredAt = md[workers.ListenerMetadataFirstFiredAt].(*time.Time)
		listener.LastFiredAt = md[workers.ListenerMetadataLastFireAt].(*time.Time)

		events := md[workers.ListenerMetadataEvents].([]workers.Event)
		for _, event := range events {
			listener.Events[event.Id()] = event.Name()
		}

		listenersListView = append(listenersListView, listener)
	}

	h.Render(r.Context(), "devices", map[string]interface{}{
		"devices":   devicesListView,
		"listeners": listenersListView,
	})
}
