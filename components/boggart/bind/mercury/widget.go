package mercury

import (
	"net/http"
	"sort"
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
	case "monthly":
		type monthly struct {
			Month                              time.Month
			Year                               int
			T1, T2, T3, T4                     uint64
			T1Delta, T2Delta, T3Delta, T4Delta uint64
			T1Trend, T2Trend, T3Trend, T4Trend int64
		}

		date, err := bind.provider.Datetime()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get datetime failed with error", "", err.Error()))

			vars["stats"] = make([]*monthly, 0)

		} else {
			vars["date"] = date

			var (
				t1, t2, t3, t4 uint64
				last           time.Month
				err            error
				monthRequest   time.Month
			)

			if r.URL().Query().Get("all") != "" {
				last = time.December
			} else {
				last = date.Month()
			}

			stats := make([]*monthly, 0, int(last))

			for i := 1; i <= int(last); i++ {
				monthRequest = time.Month(i)

				t1, t2, t3, t4, err = bind.provider.MonthlyStatByMonth(monthRequest)
				if err != nil {
					mAsString := t.Translate(r.Context(), monthRequest.String(), "")
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Get statistics for %s failed with error", "", mAsString, err.Error()))
					continue
				}

				stat := &monthly{
					T1: t1,
					T2: t2,
					T3: t3,
					T4: t4,
				}

				// счетчик отдает статистику на 1 число следующего месяца, поэтому разделяем месяца
				if monthRequest == time.January {
					stat.Month = time.December
				} else {
					stat.Month = time.Month(int(monthRequest) - 1)
				}

				if stat.Month >= date.Month() {
					stat.Year = date.Year() - 1
				} else {
					stat.Year = date.Year()
				}

				stats = append(stats, stat)
			}

			sort.SliceStable(stats, func(i, j int) bool {
				if stats[i].Year == stats[j].Year {
					return int(stats[i].Month) < int(stats[j].Month)
				}

				return stats[i].Year < stats[j].Year
			})

			// отдельно запрашиваем текущий месяц
			t1, t2, t3, t4, err = bind.provider.PowerCounters()
			if err != nil {
				mAsString := t.Translate(r.Context(), date.Month().String(), "")
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Get statistics for %s failed with error", "", mAsString, err.Error()))
			} else {
				stat := &monthly{
					Month: date.Month(),
					Year:  date.Year(),
					T1:    t1,
					T2:    t2,
					T3:    t3,
					T4:    t4,
				}

				stats = append(stats, stat)
			}

			// deltas
			for i, current := range stats {
				if i == 0 {
					continue
				}

				current.T1Delta = current.T1 - stats[i-1].T1
				current.T2Delta = current.T2 - stats[i-1].T2
				current.T3Delta = current.T3 - stats[i-1].T3
				current.T4Delta = current.T4 - stats[i-1].T4
			}

			// trends
			for i, current := range stats {
				if i < 2 {
					continue
				}

				current.T1Trend = int64(current.T1Delta) - int64(stats[i-1].T1Delta)
				current.T2Trend = int64(current.T2Delta) - int64(stats[i-1].T2Delta)
				current.T3Trend = int64(current.T3Delta) - int64(stats[i-1].T3Delta)
				current.T4Trend = int64(current.T4Delta) - int64(stats[i-1].T4Delta)
			}

			vars["stats"] = stats
		}

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
