package rm

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/broadlink"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind

	provider broadlink.Device
}

type SupportCapture interface {
	StartCaptureRemoteControlCode() error
	ReadCapturedRemoteControlCodeAsString() (broadlink.RemoteType, string, error)
}

type SupportIR interface {
	SendIRRemoteControlCodeAsString(code string, count int) error
}

type SupportRF315Mhz interface {
	SendRF315MhzRemoteControlCodeAsString(code string, count int) error
}

type SupportRF433Mhz interface {
	SendRF433MhzRemoteControlCodeAsString(code string, count int) error
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	switch cfg.Model {
	case "rm3mini":
		b.provider = broadlink.NewRMMini(cfg.MAC.HardwareAddr, cfg.Host)

	case "rm2proplus":
		b.provider = broadlink.NewRM2ProPlus3(cfg.MAC.HardwareAddr, cfg.Host)

	default:
		return errors.New("unknown model " + cfg.Model)
	}

	b.provider.SetTimeout(cfg.ConnectionTimeout)

	b.Meta().SetMAC(cfg.MAC.HardwareAddr)

	return b.MQTT().PublishAsync(context.Background(), cfg.TopicCaptureState.Format(cfg.MAC.String()), false)
}
