package internal

import (
	"net"
	"os"

	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logging"
	rsyslog "gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Component struct {
	logger   logging.Logger
	handlers []syslog.HasHandler

	server *rsyslog.Server
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
			Name: logging.ComponentName,
		},
	}
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	components, err := a.GetComponents()
	if err != nil {
		return err
	}

	c.server = rsyslog.NewServer()
	c.server.SetFormat(rsyslog.Automatic)
	c.server.SetHandler(c)

	c.logger = logging.DefaultLogger().Named(c.Name())

	c.handlers = make([]syslog.HasHandler, 0, 0)
	for _, component := range components {
		if handler, ok := component.(syslog.HasHandler); ok {
			c.handlers = append(c.handlers, handler)
		}
	}

	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	addr := net.JoinHostPort(cfg.String(syslog.ConfigHost), cfg.String(syslog.ConfigPort))
	if err := c.server.ListenUDP(addr); err != nil {
		c.logger.Fatalf("Failed to listen [%d]: %s\n", os.Getpid(), err.Error())
		return err
	}

	if err := c.server.Boot(); err != nil {
		c.logger.Fatalf("Failed to boot [%d]: %s\n", os.Getpid(), err.Error())
		return err
	}

	ready <- struct{}{}

	c.server.Wait()
	return nil
}

func (c *Component) Shutdown() error {
	if c.server != nil {
		return c.server.Kill()
	}

	return nil
}

func (c *Component) Handle(message format.LogParts, length int64, err error) {
	if err != nil {
		c.logger.Error("Handler with error",
			"error", err,
			"length", length,
			"message", message,
		)
	}

	for _, h := range c.handlers {
		go h.SyslogHandler(message)
	}
}
