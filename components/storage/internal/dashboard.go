package internal

import (
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/storage/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
	"net/http"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(storage.RouteFileStoragePrefix+":namespace/*filepath", &handlers.FSHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}),
		}
	}

	return c.routes
}
