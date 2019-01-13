package broadlink

import (
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/broadlink"
)

type BindRM struct {
	boggart.BindBase
	boggart.BindMQTT

	provider        interface{}
	mac             net.HardwareAddr
	ip              net.UDPAddr
	captureDuration time.Duration
}
type BindRMSupportCapture interface {
	StartCaptureRemoteControlCode() error
	ReadCapturedRemoteControlCodeAsString() (broadlink.RemoteType, string, error)
}

type BindRMSupportIR interface {
	SendIRRemoteControlCodeAsString(code string, count int) error
}

type BindRMSupportRF315Mhz interface {
	SendRF315MhzRemoteControlCodeAsString(code string, count int) error
}

type BindRMSupportRF433Mhz interface {
	SendRF433MhzRemoteControlCodeAsString(code string, count int) error
}
