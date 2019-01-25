package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Smart home").
		WithIcon("home").
		WithChild(dashboard.NewMenu("Manager").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Config YAML").WithUrl("/" + c.Name() + "/config/view"))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		<-c.application.ReadyComponent(c.Name())

		bindHandler := handlers.NewBindHandler(c.manager)

		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/manager/", handlers.NewManagerHandler(c.manager, c.listenersManager)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/bind/", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/bind/:id/", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/bind/:id/:action", bindHandler).
				WithMethods([]string{http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/camera/:id/:channel", handlers.NewCameraHandler(c.manager)).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
			dashboard.NewRoute("/"+c.Name()+"/config/:action", handlers.NewConfigHandler(c.manager, c)).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
		}
	}

	return c.routes
}
