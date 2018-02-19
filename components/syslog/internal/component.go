package internal

import (
	"net"
	"os"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	rsyslog "gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	handlers    []syslog.HasHandler
}

func (c *Component) Name() string {
	return syslog.ComponentName
}

func (c *Component) Version() string {
	return syslog.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name:     boggart.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.handlers = make([]syslog.HasHandler, 0, 0)

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	components, err := c.application.GetComponents()
	if err != nil {
		return err
	}

	c.logger = logger.NewOrNop(c.Name(), c.application)

	for _, component := range components {
		if handler, ok := component.(syslog.HasHandler); ok {
			c.handlers = append(c.handlers, handler)
		}
	}

	go func(l logger.Logger) {
		defer wg.Done()

		server := rsyslog.NewServer()
		server.SetFormat(rsyslog.Automatic)
		server.SetHandler(c)

		addr := net.JoinHostPort(c.config.String(syslog.ConfigHost), c.config.String(syslog.ConfigPort))
		if err := server.ListenUDP(addr); err != nil {
			l.Fatalf("Failed to listen [%d]: %s\n", os.Getpid(), err.Error())
			return
		}

		if err := server.Boot(); err != nil {
			l.Fatalf("Failed to boot [%d]: %s\n", os.Getpid(), err.Error())
			return
		}

		server.Wait()
	}(c.logger)

	return nil
}

func (c *Component) Handle(message format.LogParts, length int64, err error) {
	if err != nil {
		c.logger.Error("Handler with error", map[string]interface{}{
			"error":   err,
			"length":  length,
			"message": message,
		})
	}

	for _, h := range c.handlers {
		go h.SyslogHandler(message)
	}
}
