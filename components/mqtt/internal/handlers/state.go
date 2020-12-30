package handlers

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

type StateHandler struct {
	dashboard.Handler

	Component mqtt.Component
}

//nolint
func (h *StateHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	client := h.Component.Client()

	h.Render(r.Context(), "state", map[string]interface{}{
		"is_connected":       client.IsConnected(),
		"is_connection_open": client.IsConnectionOpen(),
	})
}
