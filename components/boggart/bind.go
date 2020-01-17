package boggart

import (
	"github.com/kihamo/go-workers/event"
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

type Bind interface{}

type BindRunner interface {
	Run() error
}
