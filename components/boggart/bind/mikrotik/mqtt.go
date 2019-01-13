package mikrotik

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTTopicWiFiMACState         mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/+/state"
	MQTTTopicWiFiConnectedMAC     mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MQTTTopicWiFiDisconnectedMAC  mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MQTTTopicVPNLoginState        mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/+/state"
	MQTTTopicVPNConnectedLogin    mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/on/login"
	MQTTTopicVPNDisconnectedLogin mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/off/login"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTTopicWiFiMACState,
		MQTTTopicWiFiConnectedMAC,
		MQTTTopicWiFiDisconnectedMAC,
		MQTTTopicVPNLoginState,
		MQTTTopicVPNConnectedLogin,
		MQTTTopicVPNDisconnectedLogin,
	}
}
