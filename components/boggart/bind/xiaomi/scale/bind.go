package scale

import (
	"context"
	"time"

	"github.com/go-ble/ble/linux/hci"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/xiaomi/scale"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.WorkersBind
	di.ProbesBind

	disconnected *atomic.BoolNull

	config               *Config
	provider             *scale.Client
	currentProfile       atomic.Value
	measureStartDatetime *atomic.Time
}

func (b *Bind) Run() error {
	b.disconnected.Nil()
	b.notifyCurrentProfile(context.Background())
	b.Meta().SetMAC(b.config.MAC.HardwareAddr)
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
	profiles := make([]*Profile, 0, len(b.config.Profiles)+1)
	guestExist := false

	for _, profile := range b.config.Profiles {
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
