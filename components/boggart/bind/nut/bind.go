package nut

import (
	"errors"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/robbiet480/go.nut"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.ProbesBind

	config    *Config
	mutex     sync.Mutex
	variables map[string]interface{}
}

func (b *Bind) connect() (client nut.Client, err error) {
	client, err = nut.Connect(b.config.Address.Host)
	if err != nil {
		return client, err
	}

	defer func() {
		if err != nil {
			_, _ = client.Disconnect()
		}
	}()

	if user := b.config.Address.User; user != nil {
		username := user.Username()
		password, _ := user.Password()

		if username != "" && password != "" {
			result, err := client.Authenticate(username, password)
			if err != nil {
				return client, err
			}

			if !result {
				return client, errors.New("authenticate failed")
			}
		}
	}

	return client, err
}

func (b *Bind) ups() (ups nut.UPS, err error) {
	client, err := b.connect()
	if err != nil {
		return ups, err
	}
	defer func() {
		_, _ = client.Disconnect()
	}()

	devices, err := client.GetUPSList()
	if err != nil {
		return ups, err
	}

	for _, device := range devices {
		if device.Name == b.config.UPS {
			for _, v := range device.Variables {
				if v.Name == "device.serial" {
					b.Meta().SetSerialNumber(strings.TrimSpace(v.Value.(string)))
					return device, nil
				}
			}

			break
		}
	}

	return ups, errors.New("device " + b.config.UPS + " not found")
}

func (b *Bind) Variables() ([]nut.Variable, error) {
	ups, err := b.ups()
	if err != nil {
		return nil, err
	}

	return ups.Variables, nil
}

func (b *Bind) SetVariable(variable, value string) (bool, error) {
	client, err := b.connect()
	if err != nil {
		return false, err
	}
	defer func() {
		_, _ = client.Disconnect()
	}()

	devices, err := client.GetUPSList()
	if err != nil {
		return false, err
	}

	for _, device := range devices {
		if device.Name == b.config.UPS {
			variable = strings.ToLower(variable)

			for _, v := range device.Variables {
				if strings.ToLower(v.Name) == variable {
					return device.SetVariable(variable, value)
				}
			}

			return false, errors.New("variable " + variable + " not found")
		}
	}

	return false, errors.New("device " + b.config.UPS + " not found")
}

func (b *Bind) SendCommand(command string) (bool, error) {
	client, err := b.connect()
	if err != nil {
		return false, err
	}
	defer func() {
		_, _ = client.Disconnect()
	}()

	devices, err := client.GetUPSList()
	if err != nil {
		return false, err
	}

	for _, device := range devices {
		if device.Name == b.config.UPS {
			command = strings.ToLower(command)

			for _, cmd := range device.Commands {
				if strings.ToLower(cmd.Name) == command {
					return device.SendCommand(cmd.Name)
				}
			}

			return false, errors.New("command " + command + " not found")
		}
	}

	return false, errors.New("device " + b.config.UPS + " not found")
}
