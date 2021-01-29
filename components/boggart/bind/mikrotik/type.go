package mikrotik

import (
	"net/url"
	"strings"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/mikrotik"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, err := url.Parse(config.Address)
	if err != nil {
		return nil, err
	}

	macMappingCase := make(map[string]string, len(config.MacAddressMapping))
	for mac, alias := range config.MacAddressMapping {
		macMappingCase[strings.ToLower(mac)] = alias
	}
	config.MacAddressMapping = macMappingCase

	username := u.User.Username()
	password, _ := u.User.Password()

	bind := &Bind{
		config:                  config,
		address:                 u,
		provider:                mikrotik.NewClient(u.Host, username, password, config.ClientTimeout),
		connectionsZombieKiller: &atomic.Once{},
	}

	return bind, nil
}
