package samsung_tizen

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT
	config *Config
	mutex  sync.RWMutex
	client *tv.ApiV2
	mac    string
}
