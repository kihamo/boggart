package mikrotik

import (
	"context"
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
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WorkersBind

	config            *Config
	address           *url.URL
	provider          *mikrotik.Client
	serialNumberLock  chan struct{}
	serialNumberReady *atomic.Bool
	clientWiFi        *PreloadMap
	clientVPN         *PreloadMap
}

func (b *Bind) Close() error {
	if b.serialNumberReady.IsFalse() {
		close(b.serialNumberLock)
	}

	return nil
}

func (b *Bind) SerialNumberWait(ctx context.Context) (string, error) {
	// максимально ждем один цикл отработки таски
	ctx, cancel := context.WithTimeout(ctx, b.config.ClientsSyncInterval+b.config.ReadinessTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("get serial number failed with error: %v", ctx.Err())

	case <-b.serialNumberLock:
		return b.Meta().SerialNumber(), nil
	}
}

func (b *Bind) SetSerialNumber(serialNumber string) {
	b.Meta().SetSerialNumber(serialNumber)

	if b.serialNumberReady.IsFalse() {
		b.serialNumberReady.True()
		close(b.serialNumberLock)
	}
}

func (b *Bind) updateWiFiClient(ctx context.Context) {
	connections, err := b.provider.InterfaceWirelessRegistrationTable(ctx)
	if err != nil {
		b.Logger().Error("Update WiFi client failed", "error", err.Error())
		return
	}

	sn := b.Meta().SerialNumber()
	active := make(map[interface{}]struct{}, len(connections))

	for _, connection := range connections {
		key := mqtt.NameReplace(connection.MacAddress.String())
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
		b.Logger().Error("Update VPN client failed", "error", err.Error())
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
