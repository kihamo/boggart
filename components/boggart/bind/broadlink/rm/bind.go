package rm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/broadlink"
)

type Bind struct {
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config *ConfigRM

	provider interface{}
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

func (b *Bind) Run() error {
	return b.MQTT().PublishAsync(context.Background(), b.config.TopicCaptureState, false)
}
