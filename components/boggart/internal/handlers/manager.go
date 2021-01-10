package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
	"go.uber.org/zap/zapcore"
)

// easyjson:json
type managerHandlerDeviceTask struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Registered bool   `json:"registered"`
}

// easyjson:json
type managerHandlerDevice struct {
	Tasks                    []managerHandlerDeviceTask `json:"tasks"`
	Tags                     []string                   `json:"tags"`
	ID                       string                     `json:"id"`
	Type                     string                     `json:"type"`
	Description              string                     `json:"description"`
	SerialNumber             string                     `json:"serial_number"`
	MAC                      string                     `json:"mac"`
	Status                   string                     `json:"status"`
	ProbeReadiness           string                     `json:"probe_readiness"`
	ProbeLiveness            string                     `json:"probe_liveness"`
	MetricsDescriptionsCount uint64                     `json:"metrics_descriptions_count"`
	MetricsCollectCount      uint64                     `json:"metrics_collect_count"`
	MetricsEmptyCount        uint64                     `json:"metrics_empty_count"`
	MQTTPublishes            int                        `json:"mqtt_publishes"`
	MQTTSubscribers          int                        `json:"mqtt_subscribers"`
	LogsCount                int                        `json:"logs_count"`
	HasMetrics               bool                       `json:"has_metrics"`
	HasWidget                bool                       `json:"has_widget"`
	LogsMaxLevel             zapcore.Level              `json:"logs_max_level"`
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

	entities := boggart.NewResponseDataTable()

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
				item.ProbeReadiness = bindSupport.ReadinessTaskID()
				item.ProbeLiveness = bindSupport.LivenessTaskID()
			}

			if bindSupport, ok := di.MetaContainerBind(bindItem.Bind()); ok {
				item.SerialNumber = bindSupport.SerialNumber()

				if mac := bindSupport.MAC(); mac != nil {
					item.MAC = mac.String()
				}
			}

			if bindSupport, ok := di.WorkersContainerBind(bindItem.Bind()); ok {
				ids := bindSupport.TasksID()
				item.Tasks = make([]managerHandlerDeviceTask, len(ids))
				for i, id := range ids {
					item.Tasks[i].ID = id[0]
					item.Tasks[i].Name = bindSupport.TaskShortName(id[1])
					item.Tasks[i].Registered = bindSupport.TaskRegisteredInQueue(id[0])
				}
			}

			if bindSupport, ok := di.MetricsContainerBind(bindItem.Bind()); ok {
				item.HasMetrics = true
				item.MetricsDescriptionsCount = bindSupport.DescriptionsCount()
				item.MetricsCollectCount = bindSupport.CollectCount()
				item.MetricsEmptyCount = bindSupport.EmptyCount()
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
