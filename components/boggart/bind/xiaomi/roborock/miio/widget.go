package miio

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()
	action := r.URL().Query().Get("action")

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	case "settings":
		carpetMode, err := bind.device.CarpetMode(ctx)
		if err != nil && !r.IsPost() {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get carpet mode failed with error %s", "", err.Error()))
		} else {
			vars["carpetMode"] = carpetMode
		}

		dnd, err := bind.device.DoNotDisturb(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get do not disturb failed with error %s", "", err.Error()))
		} else {
			now := time.Now()

			vars["dnd"] = dnd
			vars["dnd_start"] = time.Date(now.Year(), now.Month(), now.Day(), int(dnd.StartHour), int(dnd.StartMinute), 0, 0, now.Location())
			vars["dnd_end"] = time.Date(now.Year(), now.Month(), now.Day(), int(dnd.EndHour), int(dnd.EndMinute), 0, 0, now.Location()).AddDate(0, 0, 1)
		}

		if r.IsPost() {
			err = r.Original().ParseForm()
			if err == nil {
				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					switch key {
					case "carpet-mode-enabled":
						err = bind.device.SetCarpetMode(ctx,
							strings.EqualFold(value[0], "true"),
							carpetMode.CurrentIntegral,
							carpetMode.CurrentHigh,
							carpetMode.CurrentLow,
							carpetMode.StallTime)

					case "carpet-mode-integral":
						var v uint64
						v, err = strconv.ParseUint(value[0], 10, 64)

						if err == nil {
							err = bind.device.SetCarpetMode(ctx,
								carpetMode.Enabled,
								v,
								carpetMode.CurrentHigh,
								carpetMode.CurrentLow,
								carpetMode.StallTime)
						}

					case "carpet-mode-high":
						var v uint64
						v, err = strconv.ParseUint(value[0], 10, 64)

						if err == nil {
							err = bind.device.SetCarpetMode(ctx,
								carpetMode.Enabled,
								carpetMode.CurrentIntegral,
								v,
								carpetMode.CurrentLow,
								carpetMode.StallTime)
						}

					case "carpet-mode-low":
						var v uint64
						v, err = strconv.ParseUint(value[0], 10, 64)

						if err == nil {
							err = bind.device.SetCarpetMode(ctx,
								carpetMode.Enabled,
								carpetMode.CurrentIntegral,
								carpetMode.CurrentHigh,
								v,
								carpetMode.StallTime)
						}

					case "carpet-mode-stall-time":
						var v uint64
						v, err = strconv.ParseUint(value[0], 10, 64)

						if err == nil {
							err = bind.device.SetCarpetMode(ctx,
								carpetMode.Enabled,
								carpetMode.CurrentIntegral,
								carpetMode.CurrentHigh,
								carpetMode.CurrentLow,
								v)
						}

					case "fan-power":
						var v uint64

						v, err = strconv.ParseUint(value[0], 10, 64)
						if err == nil {
							err = bind.device.SetFanPower(ctx, v)
						}

					case "dnd-enabled":
						if strings.EqualFold(value[0], "true") {
							err = bind.device.SetDoNotDisturb(ctx, dnd.StartHour, dnd.StartMinute, dnd.EndHour, dnd.EndMinute)
						} else {
							err = bind.device.DisableDoNotDisturb(ctx)
						}

					case "dnd-time":
						split := strings.Split(value[0], "-")

						if len(split) == 2 {
							var (
								start time.Time
								end   time.Time
							)

							start, err = time.Parse("15:04", strings.TrimSpace(split[0]))
							if err == nil {
								end, err = time.Parse("15:04", strings.TrimSpace(split[1]))
								if err == nil {
									err = bind.device.SetDoNotDisturb(ctx, uint64(start.Hour()), uint64(start.Minute()), uint64(end.Hour()), uint64(end.Minute()))
								}
							}
						} else {
							err = errors.New("bad value")
						}

					case "volume":
						var v uint64

						v, err = strconv.ParseUint(value[0], 10, 64)
						if err == nil {
							err = bind.device.SetSoundVolume(ctx, uint32(v))
						}

					case "timezone":
						var v *time.Location

						v, err = time.LoadLocation(value[0])
						if err == nil {
							err = bind.device.SetTimezone(ctx, v)
						}
					}

					break
				}
			}

			if err != nil {
				_ = w.SendJSON(response{
					Result:  "failed",
					Message: err.Error(),
				})

			} else {
				_ = w.SendJSON(response{
					Result:  "success",
					Message: "Save success",
				})
			}

			return
		}

		fanPower, err := bind.device.FanPower(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get fan power failed with error %s", "", err.Error()))
		} else {
			vars["fanPower"] = fanPower
		}

		volume, err := bind.device.SoundVolume(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get volume failed with error %s", "", err.Error()))
		} else {
			vars["volume"] = volume
		}

		timezone, err := bind.device.Timezone(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get timezone failed with error %s", "", err.Error()))
		} else {
			vars["timezone"] = timezone.String()
		}

	default:
		status, err := bind.device.Status(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get status failed with error %s", "", err.Error()))
		}

		vars["status"] = status

		info, err := bind.device.Info(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get info failed with error %s", "", err.Error()))
		}

		vars["info"] = info

		wifi, err := bind.device.WiFiStatus(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get WiFi status failed with error %s", "", err.Error()))
		}

		vars["wifi"] = wifi
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
