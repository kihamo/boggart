package mikrotik

import (
	"context"
	"errors"
	"strings"
	"sync"

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

	cfg := b.config()
	steps := make([]installer.Step, 0, 1)

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
			WithStateTopic(cfg.TopicPackagesInstalledVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesInstalledVersion, openhab.ItemTypeString).
					WithLabel("Packages installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idPackagesLatestVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicPackagesLatestVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPackagesLatestVersion, openhab.ItemTypeString).
					WithLabel("Packages latest version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareInstalledVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFirmwareInstalledVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareInstalledVersion, openhab.ItemTypeString).
					WithLabel("Firmware installed version").
					WithIcon("text"),
			),
		openhab.NewChannel(idFirmwareLatestVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFirmwareLatestVersion.Format(sn)).
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

			if cfg.IgnoreUnknownMacAddress {
				connections = cfg.MacAddressMapping
			} else {
				connections = make(map[string]string)

				for _, i := range wirelessConnections {
					if i.Interface != iface.Name {
						continue
					}

					mac := i.MacAddress.String()
					if alias, ok := cfg.MacAddressMapping[mac]; ok {
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
						WithStateTopic(cfg.TopicInterfaceLastConnect.Format(sn, iface.Type, iface.Name)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idLastConnected, openhab.ItemTypeString).
								WithLabel("Last connected to "+iface.Name).
								WithIcon("network"),
						),
					openhab.NewChannel(id+idLastDisconnected, openhab.ChannelTypeString).
						WithStateTopic(cfg.TopicInterfaceLastDisconnect.Format(sn, iface.Type, iface.Name)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idLastDisconnected, openhab.ItemTypeString).
								WithLabel("Last disconnected to "+iface.Name).
								WithIcon("network"),
						),
				)
			}

			for mac, alias := range connections {
				clientID := "Client" + openhab.IDNormalizeCamelCase(alias)

				channels = append(channels,
					openhab.NewChannel(id+clientID, openhab.ChannelTypeContact).
						WithStateTopic(cfg.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, mac)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id+clientID, openhab.ItemTypeContact).
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
							WithStateTopic(cfg.TopicInterfaceLastConnect.Format(sn, iface.Type)).
							AddItems(
								openhab.NewItem(itemPrefix+id+idLastConnected, openhab.ItemTypeString).
									WithLabel("Last connected to L2TP").
									WithIcon("network"),
							),
						openhab.NewChannel(id+idLastDisconnected, openhab.ChannelTypeString).
							WithStateTopic(cfg.TopicInterfaceLastDisconnect.Format(sn, iface.Type)).
							AddItems(
								openhab.NewItem(itemPrefix+id+idLastDisconnected, openhab.ItemTypeString).
									WithLabel("Last disconnected to L2TP").
									WithIcon("network"),
							),
					)
				})

				clientID := "Client" + openhab.IDNormalizeCamelCase(i.User)

				channels = append(channels,
					openhab.NewChannel(id+clientID, openhab.ChannelTypeContact).
						WithStateTopic(cfg.TopicInterfaceConnect.Format(sn, iface.Type, iface.Name, i.User)).
						WithOn("true").
						WithOff("false").
						AddItems(
							openhab.NewItem(itemPrefix+id+clientID, openhab.ItemTypeContact).
								WithLabel(i.User+" connected to L2TP").
								WithIcon("boy_3"),
						),
				)

				break
			}
		}
	}

	// UPS
	workers := b.Workers()
	var hasUPS bool

	const (
		idUPSBatteryVoltage = "BatteryVoltage"
		idUPSBatteryCharge  = "BatteryCharge"
		idUPSBatteryRuntime = "BatteryRuntime"
		idUPSInputVoltage   = "InputVoltage"
		idUPSLoad           = "Load"
		idUPSStatus         = "Status"
	)

	for _, info := range workers.TasksID() {
		taskID := workers.TaskShortName(info[1])

		if !strings.HasPrefix(taskID, TaskNameUPS) {
			continue
		}

		hasUPS = true
		upsName := taskID[len(TaskNameUPS)+1:]
		upsID := "UPS_" + openhab.IDNormalizeCamelCase(upsName) + "_"

		channels = append(channels,
			openhab.NewChannel(upsID+idUPSBatteryVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicUPSBatteryVoltage.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSBatteryVoltage, openhab.ItemTypeNumber).
						WithLabel("Battery voltage [%d V]").
						WithIcon("energy"),
				),
			openhab.NewChannel(upsID+idUPSBatteryCharge, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicUPSBatteryCharge.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSBatteryCharge, openhab.ItemTypeNumber).
						WithLabel("Battery charge [%d %%]").
						WithIcon("batterylevel"),
				),
			openhab.NewChannel(upsID+idUPSBatteryRuntime, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicUPSBatteryRuntime.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSBatteryRuntime, openhab.ItemTypeNumber).
						WithLabel("Battery runtime [JS(human_seconds.js):%s]").
						WithIcon("time"),
				),
			openhab.NewChannel(upsID+idUPSInputVoltage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicUPSInputVoltage.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSInputVoltage, openhab.ItemTypeNumber).
						WithLabel("Input voltage [%d V]").
						WithIcon("energy"),
				),
			openhab.NewChannel(upsID+idUPSLoad, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicUPSLoad.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSLoad, openhab.ItemTypeNumber).
						WithLabel("Load [%d %%]").
						WithIcon("batterylevel"),
				),
			openhab.NewChannel(upsID+idUPSStatus, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicUPSStatus.Format(sn, upsName)).
				AddItems(
					openhab.NewItem(itemPrefix+upsID+idUPSStatus, openhab.ItemTypeString).
						WithLabel("Status").
						WithIcon("text"),
				),
		)
	}

	if hasUPS {
		steps = append(steps, openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds))
	}

	return openhab.StepsByBind(b, steps, channels...)
}
