package v1

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/mercury/v1"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	provider, err := b.Provider()

	if err != nil {
		widget.NotFound(w, r)
		return
	}

	vars := map[string]interface{}{}
	ctx := r.Context()

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

		date, err := provider.Datetime()
		if err != nil {
			widget.FlashError(r, "Get datetime failed with error %v", "", err)

			vars["stats"] = make([]*monthly, 0)
		}

		tariffCount, err := provider.TariffCount()
		if err != nil {
			widget.FlashError(r, "Get tariff count failed with error %v", "", err)
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

				statsValues, err := provider.MonthlyStatByMonth(monthRequest)
				if err != nil {
					mAsString := widget.Translate(ctx, monthRequest.String(), "")
					widget.FlashError(r, "Get statistics for %s failed with error %v", "", mAsString, err)

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
			powerValues, err := provider.PowerCounters()
			if err != nil {
				mAsString := widget.Translate(ctx, date.Month().String(), "")
				widget.FlashError(r, "Get statistics for %s failed with error %v", "", mAsString, err)
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
		if makeDate, err := provider.MakeDate(); err != nil {
			widget.FlashError(r, "Get make date failed with error %v", "", err)
		} else {
			type event struct {
				State bool
				Time  time.Time
			}

			events := make([]event, 0, v1.MaxEventsIndex)

			for i := uint8(0); i <= v1.MaxEventsIndex; i++ {
				state, date, err := provider.EventsPowerOnOff(i)

				if err != nil {
					widget.FlashError(r, "Get event %02x failed with error %v", "", i, err)
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
		if makeDate, err := provider.MakeDate(); err != nil {
			widget.FlashError(r, "Get make date failed with error %v", "", err)
		} else {
			type event struct {
				State bool
				Time  time.Time
			}

			events := make([]event, 0, v1.MaxEventsIndex)

			for i := uint8(0); i <= v1.MaxEventsIndex; i++ {
				state, date, err := provider.EventsOpenClose(i)

				if err != nil {
					widget.FlashError(r, "Get event %02x failed with error %v", "", i, err)
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

		mode, err := provider.DisplayMode()
		if err != nil {
			widget.FlashError(r, "Get display mode failed with error %v", "", err)
		}

		timeValues, err := provider.DisplayTime()
		if err != nil {
			widget.FlashError(r, "Get display time failed with error %v", "", err)
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
				err = provider.SetDisplayMode(mode)
				if err != nil {
					widget.FlashError(r, "Change display mode failed with error %v", "", err)
				} else {
					widget.FlashSuccess(r, "Change display mode success", "")
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
				widget.FlashError(r, "Parse value failed with error %v", "", err)
			} else if timeValues.IsChanged() {
				err = provider.SetDisplayTime(timeValues)
				if err != nil {
					widget.FlashError(r, "Change display time failed with error %v", "", err)
				} else {
					widget.FlashSuccess(r, "Change display time success", "")
				}
			}

			if err == nil {
				widget.Redirect(r.URL().Path+"?action="+action, http.StatusFound, w, r)
				return
			}
		}

		vars["mode"] = mode
		vars["time"] = timeValues

	case "holidays":
		days, err := provider.Holidays()
		if err != nil {
			widget.FlashError(r, "Get holidays failed with error %v", "", err)
		}

		vars["holidays"] = days

	default:
		// date time
		now := time.Now()
		v, err := provider.Datetime()
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
		v, err = provider.ParamLastChange()
		if !v1.CommandNotSupported(err) {
			vars["param_last_change_data"] = map[string]interface{}{
				"value": v,
				"error": err,
			}
		}

		// make date
		v, err = provider.MakeDate()
		vars["make_date"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// version
		version, _, err := provider.Version()
		vars["version"] = map[string]interface{}{
			"value": version,
			"error": err,
		}
		// last power off
		v, err = provider.LastPowerOffDatetime()
		vars["last_power_off_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// last power on
		v, err = provider.LastPowerOnDatetime()
		vars["last_power_on_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// last close cap
		v, err = provider.LastCloseCap()
		vars["last_close_cap_datetime"] = map[string]interface{}{
			"value": v,
			"error": err,
		}

		// V, A, Watts
		voltage, amperage, power, err := provider.UIPCurrent()
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
		voltageMax, voltageMaxDate, voltageMaxReset, voltageMaxDateReset, err := provider.MaximumVoltage()
		if !v1.CommandNotSupported(err) {
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
		amperageMax, amperageMaxDate, amperageMaxReset, amperageMaxDateReset, err := provider.MaximumAmperage()
		if !v1.CommandNotSupported(err) {
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
		powerMax, powerMaxDate, powerMaxReset, powerMaxDateReset, err := provider.MaximumPower()
		if !v1.CommandNotSupported(err) {
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

		// model
		twoSensors, relay, err := provider.Model()
		vars["model_two_sensors"] = map[string]interface{}{
			"value": twoSensors,
			"error": err,
		}
		vars["model_relay"] = map[string]interface{}{
			"value": relay,
			"error": err,
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
