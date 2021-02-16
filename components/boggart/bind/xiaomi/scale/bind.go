package scale

import (
	"context"
	"time"

	"github.com/go-ble/ble/linux/hci"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/bluetooth"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
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

	disconnected *atomic.BoolNull

	provider             *scale.Client
	currentProfile       atomic.Value
	measureStartDatetime *atomic.Time
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.disconnected.Nil()
	b.measureStartDatetime.Set(time.Now())

	cfg := b.config()

	b.notifyCurrentProfile(context.Background())
	b.Meta().SetMAC(cfg.MAC.HardwareAddr)

	if len(cfg.Profiles) > 0 {
		for name, profile := range cfg.Profiles {
			profile.Name = name
			profile.Age = profile.GetAge()
		}
	}

	device, err := bluetooth.NewDevice()
	if err != nil {
		return err
	}

	b.provider = scale.NewClient(device, cfg.MAC.HardwareAddr, cfg.CaptureDuration, cfg.IgnoreEmptyImpedance)

	return nil
}

func (b *Bind) CurrentProfile() *Profile {
	if profile := b.currentProfile.Load(); profile != nil {
		if p, ok := profile.(*Profile); ok {
			return p
		}
	}

	return profileGuest
}

func (b *Bind) SetProfile(name string) {
	profile := b.Profile(name)

	if profile == nil {
		profile = profileGuest
	}

	b.currentProfile.Store(profile)
	b.measureStartDatetime.Set(time.Now())
}

func (b *Bind) Profile(name string) *Profile {
	if name == "" {
		return nil
	}

	for _, p := range b.Profiles() {
		if p.Name == name {
			return p
		}
	}

	return nil
}

func (b *Bind) Profiles() []*Profile {
	cfg := b.config()

	profiles := make([]*Profile, 0, len(cfg.Profiles)+1)
	guestExist := false

	for _, profile := range cfg.Profiles {
		profiles = append(profiles, profile)

		if profile.Name == profileGuest.Name {
			guestExist = true
		}
	}

	if !guestExist {
		profiles = append(profiles, profileGuest)
	}

	return profiles
}

func (b *Bind) Measures(ctx context.Context) ([]*scale.Measure, error) {
	measures, err := b.provider.Measures(ctx)
	if err != nil && err == hci.ErrDisallowed {
		b.disconnected.True()
	}

	return measures, err
}

func (b *Bind) Close() error {
	return b.provider.Close()
}
