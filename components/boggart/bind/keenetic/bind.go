package keenetic

import (
	"encoding/json"
	"net/url"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/keenetic"
	"github.com/kihamo/boggart/providers/keenetic/models"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	client *keenetic.Client

	hotspotConnections  sync.Map
	hotspotZombieKiller *atomic.Once
}

type storeItem struct {
	version uint64
	host    *models.ShowIPHotspotResponseHostItems0
}

func (i storeItem) ID() string {
	return mqtt.NameReplace(i.host.Mac)
}

// формирует payload для MQTT топика
func (i storeItem) MarshalBinary() (data []byte, err error) {
	type payload struct {
		MAC        string `json:"mac"`
		IP         string `json:"ip"`
		Name       string `json:"name"`
		Active     bool   `json:"active"`
		Uplink     bool   `json:"uplink"`
		Registered bool   `json:"registered"`
	}

	return json.Marshal(&payload{
		MAC:        i.host.Mac,
		IP:         i.host.IP,
		Name:       i.host.Name,
		Active:     i.host.Active,
		Uplink:     i.host.Link == "up",
		Registered: i.host.Registered,
	})
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	link := &url.URL{}
	*link = cfg.Address.URL
	link.User = nil

	b.Meta().SetLinkSchemeConvert(link)

	b.client = keenetic.New(cfg.Address.Host, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	b.hotspotZombieKiller.Reset()
	b.hotspotConnections = sync.Map{}

	return nil
}
