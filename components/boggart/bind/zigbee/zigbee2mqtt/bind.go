package zigbee2mqtt

import (
	"context"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	status atomic.BoolNull

	settings     *Settings
	settingsLock sync.RWMutex

	devices     map[string]*Device
	devicesLock sync.RWMutex
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
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

	// для случая когда из /config приходит не полная инфа, сохраняем то, что пришло из /info
	if settings.Config == nil && b.settings != nil && b.settings.Config != nil {
		settings.Config = b.settings.Config
	}

	b.settings = &settings
	b.settingsLock.Unlock()

	return nil
}

func (b *Bind) setState(ctx context.Context, flag bool) error {
	if flag {
		if b.status.True() {
			return b.MQTT().PublishAsyncRawWithoutCache(ctx, b.config().TopicDevicesRequest, 1, false, true)
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

func (b *Bind) NetworkMap(ctx context.Context) (*NetworkMap, error) {
	cfg := b.config()

	message, err := b.MQTT().Request(ctx, cfg.TopicNetworkMapRequest, cfg.TopicNetworkMapResponse, "raw")
	if err != nil {
		return nil, err
	}

	m := &NetworkMap{}

	if err := message.JSONUnmarshal(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (b *Bind) SetPermitJoin(ctx context.Context, flag bool) error {
	return b.MQTT().PublishAsyncRawWithoutCache(ctx, b.config().TopicPermitJoin, 1, false, flag)
}

func (b *Bind) SetLogLevel(ctx context.Context, level string) error {
	return b.MQTT().PublishAsyncRawWithoutCache(ctx, b.config().TopicLogLevel, 1, false, level)
}
