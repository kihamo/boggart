package mikrotik

import (
	"context"
	"errors"
	"strings"

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

	var id string

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

			for mac, alias := range connections {
				id = idInterface +
					openhab.IDNormalizeCamelCase(iface.Type) + "_" +
					openhab.IDNormalizeCamelCase(iface.Name) + "_" +
					openhab.IDNormalizeCamelCase(alias)

				channels = append(channels,
					openhab.NewChannel(id, openhab.ChannelTypeContact).
						WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, mac)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
								WithLabel("Interface "+iface.Name+" connect "+mac).
								WithIcon("wireless"),
						),
				)
			}

		case InterfaceL2TPServer:
			for _, i := range interfacesL2TP {
				if i.Name != iface.Name {
					continue
				}

				id = idInterface +
					openhab.IDNormalizeCamelCase(iface.Type) + "_" +
					openhab.IDNormalizeCamelCase(i.User)

				channels = append(channels,
					openhab.NewChannel(id, openhab.ChannelTypeContact).
						WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, i.User)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
								WithLabel("Interface L2TP "+i.User).
								WithIcon("boy_3"),
						),
				)

				break
			}
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
