package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/ota/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("OTA").
		WithIcon("cloud-upload-alt").
		WithChild(dashboard.NewMenu("Upgrade").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Releases").WithRoute(routes[3]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		releasesHandler := &handlers.ReleasesHandler{
			Updater:    c.updater,
			Repository: c.uploadRepository,
		}
		upgradeHandler := &handlers.UpgradeHandler{
			Updater:    c.updater,
			Repository: c.uploadRepository,
		}

		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/upgrade/", upgradeHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/upgrade/:id/:action", upgradeHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/", releasesHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/release/:id/:action", releasesHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
		}
	}

	return c.routes
}
