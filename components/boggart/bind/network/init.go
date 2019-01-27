package network

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/network/ping"
	"github.com/kihamo/boggart/components/boggart/bind/network/service"
)

func init() {
	boggart.RegisterBindType("network:ping", ping.Type{})
	boggart.RegisterBindType("network:service", service.Type{})
}
