package internal

import (
	"net/http"

	"github.com/kihamo/boggart/components/barcode/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/decode/", &handlers.DecodeHandler{}).
				WithMethods([]string{http.MethodGet}),
		}
	}

	return c.routes
}
