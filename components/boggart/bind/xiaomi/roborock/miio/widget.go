package miio

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	action := r.URL().Query().Get("action")
	widget := b.Widget()

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	case "settings":
		carpetMode, err := b.device.CarpetMode(ctx)
		if err != nil && !r.IsPost() {
			widget.FlashError(r, "Get carpet mode failed with error %v", "", err)
		} else {
			vars["carpetMode"] = carpetMode
		}

		dnd, err := b.device.DoNotDisturb(ctx)
		if err != nil {
			widget.FlashError(r, "Get do not disturb failed with error %v", "", err)
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
						err = b.device.SetCarpetMode(ctx,
							strings.EqualFold(value[0], "true"),
							carpetMode.CurrentIntegral,
							carpetMode.CurrentHigh,
							carpetMode.CurrentLow,
							carpetMode.StallTime)

					case "carpet-mode-integral":
						var v uint64
						v, err = strconv.ParseUint(value[0], 10, 64)

						if err == nil {
							err = b.device.SetCarpetMode(ctx,
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
							err = b.device.SetCarpetMode(ctx,
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
							err = b.device.SetCarpetMode(ctx,
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
							err = b.device.SetCarpetMode(ctx,
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
							err = b.device.SetFanPower(ctx, v)
						}

					case "dnd-enabled":
						if strings.EqualFold(value[0], "true") {
							err = b.device.SetDoNotDisturb(ctx, dnd.StartHour, dnd.StartMinute, dnd.EndHour, dnd.EndMinute)
						} else {
							err = b.device.DisableDoNotDisturb(ctx)
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
									err = b.device.SetDoNotDisturb(ctx, uint64(start.Hour()), uint64(start.Minute()), uint64(end.Hour()), uint64(end.Minute()))
								}
							}
						} else {
							err = errors.New("bad value")
						}

					case "volume":
						var v uint64

						v, err = strconv.ParseUint(value[0], 10, 64)
						if err == nil {
							err = b.device.SetSoundVolume(ctx, uint32(v))
						}

					case "timezone":
						var v *time.Location

						v, err = time.LoadLocation(value[0])
						if err == nil {
							err = b.device.SetTimezone(ctx, v)
						}
					}

					break
				}
			}

			if err != nil {
				_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
			} else {
				_ = w.SendJSON(boggart.NewResponseJSON().Success("Save success"))
			}

			return
		}

		fanPower, err := b.device.FanPower(ctx)
		if err != nil {
			widget.FlashError(r, "Get fan power failed with error %v", "", err)
		} else {
			vars["fanPower"] = fanPower
		}

		volume, err := b.device.SoundVolume(ctx)
		if err != nil {
			widget.FlashError(r, "Get volume failed with error %v", "", err)
		} else {
			vars["volume"] = volume
		}

		timezone, err := b.device.Timezone(ctx)
		if err != nil {
			widget.FlashError(r, "Get timezone failed with error %v", "", err)
		} else {
			vars["timezone"] = timezone.String()
		}

	case "history":
		summary, err := b.device.CleanSummary(ctx)
		if err != nil {
			widget.FlashError(r, "Get clean summary failed with error %v", "", err)
		} else {
			details := make([]vacuum.CleanDetail, 0, len(summary.CleanupIDs))
			for _, id := range summary.CleanupIDs {
				detail, err := b.device.CleanDetails(ctx, id)
				if err != nil {
					widget.FlashError(r, "Get clean summary detail %d failed with error %v", "", id, err)
					continue
				}

				details = append(details, detail)
			}

			vars["summary"] = summary
			vars["details"] = details
		}

	default:
		vars["packets_counter"] = b.device.Client().PacketsCounter()

		status, err := b.device.Status(ctx)
		if err != nil {
			widget.FlashError(r, "Get status failed with error %v", "", err)
		}

		vars["status"] = status

		info, err := b.device.Info(ctx)
		if err != nil {
			widget.FlashError(r, "Get info failed with error %v", "", err)
		}

		vars["info"] = info

		wifi, err := b.device.WiFiStatus(ctx)
		if err != nil {
			widget.FlashError(r, "Get WiFi status failed with error %v", "", err)
		}

		vars["wifi"] = wifi
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
