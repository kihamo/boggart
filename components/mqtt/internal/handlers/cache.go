package handlers

import (
	"net/http"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

type CacheHandler struct {
	dashboard.Handler

	PayloadCache mqtt.Cache
}

func (h *CacheHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	if r.IsPost() {
		h.PayloadCache.Purge()

		h.Redirect(r.URL().Path, http.StatusFound, w, r)
		return
	}

	h.Render(r.Context(), "cache", map[string]interface{}{
		"payloads": h.PayloadCache.Payloads(),
	})
}
