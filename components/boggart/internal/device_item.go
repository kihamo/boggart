package internal

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
)

type DeviceItem struct {
	bind        boggart.DeviceBind
	id          string
	t           string
	description string
	tags        []string
	config      interface{}

	cacheMutex           sync.Mutex
	cacheTasks           []workers.Task
	cacheListeners       []workers.ListenerWithEvents
	cacheMQTTSubscribers []mqtt.Subscriber
	cacheMQTTPublishes   []mqtt.Topic
}

func (d *DeviceItem) Bind() boggart.DeviceBind {
	return d.bind
}

func (d *DeviceItem) ID() string {
	return d.id
}

func (d *DeviceItem) Type() string {
	return d.t
}

func (d *DeviceItem) Description() string {
	return d.description
}

func (d *DeviceItem) Tags() []string {
	return d.tags
}

func (d *DeviceItem) Config() interface{} {
	return d.config
}

func (d *DeviceItem) Tasks() []workers.Task {
	c, ok := d.Bind().(boggart.DeviceBindHasTasks)
	if !ok {
		return nil
	}

	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	if d.cacheTasks == nil {
		d.cacheTasks = c.Tasks()
	}

	return d.cacheTasks
}

func (d *DeviceItem) Listeners() []workers.ListenerWithEvents {
	c, ok := d.Bind().(boggart.DeviceBindHasListeners)
	if !ok {
		return nil
	}

	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	if d.cacheListeners == nil {
		d.cacheListeners = c.Listeners()
	}

	return d.cacheListeners
}

func (d *DeviceItem) MQTTSubscribers() []mqtt.Subscriber {
	c, ok := d.Bind().(boggart.DeviceBindHasMQTTSubscribers)
	if !ok {
		return nil
	}

	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	if d.cacheMQTTSubscribers == nil {
		d.cacheMQTTSubscribers = c.MQTTSubscribers()
	}

	return d.cacheMQTTSubscribers
}

func (d *DeviceItem) MQTTPublishes() []mqtt.Topic {
	c, ok := d.Bind().(boggart.DeviceBindHasMQTTPublishes)
	if !ok {
		return nil
	}

	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	if d.cacheMQTTPublishes == nil {
		d.cacheMQTTPublishes = c.MQTTPublishes()
	}

	return d.cacheMQTTPublishes
}
