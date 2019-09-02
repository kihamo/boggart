package handlers

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

type CacheHandler struct {
	dashboard.Handler

	PayloadCache mqtt.Cache
}

func (h *CacheHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "cache", map[string]interface{}{
		"payloads": h.PayloadCache.Payloads(),
	})
}
