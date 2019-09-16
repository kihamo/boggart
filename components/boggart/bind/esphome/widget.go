package esphome

import (
	"fmt"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/shadow/components/dashboard"
)

type entityRow struct {
	Key         uint32
	Type        string
	Name        string
	State       string
	StateRaw    interface{}
	stateFormat string
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	q := r.URL().Query()

	switch q.Get("action") {
	case "entity":
		id := q.Get("id")

		if id == "" {
			t.NotFound(w, r)
			return
		}

		entityID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			t.NotFound(w, r)
			return
		}

		t.handleEntity(w, r, bind, uint32(entityID))

	default:
		t.handleIndex(w, r, bind)
	}
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (t Type) handleIndex(w *dashboard.Response, r *dashboard.Request, bind *Bind) {
	ctx := r.Context()

	messages, err := bind.provider.ListEntities(ctx)
	entities := make(map[uint32]*entityRow, len(messages))

	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(),
			"Get list entities failed with error %s",
			"",
			err.Error(),
		))
	} else {
		for _, message := range messages {
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
				r.Session().FlashBag().Notice(t.Translate(r.Context(),
					"Unknown entity type %s",
					"",
					proto.MessageName(message),
				))
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

		states, err := bind.States(ctx, messages)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(),
				"Get state of entities failed with error %s",
				"",
				err.Error(),
			))
		}

		stateOn := t.Translate(r.Context(), "on", "")
		stateOff := t.Translate(r.Context(), "off", "")

		for key, entity := range entities {
			message, ok := states[key]
			if !ok {
				r.Session().FlashBag().Notice(t.Translate(r.Context(),
					"State for entity with key %d not found",
					"",
					key,
				))
				continue
			}

			var (
				found    = true
				key      uint32
				state    string
				stateRaw interface{}
			)

			switch v := message.(type) {
			case *native_api.BinarySensorStateResponse:
				key, stateRaw = v.GetKey(), v.GetState()
				if v.GetState() {
					state = stateOn
				} else {
					state = stateOff
				}
			case *native_api.CoverStateResponse:
				key = v.GetKey()
			case *native_api.FanStateResponse:
				key, stateRaw = v.GetKey(), v.GetState()
				if v.GetState() {
					state = stateOn
				} else {
					state = stateOff
				}
			case *native_api.LightStateResponse:
				key, stateRaw = v.GetKey(), v.GetState()
				if v.GetState() {
					state = stateOn
				} else {
					state = stateOff
				}
			case *native_api.SensorStateResponse:
				key, state = v.GetKey(), strconv.FormatFloat(float64(v.GetState()), 'f', -1, 64)
				stateRaw = state
			case *native_api.SwitchStateResponse:
				key, stateRaw = v.GetKey(), v.GetState()
				if v.GetState() {
					state = stateOn
				} else {
					state = stateOff
				}
			case *native_api.TextSensorStateResponse:
				key, state = v.GetKey(), v.GetState()
				stateRaw = state
			case *native_api.ClimateStateResponse:
				key, state = v.GetKey(), v.GetMode().String()
				stateRaw = state
			default:
				r.Session().FlashBag().Notice(t.Translate(r.Context(),
					"Unknown state type %s for entity with key %d",
					"",
					proto.MessageName(message),
					key,
				))
				found = false
			}

			if !found {
				continue
			}

			if entity.stateFormat != "" {
				state = fmt.Sprintf(entity.stateFormat, state)
			}

			entity.State = state
			entity.StateRaw = stateRaw
		}
	}

	t.Render(ctx, "index", map[string]interface{}{
		"entities": entities,
	})
}

func (t Type) handleEntity(w *dashboard.Response, r *dashboard.Request, bind *Bind, entityID uint32) {
	ctx := r.Context()

	list, err := bind.provider.ListEntities(ctx)
	if err != nil {
		t.InternalError(w, r, err)
	}

	for _, message := range list {
		switch v := message.(type) {
		case *native_api.ListEntitiesLightResponse:
			if entityID == v.GetKey() {
				t.handleLight(w, r, bind, v)
				return
			}
		}
	}

	t.NotFound(w, r)
}

func (t Type) handleLight(w *dashboard.Response, r *dashboard.Request, bind *Bind, entity *native_api.ListEntitiesLightResponse) {
	ctx := r.Context()

	vars := map[string]interface{}{
		"entity": entity,
	}

	states, err := bind.States(ctx, []proto.Message{entity})
	if err == nil {
		vars["state"] = states[entity.GetKey()]
	}

	vars["error"] = err

	t.Render(ctx, "light", vars)
}
