package esp

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
)

const (
	// deviceStateInit         = "init"
	deviceStateReady = "ready"
	// deviceStateDisconnected = "disconnected"
	// deviceStateSleeping     = "sleeping"
	// deviceStateLost         = "lost"
	// deviceStateAlert        = "alert"
)

func (b *Bind) registerDeviceAttributes(name string, value interface{}) {
	b.deviceAttributes.Store(name, value)
}

func (b *Bind) DeviceAttribute(key string) (interface{}, bool) {
	return b.deviceAttributes.Load(key)
}

func (b *Bind) DeviceAttributes() map[string]interface{} {
	result := make(map[string]interface{})

	b.deviceAttributes.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})

	return result
}

func (b *Bind) deviceAttributesSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	attributeName := route[len(route)-1]

	if !strings.HasPrefix(attributeName, "$") {
		return errors.New("device attribute name " + attributeName + " is wrong")
	}

	attributeName = attributeName[1:]
	b.registerDeviceAttributes(attributeName, message.String())

	switch attributeName {
	case "online": // 2.x
		if b.ProtocolVersionConstraint(">= 2.0, < 3.0") {
			b.updateStatus(message.Bool())
		}

	case "state": // 3.x
		if b.ProtocolVersionConstraint(">= 3.0") {
			switch message.String() {
			case deviceStateReady:
				b.updateStatus(true)
			default:
				b.updateStatus(false)
			}
		}
	}

	return nil
}

func (b *Bind) deviceFirmwareSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	b.registerDeviceAttributes("fw."+route[len(route)-1], message.String())

	return nil
}

func (b *Bind) deviceImplementationSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()
	name := strings.Join(route[3:], ".")

	if strings.HasPrefix(name, "ota.") {
		return nil
	}

	b.registerDeviceAttributes("implementation."+strings.Join(route[3:], "."), message.String())

	return nil
}

func (b *Bind) deviceStatsSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	b.bump()

	route := message.Topic().Split()
	attributeName := strings.Join(route[3:], ".")

	var value interface{}

	switch attributeName {
	case "interval", "uptime":
		v, _ := strconv.ParseInt(message.String(), 10, 64)
		value = time.Second * time.Duration(v)

	default:
		value = message.String()
	}

	b.registerDeviceAttributes("stats."+attributeName, value)

	return nil
}
