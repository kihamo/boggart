package openhab

import (
	"strconv"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart/di"
)

const (
	BindingID = "mqtt"

	DirectoryThings = "things/"
	DirectoryItems  = "items/"
)

var replacerID = strings.NewReplacer(
	":", "_",
	"{", "_",
	"}", "_",
	"[", "_",
	"]", "_",
	"@", "_",
	" ", "_",
	"\"", "_",
	"-", "_",
)

func IDReplace(id string) string {
	return replacerID.Replace(id)
}

func BrokerFromClientOptionsReader(ops *mqtt.ClientOptionsReader) *Broker {
	server := ops.Servers()[0]
	port, _ := strconv.Atoi(server.Port())
	tsl := ops.TLSConfig()

	return NewBroker(ops.ClientID(), server.Hostname()).
		WithLabel("Auto generate from boggart").
		// WithClientID("openhab").
		WithKeepAlive(int(ops.KeepAlive().Seconds())).
		WithUsername(ops.Username()).
		WithPassword(ops.Password()).
		WithTimeoutInMs(ops.WriteTimeout().Milliseconds()).
		WithPort(port).
		WithSecure(tsl != nil)
}

func BindStatusChannel(meta *di.MetaContainer) *Channel {
	return NewChannel("BindStatus", ChannelTypeString).
		WithStateTopic(meta.MQTTTopicStatus()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindStatus", ItemTypeString).
				WithLabel("Bind status [%s]").
				WithIcon("text"),
		)
}

func BindSerialNumberChannel(meta *di.MetaContainer) *Channel {
	return NewChannel("BindSerialNumber", ChannelTypeString).
		WithStateTopic(meta.MQTTTopicStatus()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindSerialNumber", ItemTypeString).
				WithLabel("Bind serial number [%s]").
				WithIcon("text"),
		)
}

func BindMACChannel(meta *di.MetaContainer) *Channel {
	return NewChannel("BindMAC", ChannelTypeString).
		WithStateTopic(meta.MQTTTopicStatus()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindMAC", ItemTypeString).
				WithLabel("Bind MAC address [%s]").
				WithIcon("text"),
		)
}

func GenericThingFromBindMeta(meta *di.MetaContainer) *GenericThing {
	return NewGenericThing(meta.ID()).
		WithLabel(meta.Description())
}

func ItemPrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDReplace(strings.Title(strings.ToLower(meta.ID()))) + "_"
}

func FilePrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDReplace(strings.ToLower(meta.Type()))
}
