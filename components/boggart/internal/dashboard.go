package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) GetTemplates() *assetfs.AssetFS {
	return &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "templates",
	}
}

func (c *Component) GetDashboardMenu() dashboard.Menu {
	routes := c.GetDashboardRoutes()

	return dashboard.NewMenuWithUrl(
		"Boggart",
		"/"+c.GetName()+"/",
		"home",
		[]dashboard.Menu{
			dashboard.NewMenuWithRoute("Dashboard", routes[0], "", nil, nil),
			dashboard.NewMenuWithRoute("Devices", routes[1], "", nil, nil),
		},
		nil)
}

func (c *Component) GetDashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.NewRoute(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/",
				&handlers.IndexHandler{
					Config:    c.config,
					Collector: c,
				},
				"",
				true),
			dashboard.NewRoute(
				c.GetName(),
				[]string{http.MethodGet},
				"/"+c.GetName()+"/devices/",
				&handlers.DevicesHandler{},
				"",
				true),
		}
	}

	return c.routes
}
