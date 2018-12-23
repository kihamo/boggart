package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("Smart home").
		WithIcon("home").
		WithRoute(routes[1])
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		devicesHandler := handlers.NewDevicesHandler(c.devicesManager, c.listenersManager)
		cameraHandler := &handlers.CameraHandler{
			DevicesManager: c.devicesManager,
		}

		m := c.application.GetComponent(messengers.ComponentName)
		if m != nil {
			<-c.application.ReadyComponent(messengers.ComponentName)
			cameraHandler.Messengers = m.(messengers.Component)
		}

		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/", devicesHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/devices/", devicesHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/devices/:device/:action", handlers.NewDeviceHandler(c.devicesManager)).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/camera/:sn/:channel", cameraHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
		}
	}

	return c.routes
}
