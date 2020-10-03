package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
	"go.uber.org/zap/zapcore"
)

// easyjson:json
type managerHandlerDevice struct {
	Tasks             [][]string    `json:"tasks"`
	Tags              []string      `json:"tags"`
	ID                string        `json:"id"`
	Type              string        `json:"type"`
	Description       string        `json:"description"`
	SerialNumber      string        `json:"serial_number"`
	MAC               string        `json:"mac"`
	Status            string        `json:"status"`
	MQTTPublishes     int           `json:"mqtt_publishes"`
	MQTTSubscribers   int           `json:"mqtt_subscribers"`
	LogsCount         int           `json:"logs_count"`
	HasWidget         bool          `json:"has_widget"`
	HasReadinessProbe bool          `json:"has_readiness_probe"`
	HasLivenessProbe  bool          `json:"has_liveness_probe"`
	LogsMaxLevel      zapcore.Level `json:"logs_max_level"`
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
				ID:           bindItem.ID(),
				Type:         bindItem.Type(),
				Description:  bindItem.Description(),
				Status:       bindItem.Status().String(),
				Tags:         bindItem.Tags(),
				LogsMaxLevel: zapcore.DebugLevel,
			}

			if bindSupport, ok := di.WidgetContainerBind(bindItem.Bind()); ok {
				item.HasWidget = ok && bindSupport.HandleAllowed()
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
				item.Tasks = make([][]string, 0, len(bindSupport.Tasks()))
				for _, t := range bindSupport.Tasks() {
					name := bindSupport.TaskShortName(t)
					if !bindSupport.TaskRegisteredInQueue(t) {
						name += " (!)"
					}

					item.Tasks = append(item.Tasks, []string{
						t.Id(), name,
					})
				}
			}

			if bindSupport, ok := di.MQTTContainerBind(bindItem.Bind()); ok {
				item.MQTTSubscribers = len(bindSupport.Subscribers())
				item.MQTTPublishes = len(bindSupport.Publishes())
			}

			if bindSupport, ok := bindItem.Bind().(di.LoggerContainerSupport); ok {
				records := bindSupport.LastRecords()
				item.LogsCount = len(records)

				for _, r := range records {
					if r.Level > item.LogsMaxLevel {
						item.LogsMaxLevel = r.Level
					}
				}
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
