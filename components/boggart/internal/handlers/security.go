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
	switch r.Original().FormValue("status") {
	case boggart.SecurityStatusOpen.String():
		h.securityManager.Open()

	case boggart.SecurityStatusClosed.String():
		h.securityManager.Close()

	case boggart.SecurityStatusOpenForce.String():
		h.securityManager.OpenForce()

	case boggart.SecurityStatusClosedForce.String():
		h.securityManager.CloseForce()
	}
}
