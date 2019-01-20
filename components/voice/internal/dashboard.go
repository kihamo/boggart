package internal

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/voice/internal/handlers"
	"github.com/kihamo/shadow/components/dashboard"
)

func (c *Component) DashboardTemplates() *assetfs.AssetFS {
	return dashboard.TemplatesFromAssetFS(c)
}

func (c *Component) DashboardRoutes() []dashboard.Route {
	if c.routes == nil {
		<-c.application.ReadyComponent(storage.ComponentName)

		c.routes = []dashboard.Route{
			dashboard.NewRoute("/"+c.Name()+"/file/", &handlers.FileHandler{
				Voice:   c,
				Storage: c.application.GetComponent(storage.ComponentName).(storage.Component),
			}).
				WithMethods([]string{http.MethodGet}),
		}
	}

	return c.routes
}
