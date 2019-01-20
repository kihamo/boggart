package manager

import (
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
)

type BindItem struct {
	bind        boggart.Bind
	id          string
	t           string
	description string
	tags        []string
	config      interface{}
	status      uint64

	cacheMutex           sync.Mutex
	cacheTasks           []workers.Task
	cacheListeners       []workers.ListenerWithEvents
	cacheMQTTSubscribers []mqtt.Subscriber
	cacheMQTTPublishes   []mqtt.Topic
}

type BindItemYaml struct {
	Type        string
	ID          string
	Description string
	Tags        []string
	Config      interface{}
}

func (d *BindItem) Bind() boggart.Bind {
	return d.bind
}

func (d *BindItem) ID() string {
	return d.id
}

func (d *BindItem) SetID(id string) {
	d.id = id
}

func (d *BindItem) Type() string {
	return d.t
}

func (d *BindItem) Description() string {
	return d.description
}

func (d *BindItem) Tags() []string {
	return d.tags
}

func (d *BindItem) Config() interface{} {
	return d.config
}

func (d *BindItem) Status() boggart.BindStatus {
	return boggart.BindStatus(atomic.LoadUint64(&d.status))
}

func (d *BindItem) updateStatus(status boggart.BindStatus) {
	atomic.StoreUint64(&d.status, uint64(status))
}

func (d *BindItem) Tasks() []workers.Task {
	c, ok := d.Bind().(boggart.BindHasTasks)
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

func (d *BindItem) Listeners() []workers.ListenerWithEvents {
	c, ok := d.Bind().(boggart.BindHasListeners)
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

func (d *BindItem) MQTTSubscribers() []mqtt.Subscriber {
	c, ok := d.Bind().(boggart.BindHasMQTTSubscribers)
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

func (d *BindItem) MQTTPublishes() []mqtt.Topic {
	c, ok := d.Bind().(boggart.BindHasMQTTPublishes)
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

func (d *BindItem) MarshalYAML() (interface{}, error) {
	return BindItemYaml{
		Type:        d.Type(),
		ID:          d.ID(),
		Description: d.Description(),
		Tags:        d.Tags(),
		Config:      d.Config(),
	}, nil
}
