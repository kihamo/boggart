package syslog

import (
	"context"
	"net"
	"sort"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	SystemMikrotik installer.System = "Mikrotik"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
		SystemMikrotik,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, s installer.System) ([]installer.Step, error) {
	cfg := b.config()

	if s == SystemMikrotik {
		const (
			ActionName  = "boggart"
			TagWifiName = "wifi"
			TagVPNName  = "vpn"
		)

		var (
			host = b.Config().App().String(dashboard.ConfigHost)
			port = strconv.FormatInt(cfg.Port, 10)
		)

		if host == "0.0.0.0" {
			addresses, err := net.InterfaceAddrs()
			if err == nil {
				for _, address := range addresses {
					if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
						if ipnet.IP.To4() != nil {
							host = ipnet.IP.String()
						}
					}
				}
			}
		}

		return []installer.Step{{
			Description: "Install logging from console",
			Content: `/system logging remove [/system logging find action=` + ActionName + `]
/system logging action remove [/system logging action find name=` + ActionName + `]
/system logging action add name="` + ActionName + `" target=remote remote=` + host + ` remote-port=` + port + ` src-address=0.0.0.0 bsd-syslog=yes syslog-time-format=bsd-syslog syslog-facility=daemon syslog-severity=info
/system logging add topics=wireless,info prefix="` + TagWifiName + `" action=` + ActionName + `
/system logging add topics=l2tp,info prefix="` + TagVPNName + `" action=` + ActionName,
		}, {
			Description: "Install logging from WebBox",
			Content: ` 1. Open System > Logging > Actions
 2. Press Add New button
 3. Fill form:
    Name: ` + ActionName + `
    Type: remote
    Remote Address: ` + host + `
    Remote Port: ` + port + `
    Src. Address: empty
    BSD Syslog: checked
    Syslog Facility: 3 (daemon)
    Syslog Severity: 6 (info)
 4. Press OK button to save
 5. Change tab to Rules
 6. Press Add New button
 7. Fill form:
    Enabled: checked
    Topics: wireless and info
    Prefix: ` + TagWifiName + `
    Action: ` + ActionName + `
 8. Press OK button to save
 9. Press Add New button
10. Fill form:
    Enabled: checked
    Topics: l2tp and info
    Prefix: ` + TagVPNName + `
    Action: ` + ActionName + `
11. Press OK button to save`,
		}}, nil
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	publishes := b.MQTT().Publishes()
	channels := make([]*openhab.Channel, 0, len(publishes))

	sortNames := make([]string, 0, len(publishes))
	topics := make(map[string]mqtt.Topic, len(publishes))

	parts := cfg.TopicMessage.Split()

	var nameIndex int

	for i := len(parts) - 1; i > 0; i-- {
		nameIndex = i

		if parts[i] == `+` {
			break
		}
	}

	for topic := range publishes {
		parts := topic.Split()
		if len(parts) <= nameIndex {
			continue
		}
		name := parts[nameIndex]

		if cfg.TopicMessage.Format(name) == topic {
			sortNames = append(sortNames, name)
			topics[name] = topic
		}
	}

	sort.Strings(sortNames)

	for _, name := range sortNames {
		id := openhab.IDNormalizeCamelCase(name)

		channels = append(channels, openhab.NewChannel(id, openhab.ChannelTypeString).
			WithStateTopic(topics[name]).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeString).
					WithLabel("Syslog %s").
					WithIcon("text"),
			))
	}

	return openhab.StepsByBind(b, nil, channels...)
}
