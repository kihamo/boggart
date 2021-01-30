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

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
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
						b.config.ReadinessTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.ClientsSyncInterval)),
	)
	if err != nil {
		return err
	}

	if b.config.TopicSyslog != "" {
		return b.MQTT().Subscribe(mqtt.NewSubscriber(b.config.TopicSyslog, 0, b.callbackMQTTSyslog))
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
		for _, connection := range connections {
			item := &storeItem{
				version:        version,
				interfaceType:  interfaceWirelessMQTT,
				interfaceName:  mqtt.NameReplace(connection.Interface),
				connectionName: mqtt.NameReplace(connection.MacAddress.String()),
			}

			actual, loaded := b.connectionsActive.LoadOrStore(item.String(), item)
			if loaded {
				actual.(*storeItem).version = version
			}
		}

		storedWireless = true
	} else {
		// TODO: wrap
		err = fmt.Errorf("get wireless registration table failed: %w", e)
	}

	// 1.2 L2TP

	// 1.2.1 Сначала вычитываем все интерфейсы сконфигурированные, что бы сработал механизм деактивации
	//       не активных для этого простовляем им версию 0, чтобы механиз ниже сработал
	if connections, e := b.provider.InterfaceL2TPServer(ctx); e == nil {
		for _, connection := range connections {
			item := &storeItem{
				version:        0,
				interfaceType:  interfaceWirelessMQTT,
				interfaceName:  mqtt.NameReplace(connection.Name),
				connectionName: mqtt.NameReplace(connection.User),
			}

			b.connectionsActive.LoadOrStore(item.String(), item)
		}
	}

	// 1.2.2 Получаем активные соединения и регистрируем их
	if connections, e := b.provider.PPPActive(ctx); e == nil {
		for _, connection := range connections {
			item := &storeItem{
				version:        version,
				interfaceType:  interfaceL2TPServerMQTT,
				connectionName: mqtt.NameReplace(connection.Name),
			}

			if actual, loaded := b.connectionsActive.Load(item.String()); loaded {
				actual.(*storeItem).version = version
			}
		}

		storedL2TP = true
	} else {
		// TODO: wrap
		err = fmt.Errorf("get active ppp connections failed: %w", e)
	}

	// 2. Теперь рассылаем актуальные версии в MQTT, а то, что меньше текущей версии еще и удаляем
	sn := b.Meta().SerialNumber()

	b.connectionsActive.Range(func(key, value interface{}) bool {
		item := value.(*storeItem)
		var payload bool

		if item.version == version {
			payload = true
		} else if item.interfaceType == interfaceWirelessMQTT && storedWireless || item.interfaceType == interfaceL2TPServerMQTT && storedL2TP {
			payload = false
			b.connectionsActive.Delete(key)
		} else {
			return true
		}

		_ = b.MQTT().PublishAsync(ctx, b.config.TopicInterfaceConnect.Format(sn, item.interfaceType, item.interfaceName, item.connectionName), payload)
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
					b.config.TopicInterfaceConnect.Format(b.Meta().SerialNumber()),
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

	// check upgrade
	if checkPackages, e := b.provider.SystemPackageUpdateCheck(ctx); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicPackagesInstalledVersion.Format(sn), checkPackages.InstalledVersion); e != nil {
			err = multierr.Append(err, e)
		}

		if checkPackages.LatestVersion != "" {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicPackagesLatestVersion.Format(sn), checkPackages.LatestVersion); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if checkRouterBoard, e := b.provider.SystemRouterBoard(ctx); e == nil {
		if e := b.MQTT().PublishAsync(ctx, b.config.TopicFirmwareInstalledVersion.Format(sn), checkRouterBoard.CurrentFirmware); e != nil {
			err = multierr.Append(err, e)
		}

		if checkRouterBoard.UpgradeFirmware != "" {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicFirmwareLatestVersion.Format(sn), checkRouterBoard.UpgradeFirmware); e != nil {
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
