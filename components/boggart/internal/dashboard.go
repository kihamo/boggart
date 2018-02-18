package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()
	menus := []dashboard.Menu{
		dashboard.NewMenuWithRoute("Dashboard", routes[0], "", nil, nil),
		dashboard.NewMenuWithRoute("Detect", routes[1], "", nil, nil),
		dashboard.NewMenuWithRoute("Devices", routes[2], "", nil, nil),
	}

	if u := c.config.String(boggart.ConfigMonitoringExternalURL); u != "" {
		menus = append(menus, dashboard.NewMenuWithUrl("Monitoring", u, "", nil, nil))
	}

	return dashboard.NewMenuWithUrl("Boggart", "/"+c.Name()+"/", "home", menus, nil)
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/"+c.Name()+"/",
				&handlers.IndexHandler{
					Config:    c.config,
					Component: c,
				},
				"",
				true),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/"+c.Name()+"/detect/",
				&handlers.DetectHandler{},
				"",
				true),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/"+c.Name()+"/devices/",
				&handlers.DevicesHandler{
					DevicesManager: c.devicesManager,
				},
				"",
				true),
			dashboard.NewRoute(
				c.Name(),
				[]string{http.MethodGet},
				"/"+c.Name()+"/camera/:place/:action/",
				&handlers.CameraHandler{
					DevicesManager: c.devicesManager,
				},
				"",
				true),
		}
	}

	return c.routes
}
