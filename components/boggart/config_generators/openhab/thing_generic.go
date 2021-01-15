package openhab

const (
	GenericThingTypeID = "topic"
)

type GenericThings []*GenericThing

func (t GenericThings) String() string {
	var s string

	for i, thing := range t {
		if i > 0 {
			s += "\n"
		}

		s += thing.String()
	}

	return s
}

type GenericThing struct {
	*thing

	channels []*Channel
}

func NewGenericThing(id string) *GenericThing {
	return &GenericThing{
		thing: newThing(id, BindingID, GenericThingTypeID),
	}
}

func (t *GenericThing) GenericThingID() string {
	return t.thing.ThingID()
}

func (t *GenericThing) WithLabel(label string) *GenericThing {
	t.thing.withLabel(label)
	return t
}

func (t *GenericThing) WithLocation(location string) *GenericThing {
	t.thing.withLocation(location)
	return t
}

func (t *GenericThing) WithBroker(broker *Broker) *GenericThing {
	t.thing.withBridge(broker.BrokerID())

	for _, channel := range t.channels {
		channel.WithGenericThing(t)
	}

	return t
}

func (t *GenericThing) WithAvailabilityTopic(availabilityTopic string) *GenericThing {
	const key = "availabilityTopic"

	if availabilityTopic == "" {
		t.parameters.Delete(key)
		return t
	}

	t.parameters.Set(key, availabilityTopic)
	return t
}

func (t *GenericThing) WithPayloadAvailable(payloadAvailable string) *GenericThing {
	const key = "payloadAvailable"

	if payloadAvailable == "" {
		t.parameters.Delete(key)
		return t
	}

	t.parameters.Set(key, payloadAvailable)
	return t
}

func (t *GenericThing) WithPayloadNotAvailable(payloadNotAvailable string) *GenericThing {
	const key = "payloadNotAvailable"

	if payloadNotAvailable == "" {
		t.parameters.Delete(key)
		return t
	}

	t.parameters.Set(key, payloadNotAvailable)
	return t
}

func (t *GenericThing) AddChannels(channels ...*Channel) *GenericThing {
	for _, channel := range channels {
		channel.WithGenericThing(t)
		t.channels = append(t.channels, channel)
	}

	return t
}

func (t *GenericThing) Items() Items {
	items := make(Items, 0)

	for _, channel := range t.channels {
		items = append(items, channel.Items()...)
	}

	return items
}

func (t *GenericThing) String() string {
	s := "Thing " + t.thing.String()

	if len(t.channels) > 0 {
		s += " {\n  Channels:"

		for _, channel := range t.channels {
			s += "\n    " + channel.String()
		}

		s += "\n}"
	}

	return s
}
