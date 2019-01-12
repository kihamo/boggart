package nut

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	config *Config

	mutex     sync.Mutex
	variables map[string]interface{}
}
