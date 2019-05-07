package devices

import (
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
)

/*
Type	Command	Description
START_VACUUM	app_start	Start vacuuming
STOP_VACUUM	app_stop	Stop vacuuming
START_SPOT	app_spot	Start spot cleaning
PAUSE	app_pause	Pause cleaning
CHARGE	app_charge	Start charging
FIND_ME	find_me	Send findme
// CONSUMABLES_GET	get_consumable	Get consumables status
// CONSUMABLES_RESET	reset_consumable	Reset consumables
// CLEAN_SUMMARY_GET	get_clean_summary	Cleaning details
// CLEAN_RECORD_GET	get_clean_record	Cleaning details
CLEAN_RECORD_MAP_GET	get_clean_record_map	Get the map reference of a historical cleaning
GET_MAP	get_map_v1	Get Map
GET_STATUS	get_status	Get Status information
// GET_SERIAL_NUMBER	get_serial_number	Get Serial #
DND_GET	get_dnd_timer	Do Not Disturb Settings
DND_SET	set_dnd_timer	Set the do not disturb timings
DND_CLOSE	close_dnd_timer	Disable the do not disturb function
TIMER_SET	set_timer	Add a timer
TIMER_UPDATE	upd_timer	Activate/deactivate a timer
TIMER_GET	get_timer	Get Timers
TIMER_DEL	del_timer	Remove a timer
TIMERZONE_GET	get_timezone	Get timezone
TIMERZONE_SET	set_timezone	Set timezone
SOUND_INSTALL	dnld_install_sound	Voice pack installation
SOUND_GET_CURRENT	get_current_sound	Current voice
// SOUND_GET_VOLUME	get_sound_volume	-
LOG_UPLOAD_GET	get_log_upload_status	-
LOG_UPLOAD_ENABLE	enable_log_upload	-
SET_MODE	set_custom_mode	Set the vacuum level
GET_MODE	get_custom_mode	Get the vacuum level
REMOTE_START	app_rc_start	Start remote control
REMOTE_END	app_rc_end	End remote control
REMOTE_MOVE	app_rc_move	Remote control move command
GET_GATEWAY	get_gateway	Get current gatway

Robo Vacuum v2 and v1 with firmware versions 3.3.9_003194 or newer
Type	Command	Description
START_ZONE	app_zoned_clean	Start zone vacuum
GOTO_TARGET	app_goto_target	Send vacuum to coordinates

Generic MiIO Commands
Type	Command	Description
INFO	miIO.info	Get device info
ROUTER	miIO.config_router	Set Wifi settings of the device
OTA	miIO.ota	Update firmware over air
OTA_PROG	miIO.get_ota_progress	Update firmware over air Progress
OTA_STATE	miIO.get_ota_state	Update firmware over air Status
*/

const (
	ConsumableFilter    vacuumConsumable = "filter_work_time"
	ConsumableBrushMain vacuumConsumable = "main_brush_work_time"
	ConsumableBrushSide vacuumConsumable = "side_brush_work_time"
	ConsumableSensor    vacuumConsumable = "sensor_dirty_time"

	ConsumableLifetimeFilter    time.Duration = 150
	ConsumableLifetimeBrushMain time.Duration = 300
	ConsumableLifetimeBrushSide time.Duration = 200
	ConsumableLifetimeSensor    time.Duration = 30
)

// https://github.com/marcelrv/XiaomiRobotVacuumProtocol
type vacuumConsumable string

type CleanSummary struct {
	TotalTime     time.Duration
	TotalArea     uint64 // mm2
	TotalCleanups uint64
	CleanupIDs    []uint64
}

type CleanDetail struct {
	StartTime        time.Time
	EndTime          time.Time
	CleaningDuration time.Duration
	Area             uint64 // mm2
	Completed        bool
}

type Vacuum struct {
	miio.Device
}

func NewVacuum(address, token string) *Vacuum {
	d := &Vacuum{
		Device: *miio.NewDevice(address, token),
	}

	return d
}

func (d *Vacuum) SerialNumber() (string, error) {
	type response struct {
		miio.Response

		Result []struct {
			SerialNumber string `json:"serial_number"`
		} `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_serial_number", nil, &reply)
	if err != nil {
		return "", err
	}

	return reply.Result[0].SerialNumber, nil
}

func (d *Vacuum) Consumables() (map[vacuumConsumable]time.Duration, error) {
	type response struct {
		miio.Response

		Result []map[string]uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_consumable", nil, &reply)
	if err != nil {
		return nil, err
	}

	consumables := make(map[vacuumConsumable]time.Duration, len(reply.Result[0]))
	for n, v := range reply.Result[0] {
		consumables[vacuumConsumable(n)] = time.Duration(v/60/60) * time.Hour
	}

	return consumables, nil
}

func (d *Vacuum) ResetConsumable(consumable vacuumConsumable) error {
	var reply miio.ResponseOK

	err := d.Client().Send("reset_consumable", []string{string(consumable)}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Vacuum) CleanSummary() (CleanSummary, error) {
	type response struct {
		miio.Response

		Result []interface{} `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_clean_summary", nil, &reply)
	if err != nil {
		return CleanSummary{}, err
	}

	result := CleanSummary{}

	for i, v := range reply.Result {
		switch i {
		case 0:
			result.TotalTime = time.Duration(v.(float64)) * time.Second

		case 1:
			result.TotalArea = uint64(v.(float64))

		case 2:
			result.TotalCleanups = uint64(v.(float64))

		case 3:
			values := v.([]interface{})
			result.CleanupIDs = make([]uint64, len(values), len(values))

			for i2, v2 := range values {
				result.CleanupIDs[i2] = uint64(v2.(float64))
			}
		}
	}

	return result, nil
}

func (d *Vacuum) CleanDetails(id uint64) (CleanDetail, error) {
	type response struct {
		miio.Response

		Result [][]int64 `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_clean_record", []uint64{id}, &reply)
	if err != nil {
		return CleanDetail{}, err
	}

	result := CleanDetail{}

	for i, v := range reply.Result[0] {
		switch i {
		case 0:
			result.StartTime = time.Unix(v, 0)

		case 1:
			result.EndTime = time.Unix(v, 0)

		case 2:
			result.CleaningDuration = time.Duration(v) * time.Second

		case 3:
			result.Area = uint64(v)

		case 5:
			result.Completed = v == 1
		}
	}

	return result, nil
}

func (d *Vacuum) SoundVolume() (uint64, error) {
	type response struct {
		miio.Response

		Result []uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_sound_volume", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SetSoundVolume(volume uint64) error {
	if volume > 100 {
		volume = 100
	}

	var reply miio.ResponseOK

	err := d.Client().Send("change_sound_volume", []uint64{volume}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Vacuum) Timezone() (*time.Location, error) {
	type response struct {
		miio.Response

		Result []string `json:"result"`
	}

	var reply response

	err := d.Client().Send("get_timezone", nil, &reply)
	if err != nil {
		return nil, err
	}

	return time.LoadLocation(reply.Result[0])
}

func (d *Vacuum) SetTimezone(zone time.Location) error {
	var reply miio.ResponseOK

	err := d.Client().Send("set_timezone", []string{zone.String()}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
