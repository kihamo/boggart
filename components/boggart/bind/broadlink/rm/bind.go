package rm

import (
	"context"
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	provider        interface{}
	mac             net.HardwareAddr
	ip              net.UDPAddr
	captureDuration time.Duration

	livenessInterval time.Duration
	livenessTimeout  time.Duration
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
	b.MQTTPublishAsync(context.Background(), MQTTPublishTopicCaptureState.Format(mqtt.NameReplace(b.SerialNumber())), false)

	return nil
}
