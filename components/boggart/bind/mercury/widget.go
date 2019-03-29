package mercury

import (
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{}

	action := r.URL().Query().Get("action")
	vars["action"] = action

	switch action {
	case "display":
		var (
			modeT1, modeT2, modeT3, modeT4, modeAmount, modePower, modeTime, modeDate bool
			timeT1, timeT2, timeT3, timeT4                                            uint64
			err                                                                       error
		)

		if r.IsPost() {
			modeT1 = r.Original().FormValue("mode_t1") != ""
			modeT2 = r.Original().FormValue("mode_t2") != ""
			modeT3 = r.Original().FormValue("mode_t3") != ""
			modeT4 = r.Original().FormValue("mode_t4") != ""
			modeAmount = r.Original().FormValue("mode_amount") != ""
			modePower = r.Original().FormValue("mode_power") != ""
			modeTime = r.Original().FormValue("mode_time") != ""
			modeDate = r.Original().FormValue("mode_date") != ""

			err = bind.provider.SetDisplayMode(modeT1, modeT2, modeT3, modeT4, modeAmount, modePower, modeTime, modeDate)
			if err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Change display mode failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Change display mode success", ""))
			}

			timeT1, err = strconv.ParseUint(r.Original().FormValue("time_t1"), 10, 64)

			if err == nil {
				timeT2, err = strconv.ParseUint(r.Original().FormValue("time_t2"), 10, 64)
			}

			if err == nil {
				timeT3, err = strconv.ParseUint(r.Original().FormValue("time_t3"), 10, 64)
			}

			if err == nil {
				timeT4, err = strconv.ParseUint(r.Original().FormValue("time_t4"), 10, 64)
			}

			if err == nil {
				err = bind.provider.SetDisplayTime(timeT1, timeT2, timeT3, timeT4)
			}

			if err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Change display time failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Change display time success", ""))
				t.Redirect(r.URL().Path+"?action="+action, http.StatusFound, w, r)
				return
			}
		}

		modeT1, modeT2, modeT3, modeT4, modeAmount, modePower, modeTime, modeDate, err = bind.provider.DisplayMode()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get display mode failed with error %s", "", err.Error()))
		}

		vars["mode_t1"] = modeT1
		vars["mode_t2"] = modeT2
		vars["mode_t3"] = modeT3
		vars["mode_t4"] = modeT4
		vars["mode_amount"] = modeAmount
		vars["mode_power"] = modePower
		vars["mode_time"] = modeTime
		vars["mode_date"] = modeDate

		timeT1, timeT2, timeT3, timeT4, err = bind.provider.DisplayTime()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get display time failed with error %s", "", err.Error()))
		}

		vars["time_t1"] = timeT1
		vars["time_t2"] = timeT2
		vars["time_t3"] = timeT3
		vars["time_t4"] = timeT4

	default:
		// date time
		now := time.Now()
		v, err := bind.provider.Datetime()
		variable := map[string]interface{}{
			"value": v,
			"delta": 0,
			"error": err,
		}
		if err == nil {
			variable["delta"] = int64(now.Sub(v).Seconds())
		}

		vars["datetime"] = variable

		// make date
		v, err = bind.provider.MakeDate()
		vars["make_date"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// last power off
		v, err = bind.provider.LastPowerOffDatetime()
		vars["last_power_off_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// last power on
		v, err = bind.provider.LastPowerOnDatetime()
		vars["last_power_on_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
