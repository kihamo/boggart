package mercury

import (
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
)

type Bind struct {
	tariff1         uint64
	tariff2         uint64
	tariff3         uint64
	tariff4         uint64
	voltage         uint64
	amperage        uint64
	power           uint64
	batteryVoltage  uint64
	lastPowerOff    int64
	lastPowerOn     int64
	makeDate        int64
	firmwareDate    int64
	firmwareVersion string

	boggart.BindBase
	boggart.BindMQTT

	mutex    sync.Mutex
	provider *mercury.ElectricityMeter200

	updaterInterval time.Duration
}
