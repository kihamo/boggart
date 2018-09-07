package internal

import (
	"net/url"
	"sync"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/kihamo/boggart/components/openhab"
	apiclient "github.com/kihamo/boggart/components/openhab/client"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

type Component struct {
	application shadow.Application
	config      config.Component

	mutex  sync.RWMutex
	apiUrl *url.URL
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
			Name: dashboard.ComponentName,
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

func (c *Component) Run() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.apiUrl, err = url.Parse(c.config.String(openhab.ConfigAPIURL))
	if err != nil {
		return err
	}

	return nil
}

func (c *Component) Client() *apiclient.OpenHABREST {
	c.mutex.RLock()
	transport := httptransport.New(c.apiUrl.Host, "/rest", []string{c.apiUrl.Scheme})
	c.mutex.RUnlock()

	transport.Debug = c.config.Bool(config.ConfigDebug)

	return apiclient.New(transport, nil)
}
