package esphome

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/shadow/components/dashboard"
)

type entityRow struct {
	ObjectID string
	Type     string
	Name     string
	State    string
	Entity   proto.Message
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	q := r.URL().Query()

	switch q.Get("action") {
	case "entity":
		objectID := q.Get("object")

		if objectID == "" {
			t.NotFound(w, r)
			return
		}

		entity, err := bind.EntityByObjectID(r.Context(), objectID)
		if err != nil {
			t.NotFound(w, r)
		}

		if native_api.EntityType(entity) == native_api.EntityTypeLight {
			t.handleLight(w, r, bind, entity.(*native_api.ListEntitiesLightResponse))
		}

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
		r.Session().FlashBag().Error(t.Translate(ctx,
			"Get list entities failed with error %s",
			"",
			err.Error(),
		))
	} else {
		states, err := bind.States(ctx, messages...)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx,
				"Get state of entities failed with error: %s",
				"",
				err.Error(),
			))
		} else {
			for _, message := range messages {
				e, ok := message.(native_api.MessageEntity)
				if !ok {
					continue
				}

				entities[e.GetKey()] = &entityRow{
					ObjectID: e.GetObjectId(),
					Name:     e.GetName(),
					Type:     native_api.EntityType(message),
					Entity:   message,
				}
			}

			for _, message := range states {
				s, ok := message.(native_api.MessageState)
				if !ok {
					continue
				}

				var row *entityRow

				row, ok = entities[s.GetKey()]
				if !ok {
					continue
				}

				row.State, err = native_api.State(row.Entity, message, true)
				if err != nil {
					r.Session().FlashBag().Notice(t.Translate(ctx,
						"Unknown state type %s for entity with key %d",
						"",
						proto.MessageName(message),
						s.GetKey(),
					))
				}
			}
		}
	}

	t.Render(ctx, "index", map[string]interface{}{
		"entities": entities,
	})
}

func (t Type) handleLight(w *dashboard.Response, r *dashboard.Request, bind *Bind, entity *native_api.ListEntitiesLightResponse) {
	ctx := r.Context()

	if r.IsPost() {
		err := r.Original().ParseForm()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Parse form failed with error %s", "", err.Error()))
		} else {
			cmd := &native_api.LightCommandRequest{
				Key: entity.Key,
			}

			for key, value := range r.Original().PostForm {
				if len(value) == 0 {
					continue
				}

				switch key {
				case "state":
					cmd.State = value[0] == "on"
					cmd.HasState = true
				case "brightness":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.Brightness = float32(val)
						cmd.HasBrightness = true
					} else {
						err = e
					}
				case "red":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.Red = float32(val)
						cmd.HasRgb = true
					} else {
						err = e
					}
				case "green":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.Green = float32(val)
						cmd.HasRgb = true
					} else {
						err = e
					}
				case "blue":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.Blue = float32(val)
						cmd.HasRgb = true
					} else {
						err = e
					}
				case "white":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.White = float32(val)
						cmd.HasWhite = true
					} else {
						err = e
					}
				case "color-temperature":
					if val, e := strconv.ParseFloat(value[0], 64); e == nil {
						cmd.ColorTemperature = float32(val)
						cmd.HasColorTemperature = true
					} else {
						err = e
					}
				case "effect":
					cmd.Effect = value[0]
					cmd.HasEffect = true
				case "flash-length":
					if val, e := time.ParseDuration(value[0]); e == nil {
						if val > 0 {
							cmd.FlashLength = uint32(val.Seconds())
							cmd.HasFlashLength = true
						}
					} else {
						err = e
					}
				case "transition-length":
					if val, e := time.ParseDuration(value[0]); e == nil {
						cmd.TransitionLength = uint32(val.Seconds())
						cmd.HasTransitionLength = true
					} else {
						err = e
					}
				}

				if err != nil {
					r.Session().FlashBag().Error(t.Translate(ctx, "Parse param %s failed with error %s", "", key, err.Error()))
					err = nil
				}
			}

			if !cmd.HasState {
				cmd.State = false
				cmd.HasState = true
			}

			err = bind.provider.LightCommand(ctx, cmd)
			if err != nil {
				r.Session().FlashBag().Error(t.Translate(ctx, "Execute command failed with error %s", "", err.Error()))
			} else {
				go bind.syncState(context.Background(), entity)
				t.Redirect(r.URL().Path+"?action=entity&object="+entity.GetObjectId(), http.StatusFound, w, r)
				return
			}
		}
	}

	vars := map[string]interface{}{
		"entity": entity,
	}

	states, err := bind.States(ctx, entity)
	if err == nil {
		vars["state"] = states[entity.GetKey()]
	}

	vars["error"] = err

	t.Render(ctx, "light", vars)
}
