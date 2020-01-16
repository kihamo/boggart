package boggart

import (
	"context"

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
}

type BindStatusManager func() BindStatus

type Bind interface {
	Run() error
	SetID(string)
	SetStatusManager(BindStatusManager)
	SerialNumber() string
}

type BindLogger interface {
	SetLogger(logging.Logger)
}

type BindHasReadinessProbe interface {
	ReadinessProbe(context.Context) error
}

type BindHasLivenessProbe interface {
	LivenessProbe(context.Context) error
}
