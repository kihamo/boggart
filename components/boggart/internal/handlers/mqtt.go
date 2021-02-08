package handlers

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

type MQTTHandler struct {
	dashboard.Handler

	componentBoggart boggart.Component
	componentMQTT    mqtt.Component
}

func NewMQTTHandler(b boggart.Component, m mqtt.Component) *MQTTHandler {
	return &MQTTHandler{
		componentBoggart: b,
		componentMQTT:    m,
	}
}

func (h *MQTTHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	id := r.URL().Query().Get(":id")

	if id == "" {
		h.NotFound(w, r)
		return
	}

	bindItem := h.componentBoggart.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	bindSupport, ok := bindItem.Bind().(di.MQTTContainerSupport)
	if !ok {
		h.NotFound(w, r)
		return
	}

	type itemView struct {
		Topic    string
		Calls    uint64
		Datetime time.Time
		Payload  interface{}
	}

	publishes := bindSupport.MQTT().Publishes()
	publishesItems := make([]itemView, 0, len(publishes))

	for _, item := range h.componentMQTT.CacheItems() {
		for sent, count := range publishes {
			if item.Topic().String() == sent.String() {
				view := itemView{
					Topic:    sent.String(),
					Calls:    count,
					Datetime: item.Datetime(),
					Payload:  item.Payload(),
				}

				publishesItems = append(publishesItems, view)

				break
			}
		}
	}

	subscribers := bindSupport.MQTT().Subscribers()

	h.Render(r.Context(), "mqtt", map[string]interface{}{
		"bind":        bindItem,
		"publishes":   publishesItems,
		"subscribers": subscribers,
	})
}
