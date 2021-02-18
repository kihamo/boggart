package boggart

import (
	"errors"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kihamo/boggart/types"
	"github.com/mitchellh/mapstructure"
)

type BindTypeItem struct {
	typ     BindType
	aliases []string
	isAlias bool
}

func (t *BindTypeItem) Type() BindType {
	return t.typ
}

func (t *BindTypeItem) Aliases() []string {
	return t.aliases
}

func (t *BindTypeItem) IsAlias() bool {
	return t.isAlias
}

var (
	bindTypesMutex sync.RWMutex
	bindTypes      = make(map[string]*BindTypeItem)
)

func RegisterBindType(name string, kind BindType, aliases ...string) {
	bindTypesMutex.Lock()
	defer bindTypesMutex.Unlock()

	if kind == nil {
		panic("Bind type name is nil")
	}

	if _, dup := bindTypes[name]; dup {
		panic("Register called twice for bind type " + name)
	}

	bindTypes[name] = &BindTypeItem{
		typ:     kind,
		aliases: aliases,
	}

	for _, name := range aliases {
		if _, dup := bindTypes[name]; dup {
			panic("Register called twice for bind type " + name)
		}

		bindTypes[name] = &BindTypeItem{
			typ:     kind,
			isAlias: true,
		}
	}
}

func GetBindType(name string) (BindType, error) {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	kind, ok := bindTypes[name]
	if !ok {
		return nil, errors.New("bind type " + name + " isn't register")
	}

	return kind.typ, nil
}

func GetBindTypes() map[string]*BindTypeItem {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	return bindTypes
}

type BindType interface {
	ConfigDefaults() interface{}
	CreateBind() Bind
}

func ValidateBindConfig(t BindType, config interface{}) (cfg interface{}, md *mapstructure.Metadata, err error) {
	if prepare := t.ConfigDefaults(); prepare != nil {
		md = new(mapstructure.Metadata)

		mapStructureDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata: md,
			Result:   &prepare,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				types.StringToTimeHookFunc(time.RFC3339),
				types.StringToTimeDurationHookFunc(),
				types.StringToIPNetHookFunc(),
				types.StringToIPHookFunc(),
				types.StringToMACHookFunc(),
				types.StringToURLHookFunc(),
				types.StringToFileModeHookFunc(),
				types.StringToSliceHookFunc(","),
				types.StringToSliceHookFunc(";"),
			),
		})

		if err != nil {
			return cfg, md, err
		}

		if err := mapStructureDecoder.Decode(config); err != nil {
			return cfg, md, err
		}

		if _, err = govalidator.ValidateStruct(prepare); err != nil {
			return cfg, md, err
		}

		cfg = prepare
	} else {
		cfg = config
	}

	return cfg, md, err
}
