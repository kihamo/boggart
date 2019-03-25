package homie

import (
	"context"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	a "github.com/kihamo/boggart/components/boggart/atomic"
)

const (
	configNameSeparator = "."
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config     *Config
	lastUpdate *a.TimeNull

	deviceAttributes *sync.Map

	otaRun      *a.Bool
	otaWritten  *a.Uint32
	otaTotal    *a.Uint32
	otaChecksum *a.String
	otaFlash    chan struct{}

	settings *sync.Map
}

func (b *Bind) UpdateStatus(status boggart.BindStatus) {
	b.BindBase.UpdateStatus(status)

	if status == boggart.BindStatusOnline && b.OTAIsRunning() {
		b.otaFlash <- struct{}{}
	}
}

func (b *Bind) Broadcast(ctx context.Context, level string, payload interface{}) error {
	return b.MQTTPublishRaw(ctx, MQTTPublishTopicBroadcast.Format(b.config.BaseTopic, level), 1, false, payload)
}

func (b *Bind) Restart(ctx context.Context) error {
	return b.MQTTPublishRaw(ctx, MQTTPublishTopicRestart.Format(b.config.BaseTopic, b.SerialNumber()), 1, false, true)
}

func (b *Bind) Reset(ctx context.Context) error {
	return b.MQTTPublish(ctx, MQTTPublishTopicReset.Format(b.config.BaseTopic, b.SerialNumber()), true)
}

func (b *Bind) bump() {
	b.lastUpdate.Set(time.Now())
}

func (b *Bind) LastUpdate() *time.Time {
	if b.lastUpdate.IsNil() {
		return nil
	}

	t := b.lastUpdate.Load()
	return &t
}
