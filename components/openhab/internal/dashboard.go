package internal

import (
	"net/http"

	"github.com/kihamo/boggart/components/openhab/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	return []dashboard.Route{
		dashboard.NewRoute("/"+c.Name()+"/proxy/:type/:id/:key/", &handlers.ProxyHandler{}).
			WithMethods([]string{http.MethodGet}),
	}
}
