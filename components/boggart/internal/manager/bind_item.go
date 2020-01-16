package manager

import (
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
)

type BindItem struct {
	status uint64

	bind           boggart.Bind
	bindType       boggart.BindType
	id             string
	t              string
	description    string
	tags           []string
	config         interface{}
	probeReadiness workers.Task
	probeLiveness  workers.Task
}

type BindItemYaml struct {
	Type        string
	ID          string
	Description string
	Tags        []string
	Config      interface{}
}

func (i *BindItem) Bind() boggart.Bind {
	return i.bind
}

func (i *BindItem) BindType() boggart.BindType {
	return i.bindType
}

func (i *BindItem) ID() string {
	return i.id
}

func (i *BindItem) SetID(id string) {
	i.id = id
}

func (i *BindItem) Type() string {
	return i.t
}

func (i *BindItem) Description() string {
	return i.description
}

func (i *BindItem) Tags() []string {
	return i.tags
}

func (i *BindItem) Config() interface{} {
	return i.config
}

func (i *BindItem) Status() boggart.BindStatus {
	return boggart.BindStatus(atomic.LoadUint64(&i.status))
}

func (i *BindItem) updateStatus(status boggart.BindStatus) bool {
	value := uint64(status)
	old := atomic.SwapUint64(&i.status, value)

	return old != value
}

func (i *BindItem) MarshalYAML() (interface{}, error) {
	return BindItemYaml{
		Type:        i.Type(),
		ID:          i.ID(),
		Description: i.Description(),
		Tags:        i.Tags(),
		Config:      i.Config(),
	}, nil
}
