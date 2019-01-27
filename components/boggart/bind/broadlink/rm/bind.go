package rm

import (
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
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
