package mikrotik

import (
	"regexp"
	"strings"
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

	provider *mikrotik.Client

	connectionsActive       sync.Map
	connectionsZombieKiller *atomic.Once
	connectionsFirstLoad    map[string]*atomic.Once
	macMappingCase          map[string]string
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

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	username := cfg.Address.User.Username()
	password, _ := cfg.Address.User.Password()

	b.provider = mikrotik.NewClient(cfg.Address.Host, username, password, cfg.ClientTimeout)

	b.macMappingCase = make(map[string]string, len(cfg.MacAddressMapping))
	for mac, alias := range cfg.MacAddressMapping {
		mac := strings.ReplaceAll(mac, "-", ":")
		mac = strings.ToLower(mac)

		b.macMappingCase[mac] = alias
	}

	for _, o := range b.connectionsFirstLoad {
		o.Reset()
	}

	b.connectionsZombieKiller.Reset()
	b.connectionsActive = sync.Map{}

	return nil
}

func (b *Bind) loadOrStoreItem(item *storeItem) {
	actual, loaded := b.connectionsActive.LoadOrStore(item.String(), item)
	if loaded {
		actual.(*storeItem).version = item.version
		actual.(*storeItem).isUpdated = true
	}
}
