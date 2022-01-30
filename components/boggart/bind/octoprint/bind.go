package octoprint

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
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
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Settings() *models.Settings {
	b.settingsMutex.RLock()
	defer b.settingsMutex.RUnlock()

	if b.settings == nil {
		if response, err := b.provider.Settings.GetSettings(settings.NewGetSettingsParams(), nil); err == nil {
			b.settings = response.GetPayload()
		}
	}

	return b.settings
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

	return nil
}
