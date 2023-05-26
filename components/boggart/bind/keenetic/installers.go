package keenetic

import (
	"context"
	"errors"
	"strings"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	cfg := b.config()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const (
		idHotspotConnectLast       = "HotspotConnect_Last"
		idHotspotConnectRegistered = "HotspotConnect_"
		idHostMAC                  = "MAC"
		idHostIP                   = "IP"
		idHostName                 = "Name"
		idHostActive               = "Active"
		idHostUplink               = "Uplink"
		idHostRegistered           = "Registered"
	)

	channelId := idHotspotConnectLast
	subItemPrefix := itemPrefix + "Last_"

	channels := []*openhab.Channel{
		openhab.NewChannel(channelId, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicHotspotState.Format(sn)).
			AddItems(
				openhab.NewItem(subItemPrefix+idHostMAC, openhab.ItemTypeString).
					WithLabel("MAC address").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostIP, openhab.ItemTypeString).
					WithLabel("IP address").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostName, openhab.ItemTypeString).
					WithLabel("Name").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostActive, openhab.ItemTypeContact).
					WithLabel("Active").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostUplink, openhab.ItemTypeContact).
					WithLabel("Uplink").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostRegistered, openhab.ItemTypeContact).
					WithLabel("Registered").
					WithIcon("text"),
			),
	}

	var macNormalize string

	b.hotspotConnections.Range(func(key, value interface{}) bool {
		si := value.(*storeItem)

		if !si.host.Registered {
			return true
		}

		macNormalize = openhab.IDNormalize(strings.ReplaceAll(si.host.Mac, ":", ""))
		subItemPrefix = itemPrefix + macNormalize + "_"
		channelId = idHotspotConnectRegistered + macNormalize

		channels = append(channels,
			openhab.NewChannel(channelId, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				AddItems(
					openhab.NewItem(subItemPrefix+idHostMAC, openhab.ItemTypeString).
						WithLabel("MAC address").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostIP, openhab.ItemTypeString).
						WithLabel("IP address").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostName, openhab.ItemTypeString).
						WithLabel("Name").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostActive, openhab.ItemTypeContact).
						WithLabel("Active").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostUplink, openhab.ItemTypeContact).
						WithLabel("Uplink").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostRegistered, openhab.ItemTypeContact).
						WithLabel("Registered").
						WithIcon("text"),
				),
		)

		return true
	})

	return openhab.StepsByBind(b, nil, channels...)
}
