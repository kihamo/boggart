package handlers

import (
	"github.com/kihamo/shadow/components/dashboard"
)

type DetectHandler struct {
	dashboard.Handler
}

func (h *DetectHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "detect", nil)
}
