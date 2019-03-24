package internal

import (
	"net/http"
	"strings"

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
			dashboard.NewRoute("/"+c.Name()+"/bind/:id/:action/*path", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/config/:action", handlers.NewConfigHandler(c.manager, c)).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/widget/:id", handlers.NewWidgetHandler(c.manager)).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
		}
	}

	return c.routes
}

func (c *Component) DashboardMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				if dashboard.TemplateNamespaceFromContext(ctx) != c.Name() {
					next.ServeHTTP(w, r)
					return
				}

				// уже авторизован общей авторизацией
				if request := dashboard.RequestFromContext(ctx); request != nil && request.User().IsAuthorized() {
					next.ServeHTTP(w, r)
					return
				}

				// для защищенных маршрутов сработает мидлваря с общей авторизацией
				if route := dashboard.RouteFromContext(ctx); route != nil && (route.Auth() || route.HandlerName() == "AssetsHandler") {
					next.ServeHTTP(w, r)
					return
				}

				keysConfig := c.config.String(boggart.ConfigAccessKeys)
				if keysConfig == "" {
					next.ServeHTTP(w, r)
					return
				}

				keys := strings.Split(keysConfig, ",")
				if len(keys) == 0 {
					next.ServeHTTP(w, r)
					return
				}

				key := strings.TrimSpace(r.URL.Query().Get("key"))
				if key == "" {
					http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
					return
				}

				for _, k := range keys {
					k = strings.TrimSpace(k)

					if k == key {
						next.ServeHTTP(w, r)
						return
					}
				}

				http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			})
		},
	}
}
