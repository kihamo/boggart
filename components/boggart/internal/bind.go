package internal

import (
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
)

type BindItem struct {
	status uint64

	bind        boggart.Bind
	bindType    boggart.BindType
	id          string
	t           string
	description string
	tags        []string
}

type BindItemYaml struct {
	Type        string
	ID          string
	Description string      `yaml:",omitempty"`
	Tags        []string    `yaml:",omitempty"`
	Config      interface{} `yaml:",omitempty"`
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

func (i *BindItem) Status() boggart.BindStatus {
	return boggart.BindStatus(atomic.LoadUint64(&i.status))
}

func (i *BindItem) updateStatus(status boggart.BindStatus) bool {
	value := uint64(status)
	old := atomic.SwapUint64(&i.status, value)

	return old != value
}

func (i *BindItem) MarshalShortYAML() (interface{}, error) {
	return i.marshalYAML(true)
}

func (i *BindItem) MarshalYAML() (interface{}, error) {
	return i.marshalYAML(false)
}

func (i *BindItem) marshalYAML(short bool) (interface{}, error) {
	config, ok := di.ConfigForBind(i.bind)
	if short && ok {
		defaults := i.bindType.ConfigDefaults()

		originalV := reflect.Indirect(reflect.ValueOf(config))
		defaultsV := reflect.Indirect(reflect.ValueOf(defaults))

		shortConfig := make(map[string]interface{}, defaultsV.NumField())

		if defaultsV.Kind() == reflect.Struct && defaultsV.Kind() == originalV.Kind() && defaultsV.NumField() == originalV.NumField() {
			for i := 0; i < defaultsV.NumField(); i++ {
				originalF := originalV.Type().Field(i)
				defaultsF := defaultsV.Type().Field(i)

				if originalF.Name != defaultsF.Name {
					continue
				}

				name := strings.ToLower(defaultsF.Name)
				if tag := defaultsF.Tag.Get("yaml"); tag != "" {
					if val := strings.Split(tag, ",")[0]; val != "" {
						// FIXME: inline
						name = val
					}
				}

				value := originalV.FieldByName(originalF.Name).Interface()

				if reflect.DeepEqual(value, defaultsV.FieldByName(defaultsF.Name).Interface()) {
					continue
				}

				shortConfig[name] = value
			}
		}

		if len(shortConfig) > 0 {
			config = shortConfig
		} else {
			config = nil
		}
	}

	return BindItemYaml{
		Type:        i.Type(),
		ID:          i.ID(),
		Description: i.Description(),
		Tags:        i.Tags(),
		Config:      config,
	}, nil
}
