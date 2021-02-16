package boggart

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

func (i BindStatus) IsStatus(status BindStatus) bool {
	return i == status
}

func (i BindStatus) IsStatusUnknown() bool {
	return i.IsStatus(BindStatusUnknown)
}

func (i BindStatus) IsStatusUninitialized() bool {
	return i.IsStatus(BindStatusUninitialized)
}

func (i BindStatus) IsStatusInitializing() bool {
	return i.IsStatus(BindStatusInitializing)
}

func (i BindStatus) IsStatusOnline() bool {
	return i.IsStatus(BindStatusOnline)
}

func (i BindStatus) IsStatusOffline() bool {
	return i.IsStatus(BindStatusOffline)
}

func (i BindStatus) IsStatusRemoving() bool {
	return i.IsStatus(BindStatusRemoving)
}

func (i BindStatus) IsStatusRemoved() bool {
	return i.IsStatus(BindStatusRemoved)
}

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

type BindRunnable interface {
	Run() error
}
