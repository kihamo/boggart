package v1

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	mercury "github.com/kihamo/boggart/components/boggart/providers/mercury/v1"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	provider *mercury.MercuryV1

	tariff1          *atomic.Uint32Null
	tariff2          *atomic.Uint32Null
	tariff3          *atomic.Uint32Null
	tariff4          *atomic.Uint32Null
	voltage          *atomic.Uint32Null
	amperage         *atomic.Float32Null
	power            *atomic.Uint32Null
	batteryVoltage   *atomic.Float32Null
	lastPowerOffDate *atomic.Uint32Null
	lastPowerOnDate  *atomic.Uint32Null
	makeDate         *atomic.Uint32Null
	firmwareDate     *atomic.Uint32Null
	firmwareVersion  *atomic.String

	updaterInterval time.Duration
}
