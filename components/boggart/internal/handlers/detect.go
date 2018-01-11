package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type DetectHandler struct {
	dashboard.Handler
}

func (h *DetectHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), boggart.ComponentName, "detect", nil)
}
