package mikrotik

import (
	"context"
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/listener"
)

type Listener struct {
	listener.BaseListener

	router *Bind
}

func NewListener(router *Bind) *Listener {
	l := &Listener{
		router: router,
	}
	l.Init()

	return l
}

func (b *Bind) Listeners() []workers.ListenerWithEvents {
	return []workers.ListenerWithEvents{
		NewListener(b),
	}
}

func (l *Listener) Events() []workers.Event {
	return []workers.Event{
		boggart.DeviceEventSyslogReceive,
	}
}

func (l *Listener) Name() string {
	return "bind-mikrotik-" + l.router.host
}

func (l *Listener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.DeviceEventSyslogReceive:
		message := args[0].(map[string]interface{})

		client, ok := message["client"]
		if !ok || client != l.router.syslogClient {
			return
		}

		tag, ok := message["tag"]
		if !ok {
			return
		}

		content, ok := message["content"]
		if !ok {
			return
		}

		switch tag {
		case "wifi":
			check := wifiClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 4 {
				return
			}

			if _, err := net.ParseMAC(check[1]); err != nil {
				return
			}

			mac, err := l.router.Mac(ctx, check[1])
			if err != nil {
				return
			}

			sn := l.router.SerialNumber()
			login := mqtt.NameReplace(mac.Address)

			switch check[3] {
			case "connected":
				l.router.MQTTPublishAsync(ctx, MQTTTopicWiFiConnectedMAC.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MQTTTopicWiFiMACState.Format(sn, login), 0, false, []byte(`1`))
			case "disconnected":
				l.router.MQTTPublishAsync(ctx, MQTTTopicWiFiDisconnectedMAC.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MQTTTopicWiFiMACState.Format(sn, login), 0, false, []byte(`0`))
			}

		case "vpn":
			check := vpnClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 2 {
				return
			}

			sn := l.router.SerialNumber()
			login := mqtt.NameReplace(check[1])

			switch check[2] {
			case "in":
				l.router.MQTTPublishAsync(ctx, MQTTTopicVPNConnectedLogin.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MQTTTopicVPNLoginState.Format(sn, login), 0, false, []byte(`1`))
			case "out":
				l.router.MQTTPublishAsync(ctx, MQTTTopicVPNDisconnectedLogin.Format(sn), 0, false, login)
				l.router.MQTTPublishAsync(ctx, MQTTTopicVPNLoginState.Format(sn, login), 0, false, []byte(`0`))
			}
		}
	}
}
