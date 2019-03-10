package mikrotik

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPublishTopicWiFiMACState         mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/+/state"
	MQTTPublishTopicWiFiConnectedMAC     mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/on/mac"
	MQTTPublishTopicWiFiDisconnectedMAC  mqtt.Topic = boggart.ComponentName + "/router/+/wifi/clients/last/off/mac"
	MQTTPublishTopicVPNLoginState        mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/+/state"
	MQTTPublishTopicVPNConnectedLogin    mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/on/login"
	MQTTPublishTopicVPNDisconnectedLogin mqtt.Topic = boggart.ComponentName + "/router/+/vpn/clients/last/off/login"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicWiFiMACState,
		MQTTPublishTopicWiFiConnectedMAC,
		MQTTPublishTopicWiFiDisconnectedMAC,
		MQTTPublishTopicVPNLoginState,
		MQTTPublishTopicVPNConnectedLogin,
		MQTTPublishTopicVPNDisconnectedLogin,
	}
}
