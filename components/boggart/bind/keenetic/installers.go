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
		idHotspotConnectLastState      = "HotspotConnect_Last"
		idHotspotConnectLastActive     = "HotspotConnect_Last_Active"
		idHotspotConnectLastUplink     = "HotspotConnect_Last_Uplink"
		idHotspotConnectLastRegistered = "HotspotConnect_Last_Registered"
		idHotspotConnectRegistered     = "HotspotConnect_"
		idHostMAC                      = "MAC"
		idHostIP                       = "IP"
		idHostName                     = "Name"
		idHostActive                   = "Active"
		idHostUplink                   = "Uplink"
		idHostRegistered               = "Registered"
		idHostSpeed                    = "Speed"
	)

	subItemPrefix := itemPrefix + "Last_"

	channels := []*openhab.Channel{
		openhab.NewChannel(idHotspotConnectLastState, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicHotspotState.Format(sn)).
			AddItems(
				openhab.NewItem(subItemPrefix+idHostMAC, openhab.ItemTypeString).
					WithLabel("MAC address").
					WithProfile("transform:JSONPATH", "function", "$.mac").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostIP, openhab.ItemTypeString).
					WithLabel("IP address").
					WithProfile("transform:JSONPATH", "function", "$.ip").
					WithIcon("text"),
				openhab.NewItem(subItemPrefix+idHostName, openhab.ItemTypeString).
					WithLabel("Name").
					WithProfile("transform:JSONPATH", "function", "$.name").
					WithIcon("text"),
			),
		openhab.NewChannel(idHotspotConnectLastActive, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHotspotState.Format(sn)).
			WithTransformationPattern("JSONPATH:$.active").
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(subItemPrefix+idHostActive, openhab.ItemTypeContact).
					WithLabel("Active").
					WithIcon("text"),
			),
		openhab.NewChannel(idHotspotConnectLastUplink, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHotspotState.Format(sn)).
			WithTransformationPattern("JSONPATH:$.uplink").
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(subItemPrefix+idHostUplink, openhab.ItemTypeContact).
					WithLabel("Uplink").
					WithIcon("text"),
			),
		openhab.NewChannel(idHotspotConnectLastRegistered, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHotspotState.Format(sn)).
			WithTransformationPattern("JSONPATH:$.registered").
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(subItemPrefix+idHostRegistered, openhab.ItemTypeContact).
					WithLabel("Registered").
					WithIcon("text"),
			),
	}

	var channelId, macNormalize string

	b.hotspotConnections.Range(func(key, value interface{}) bool {
		si := value.(*storeItem)

		if !si.host.Registered {
			return true
		}

		macNormalize = openhab.IDNormalize(strings.ReplaceAll(si.host.Mac, ":", ""))
		subItemPrefix = itemPrefix + macNormalize + "_"
		channelId = idHotspotConnectRegistered + macNormalize

		channels = append(channels,
			openhab.NewChannel(channelId+"_State", openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				AddItems(
					openhab.NewItem(subItemPrefix+idHostMAC, openhab.ItemTypeString).
						WithLabel("MAC address").
						WithProfile("transform:JSONPATH", "function", "$.mac").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostIP, openhab.ItemTypeString).
						WithLabel("IP address").
						WithProfile("transform:JSONPATH", "function", "$.ip").
						WithIcon("text"),
					openhab.NewItem(subItemPrefix+idHostName, openhab.ItemTypeString).
						WithLabel("Name").
						WithProfile("transform:JSONPATH", "function", "$.name").
						WithIcon("text"),
				),
			openhab.NewChannel(channelId+"_"+idHostActive, openhab.ChannelTypeContact).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				WithTransformationPattern("JSONPATH:$.active").
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(subItemPrefix+idHostActive, openhab.ItemTypeContact).
						WithLabel("Active").
						WithIcon("siren"),
				),
			openhab.NewChannel(channelId+"_"+idHostUplink, openhab.ChannelTypeContact).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				WithTransformationPattern("JSONPATH:$.uplink").
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(subItemPrefix+idHostUplink, openhab.ItemTypeContact).
						WithLabel("Uplink").
						WithIcon("network"),
				),
			openhab.NewChannel(channelId+"_"+idHostRegistered, openhab.ChannelTypeContact).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				WithTransformationPattern("JSONPATH:$.registered").
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(subItemPrefix+idHostRegistered, openhab.ItemTypeContact).
						WithLabel("Registered").
						WithIcon("keyring"),
				),
			openhab.NewChannel(channelId+"_"+idHostSpeed, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicHotspotState.Format(sn, si.ID())).
				WithTransformationPattern("JSONPATH:$.speed").
				AddItems(
					openhab.NewItem(subItemPrefix+idHostSpeed, openhab.ItemTypeNumber).
						WithLabel("Host speed [%d Mbit/sec]").
						WithIcon("returnpipe"),
				),
		)

		return true
	})

	return openhab.StepsByBind(b, nil, channels...)
}
