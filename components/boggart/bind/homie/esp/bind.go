package esp

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

const (
	configNameSeparator = "."
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config     *Config
	lastUpdate atomic.TimeNull

	deviceAttributes sync.Map
	nodes            sync.Map
	settings         sync.Map

	otaEnabled  atomic.Bool
	otaRun      atomic.Bool
	otaWritten  atomic.Uint32
	otaTotal    atomic.Uint32
	otaChecksum atomic.String
	otaFlash    chan struct{}

	status atomic.BoolNull
}

func (b *Bind) updateStatus(status bool) {
	b.status.Set(status)

	if status && b.OTAIsRunning() {
		b.otaFlash <- struct{}{}
	}
}

func (b *Bind) Broadcast(ctx context.Context, level string, payload interface{}) error {
	return b.MQTTPublishRaw(ctx, b.config.TopicBroadcast.Format(level), 1, false, payload)
}

func (b *Bind) Restart(ctx context.Context) error {
	return b.MQTTPublishRaw(ctx, b.config.TopicRestart, 1, false, true)
}

func (b *Bind) Reset(ctx context.Context) error {
	return b.MQTTPublish(ctx, b.config.TopicReset, true)
}

func (b *Bind) ProtocolVersion() string {
	v, ok := b.DeviceAttribute("homie")
	if ok {
		return v.(string)
	}

	return ""
}

func (b *Bind) ProtocolVersionConstraint(constraint string) bool {
	current, err := version.NewVersion(b.ProtocolVersion())
	if err != nil {
		return false
	}

	constraints, err := version.NewConstraint(constraint)
	if err != nil {
		return false
	}

	return constraints.Check(current)
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
