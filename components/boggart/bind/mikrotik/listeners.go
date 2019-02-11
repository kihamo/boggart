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

	bind *Bind
}

func (b *Bind) Listeners() []workers.ListenerWithEvents {
	if b.syslogClient == "" {
		return nil
	}

	return []workers.ListenerWithEvents{
		NewListener(b),
	}
}

func NewListener(bind *Bind) *Listener {
	l := &Listener{
		bind: bind,
	}
	l.Init()

	return l
}

func (l *Listener) Events() []workers.Event {
	return []workers.Event{
		boggart.BindEventSyslogReceive,
	}
}

func (l *Listener) Name() string {
	return "bind-mikrotik-" + l.bind.syslogClient
}

func (l *Listener) Run(ctx context.Context, event workers.Event, t time.Time, args ...interface{}) {
	switch event {
	case boggart.BindEventSyslogReceive:
		message := args[0].(map[string]interface{})

		client, ok := message["client"]
		if !ok || client != l.bind.syslogClient {
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

			mac, err := l.bind.Mac(ctx, check[1])
			if err != nil {
				return
			}

			sn := l.bind.SerialNumber()
			login := mqtt.NameReplace(mac.Address)

			switch check[3] {
			case "connected":
				// TODO:
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiConnectedMAC.Format(sn), login)
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiMACState.Format(sn, login), true)
			case "disconnected":
				// TODO:
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiDisconnectedMAC.Format(sn), login)
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicWiFiMACState.Format(sn, login), false)
			}

		case "vpn":
			check := vpnClientRegexp.FindStringSubmatch(content.(string))
			if len(check) < 2 {
				return
			}

			sn := l.bind.SerialNumber()
			login := mqtt.NameReplace(check[1])

			switch check[2] {
			case "in":
				// TODO:
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicVPNConnectedLogin.Format(sn), login)
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicVPNLoginState.Format(sn, login), true)
			case "out":
				// TODO:
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicVPNDisconnectedLogin.Format(sn), login)
				_ = l.bind.MQTTPublishAsync(ctx, MQTTPublishTopicVPNLoginState.Format(sn, login), false)
			}
		}
	}
}
