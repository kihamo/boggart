package internal

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
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
		WithChild(dashboard.NewMenu("Manager").WithRoute(routes[9])).
		WithChild(dashboard.NewMenu("Workers").WithRoute(routes[13])).
		WithChild(dashboard.NewMenu("Config YAML").WithURL("/" + c.Name() + "/config/view"))
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		<-c.application.ReadyComponent(c.Name())

		bindHandler := handlers.NewBindHandler(c)
		configHandler := handlers.NewConfigHandler(c)
		logsHandler := handlers.NewLogsHandler(c)

		c.routes = []dashboard.Route{
			dashboard.RouteFromAssetFS(c),
			dashboard.NewRoute("/"+c.Name()+"/bind/", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/bind/:id/", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/bind/:id/:action/*path", bindHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/config/:action/", configHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/config/:action/:id/", configHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/installer/:id/:system/", handlers.NewInstallerHandler(c)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/logs/:id/", logsHandler).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/logs/:id/:action/", logsHandler).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/manager/", handlers.NewManagerHandler(c, c.tasksManager)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/metrics/:id/", handlers.NewMetricsHandler(c)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/mqtt/:id/", handlers.NewMQTTHandler(c, c.mqtt)).
				WithMethods([]string{http.MethodGet}).
				WithAuth(true),
			dashboard.NewRoute("/"+c.Name()+"/widget/:id/", handlers.NewWidgetHandler(c)).
				WithMethods([]string{http.MethodGet, http.MethodPost}),
			dashboard.NewRoute("/"+c.Name()+"/workers/", handlers.NewWorkersHandler(c.tasksManager)).
				WithMethods([]string{http.MethodGet, http.MethodPost}).
				WithAuth(true),
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

				key := strings.TrimSpace(r.URL.Query().Get(boggart.AccessKeyName))
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

func (c *Component) DashboardTemplateFunctions() map[string]interface{} {
	funcMap := template.FuncMap{
		"human_bytes": templateFunctionHumanBytes,
		"widget_url":  templateFunctionWidgetURL,
	}

	for prefix, t := range boggart.GetBindTypes() {
		buf := strings.Builder{}
		for i, r := range prefix {
			switch {
			case r == '_':
				buf.WriteRune(r)
			case i == 0 && !unicode.IsLetter(r):
				buf.WriteString("_")
			case !unicode.IsLetter(r) && !unicode.IsDigit(r):
				buf.WriteString("_")
			default:
				buf.WriteRune(r)
			}
		}
		prefix = buf.String()

		if cmp, ok := t.Type().(dashboard.HasTemplateFunctions); ok {
			for alias, f := range cmp.DashboardTemplateFunctions() {
				alias = prefix + "_" + alias

				if _, ok := funcMap[alias]; !ok {
					funcMap[alias] = f
				}
			}
		}
	}

	return funcMap
}

func templateFunctionHumanBytes(size interface{}) string {
	s := fmt.Sprintf("%v", size)
	b, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return s
	}

	const unit = 1024

	if b < unit {
		return s + " B"
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func templateFunctionWidgetURL(opts ...interface{}) template.URL {
	u := "/"
	ctx := context.Background()

	if len(opts) > 0 {
		if templateCtx, ok := opts[0].(map[string]interface{}); ok {
			if requestCtx, ok := templateCtx["Request"]; ok {
				if request, ok := requestCtx.(*dashboard.Request); ok {
					ctx = request.Context()
					opts = opts[1:]
				}
			}
		}
	}

	if widget := di.WidgetFromContext(ctx); widget != nil {
		values := make(map[string]string, len(opts)/2)

		for i := 0; i < len(opts); i += 2 {
			cur := fmt.Sprintf("%v", opts[i])

			if i+1 < len(opts) {
				values[cur] = fmt.Sprintf("%v", opts[i+1])
			} else {
				values[cur] = ""
			}
		}

		if generateURL, err := widget.URL(values); err == nil {
			u = generateURL.String()
		}
	}

	return template.URL(u)
}
