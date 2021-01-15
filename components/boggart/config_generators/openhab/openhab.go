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
