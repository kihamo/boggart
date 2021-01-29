package mikrotik

import (
	"net/url"
	"regexp"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/mikrotik"
)

var (
	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
	vpnClientRegexp  = regexp.MustCompile(`^(\S+) logged (in|out), (.+?)$`)
)

const (
	InterfaceWireless   = "wlan"
	InterfaceL2TPServer = "l2tp-in"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	config   *Config
	address  *url.URL
	provider *mikrotik.Client

	connectionsActive       sync.Map
	connectionsZombieKiller *atomic.Once
}

type storeItem struct {
	version       uint64
	name          string
	interfaceType string
}

func (i storeItem) String() string {
	return i.interfaceType + "/" + i.name
}
