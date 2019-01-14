package mikrotik

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("bind-mikrotik-liveness-" + b.host)

	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.updaterInterval)
	taskStateUpdater.SetName("bind-mikrotik-updater-" + b.host)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	system, err := b.provider.SystemRouterboard(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	if system.SerialNumber == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, errors.New("serial number is empty")
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	if b.SerialNumber() != "" {
		return nil, nil
	}

	b.SetSerialNumber(system.SerialNumber)
	sn := system.SerialNumber

	// wifi clients
	clients, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		return nil, err
	}

	for _, connection := range clients {
		mac, err := b.Mac(ctx, connection.MacAddress)
		if err != nil {
			return nil, err
		}

		login := mqtt.NameReplace(mac.Address)

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiConnectedMAC.Format(sn), 0, false, login)
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiMACState.Format(sn, login), 0, false, true)
	}

	// vpn clients
	connections, err := b.provider.PPPActive(ctx)
	if err != nil {
		return nil, err
	}

	for _, connection := range connections {
		login := mqtt.NameReplace(connection.Name)

		// TODO:
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicVPNConnectedLogin.Format(sn), 0, false, login)
		_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicVPNLoginState.Format(sn, login), 0, false, true)
	}

	return nil, nil
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	serialNumber := b.SerialNumber()
	if serialNumber == "" {
		return nil, nil
	}

	arp, err := b.provider.IPARP(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	dns, err := b.provider.IPDNSStatic(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	leases, err := b.provider.IPDHCPServerLease(ctx)
	if err != nil && !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	// Wifi clients
	clients, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err == nil {
		metricWifiClients.With("serial_number", serialNumber).Set(float64(len(clients)))

		for _, client := range clients {
			bytes := strings.Split(client.Bytes, ",")
			if len(bytes) != 2 {
				return nil, err
			}

			name := mikrotik.GetNameByMac(client.MacAddress, arp, dns, leases)

			sent, err := strconv.ParseFloat(bytes[0], 64)
			if err != nil {
				return nil, err
			}

			received, err := strconv.ParseFloat(bytes[1], 64)
			if err == nil {
				metricTrafficReceivedBytes.With("serial_number", serialNumber).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(received)
				metricTrafficSentBytes.With("serial_number", serialNumber).With(
					"interface", client.Interface,
					"mac", client.MacAddress,
					"name", name).Set(sent)
			} else if !mikrotik.IsEmptyResponse(err) {
				return nil, err
			}
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	// Ports on mikrotik
	stats, err := b.provider.InterfaceStats(ctx)
	if err == nil {
		for _, stat := range stats {
			metricTrafficReceivedBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.RXByte))
			metricTrafficSentBytes.With("serial_number", serialNumber).With(
				"interface", stat.Name,
				"mac", stat.MacAddress).Set(float64(stat.TXByte))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	resource, err := b.provider.SystemResource(ctx)
	if err == nil {
		metricCPULoad.With("serial_number", serialNumber).Set(float64(resource.CPULoad))
		metricMemoryAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeMemory))
		metricMemoryUsage.With("serial_number", serialNumber).Set(float64(resource.TotalMemory - resource.FreeMemory))
		metricStorageAvailable.With("serial_number", serialNumber).Set(float64(resource.FreeHDDSpace))
		metricStorageUsage.With("serial_number", serialNumber).Set(float64(resource.TotalHDDSpace - resource.FreeHDDSpace))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	disks, err := b.provider.SystemDisk(ctx)
	if err == nil {
		for _, disk := range disks {
			metricDiskUsage.With("serial_number", serialNumber).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Size - disk.Free))
			metricDiskAvailable.With("serial_number", serialNumber).With(
				"name", disk.Name,
				"label", disk.Label,
			).Set(float64(disk.Free))
		}
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	health, err := b.provider.SystemHealth(ctx)
	if err == nil {
		metricVoltage.With("serial_number", serialNumber).Set(health.Voltage)
		metricTemperature.With("serial_number", serialNumber).Set(float64(health.Temperature))
	} else if !mikrotik.IsEmptyResponse(err) {
		return nil, err
	}

	return nil, nil
}
