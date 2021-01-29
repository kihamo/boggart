package mikrotik

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
	"github.com/kihamo/boggart/components/mqtt"
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

	const (
		idPackagesInstalledVersion = "PackagesInstalledVersion"
		idPackagesLatestVersion    = "PackagesLatestVersion"
		idFirmwareInstalledVersion = "FirmwareInstalledVersion"
		idFirmwareLatestVersion    = "FirmwareLatestVersion"
		idInterface                = "Interface"
	)

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	channels := []*openhab.Channel{
		openhab.NewChannel(idPackagesInstalledVersion, openhab.ItemTypeString).
			WithStateTopic(b.config.TopicPackagesInstalledVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesInstalledVersion, openhab.ItemTypeString).
					WithLabel("Packages installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idPackagesLatestVersion, openhab.ItemTypeString).
			WithStateTopic(b.config.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesLatestVersion, openhab.ItemTypeString).
					WithLabel("Packages latest version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareInstalledVersion, openhab.ItemTypeString).
			WithStateTopic(b.config.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareInstalledVersion, openhab.ItemTypeString).
					WithLabel("Firmware installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareLatestVersion, openhab.ItemTypeString).
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

		id = idInterface + openhab.IDNormalizeCamelCase(iface.Name)
		var channel *openhab.Channel

		switch iface.Type {
		case InterfaceWireless:
			channel = openhab.NewChannel(id, openhab.ItemTypeSwitch).
				WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, mqtt.NameReplace(iface.Type))).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
						WithLabel("Interface " + iface.Name + " connect").
						WithIcon("wireless"),
				)

		case InterfaceL2TPServer:
			for _, i := range interfacesL2TP {
				if i.Name != iface.Name {
					continue
				}

				channel = openhab.NewChannel(id, openhab.ItemTypeSwitch).
					WithStateTopic(b.config.TopicInterfaceConnect.Format(sn, InterfaceL2TPServer, mqtt.NameReplace(i.User))).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
							WithLabel("Interface L2TP " + i.User).
							WithIcon("boy_3"),
					)

				break
			}
		}

		if channel != nil {
			channel.
				WithOn("true").
				WithOff("false")

			channels = append(channels, channel)
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
