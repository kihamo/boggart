package internal

import (
	"errors"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/mqtt/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardMenu() dashboard.Menu {
	routes := c.DashboardRoutes()

	return dashboard.NewMenu("MQTT").
		WithIcon("list").
		WithChild(dashboard.NewMenu("State").WithRoute(routes[1])).
		WithChild(dashboard.NewMenu("Subscriptions").WithRoute(routes[2])).
		WithChild(dashboard.NewMenu("Cache").WithRoute(routes[3]))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/state/", &handlers.StateHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/subscriptions/", &handlers.SubscriptionsHandler{
				Component: c,
			}).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/cache/", &handlers.CacheHandler{
				PayloadCache: c.payloadCache,
			}).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
		}
	}

	return c.routes
}

func (c *Component) LivenessCheck() map[string]dashboard.HealthCheck {
	return map[string]dashboard.HealthCheck{
		c.Name() + "_client_is_connect": func() error {
			if c.Client().IsConnectionOpen() {
				return nil
			}

			return errors.New("connection is closed")
		},
	}
}
