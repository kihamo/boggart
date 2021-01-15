package openhab

import (
	"strconv"

	"github.com/eclipse/paho.mqtt.golang"
)

const (
	BrokerTypeID        = "broker"
	BrokerReconnectTime = 60000
	BrokerKeepAlive     = 60
)

type Broker struct {
	*thing

	things GenericThings
}

func NewBroker(id, host string) *Broker {
	return (&Broker{
		thing:  newThing(id, BindingID, BrokerTypeID),
		things: make(GenericThings, 0),
	}).
		WithHost(host).
		WithSecure(true).
		WithEnableDiscovery(true).
		WithReconnectTime(BrokerReconnectTime).
		WithKeepAlive(BrokerKeepAlive)
}

func (b *Broker) BrokerID() string {
	return b.thing.ThingID()
}

func (b *Broker) WithLabel(label string) *Broker {
	b.thing.withLabel(label)
	return b
}

func (b *Broker) WithLocation(location string) *Broker {
	b.thing.withLocation(location)
	return b
}

func (b *Broker) WithHost(host string) *Broker {
	const key = "host"

	if host == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, host)
	return b
}

func (b *Broker) WithPort(port int) *Broker {
	const key = "port"

	if port <= 0 {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, port)
	return b
}

func (b *Broker) WithSecure(secure bool) *Broker {
	const key = "secure"

	if secure {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, secure)
	return b
}

func (b *Broker) WithUsername(username string) *Broker {
	const key = "username"

	if username == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, username)
	return b
}

func (b *Broker) WithPassword(password string) *Broker {
	const key = "password"

	if password == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, password)
	return b
}

func (b *Broker) WithClientID(clientID string) *Broker {
	const key = "clientID"

	if clientID == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, clientID)
	return b
}

func (b *Broker) WithQOS(qos int) *Broker {
	const key = "qos"

	if qos <= 0 || qos > 2 {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, qos)
	return b
}

func (b *Broker) WithKeepAlive(keepAlive int) *Broker {
	const key = "keepAlive"

	if keepAlive == BrokerKeepAlive {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, keepAlive)
	return b
}

func (b *Broker) WithLWTTopic(lwtTopic string) *Broker {
	const key = "lwtTopic"

	if lwtTopic == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, lwtTopic)
	return b
}

func (b *Broker) WithLWTMessage(lwtMessage string) *Broker {
	const key = "lwtMessage"

	if lwtMessage == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, lwtMessage)
	return b
}

func (b *Broker) WithLWTQOS(lwtQos int) *Broker {
	const key = "lwtQos"

	if lwtQos <= 0 || lwtQos > 2 {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, lwtQos)
	return b
}

func (b *Broker) WithLWTRetain(lwtRetain bool) *Broker {
	const key = "lwtRetain"

	if !lwtRetain {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, lwtRetain)
	return b
}

func (b *Broker) WithReconnectTime(reconnectTime int) *Broker {
	const key = "reconnectTime"

	if reconnectTime == BrokerReconnectTime {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, reconnectTime)
	return b
}

func (b *Broker) WithTimeoutInMs(timeoutInMs int64) *Broker {
	const key = "timeoutInMs"

	if timeoutInMs <= 0 {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, timeoutInMs)
	return b
}

func (b *Broker) WithCertificate(certificate string) *Broker {
	const key = "certificate"

	if certificate == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, certificate)
	return b
}

func (b *Broker) WithPublicKey(publicKey string) *Broker {
	const key = "publickey"

	if publicKey == "" {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, publicKey)
	return b
}

func (b *Broker) WithCertificatePin(certificatePin bool) *Broker {
	const key = "certificatepin"

	if !certificatePin {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, certificatePin)
	return b
}

func (b *Broker) WithPublicKeyPin(publicKeyPin bool) *Broker {
	const key = "publickeypin"

	if !publicKeyPin {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, publicKeyPin)
	return b
}

func (b *Broker) WithEnableDiscovery(enableDiscovery bool) *Broker {
	const key = "enableDiscovery"

	if enableDiscovery {
		b.parameters.Delete(key)
		return b
	}

	b.parameters.Set(key, enableDiscovery)
	return b
}

func (b *Broker) AddThings(things ...*GenericThing) *Broker {
	for _, thing := range things {
		thing.WithBroker(b)
		b.things = append(b.things, thing)
	}

	return b
}

func (b *Broker) Things() GenericThings {
	return b.things
}

func (b *Broker) Items() Items {
	items := make(Items, 0)

	for _, thing := range b.Things() {
		items = append(items, thing.Items()...)
	}

	return items
}

func (b *Broker) String() string {
	return "Bridge " + b.thing.String()
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
