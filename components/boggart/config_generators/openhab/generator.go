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

func IDNormalize(id string) string {
	return replacerID.Replace(id)
}

func IDNormalizeCamelCase(id string) string {
	id = IDNormalize(id)
	id = strings.ReplaceAll(id, "_", " ")
	id = strings.Title(id)
	id = strings.ReplaceAll(id, " ", "")

	return id
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
		WithStateTopic(meta.MQTTTopicSerialNumber()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindSerialNumber", ItemTypeString).
				WithLabel("Bind serial number [%s]").
				WithIcon("text"),
		)
}

func BindMACChannel(meta *di.MetaContainer) *Channel {
	return NewChannel("BindMAC", ChannelTypeString).
		WithStateTopic(meta.MQTTTopicMAC()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindMAC", ItemTypeString).
				WithLabel("Bind MAC address [%s]").
				WithIcon("text"),
		)
}

func GenericThingFromBindMeta(meta *di.MetaContainer) *GenericThing {
	return NewGenericThing(strings.ToLower(meta.ID())).
		WithLabel(meta.Description())
}

func ItemPrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDNormalizeCamelCase(meta.ID()) + "_"
}

func FilePrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDNormalize(strings.ToLower(meta.Type()))
}
