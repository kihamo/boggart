package v1

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/mercury/v1"
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
			Values                             *v1.TariffValues
			T1Delta, T2Delta, T3Delta, T4Delta uint64
			T1Trend, T2Trend, T3Trend, T4Trend int64
		}

		date, err := bind.provider.Datetime()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get datetime failed with error %s", "", err.Error()))

			vars["stats"] = make([]*monthly, 0)
		}

		tariffCount, err := bind.provider.TariffCount()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get tariff count failed with error %s", "", err.Error()))
		}

		if err == nil {
			vars["date"] = date
			vars["tariff_count"] = tariffCount

			var (
				last         time.Month
				err          error
				monthRequest time.Month
			)

			if r.URL().Query().Get("all") != "" {
				last = time.December
			} else {
				last = date.Month()
			}

			stats := make([]*monthly, 0, int(last))

			for i := 1; i <= int(last); i++ {
				monthRequest = time.Month(i)

				statsValues, err := bind.provider.MonthlyStatByMonth(monthRequest)
				if err != nil {
					mAsString := t.Translate(r.Context(), monthRequest.String(), "")
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Get statistics for %s failed with error %s", "", mAsString, err.Error()))
					continue
				}

				stat := &monthly{
					Values: statsValues,
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
			powerValues, err := bind.provider.PowerCounters()
			if err != nil {
				mAsString := t.Translate(r.Context(), date.Month().String(), "")
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Get statistics for %s failed with error %s", "", mAsString, err.Error()))
			} else {
				stat := &monthly{
					Month:  date.Month(),
					Year:   date.Year(),
					Values: powerValues,
				}

				stats = append(stats, stat)
			}

			// deltas
			for i, current := range stats {
				if i == 0 {
					continue
				}

				current.T1Delta = current.Values.Tariff1() - stats[i-1].Values.Tariff1()
				current.T2Delta = current.Values.Tariff2() - stats[i-1].Values.Tariff2()
				current.T3Delta = current.Values.Tariff3() - stats[i-1].Values.Tariff3()
				current.T4Delta = current.Values.Tariff4() - stats[i-1].Values.Tariff4()
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

	case "events-on-off":
		if makeDate, err := bind.provider.MakeDate(); err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get make date failed with error %s", "", err.Error()))
		} else {
			type event struct {
				State bool
				Time  time.Time
			}

			events := make([]event, 0, v1.MaxEventsIndex)

			for i := uint8(0); i <= v1.MaxEventsIndex; i++ {
				state, date, err := bind.provider.EventsPowerOnOff(i)

				if err != nil {
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Get event %02x failed with error %s", "", i, err.Error()))
					break
				}

				if date.Before(makeDate) {
					break
				}

				events = append(events, event{
					State: state,
					Time:  date,
				})
			}

			vars["events"] = events
		}

	case "events-open-close":
		if makeDate, err := bind.provider.MakeDate(); err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get make date failed with error %s", "", err.Error()))
		} else {

			type event struct {
				State bool
				Time  time.Time
			}

			events := make([]event, 0, v1.MaxEventsIndex)

			for i := uint8(0); i <= v1.MaxEventsIndex; i++ {
				state, date, err := bind.provider.EventsOpenClose(i)

				if err != nil {
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Get event %02x failed with error %s", "", i, err.Error()))
					break
				}

				if date.Before(makeDate) {
					break
				}

				events = append(events, event{
					State: state,
					Time:  date,
				})
			}

			vars["events"] = events
		}

	case "display":
		var err error

		mode, err := bind.provider.DisplayMode()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get display mode failed with error %s", "", err.Error()))
		}

		timeValues, err := bind.provider.DisplayTime()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get display time failed with error %s", "", err.Error()))
		}

		if r.IsPost() {
			mode.SetTariff1(r.Original().FormValue("mode_t1") != "")
			mode.SetTariff2(r.Original().FormValue("mode_t2") != "")
			mode.SetTariff3(r.Original().FormValue("mode_t3") != "")
			mode.SetTariff4(r.Original().FormValue("mode_t4") != "")
			mode.SetAmount(r.Original().FormValue("mode_amount") != "")
			mode.SetPower(r.Original().FormValue("mode_power") != "")
			mode.SetTime(r.Original().FormValue("mode_time") != "")
			mode.SetDate(r.Original().FormValue("mode_date") != "")

			if mode.IsChanged() {
				err = bind.provider.SetDisplayMode(mode)
				if err != nil {
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Change display mode failed with error %s", "", err.Error()))
				} else {
					r.Session().FlashBag().Success(t.Translate(r.Context(), "Change display mode success", ""))
				}
			}

			var timeValue uint64

			if timeValue, err = strconv.ParseUint(r.Original().FormValue("time_t1"), 10, 64); err == nil {
				timeValues.SetTariff1(timeValue)
			}

			if err == nil {
				if timeValue, err = strconv.ParseUint(r.Original().FormValue("time_t2"), 10, 64); err == nil {
					timeValues.SetTariff2(timeValue)
				}
			}

			if err == nil {
				if timeValue, err = strconv.ParseUint(r.Original().FormValue("time_t3"), 10, 64); err == nil {
					timeValues.SetTariff3(timeValue)
				}
			}

			if err == nil {
				if timeValue, err = strconv.ParseUint(r.Original().FormValue("time_t4"), 10, 64); err == nil {
					timeValues.SetTariff4(timeValue)
				}
			}

			if err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse value failed with error %s", "", err.Error()))
			} else if timeValues.IsChanged() {
				err = bind.provider.SetDisplayTime(timeValues)
				if err != nil {
					r.Session().FlashBag().Error(t.Translate(r.Context(), "Change display time failed with error %s", "", err.Error()))
				} else {
					r.Session().FlashBag().Success(t.Translate(r.Context(), "Change display time success", ""))
				}
			}

			if err == nil {
				t.Redirect(r.URL().Path+"?action="+action, http.StatusFound, w, r)
				return
			}
		}

		vars["mode"] = mode
		vars["time"] = timeValues

	case "holidays":
		days, err := bind.provider.Holidays()
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get holidays failed with error %s", "", err.Error()))
		}

		vars["holidays"] = days

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

		// param last change
		v, err = bind.provider.ParamLastChange()
		if err != v1.ErrCommandNotSupported {
			vars["param_last_change_data"] = map[string]interface{}{
				"value": v,
				"error": err,
			}
		}

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

		// last close cap
		v, err = bind.provider.LastCloseCap()
		vars["last_close_cap_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// V, A, Watts
		voltage, amperage, power, err := bind.provider.UIPCurrent()
		vars["voltage"] = map[string]interface{}{
			"value": voltage,
			"error": err,
		}
		vars["amperage"] = map[string]interface{}{
			"value": amperage,
			"error": err,
		}
		vars["power"] = map[string]interface{}{
			"value": power,
			"error": err,
		}

		// V max
		voltageMax, voltageMaxDate, voltageMaxReset, voltageMaxDateReset, err := bind.provider.MaximumVoltage()
		if err != v1.ErrCommandNotSupported {
			vars["voltage_max"] = map[string]interface{}{
				"value": voltageMax,
				"date":  voltageMaxDate,
				"error": err,
			}

			vars["voltage_max_reset"] = map[string]interface{}{
				"value": voltageMaxReset,
				"date":  voltageMaxDateReset,
				"error": err,
			}
		}

		// A max
		amperageMax, amperageMaxDate, amperageMaxReset, amperageMaxDateReset, err := bind.provider.MaximumAmperage()
		if err != v1.ErrCommandNotSupported {
			vars["amperage_max"] = map[string]interface{}{
				"value": amperageMax,
				"date":  amperageMaxDate,
				"error": err,
			}

			vars["amperage_max_reset"] = map[string]interface{}{
				"value": amperageMaxReset,
				"date":  amperageMaxDateReset,
				"error": err,
			}
		}

		// Watts max
		powerMax, powerMaxDate, powerMaxReset, powerMaxDateReset, err := bind.provider.MaximumPower()
		if err != v1.ErrCommandNotSupported {
			vars["power_max"] = map[string]interface{}{
				"value": powerMax,
				"date":  powerMaxDate,
				"error": err,
			}

			vars["power_max_reset"] = map[string]interface{}{
				"value": powerMaxReset,
				"date":  powerMaxDateReset,
				"error": err,
			}
		}
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
