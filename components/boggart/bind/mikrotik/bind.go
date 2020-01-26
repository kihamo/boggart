package mikrotik

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/mikrotik"
)

var (
	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
	vpnClientRegexp  = regexp.MustCompile(`^(\S+) logged (in|out), (.+?)$`)
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config            *Config
	address           *url.URL
	provider          *mikrotik.Client
	serialNumberLock  chan struct{}
	serialNumberReady *atomic.Bool
	clientWiFi        *PreloadMap
	clientVPN         *PreloadMap
}

type Mac struct {
	Address string
	ARP     struct {
		IP      string
		Comment string
	}
	DHCP struct {
		Hostname string
	}
}

func (b *Bind) Close() error {
	if b.serialNumberReady.IsFalse() {
		close(b.serialNumberLock)
	}

	return nil
}

func (b *Bind) SerialNumberWait() string {
	<-b.serialNumberLock

	return b.Meta().SerialNumber()
}

func (b *Bind) SetSerialNumber(serialNumber string) {
	b.Meta().SetSerialNumber(serialNumber)

	if b.serialNumberReady.IsFalse() {
		b.serialNumberReady.True()
		close(b.serialNumberLock)
	}
}

func (b *Bind) Mac(ctx context.Context, mac string) (*Mac, error) {
	if b.Meta().SerialNumber() == "" {
		return nil, errors.New("serial number is empty")
	}

	info := &Mac{
		Address: mac,
	}

	if table, err := b.provider.IPARP(ctx); err == nil {
		for _, row := range table {
			if row.MacAddress == mac {
				info.ARP.IP = row.Address
				info.ARP.Comment = row.Comment
				break
			}
		}
	} else {
		return nil, err
	}

	if leases, err := b.provider.IPDHCPServerLease(ctx); err == nil {
		for _, lease := range leases {
			if lease.MacAddress == mac {
				info.DHCP.Hostname = lease.MacAddress
				break
			}
		}
	} else {
		return nil, err
	}

	return info, nil
}

func (b *Bind) updateWiFiClient(ctx context.Context) {
	connections, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		// TODO: log
		return
	}

	sn := b.Meta().SerialNumber()
	active := make(map[interface{}]struct{}, len(connections))

	for _, connection := range connections {
		mac, err := b.Mac(ctx, connection.MacAddress)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		key := mqtt.NameReplace(mac.Address)
		active[key] = struct{}{}

		if _, ok := b.clientWiFi.Load(key); !ok {
			b.clientWiFi.Store(key, true)
			_ = b.MQTT().PublishAsync(ctx, b.config.TopicWiFiMACState.Format(sn, key), true)
		}
	}

	if b.clientWiFi.IsReady() {
		b.clientWiFi.Range(func(key, value interface{}) bool {
			if _, ok := active[key]; !ok {
				b.clientWiFi.Delete(key)
				_ = b.MQTT().PublishAsync(ctx, b.config.TopicWiFiMACState.Format(sn, key), false)
			}

			return true
		})
	}
}

func (b *Bind) updateVPNClient(ctx context.Context) {
	connections, err := b.provider.PPPActive(ctx)
	if err != nil {
		// TODO: log
		return
	}

	sn := b.Meta().SerialNumber()
	active := make(map[interface{}]struct{}, len(connections))

	for _, connection := range connections {
		key := mqtt.NameReplace(connection.Name)
		active[key] = struct{}{}

		if _, ok := b.clientVPN.Load(key); !ok {
			b.clientVPN.Store(key, true)
			_ = b.MQTT().PublishAsync(ctx, b.config.TopicVPNLoginState.Format(sn, key), true)
		}
	}

	if b.clientVPN.IsReady() {
		b.clientVPN.Range(func(key, value interface{}) bool {
			if _, ok := active[key]; !ok {
				b.clientVPN.Delete(key)
				_ = b.MQTT().PublishAsync(ctx, b.config.TopicVPNLoginState.Format(sn, key), false)
			}

			return true
		})
	}
}
