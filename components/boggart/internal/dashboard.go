package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	menus := dashboard.NewMenu("Smart home").
		WithIcon("home").
		WithChild(dashboard.NewMenu("Dashboard").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Security").WithRoute(routes[2])).
		WithChild(dashboard.NewMenu("Devices").WithRoute(routes[3]))

	if u := c.config.String(boggart.ConfigMonitoringExternalURL); u != "" {
		menus = menus.WithChild(dashboard.NewMenu("Monitoring").WithUrl(u))
	}

	return menus
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/", &handlers.IndexHandler{}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/security/", handlers.NewSecurityHandler(c.securityManager)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/devices/", handlers.NewDevicesHandler(c.devicesManager, c.listenersManager)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/devices/:device/:action", handlers.NewDeviceHandler(c.devicesManager)).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/camera/:place/:action/", handlers.NewCameraHandler(c.devicesManager)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
		}
	}

	return c.routes
}
