package internal

import (
	"net/http"

	"github.com/kihamo/boggart/components/openhab/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	proxyHandler := &handlers.ProxyHandler{
		Component: c,
	}

	m := c.application.GetComponent(messengers.ComponentName)
	if m != nil {
		proxyHandler.Messengers = m.(messengers.Component)
	}

	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/proxy/:type/:id/:key/", proxyHandler).
			WithMethods([]string{http.MethodGet, http.MethodPost}),
	}
}
