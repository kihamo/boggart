package nut

import (
	"errors"
	"strings"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/nut"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	provider        *nut.Client
	updaterInterval *atomic.Duration
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	var username, password string

	cfg := b.config()

	if cfg.Address.User != nil {
		username = cfg.Address.User.Username()
		password, _ = cfg.Address.User.Password()
	}

	b.provider = nut.New(cfg.Address.Host, username, password)
	b.updaterInterval.Set(cfg.UpdaterInterval)

	return nil
}

func (b *Bind) ups() (ups nut.UPS, err error) {
	devices, err := b.provider.UPS()
	if err != nil {
		return ups, err
	}

	cfg := b.config()

	for _, device := range devices {
		if device.Name == cfg.UPS {
			return device, nil
		}
	}

	return ups, errors.New("device " + cfg.UPS + " not found")
}

func (b *Bind) Variables() ([]nut.Variable, error) {
	ups, err := b.ups()
	if err != nil {
		return nil, err
	}

	return ups.Variables()
}

func (b *Bind) Commands() ([]nut.Command, error) {
	ups, err := b.ups()
	if err != nil {
		return nil, err
	}

	return ups.Commands()
}

func (b *Bind) SetVariable(name, value string) error {
	variables, err := b.Variables()
	if err != nil {
		return err
	}

	for _, variable := range variables {
		if strings.EqualFold(variable.Name, name) {
			return variable.Set(value)
		}
	}

	return errors.New("variable " + name + " not found")
}

func (b *Bind) SendCommand(name string) error {
	commands, err := b.Commands()
	if err != nil {
		return err
	}

	for _, command := range commands {
		if strings.EqualFold(command.Name, name) {
			return command.Call()
		}
	}

	return errors.New("command " + name + " not found")
}
