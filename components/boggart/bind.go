package boggart

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/event"
)

var (
	BindEventSyslogReceive = event.NewBaseEvent("SyslogReceive")
	BindEventManagerReady  = event.NewBaseEvent("ManagerReady")
)

type BindStatus uint64

const (
	BindStatusUnknown BindStatus = iota
	BindStatusUninitialized
	BindStatusInitializing
	BindStatusOnline
	BindStatusOffline
	BindStatusRemoving
	BindStatusRemoved
)

type BindItem interface {
	Bind() Bind
	ID() string
	Type() string
	Description() string
	Tags() []string
	Config() interface{}
	Status() BindStatus
	Tasks() []workers.Task
	Listeners() []workers.ListenerWithEvents
	MQTTSubscribers() []mqtt.Subscriber
	MQTTPublishes() []mqtt.Topic
}

type BindStatusGetter func() BindStatus
type BindStatusSetter func(BindStatus)

type Bind interface {
	SetStatusManager(BindStatusGetter, BindStatusSetter)
	SerialNumber() string
}

type BindCloser interface {
	Close() error
}

type BindHasTasks interface {
	Tasks() []workers.Task
}

type BindHasListeners interface {
	Listeners() []workers.ListenerWithEvents
}

type BindHasMQTTClient interface {
	SetMQTTClient(mqtt.Component)
}

type BindHasMQTTSubscribers interface {
	MQTTSubscribers() []mqtt.Subscriber
}

type BindHasMQTTPublishes interface {
	MQTTPublishes() []mqtt.Topic
}
