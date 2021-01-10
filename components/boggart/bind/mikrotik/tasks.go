package mikrotik

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/mikrotik"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
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

	b.SetSerialNumber(system.SerialNumber)

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
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
			WithName("clients-sync").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskClientsSyncHandler),
						b.config.ReadinessTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.ClientsSyncInterval)),
	)

	return err
}

func (b *Bind) taskClientsSyncHandler(ctx context.Context) error {
	b.updateWiFiClient(ctx)
	b.clientWiFi.Ready()

	b.updateVPNClient(ctx)
	b.clientVPN.Ready()

	return nil
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
