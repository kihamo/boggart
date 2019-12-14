package mikrotik

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/providers/mikrotik"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillSuccessTask(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	taskClientsSync := b.WrapTaskIsOnline(b.taskClientsSync)
	taskClientsSync.SetTimeout(b.config.LivenessTimeout)
	taskClientsSync.SetRepeats(-1)
	taskClientsSync.SetRepeatInterval(b.config.ClientsSyncInterval)
	taskClientsSync.SetName("clients-sync")

	taskStateUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskSerialNumber,
		taskClientsSync,
		taskStateUpdater,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, errors.New("bind isn't online")
	}

	system, err := b.provider.SystemRouterboard(ctx)
	if err != nil {
		return nil, err
	}

	if system.SerialNumber == "" {
		return nil, errors.New("serial number is empty")
	}

	b.SetSerialNumber(system.SerialNumber)
	return nil, nil
}

func (b *Bind) taskClientsSync(ctx context.Context) error {
	b.updateWiFiClient(ctx)
	b.clientWiFi.Ready()

	b.updateVPNClient(ctx)
	b.clientVPN.Ready()

	return nil
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	arp, err := b.provider.IPARP(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return err
	}

	dns, err := b.provider.IPDNSStatic(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return err
	}

	leases, err := b.provider.IPDHCPServerLease(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return err
	}

	// Wifi clients
	clients, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err == nil {
		metricWifiClients.With("serial_number", sn).Set(float64(len(clients)))

		for _, client := range clients {
			bytes := strings.Split(client.Bytes, ",")
			if len(bytes) != 2 {
				return err
			}

			name := mikrotik.GetNameByMac(client.MacAddress, arp, dns, leases)

			sent, err := strconv.ParseFloat(bytes[0], 64)
			if err != nil {
				return err
			}

			received, err := strconv.ParseFloat(bytes[1], 64)
			if err == nil {
				metricTrafficReceivedBytes.With("serial_number", sn).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(received)
				metricTrafficSentBytes.With("serial_number", sn).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(sent)
			} else if !mikrotik.IsEmptyResponse(err) {
				return err
			}
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return err
	}

	// Ports on mikrotik
	stats, err := b.provider.InterfaceStats(ctx)
	if err == nil {
		for _, stat := range stats {
			metricTrafficReceivedBytes.With("serial_number", sn).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.RXByte))
			metricTrafficSentBytes.With("serial_number", sn).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.TXByte))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return err
	}

	resource, err := b.provider.SystemResource(ctx)
	if err == nil {
		metricCPULoad.With("serial_number", sn).Set(float64(resource.CPULoad))
		metricMemoryAvailable.With("serial_number", sn).Set(float64(resource.FreeMemory))
		metricMemoryUsage.With("serial_number", sn).Set(float64(resource.TotalMemory - resource.FreeMemory))
		metricStorageAvailable.With("serial_number", sn).Set(float64(resource.FreeHDDSpace))
		metricStorageUsage.With("serial_number", sn).Set(float64(resource.TotalHDDSpace - resource.FreeHDDSpace))
	} else if !mikrotik.IsEmptyResponse(err) {
		return err
	}

	disks, err := b.provider.SystemDisk(ctx)
	if err == nil {
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
	} else if !mikrotik.IsEmptyResponse(err) {
		return err
	}

	health, err := b.provider.SystemHealth(ctx)
	if err == nil {
		metricVoltage.With("serial_number", sn).Set(health.Voltage)
		metricTemperature.With("serial_number", sn).Set(float64(health.Temperature))
	} else if !mikrotik.IsEmptyResponse(err) {
		return err
	}

	return nil
}
