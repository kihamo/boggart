package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/shadow/components/dashboard"
)

type WidgetHandler struct {
	dashboard.Handler

	manager *manager.Manager
}

func NewWidgetHandler(manager *manager.Manager) *WidgetHandler {
	return &WidgetHandler{
		manager: manager,
	}
}

func (h *WidgetHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	id := r.URL().Query().Get(":id")

	if id == "" {
		h.NotFound(w, r)
		return
	}

	bind := h.manager.Bind(id)
	if bind == nil {
		h.NotFound(w, r)
		return
	}

	widget, ok := bind.BindType().(boggart.BindTypeHasWidget)
	if !ok {
		h.NotFound(w, r)
		return
	}

	r = r.WithContext(dashboard.ContextWithTemplateNamespace(r.Context(), boggart.ComponentName+"-bind-"+bind.Type()))
	widget.Widget(w, r, bind)
}
