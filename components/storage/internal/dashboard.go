package internal

import (
	"net/http"

	"github.com/kihamo/boggart/components/storage/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/:namespace/*filepath", handlers.NewFSHandler()).
				WithMethods([]string{http.MethodGet}),
		}
	}

	return c.routes
}
