package zigbee2mqtt

import (
	"context"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	config *Config
	status atomic.BoolNull

	settings     *Settings
	settingsLock sync.RWMutex

	devices     map[string]*Device
	devicesLock sync.RWMutex
}

func (b *Bind) Run() error {
	b.status.Nil()

	return nil
}

func (b *Bind) Settings() *Settings {
	b.settingsLock.RLock()
	defer b.settingsLock.RUnlock()

	return b.settings
}

func (b *Bind) setSettingsFromMessage(message mqtt.Message) error {
	var settings Settings

	if err := message.JSONUnmarshal(&settings); err != nil {
		return err
	}

	b.settingsLock.Lock()
	b.settings = &settings
	b.settingsLock.Unlock()

	return nil
}

func (b *Bind) setState(ctx context.Context, flag bool) error {
	if flag {
		if b.status.True() {
			return b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicDeviceGet, true)
		}

		return nil
	}

	b.status.False()
	return nil
}

func (b *Bind) Devices() []*Device {
	b.devicesLock.RLock()
	defer b.devicesLock.RUnlock()

	devices := make([]*Device, 0, len(b.devices))

	for _, d := range b.devices {
		devices = append(devices, d)
	}

	return devices
}

func (b *Bind) SetPermitJoin(ctx context.Context, flag bool) error {
	return b.MQTT().PublishWithoutCache(ctx, b.config.TopicPermitJoin, flag)
}

func (b *Bind) SetLogLevel(ctx context.Context, level string) error {
	return b.MQTT().PublishWithoutCache(ctx, b.config.TopicLogLevel, level)
}
