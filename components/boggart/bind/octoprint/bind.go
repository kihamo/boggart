package octoprint

import (
	"context"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/octoprint"
	"github.com/kihamo/boggart/providers/octoprint/client/settings"
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

func (b *Bind) SystemSettingsUpdate(ctx context.Context) error {
	response, err := b.provider.Settings.GetSettings(settings.NewGetSettingsParamsWithContext(ctx), nil)
	if err != nil {
		return err
	}

	b.settingsMutex.Lock()
	b.settings = response.GetPayload()
	b.settingsMutex.Unlock()

	return nil
}

func (b *Bind) SystemSettings() *models.Settings {
	b.settingsMutex.RLock()
	defer b.settingsMutex.RUnlock()

	return b.settings
}

func (b *Bind) PluginMQTTSettings() *models.SettingsPluginsMqtt {
	settings := b.SystemSettings()
	if settings == nil {
		return nil
	}

	local := settings.Plugins.Mqtt
	if local == nil || local.Broker == nil || local.Publish == nil || local.Publish.Events == nil || local.Broker.URL == "" {
		return nil
	}

	return local
}

func (b *Bind) DisplayLayerProgressEnabled() bool {
	settings := b.SystemSettings()
	if settings == nil {
		return false
	}

	return settings != nil
}

func (b *Bind) TemperatureFromMQTT() bool {
	cfg := b.PluginMQTTSettings()
	return cfg != nil && cfg.Publish.TemperatureActive
}

func (b *Bind) JobFromMQTT() bool {
	cfg := b.PluginMQTTSettings()
	return cfg != nil && cfg.Publish.ProgressActive && cfg.Publish.PrinterData
}

func (b *Bind) TemperatureTopic() mqtt.Topic {
	cfg := b.PluginMQTTSettings()
	if cfg == nil {
		return ""
	}

	return mqtt.Topic(cfg.Publish.BaseTopic +
		strings.Replace(cfg.Publish.TemperatureTopic, "{temp}", "+", 1))
}

func (b *Bind) JobTopic() mqtt.Topic {
	cfg := b.PluginMQTTSettings()
	if cfg == nil {
		return ""
	}

	return mqtt.Topic(cfg.Publish.BaseTopic +
		strings.Replace(cfg.Publish.ProgressTopic, "{progress}", "printing", 1))
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

	b.Meta().SetLink(&cfg.Address.URL)

	return nil
}
