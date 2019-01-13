package nut

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/robbiet480/go.nut"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	host            string
	username        string
	password        string
	ups             string
	updaterInterval time.Duration

	mutex     sync.Mutex
	variables map[string]interface{}
}

func (b *Bind) connect() (client nut.Client, err error) {
	client, err = nut.Connect(b.host)
	if err != nil {
		return client, err
	}

	defer func() {
		if err != nil {
			client.Disconnect()
		}
	}()

	if b.username != "" && b.password != "" {
		result, err := client.Authenticate(b.username, b.password)
		if err != nil {
			return client, err
		}

		if !result {
			return client, errors.New("authenticate failed")
		}
	}

	return client, err
}

func (b *Bind) GetUPS() (ups nut.UPS, err error) {
	client, err := b.connect()
	if err != nil {
		return ups, err
	}
	defer client.Disconnect()

	devices, err := client.GetUPSList()
	if err != nil {
		return ups, err
	}

	for _, device := range devices {
		if device.Name == b.ups {
			for _, v := range device.Variables {
				if v.Name == "device.serial" {
					b.SetSerialNumber(strings.TrimSpace(v.Value.(string)))
					return device, nil
				}
			}

			break
		}
	}

	return ups, errors.New("device " + b.ups + " not found")
}

func (b *Bind) SendCommand(command string) (bool, error) {
	client, err := b.connect()
	if err != nil {
		return false, err
	}
	defer client.Disconnect()

	devices, err := client.GetUPSList()
	if err != nil {
		return false, err
	}

	for _, device := range devices {
		if device.Name == b.ups {
			command = strings.ToLower(command)

			for _, cmd := range device.Commands {
				if strings.ToLower(cmd.Name) == command {
					return device.SendCommand(cmd.Name)
				}
			}

			return false, errors.New("command " + command + " not found")
		}
	}

	return false, errors.New("device " + b.ups + " not found")
}
