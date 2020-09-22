package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type managerHandlerDevice struct {
	ID                string            `json:"id"`
	Type              string            `json:"type"`
	Description       string            `json:"description"`
	SerialNumber      string            `json:"serial_number"`
	MAC               string            `json:"mac"`
	Status            string            `json:"status"`
	Tasks             map[string]string `json:"tasks"`
	MQTTPublishes     int               `json:"mqtt_publishes"`
	MQTTSubscribers   int               `json:"mqtt_subscribers"`
	Tags              []string          `json:"tags"`
	HasWidget         bool              `json:"has_widget"`
	HasReadinessProbe bool              `json:"has_readiness_probe"`
	HasLivenessProbe  bool              `json:"has_liveness_probe"`
	LogsCount         int               `json:"logs_count"`
}

type ManagerHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewManagerHandler(component boggart.Component) *ManagerHandler {
	return &ManagerHandler{
		component: component,
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
		list := make([]managerHandlerDevice, 0)

		for _, bindItem := range h.component.BindItems() {
			item := managerHandlerDevice{
				ID:          bindItem.ID(),
				Type:        bindItem.Type(),
				Description: bindItem.Description(),
				Status:      bindItem.Status().String(),
				Tags:        bindItem.Tags(),
			}

			if _, ok := di.WidgetContainerBind(bindItem.Bind()); ok {
				item.HasWidget = ok
			}

			if bindSupport, ok := di.ProbesContainerBind(bindItem.Bind()); ok {
				item.HasReadinessProbe = bindSupport.Readiness() != nil
				item.HasLivenessProbe = bindSupport.Liveness() != nil
			}

			if bindSupport, ok := di.MetaContainerBind(bindItem.Bind()); ok {
				item.SerialNumber = bindSupport.SerialNumber()

				if mac := bindSupport.MAC(); mac != nil {
					item.MAC = mac.String()
				}
			}

			if bindSupport, ok := di.WorkersContainerBind(bindItem.Bind()); ok {
				item.Tasks = make(map[string]string, len(bindSupport.Tasks()))
				for _, t := range bindSupport.Tasks() {
					item.Tasks[t.Id()] = bindSupport.TaskShortName(t)
				}
			}

			if bindSupport, ok := di.MQTTContainerBind(bindItem.Bind()); ok {
				item.MQTTSubscribers = len(bindSupport.Subscribers())
				item.MQTTPublishes = len(bindSupport.Publishes())
			}

			if bindSupport, ok := bindItem.Bind().(di.LoggerContainerSupport); ok {
				item.LogsCount = len(bindSupport.LastRecords())
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
	_ = w.SendJSON(entities)
}
