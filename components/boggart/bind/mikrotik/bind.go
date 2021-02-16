package mikrotik

import (
	"net/url"
	"regexp"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
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

var (
	interfaceWirelessMQTT   = mqtt.NameReplace(InterfaceWireless)
	interfaceL2TPServerMQTT = mqtt.NameReplace(InterfaceL2TPServer)
)

type Bind struct {
	di.ConfigBind
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
	connectionsFirstLoad    map[string]*atomic.Once
}

type storeItem struct {
	version        uint64
	isUpdated      bool
	connectionName string
	interfaceType  string
	interfaceName  string // is unique
}

func (i storeItem) String() string {
	return i.interfaceType + "/" + i.connectionName
}

func (b *Bind) Run() error {
	for _, o := range b.connectionsFirstLoad {
		o.Reset()
	}

	return nil
}

func (b *Bind) loadOrStoreItem(item *storeItem) {
	actual, loaded := b.connectionsActive.LoadOrStore(item.String(), item)
	if loaded {
		actual.(*storeItem).version = item.version
		actual.(*storeItem).isUpdated = true
	}
}
