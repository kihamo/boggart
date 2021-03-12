package mikrotik

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/mikrotik"
	"go.uber.org/multierr"
)

const (
	TaskNameSerialNumber        = "serial-number"
	TaskNameUpdater             = "updater"
	TaskNameInterfaceConnection = "connections"
	TaskNameUPS                 = "ups"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName(TaskNameSerialNumber).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumberHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
	system, err := b.provider.SystemRouterBoard(ctx)
	if err != nil {
		return err
	}

	if system.SerialNumber == "" {
		return errors.New("serial number is empty")
	}

	b.Meta().SetSerialNumber(system.SerialNumber)
	cfg := b.config()

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)
	if err != nil {
		return err
	}

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameInterfaceConnection).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFunc(b.taskInterfaceConnectionHandler),
						cfg.ReadinessTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.ClientsSyncInterval)),
	)
	if err != nil {
		return err
	}

	// ups updater
	if devices, err := b.provider.SystemUPS(ctx); err == nil {
		for _, device := range devices {
			if device.Disabled || device.Invalid {
				continue
			}

			_, err = b.Workers().RegisterTask(
				tasks.NewTask().
					WithName(TaskNameUPS + "-" + device.Name).
					WithHandler(
						b.Workers().WrapTaskHandlerIsOnline(
							tasks.HandlerFuncFromShortToLong(b.taskUPSHandler(device)),
						),
					).
					WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UPSInterval)),
			)
			if err != nil {
				return err
			}
		}
	}

	if cfg.TopicSyslog != "" {
		return b.MQTT().Subscribe(mqtt.NewSubscriber(cfg.TopicSyslog, 0, b.callbackMQTTSyslog))
	}

	return nil
}

func (b *Bind) taskInterfaceConnectionHandler(ctx context.Context, meta tasks.Meta, _ tasks.Task) (err error) {
	var (
		storedWireless bool
		storedL2TP     bool
	)

	// 0. Формируем индекс новой версии
	version := meta.Attempts()

	// 1. Формируем новую версию данных, получая активные соединения по каждому из типов

	// 1.1 Wireless
	if connections, e := b.provider.InterfaceWirelessRegistrationTable(ctx); e == nil {
		var isUpdated bool

		// чтобы в первую загрузку пометить все элементы как старые
		b.connectionsFirstLoad[InterfaceWireless].Do(func() {
			isUpdated = true
		})

		for _, connection := range connections {
			b.loadOrStoreItem(&storeItem{
				version:        version,
				isUpdated:      isUpdated,
				interfaceType:  interfaceWirelessMQTT,
				interfaceName:  mqtt.NameReplace(connection.Interface),
				connectionName: mqtt.NameReplace(connection.MacAddress.String()),
			})
		}

		storedWireless = true
	} else {
		err = fmt.Errorf("get wireless registration table failed: %w", e)
	}

	// 1.2 L2TP
	if connections, e := b.provider.InterfaceL2TPServer(ctx); e == nil {
		var isUpdated bool

		// чтобы в первую загрузку пометить все элементы как старые
		b.connectionsFirstLoad[InterfaceL2TPServer].Do(func() {
			isUpdated = true
		})

		for _, connection := range connections {
			if !connection.Running {
				continue
			}

			b.loadOrStoreItem(&storeItem{
				version:        version,
				isUpdated:      isUpdated,
				interfaceType:  interfaceL2TPServerMQTT,
				interfaceName:  mqtt.NameReplace(connection.Name),
				connectionName: mqtt.NameReplace(connection.User),
			})
		}

		storedL2TP = true
	} else {
		err = fmt.Errorf("get ppp connections failed: %w", e)
	}

	// 2. Теперь рассылаем актуальные версии в MQTT, а то, что меньше текущей версии еще и удаляем
	sn := b.Meta().SerialNumber()
	cfg := b.config()

	b.connectionsActive.Range(func(key, value interface{}) bool {
		item := value.(*storeItem)
		var payload bool

		if item.version != version {
			// среди активных не обнаружено, принимаем решение об исключении
			if item.interfaceType == interfaceWirelessMQTT && storedWireless {
				payload = false
			} else if item.interfaceType == interfaceL2TPServerMQTT && storedL2TP {
				payload = false
			} else {
				return true
			}
		} else if !item.isUpdated {
			// значит пришел новый элемент и это новый коннект
			payload = true
		} else {
			return true
		}

		_ = b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicInterfaceConnect.Format(sn, item.interfaceType, item.interfaceName, item.connectionName), payload)

		if payload {
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicInterfaceLastConnect.Format(sn, item.interfaceType, item.interfaceName), item.connectionName)
		} else {
			b.connectionsActive.Delete(key)

			_ = b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicInterfaceLastDisconnect.Format(sn, item.interfaceType, item.interfaceName), item.connectionName)
		}

		return true
	})

	// 3. После первого успешного обновления справочников, необходимо запустить разово подписчика MQTT
	//    который зачистит записи зобми (те, которые из-за retained остались в MQTT, но из роутера мы уже
	//    их не получаем, так как например их удалили)
	if err == nil {
		b.connectionsZombieKiller.Do(func() {
			// TODO: придумать как в конце удалить эту подписку, так как операция нужна разовая,
			// а с таким подходом после каждой публикации будет прогонять еще раз

			err = b.MQTT().Subscribe(
				mqtt.NewSubscriber(
					cfg.TopicInterfaceConnect.Format(b.Meta().SerialNumber()),
					0,
					b.callbackMQTTInterfacesZombies,
				),
			)
		})

		if err != nil {
			b.connectionsZombieKiller.Reset()
		}
	}

	return err
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	// Ports on mikrotik
	if stats, e := b.provider.InterfaceStats(ctx); e == nil {
		for _, stat := range stats {
			metricTrafficReceivedBytes.With("serial_number", sn).With(
				"interface", stat.Name,
				"mac", stat.MacAddress.String()).Set(float64(stat.RXByte))
			metricTrafficSentBytes.With("serial_number", sn).With(
				"interface", stat.Name,
				"mac", stat.MacAddress.String()).Set(float64(stat.TXByte))
		}
	} else if !mikrotik.IsEmptyResponse(e) {
		err = multierr.Append(err, e)
	}

	if resource, e := b.provider.SystemResource(ctx); e == nil {
		metricCPULoad.With("serial_number", sn).Set(float64(resource.CPULoad))
		metricMemoryAvailable.With("serial_number", sn).Set(float64(resource.FreeMemory))
		metricMemoryUsage.With("serial_number", sn).Set(float64(resource.TotalMemory - resource.FreeMemory))
		metricStorageAvailable.With("serial_number", sn).Set(float64(resource.FreeHDDSpace))
		metricStorageUsage.With("serial_number", sn).Set(float64(resource.TotalHDDSpace - resource.FreeHDDSpace))
	} else if !mikrotik.IsEmptyResponse(e) {
		err = multierr.Append(err, e)
	}

	if disks, e := b.provider.SystemDisk(ctx); e == nil {
		for _, disk := range disks {
			metricDiskUsage.With("serial_number", sn).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Size - disk.Free))
			metricDiskAvailable.With("serial_number", sn).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Free))
		}
	} else if !mikrotik.IsEmptyResponse(e) {
		err = multierr.Append(err, e)
	}

	if health, e := b.provider.SystemHealth(ctx); e == nil {
		metricVoltage.With("serial_number", sn).Set(health.Voltage)
		metricTemperature.With("serial_number", sn).Set(float64(health.Temperature))
	} else if !mikrotik.IsEmptyResponse(e) {
		err = multierr.Append(err, e)
	}

	cfg := b.config()

	// check upgrade
	if checkPackages, e := b.provider.SystemPackageUpdateCheck(ctx); e == nil {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicPackagesInstalledVersion.Format(sn), checkPackages.InstalledVersion); e != nil {
			err = multierr.Append(err, e)
		}

		if checkPackages.LatestVersion != "" {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicPackagesLatestVersion.Format(sn), checkPackages.LatestVersion); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if checkRouterBoard, e := b.provider.SystemRouterBoard(ctx); e == nil {
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicFirmwareInstalledVersion.Format(sn), checkRouterBoard.CurrentFirmware); e != nil {
			err = multierr.Append(err, e)
		}

		if checkRouterBoard.UpgradeFirmware != "" {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicFirmwareLatestVersion.Format(sn), checkRouterBoard.UpgradeFirmware); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	// Wifi clients
	if clients, e := b.provider.InterfaceWirelessRegistrationTable(ctx); e == nil {
		metricWifiClients.With("serial_number", sn).Set(float64(len(clients)))

		var hasError bool

		arp, e := b.provider.IPARP(ctx)
		if e != nil && !mikrotik.IsEmptyResponse(e) {
			err = multierr.Append(err, e)
			hasError = true
		}

		dns, e := b.provider.IPDNSStatic(ctx)
		if e != nil && !mikrotik.IsEmptyResponse(e) {
			err = multierr.Append(err, e)
			hasError = true
		}

		leases, e := b.provider.IPDHCPServerLease(ctx)
		if e != nil && !mikrotik.IsEmptyResponse(e) {
			err = multierr.Append(err, e)
			hasError = true
		}

		if !hasError {
			for _, client := range clients {
				bytes := strings.Split(client.Bytes, ",")
				if len(bytes) != 2 {
					continue
				}

				name := mikrotik.GetNameByMac(client.MacAddress, arp, dns, leases)

				sent, e := strconv.ParseFloat(bytes[0], 64)
				if e != nil {
					err = multierr.Append(err, e)
					continue
				}

				received, e := strconv.ParseFloat(bytes[1], 64)
				if e != nil {
					err = multierr.Append(err, e)
					continue
				}

				metricTrafficReceivedBytes.With("serial_number", sn).With(
					"interface", client.Interface,
					"mac", client.MacAddress.String(),
					"name", name).Set(received)
				metricTrafficSentBytes.With("serial_number", sn).With(
					"interface", client.Interface,
					"mac", client.MacAddress.String(),
					"name", name).Set(sent)
			}
		}
	} else if !mikrotik.IsEmptyResponse(e) {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) taskUPSHandler(device mikrotik.SystemUPS) func(context.Context) error {
	return func(ctx context.Context) error {
		result, err := b.provider.SystemUPSMonitor(ctx, device.Name)
		if err != nil {
			return err
		}

		runtimeLeft, err := time.ParseDuration(result.RuntimeLeft)
		if err != nil {
			return err
		}

		cfg := b.config()
		sn := b.Meta().SerialNumber()

		metricUPSBatteryVoltage.With("serial_number", sn).With("ups", device.Name).Set(result.BatteryVoltage)
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSBatteryVoltage.Format(sn, device.Name), result.BatteryVoltage); e != nil {
			err = multierr.Append(err, e)
		}

		metricUPSBatteryCharge.With("serial_number", sn).With("ups", device.Name).Set(float64(result.BatteryCharge))
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSBatteryCharge.Format(sn, device.Name), result.BatteryCharge); e != nil {
			err = multierr.Append(err, e)
		}

		metricUPSInputVoltage.With("serial_number", sn).With("ups", device.Name).Set(float64(result.LineVoltage))
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSInputVoltage.Format(sn, device.Name), result.LineVoltage); e != nil {
			err = multierr.Append(err, e)
		}

		metricUPSLoad.With("serial_number", sn).With("ups", device.Name).Set(float64(result.Load))
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSLoad.Format(sn, device.Name), result.Load); e != nil {
			err = multierr.Append(err, e)
		}

		metricUPSBatteryRuntime.With("serial_number", sn).With("ups", device.Name).Set(runtimeLeft.Seconds())
		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSBatteryRuntime.Format(sn, device.Name), runtimeLeft.Seconds()); e != nil {
			err = multierr.Append(err, e)
		}

		var status string

		if result.Online {
			status = "OL"
		} else if result.OnBattery {
			status = "OB"
		}

		if result.LowBattery {
			status += " LB"
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicUPSStatus.Format(sn, device.Name), status); e != nil {
			err = multierr.Append(err, e)
		}

		return err
	}
}
