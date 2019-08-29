package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	mercury "github.com/kihamo/boggart/components/boggart/providers/mercury/v3"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	tariff1   *atomic.Uint32Null
	voltage1  *atomic.Float32Null
	voltage2  *atomic.Float32Null
	voltage3  *atomic.Float32Null
	amperage1 *atomic.Float32Null
	amperage2 *atomic.Float32Null
	amperage3 *atomic.Float32Null
	power1    *atomic.Float32Null
	power2    *atomic.Float32Null
	power3    *atomic.Float32Null

	provider *mercury.MercuryV3
	config   *Config
}
