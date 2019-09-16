package esphome

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	subscribeStateTimeoutByEntity = time.Millisecond * 200
)

type entityRow struct {
	Key         uint32
	Type        string
	Name        string
	State       string
	stateFormat string
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()

	list, err := bind.provider.ListEntities(ctx)
	entities := make(map[uint32]*entityRow, len(list))

	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get list entities failed with error %s", "", err.Error()))
	} else {
		for _, message := range list {
			var (
				key         uint32
				typeName    string
				name        string
				stateFormat string
			)

			switch v := message.(type) {
			case *native_api.ListEntitiesBinarySensorResponse:
				typeName, key, name = "binary_sensor", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesCoverResponse:
				typeName, key, name = "cover", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesFanResponse:
				typeName, key, name = "fan", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesLightResponse:
				typeName, key, name = "light", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesSensorResponse:
				typeName, key, name = "sensor", v.GetKey(), v.GetName()
				stateFormat = "%s " + t.Translate(r.Context(), v.GetUnitOfMeasurement(), "")
			case *native_api.ListEntitiesSwitchResponse:
				typeName, key, name = "switch", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesTextSensorResponse:
				typeName, key, name = "text_sensor", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesServicesResponse:
				typeName, key, name = "services", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesCameraResponse:
				typeName, key, name = "camera", v.GetKey(), v.GetName()
			case *native_api.ListEntitiesClimateResponse:
				typeName, key, name = "climate", v.GetKey(), v.GetName()
			default:
				r.Session().FlashBag().Notice(t.Translate(r.Context(), "Unknown entity type %s", "", proto.MessageName(message)))
			}

			if key > 0 {
				entities[key] = &entityRow{
					Key:         key,
					Type:        typeName,
					Name:        name,
					stateFormat: stateFormat,
				}
			}
		}

		i := len(entities)

		timeout := subscribeStateTimeoutByEntity*time.Duration(i) + time.Second
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		chMessage, err := bind.provider.SubscribeStates(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get state of entities failed with error %s", "", err.Error()))
		}

		stateOn := t.Translate(r.Context(), "on", "")
		stateOff := t.Translate(r.Context(), "off", "")

	loop:
		for {
			select {
			case message := <-chMessage:
				var (
					key   uint32
					state string
				)

				switch v := message.(type) {
				case *native_api.BinarySensorStateResponse:
					key = v.GetKey()
					if v.GetState() {
						state = stateOn
					} else {
						state = stateOff
					}
				case *native_api.CoverStateResponse:
					key, state = v.GetKey(), ""
				case *native_api.FanStateResponse:
					key = v.GetKey()
					if v.GetState() {
						state = stateOn
					} else {
						state = stateOff
					}
				case *native_api.LightStateResponse:
					key = v.GetKey()
					if v.GetState() {
						state = stateOn
					} else {
						state = stateOff
					}
				case *native_api.SensorStateResponse:
					key, state = v.GetKey(), strconv.FormatFloat(float64(v.GetState()), 'f', -1, 64)
				case *native_api.SwitchStateResponse:
					key = v.GetKey()
					if v.GetState() {
						state = stateOn
					} else {
						state = stateOff
					}
				case *native_api.TextSensorStateResponse:
					key, state = v.GetKey(), v.GetState()
				case *native_api.ClimateStateResponse:
					key, state = v.GetKey(), v.GetMode().String()
				default:
					r.Session().FlashBag().Notice(t.Translate(r.Context(), "Unknown entity state type %s", "", proto.MessageName(message)))
					break loop
				}

				entity, ok := entities[key]
				if !ok {
					break loop
				}

				if entity.stateFormat != "" {
					state = fmt.Sprintf(entity.stateFormat, state)
				}

				entity.State = state
				i--

				if i == 0 {
					break loop
				}

			case <-ctx.Done():
				break loop
			}
		}
	}

	t.Render(r.Context(), "widget", map[string]interface{}{
		"entities": entities,
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
