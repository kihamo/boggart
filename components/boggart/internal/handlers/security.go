package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type SecurityHandler struct {
	dashboard.Handler

	securityManager boggart.SecurityManager
}

func NewSecurityHandler(securityManager boggart.SecurityManager) *SecurityHandler {
	return &SecurityHandler{
		securityManager: securityManager,
	}
}

func (h *SecurityHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	h.Render(r.Context(), "security", map[string]interface{}{
		"status": h.securityManager.Status(),
	})
}
