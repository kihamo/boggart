package boggart

import (
	"errors"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kihamo/boggart/types"
	"github.com/mitchellh/mapstructure"
)

var (
	bindTypesMutex sync.RWMutex
	bindTypes      = make(map[string]BindType)
)

func RegisterBindType(name string, kind BindType, aliases ...string) {
	bindTypesMutex.Lock()
	defer bindTypesMutex.Unlock()

	if kind == nil {
		panic("Bind type name is nil")
	}

	for _, name := range append([]string{name}, aliases...) {
		if _, dup := bindTypes[name]; dup {
			panic("Register called twice for bind type " + name)
		}

		bindTypes[name] = kind
	}
}

func GetBindType(name string) (BindType, error) {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	kind, ok := bindTypes[name]
	if !ok {
		return nil, errors.New("bind type " + name + " isn't register")
	}

	return kind, nil
}

func GetBindTypes() map[string]BindType {
	bindTypesMutex.RLock()
	defer bindTypesMutex.RUnlock()

	return bindTypes
}

type BindType interface {
	ConfigDefaults() interface{}
	CreateBind(config interface{}) (Bind, error)
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
