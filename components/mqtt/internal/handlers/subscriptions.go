package handlers

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

type SubscriptionsHandler struct {
	dashboard.Handler

	Component mqtt.Component
}

func (h *SubscriptionsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "subscriptions", map[string]interface{}{
		"subscriptions": h.Component.Subscriptions(),
	})
}
