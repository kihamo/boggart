package octoprint

import (
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/octoprint"
	"github.com/kihamo/boggart/providers/octoprint/models"
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

	provider *octoprint.Client

	settings      *models.Settings
	settingsMutex sync.RWMutex

	devices      map[string]bool
	devicesMutex sync.RWMutex
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) SystemSettings() *models.Settings {
	b.settingsMutex.RLock()
	defer b.settingsMutex.RUnlock()

	return b.settings
}

func (b *Bind) PluginMQTTSettings() *models.SettingsPluginsMqtt {
	b.settingsMutex.RLock()
	defer b.settingsMutex.RUnlock()

	local := b.settings.Plugins.Mqtt
	if local == nil || local.Broker == nil || local.Publish == nil || local.Publish.Events == nil || local.Broker.URL == "" {
		return nil
	}

	return local
}

func (b *Bind) TemperatureFromMQTT() bool {
	cfg := b.PluginMQTTSettings()
	return cfg != nil && cfg.Publish.TemperatureActive
}

func (b *Bind) TemperatureTopic() mqtt.Topic {
	cfg := b.PluginMQTTSettings()
	if cfg == nil {
		return ""
	}

	return mqtt.Topic(cfg.Publish.BaseTopic +
		strings.Replace(cfg.Publish.TemperatureTopic, "{temp}", "+", 1))
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.provider = octoprint.New(cfg.Address.Host, cfg.APIKey, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	b.settingsMutex.Lock()
	b.settings = nil
	b.settingsMutex.Unlock()

	b.devicesMutex.Lock()
	b.devices = make(map[string]bool, 0)
	b.devicesMutex.Unlock()

	return nil
}
