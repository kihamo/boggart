package subscribes

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/logger"
)

const (
	GPIOTopicPrefix = "boggart/gpio/"
)

type GPIOSubscribe struct {
	devicesManager boggart.DevicesManager
	logger         logger.Logger
}

func NewGPIOSubscribe(m boggart.DevicesManager, l logger.Logger) *GPIOSubscribe {
	return &GPIOSubscribe{
		devicesManager: m,
		logger:         l,
	}
}

func (s *GPIOSubscribe) Filters() map[string]byte {
	return map[string]byte{
		GPIOTopicPrefix + "#": 0,
	}
}

func (s *GPIOSubscribe) Callback(client mqtt.Component, message m.Message) {
	t := message.Topic()
	if !strings.HasPrefix(t, GPIOTopicPrefix) {
		return
	}

	pin := t[len(GPIOTopicPrefix):]
	number, err := strconv.ParseInt(pin, 10, 64)
	if err != nil {
		return
	}

	device := s.devicesManager.Device(fmt.Sprintf("pin.%d", number))
	if device == nil {
		return
	}

	deviceGPIO, ok := device.(*devices.GPIOPin)
	if !ok || !device.IsEnabled() {
		return
	}

	if bytes.Equal(message.Payload(), []byte(`1`)) {
		err = deviceGPIO.High()
	} else {
		err = deviceGPIO.Low()
	}

	if err != nil {
		s.logger.Errorf("Out %s failed with error %s", message.Payload(), err.Error())
	}
}
