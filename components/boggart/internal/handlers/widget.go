package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
)

type WidgetHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewWidgetHandler(component boggart.Component) *WidgetHandler {
	return &WidgetHandler{
		component: component,
	}
}

func (h *WidgetHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	id := r.URL().Query().Get(":id")

	if id == "" {
		h.NotFound(w, r)
		return
	}

	bind := h.component.Bind(id)
	if bind == nil {
		h.NotFound(w, r)
		return
	}

	if widget, ok := di.WidgetContainerBind(bind.Bind()); ok {
		widget.Handle(w, r)
	} else {
		h.NotFound(w, r)
	}
}
