package boggart

import (
	"errors"
	"sync"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/mitchellh/mapstructure"
)

var (
	bindTypesMutex sync.RWMutex
	bindTypes      = make(map[string]BindType)
)

func RegisterBindType(name string, kind BindType) {
	bindTypesMutex.Lock()
	defer bindTypesMutex.Unlock()

	if kind == nil {
		panic("Bind type name is nil")
	}

	if _, dup := bindTypes[name]; dup {
		panic("Register called twice for bind type " + name)
	}

	bindTypes[name] = kind
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
	Config() interface{}
	CreateBind(config interface{}) (Bind, error)
}

type BindTypeHasWidget interface {
	Widget(*dashboard.Response, *dashboard.Request, BindItem)
	WidgetAssetFS() *assetfs.AssetFS
}

func ValidateBindConfig(t BindType, config map[string]interface{}) (cfg interface{}, err error) {
	if prepare := t.Config(); prepare != nil {
		mapStructureDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &prepare,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeHookFunc(time.RFC3339),
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToIPNetHookFunc(),
				StringToIPHookFunc(),
				StringToMACHookFunc(),
				StringToURLHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				mapstructure.StringToSliceHookFunc(";"),
			),
		})

		if err != nil {
			return cfg, err
		}

		if err := mapStructureDecoder.Decode(config); err != nil {
			return cfg, err
		}

		if _, err = govalidator.ValidateStruct(prepare); err != nil {
			return cfg, err
		}

		cfg = prepare
	} else {
		cfg = config
	}

	return cfg, err
}
