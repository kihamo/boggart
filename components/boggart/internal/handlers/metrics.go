package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
)

type MetricsHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewMetricsHandler(b boggart.Component) *MetricsHandler {
	return &MetricsHandler{
		component: b,
	}
}

func (h *MetricsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	id := r.URL().Query().Get(":id")

	if id == "" {
		h.NotFound(w, r)
		return
	}

	bindItem := h.component.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	bindSupport, ok := bindItem.Bind().(di.MetricsContainerSupport)
	if !ok {
		h.NotFound(w, r)
		return
	}

	measures, err := bindSupport.Metrics().Gather()
	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	h.Render(r.Context(), "metrics", map[string]interface{}{
		"bind":     bindItem,
		"measures": measures,
	})
}
