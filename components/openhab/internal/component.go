package internal

import (
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/kihamo/boggart/components/openhab"
	apiclient "github.com/kihamo/boggart/components/openhab/client"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

type Component struct {
	application shadow.Application
	config      config.Component
}

func (c *Component) Name() string {
	return openhab.ComponentName
}

func (c *Component) Version() string {
	return openhab.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name:     dashboard.ComponentName,
			Required: true,
		},
		{
			Name: messengers.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Client() *apiclient.OpenHABREST {
	u, err := url.Parse(c.config.String(openhab.ConfigAPIURL))
	if err != nil {
		return nil
	}

	httpClient := &http.Client{
		Transport: &nethttp.Transport{},
	}

	transport := httptransport.NewWithClient(u.Host, "/rest", []string{u.Scheme}, httpClient)
	transport.Debug = c.config.Bool(config.ConfigDebug)

	return apiclient.New(transport, nil)
}
