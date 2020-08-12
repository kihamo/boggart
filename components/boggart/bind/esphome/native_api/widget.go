package nativeapi

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/protobuf/proto"
	api "github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/shadow/components/dashboard"
)

type entityRow struct {
	ObjectID string
	Type     string
	Name     string
	State    string
	Entity   proto.Message
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()

	switch q.Get("action") {
	case "entity":
		objectID := q.Get("object")

		if objectID == "" {
			widget.NotFound(w, r)
			return
		}

		entity, err := b.EntityByObjectID(r.Context(), objectID)
		if err != nil {
			widget.NotFound(w, r)
		}

		if api.EntityType(entity) == api.EntityTypeLight {
			b.handleLight(w, r, entity.(*api.ListEntitiesLightResponse))
		}

	case "state":
		b.handleState(w, r)

	case "ota":
		b.handleOTA(w, r)

	default:
		b.handleIndex(w, r)
	}
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) handleIndex(_ *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	widget := b.Widget()
	otaWritten, otaTotal := b.ota.Progress()

	vars := map[string]interface{}{
		"ota_running":  b.ota.IsRunning(),
		"ota_written":  otaWritten,
		"ota_total":    otaTotal,
		"ota_checksum": b.ota.Checksum(),
		"ota_progress": (float64(otaWritten) * float64(100)) / float64(otaTotal),
		"ota_error":    b.ota.LastError(),
	}

	if !b.ota.IsRunning() {
		messages, err := b.provider.ListEntities(ctx)
		entities := make(map[uint32]*entityRow, len(messages))

		if err != nil {
			r.Session().FlashBag().Error(widget.Translate(ctx, "Get list entities failed with error %s", "", err.Error()))
		} else {
			states, err := b.States(ctx, messages...)
			if err != nil {
				r.Session().FlashBag().Error(widget.Translate(ctx, "Get state of entities failed with error: %s", "", err.Error()))
			} else {
				for _, message := range messages {
					e, ok := message.(api.MessageEntity)
					if !ok {
						continue
					}

					entities[e.GetKey()] = &entityRow{
						ObjectID: e.GetObjectId(),
						Name:     e.GetName(),
						Type:     api.EntityType(message),
						Entity:   message,
					}
				}

				for _, message := range states {
					s, ok := message.(api.MessageState)
					if !ok {
						continue
					}

					var row *entityRow

					row, ok = entities[s.GetKey()]
					if !ok {
						continue
					}

					row.State, err = api.State(row.Entity, message, true)
					if err != nil {
						r.Session().FlashBag().Notice(widget.Translate(ctx, "Unknown state type %s for entity with key %d", "", proto.MessageName(message), s.GetKey()))
					}
				}
			}
		}

		vars["entities"] = entities
	}

	widget.Render(ctx, "index", vars)
}

func (b *Bind) handleLight(w *dashboard.Response, r *dashboard.Request, entity *api.ListEntitiesLightResponse) {
	ctx := r.Context()
	widget := b.Widget()

	if r.IsPost() {
		err := r.Original().ParseForm()
		if err != nil {
			r.Session().FlashBag().Error(widget.Translate(ctx, "Parse form failed with error %s", "", err.Error()))
		} else {
			cmd := &api.LightCommandRequest{
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
					r.Session().FlashBag().Error(widget.Translate(ctx, "Parse param %s failed with error %s", "", key, err.Error()))
					err = nil
				}
			}

			if !cmd.HasState {
				cmd.State = false
				cmd.HasState = true
			}

			err = b.provider.LightCommand(ctx, cmd)
			if err != nil {
				r.Session().FlashBag().Error(widget.Translate(ctx, "Execute command failed with error %s", "", err.Error()))
			} else {
				go b.syncState(context.Background(), entity)
				widget.Redirect(r.URL().Path+"?action=entity&object="+entity.GetObjectId(), http.StatusFound, w, r)
				return
			}
		}
	}

	vars := map[string]interface{}{
		"entity": entity,
	}

	states, err := b.States(ctx, entity)
	if err == nil {
		vars["state"] = states[entity.GetKey()]
	}

	vars["error"] = err

	widget.Render(ctx, "light", vars)
}

func (b *Bind) handleState(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()

	objectID := q.Get("object")
	state := q.Get("state")

	if objectID == "" || state == "" {
		widget.NotFound(w, r)
		return
	}

	entity, e := b.EntityByObjectID(r.Context(), objectID)
	if e != nil {
		widget.NotFound(w, r)
		return
	}

	ctx := r.Context()

	var err error

	switch api.EntityType(entity) {
	case api.EntityTypeClimate:
		v, e := strconv.ParseUint(state, 10, 64)
		if e != nil {
			widget.NotFound(w, r)
			return
		}

		s := api.ClimateMode(v)

		switch s {
		case api.ClimateMode_CLIMATE_MODE_OFF:
		case api.ClimateMode_CLIMATE_MODE_AUTO:
		case api.ClimateMode_CLIMATE_MODE_COOL:
		case api.ClimateMode_CLIMATE_MODE_HEAT:
			// skip
		default:
			widget.NotFound(w, r)
			return
		}

		err = b.provider.ClimateCommand(ctx, &api.ClimateCommandRequest{
			Key:     entity.(*api.ListEntitiesClimateResponse).Key,
			HasMode: true,
			Mode:    s,
		})

	case api.EntityTypeFan:
		s, e := strconv.ParseBool(state)
		if e != nil {
			widget.NotFound(w, r)
			return
		}

		err = b.provider.FanCommand(ctx, &api.FanCommandRequest{
			Key:      entity.(*api.ListEntitiesBinarySensorResponse).Key,
			HasState: true,
			State:    s,
		})

	case api.EntityTypeLight:
		s, e := strconv.ParseBool(state)
		if e != nil {
			widget.NotFound(w, r)
			return
		}

		err = b.provider.LightCommand(ctx, &api.LightCommandRequest{
			Key:      entity.(*api.ListEntitiesLightResponse).Key,
			HasState: true,
			State:    s,
		})

	case api.EntityTypeSwitch:
		s, e := strconv.ParseBool(state)
		if e != nil {
			widget.NotFound(w, r)
			return
		}

		err = b.provider.SwitchCommand(ctx, &api.SwitchCommandRequest{
			Key:   entity.(*api.ListEntitiesSwitchResponse).Key,
			State: s,
		})

	default:
		widget.NotFound(w, r)
		return
	}

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	} else {
		r.Session().FlashBag().Success(widget.Translate(ctx, "Success toggle", ""))
	}

	redirectURL := &url.URL{}
	*redirectURL = *r.Original().URL
	redirectURL.RawQuery = ""

	widget.Redirect(redirectURL.String(), http.StatusFound, w, r)
}

func (b *Bind) handleOTA(w *dashboard.Response, r *dashboard.Request) {
	file, header, err := r.Original().FormFile("firmware")

	if err == nil {
		defer file.Close()

		t := header.Header.Get("Content-Type")

		switch t {
		case "application/macbinary":
			buf := bytes.NewBuffer(nil)
			_, err = buf.ReadFrom(file)

			if err == nil {
				err = b.ota.UploadAsync(buf)
			}

		default:
			err = errors.New("unknown content type " + t)
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		_, _ = w.Write([]byte("success"))
	}
}
