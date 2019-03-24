package boggart

import (
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/event"
	"github.com/kihamo/shadow/components/logging"
)

var (
	BindEventSyslogReceive = event.NewBaseEvent("SyslogReceive")
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
	BindType() BindType
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
	Run() error
	SetStatusManager(BindStatusGetter, BindStatusSetter)
	SerialNumber() string
}

type BindLogger interface {
	SetLogger(logging.Logger)
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
