package mikrotik

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	ctx := context.Background()

	interfaces, err := b.provider.Interfaces(ctx)
	if err != nil {
		return nil, err
	}

	interfacesL2TP, err := b.provider.InterfaceL2TPServer(ctx)
	if err != nil {
		return nil, err
	}

	wirelessConnections, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		return nil, err
	}

	const (
		idPackagesInstalledVersion = "PackagesInstalledVersion"
		idPackagesLatestVersion    = "PackagesLatestVersion"
		idFirmwareInstalledVersion = "FirmwareInstalledVersion"
		idFirmwareLatestVersion    = "FirmwareLatestVersion"
		idInterface                = "Interface"
		idLastConnected            = "LastConnected"
		idLastDisconnected         = "LastDisconnected"
	)

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	channels := []*openhab.Channel{
		openhab.NewChannel(idPackagesInstalledVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicPackagesInstalledVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesInstalledVersion, openhab.ItemTypeString).
					WithLabel("Packages installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idPackagesLatestVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesLatestVersion, openhab.ItemTypeString).
					WithLabel("Packages latest version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareInstalledVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareInstalledVersion, openhab.ItemTypeString).
					WithLabel("Firmware installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareLatestVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareLatestVersion, openhab.ItemTypeString).
					WithLabel("Firmware latest version").
					WithIcon("text"),
			),
	}

	var (
		id       string
		lastL2TP sync.Once
	)

	for _, iface := range interfaces {
		if iface.Disabled {
			continue
		}

		switch iface.Type {
		case InterfaceWireless:
			var connections map[string]string

			if b.config.IgnoreUnknownMacAddress {
				connections = b.config.MacAddressMapping
			} else {
				connections = make(map[string]string)

				for _, i := range wirelessConnections {
					if i.Interface != iface.Name {
						continue
					}

					mac := i.MacAddress.String()
					if alias, ok := b.config.MacAddressMapping[mac]; ok {
						connections[mac] = alias
					} else {
						connections[mac] = strings.ReplaceAll(mac, ":", "")
					}
				}
			}

			id = idInterface +
				openhab.IDNormalizeCamelCase(iface.Type) + "_" +
				openhab.IDNormalizeCamelCase(iface.Name) + "_"

			if len(connections) > 0 {
				channels = append(channels,
					openhab.NewChannel(id+idLastConnected, openhab.ChannelTypeString).
						WithStateTopic(b.config.TopicInterfaceLastConnect.Format(sn, iface.Type, iface.Name)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idLastConnected, openhab.ItemTypeString).
								WithLabel("Last connected to "+iface.Name).
								WithIcon("network"),
						),
					openhab.NewChannel(id+idLastDisconnected, openhab.ChannelTypeString).
						WithStateTopic(b.config.TopicInterfaceLastDisconnect.Format(sn, iface.Type, iface.Name)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idLastDisconnected, openhab.ItemTypeString).
								WithLabel("Last disconnected to "+iface.Name).
								WithIcon("network"),
						),
				)
			}

			for mac, alias := range connections {
				clientId := "Client" + openhab.IDNormalizeCamelCase(alias)

				channels = append(channels,
					openhab.NewChannel(id+clientId, openhab.ChannelTypeContact).
						WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, mac)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id+clientId, openhab.ItemTypeContact).
								WithLabel(alias+" connected to "+iface.Name).
								WithIcon("network"),
						),
				)
			}

		case InterfaceL2TPServer:
			id = idInterface +
				openhab.IDNormalizeCamelCase(iface.Type) + "_"

			for _, i := range interfacesL2TP {
				if i.Name != iface.Name {
					continue
				}

				lastL2TP.Do(func() {
					channels = append(channels,
						openhab.NewChannel(id+idLastConnected, openhab.ChannelTypeString).
							WithStateTopic(b.config.TopicInterfaceLastConnect.Format(sn, iface.Type)).
							AddItems(
								openhab.NewItem(itemPrefix+id+idLastConnected, openhab.ItemTypeString).
									WithLabel("Last connected to L2TP").
									WithIcon("network"),
							),
						openhab.NewChannel(id+idLastDisconnected, openhab.ChannelTypeString).
							WithStateTopic(b.config.TopicInterfaceLastDisconnect.Format(sn, iface.Type)).
							AddItems(
								openhab.NewItem(itemPrefix+id+idLastDisconnected, openhab.ItemTypeString).
									WithLabel("Last disconnected to L2TP").
									WithIcon("network"),
							),
					)
				})

				clientId := "Client" + openhab.IDNormalizeCamelCase(i.User)

				channels = append(channels,
					openhab.NewChannel(id+clientId, openhab.ChannelTypeContact).
						WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, i.User)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id+clientId, openhab.ItemTypeContact).
								WithLabel(i.User+" connected to L2TP").
								WithIcon("boy_3"),
						),
				)

				break
			}
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
