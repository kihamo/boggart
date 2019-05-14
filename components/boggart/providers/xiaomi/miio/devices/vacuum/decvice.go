package vacuum

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio/internal"
)

/*
Type	Command	Description
// START_VACUUM	app_start	Start vacuuming
// STOP_VACUUM	app_stop	Stop vacuuming
// START_SPOT	app_spot	Start spot cleaning
// PAUSE	app_pause	Pause cleaning
// CHARGE	app_charge	Start charging
// FIND_ME	find_me	Send findme
// CONSUMABLES_GET	get_consumable	Get consumables status
// CONSUMABLES_RESET	reset_consumable	Reset consumables
// CLEAN_SUMMARY_GET	get_clean_summary	Cleaning details
// CLEAN_RECORD_GET	get_clean_record	Cleaning details
CLEAN_RECORD_MAP_GET	get_clean_record_map	Get the map reference of a historical cleaning
GET_MAP	get_map_v1	Get Map
// GET_STATUS	get_status	Get Status information
// GET_SERIAL_NUMBER	get_serial_number	Get Serial #
// DND_GET	get_dnd_timer	Do Not Disturb Settings
// DND_SET	set_dnd_timer	Set the do not disturb timings
// DND_CLOSE	close_dnd_timer	Disable the do not disturb function
TIMER_SET	set_timer	Add a timer
TIMER_UPDATE	upd_timer	Activate/deactivate a timer
TIMER_GET	get_timer	Get Timers
TIMER_DEL	del_timer	Remove a timer
// TIMERZONE_GET	get_timezone	Get timezone
// TIMERZONE_SET	set_timezone	Set timezone
// SOUND_INSTALL	dnld_install_sound	Voice pack installation
// SOUND_GET_CURRENT	get_current_sound	Current voice
// SOUND_GET_VOLUME	get_sound_volume	-
// LOG_UPLOAD_GET	get_log_upload_status	-
LOG_UPLOAD_ENABLE	enable_log_upload	-
// SET_MODE	set_custom_mode	Set the vacuum level
// GET_MODE	get_custom_mode	Get the vacuum level
REMOTE_START	app_rc_start	Start remote control
REMOTE_END	app_rc_end	End remote control
REMOTE_MOVE	app_rc_move	Remote control move command
// GET_GATEWAY	get_gateway	Get current gatway

Robo Vacuum v2 and v1 with firmware versions 3.3.9_003194 or newer
Type	Command	Description
START_ZONE	app_zoned_clean	Start zone vacuum
GOTO_TARGET	app_goto_target	Send vacuum to coordinates

Generic MiIO Commands
Type	Command	Description
// INFO	miIO.info	Get device info
ROUTER	miIO.config_router	Set Wifi settings of the device
OTA	miIO.ota	Update firmware over air
OTA_PROG	miIO.get_ota_progress	Update firmware over air Progress
// OTA_STATE	miIO.get_ota_state	Update firmware over air Status
*/

const (
	StatusUnknown uint64 = iota
	StatusInitiating
	StatusSleeping
	StatusWaiting
	StatusUnknown4
	StatusCleaning
	StatusReturningHome
	StatusRemoteControl
	StatusCharging
	StatusChargingError
	StatusPause
	StatusSpotCleaning
	StatusInError
	StatusShuttingDown
	StatusUpdating
	StatusDocking
	StatusGoTo
	StatusZoneCleaning
	StatusFull uint64 = 100

	FanPowerQuiet    uint64 = 38
	FanPowerBalanced uint64 = 60
	FanPowerTurbo    uint64 = 75
	FanPowerMax      uint64 = 100
	FanPowerMob      uint64 = 105

	ConsumableFilter    consumable = "filter_work_time"
	ConsumableBrushMain consumable = "main_brush_work_time"
	ConsumableBrushSide consumable = "side_brush_work_time"
	ConsumableSensor    consumable = "sensor_dirty_time"

	ConsumableLifetimeFilter    time.Duration = 150
	ConsumableLifetimeBrushMain time.Duration = 300
	ConsumableLifetimeBrushSide time.Duration = 200
	ConsumableLifetimeSensor    time.Duration = 30
)

const (
	SoundInstallStateUnknown uint64 = iota
	SoundInstallStateDownloading
	SoundInstallStateInstalling
	SoundInstallStateInstalled
	SoundInstallStateError
)

const (
	SoundInstallErrorNo uint64 = iota
	SoundInstallErrorUnknown1
	SoundInstallErrorFailedDownload
	SoundInstallErrorWrongChecksum
	SoundInstallErrorUnknown4
	SoundInstallErrorUnknown5

/*
   0 => 'None',
       1 => 'Laser sensor fault',
       2 => 'Collision sensor error',
       3 => 'Wheel floating',
       4 => 'Cliff sensor fault',
       5 => 'Main brush blocked',
       6 => 'Side brush blocked',
       7 => 'Wheel blocked',
       8 => 'Device stuck',
       9 => 'Dust bin missing',
       10 => 'Filter blocked',
       11 => 'Magnetic field detected',
       12 => 'Low battery',
       13 => 'Charging problem',
       14 => 'Battery failure',
       15 => 'Wall sensor fault',
       16 => 'Uneven surface',
       17 => 'Side brush failure',
       18 => 'Suction fan failure',
       19 => 'Unpowered charging station',
       20 => 'Unknown'
*/
)

// https://github.com/marcelrv/XiaomiRobotVacuumProtocol
type Status struct {
	MessageVersion  uint32        `json:"msg_ver"`
	MessageSequence uint32        `json:"msg_seq"`
	State           uint32        `json:"state"`
	Battery         uint32        `json:"battery"`
	CleanTime       time.Duration `json:"clean_time"`
	CleanArea       uint32        `json:"clean_area"` // mm2
	ErrorCode       uint64        `json:"error_code"`
	MapPresent      bool          `json:"map_present"`
	InCleaning      bool          `json:"in_cleaning"`
	InReturning     bool          `json:"in_returning"`
	InFreshState    bool          `json:"in_fresh_state"`
	LabStatus       bool          `json:"lab_status"`
	FanPower        uint32        `json:"fan_power"`
	DNDEnabled      bool          `json:"dnd_enabled"`
}

type CarpetMode struct {
	Enabled         bool   `json:"enabled"`
	CurrentIntegral uint64 `json:"current_integral"`
	CurrentHigh     uint64 `json:"current_high"`
	CurrentLow      uint64 `json:"current_low"`
	StallTime       uint64 `json:"stall_time"`
}

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

type DoNotDisturb struct {
	Enabled     bool   `json:"enabled"`
	StartHour   uint64 `json:"start_hour"`
	StartMinute uint64 `json:"start_minute"`
	EndHour     uint64 `json:"end_hour"`
	EndMinute   uint64 `json:"end_minute"`
}

type Locale struct {
	Name     string           `json:"name"`
	Bom      string           `json:"bom"`
	Location string           `json:"location"`
	Language string           `json:"language"`
	WiFiPlan string           `json:"wifiplan"`
	Timezone boggart.Location `json:"timezone"`
}

type Sound struct {
	SIDInUse       uint64 `json:"sid_in_use"`
	SIDVersion     uint64 `json:"sid_version"`
	SIDInProgress  uint64 `json:"sid_in_progress"`
	Location       string `json:"location"`
	Bom            string `json:"bom"`
	Language       string `json:"language"`
	MessageVersion uint64 `json:"msg_ver"`
}

type SoundInstallStatus struct {
	Progress      uint64 `json:"progress"`
	State         uint64 `json:"state"`
	Error         uint64 `json:"error"`
	SIDInProgress uint64 `json:"sid_in_progress"`
}

type Gateway struct {
	IP  boggart.IP           `json:"gateway_ip"`
	MAC boggart.HardwareAddr `json:"gateway_mac"`
}

type LogUploadStatus struct {
	Status     uint64 `json:"log_upload_status"`
	Location   string `json:"location"`
	PolicyName uint64 `json:"policy_name"`
}

type consumable string

type Device struct {
	miio.Device
}

func New(address, token string) *Device {
	d := &Device{
		Device: *miio.NewDevice(address, token),
	}

	return d
}

func (d *Device) SerialNumber(ctx context.Context) (string, error) {
	type response struct {
		miio.Response

		Result []struct {
			SerialNumber string `json:"serial_number"`
		} `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_serial_number", nil, &reply)
	if err != nil {
		return "", err
	}

	return reply.Result[0].SerialNumber, nil
}

func (d *Device) Status(ctx context.Context) (result Status, err error) {
	type response struct {
		miio.Response

		Result []struct {
			Status

			MapPresent   uint64 `json:"map_present"`
			InCleaning   uint64 `json:"in_cleaning"`
			InReturning  uint64 `json:"in_returning"`
			InFreshState uint64 `json:"in_fresh_state"`
			LabStatus    uint64 `json:"lab_status"`
			DNDEnabled   uint64 `json:"dnd_enabled"`
		} `json:"result"`
	}

	var reply response

	err = d.Client().Send(ctx, "get_status", nil, &reply)
	if err == nil {
		r := &reply.Result[0]
		result.MessageVersion = r.MessageVersion
		result.MessageSequence = r.MessageSequence
		result.State = r.State
		result.Battery = r.Battery
		result.CleanTime = r.CleanTime * time.Second
		result.CleanArea = r.CleanArea
		result.ErrorCode = r.ErrorCode
		result.MapPresent = r.MapPresent == 1
		result.InCleaning = r.InCleaning == 1
		result.InReturning = r.InReturning == 1
		result.InFreshState = r.InFreshState == 1
		result.LabStatus = r.LabStatus == 1
		result.FanPower = r.FanPower
		result.DNDEnabled = r.DNDEnabled == 1
	}

	return result, err
}

func (d *Device) Start(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "app_start", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Spot(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "app_spot", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Stop(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "app_stop", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Pause(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "app_pause", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Home(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "app_charge", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) FanPower(ctx context.Context) (uint32, error) {
	type response struct {
		miio.Response

		Result []uint32 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_custom_mode", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Device) SetFanPower(ctx context.Context, power uint64) error {
	if power > 105 {
		power = 105
	} else if power < 1 {
		power = 1
	}

	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "set_custom_mode", []uint64{power}, &reply)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) CarpetMode(ctx context.Context) (result CarpetMode, err error) {
	type response struct {
		miio.Response

		Result []struct {
			CarpetMode

			Enabled uint64 `json:"enable"`
		} `json:"result"`
	}

	var reply response

	err = d.Client().Send(ctx, "get_carpet_mode", nil, &reply)
	if err == nil {
		r := &reply.Result[0]
		result.Enabled = r.Enabled == 1
		result.CurrentIntegral = r.CurrentIntegral
		result.CurrentHigh = r.CurrentHigh
		result.CurrentLow = r.CurrentLow
		result.StallTime = r.StallTime
	}

	return result, nil
}

func (d *Device) SetCarpetMode(ctx context.Context, enabled bool, integral, high, low, stallTime uint64) error {
	var reply miio.ResponseOK

	request := map[string]uint64{
		"enable":           0,
		"current_integral": integral,
		"current_high":     high,
		"current_low":      low,
		"stall_time":       stallTime,
	}

	if enabled {
		request["enable"] = 1
	}

	err := d.Client().Send(ctx, "set_carpet_mode", []interface{}{request}, &reply)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) Consumables(ctx context.Context) (map[consumable]time.Duration, error) {
	type response struct {
		miio.Response

		Result []map[string]uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_consumable", nil, &reply)
	if err != nil {
		return nil, err
	}

	consumables := make(map[consumable]time.Duration, len(reply.Result[0]))
	for n, v := range reply.Result[0] {
		consumables[consumable(n)] = time.Duration(v) * time.Second
	}

	return consumables, nil
}

func (d *Device) ConsumableReset(ctx context.Context, consumable consumable) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "reset_consumable", []string{string(consumable)}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) CleanSummary(ctx context.Context) (CleanSummary, error) {
	type response struct {
		miio.Response

		Result []interface{} `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_summary", nil, &reply)
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

func (d *Device) CleanDetails(ctx context.Context, id uint64) (CleanDetail, error) {
	type response struct {
		miio.Response

		Result [][]int64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_record", []uint64{id}, &reply)
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

func (d *Device) SoundVolumeTest(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "test_sound_volume", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) SoundVolume(ctx context.Context) (uint32, error) {
	type response struct {
		miio.Response

		Result []uint32 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_volume", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Device) SetSoundVolume(ctx context.Context, volume uint32) error {
	if volume > 100 {
		volume = 100
	}

	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "change_sound_volume", []uint32{volume}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

// Мой кастомный {"result":[{"sid_in_use":10000,"sid_version":1,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236719}
// English       {"result":[{"sid_in_use":3,"sid_version":2,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236769}
// По-умолчанию  {"result":[{"sid_in_use":1,"sid_version":2,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236821}

func (d *Device) SoundCurrent(ctx context.Context) (Sound, error) {
	type response struct {
		miio.Response

		Result []Sound `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_current_sound", nil, &reply)
	if err != nil {
		return Sound{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstall(ctx context.Context, url, md5sum string, sid uint64) (SoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []SoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "dnld_install_sound", map[string]interface{}{
		"md5": md5sum,
		"url": url,
		"sid": sid,
		//"sver": 2,
	}, &reply)
	if err != nil {
		return SoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstallProgress(ctx context.Context) (SoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []SoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_progress", nil, &reply)
	if err != nil {
		return SoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstallLocalServer(ctx context.Context, file io.ReadSeeker, sid uint64) error {
	server, err := internal.NewServer(file, d.HostnameForLocalServer())
	if err != nil {
		return err
	}
	defer server.Close()

	var status SoundInstallStatus
	status, err = d.SoundInstall(ctx, server.URL().String(), server.MD5(), sid)
	if err == nil {
		if status.Error != SoundInstallErrorNo {
			return errors.New("return error code " + strconv.FormatUint(status.Error, 10))
		}

		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			if status, err := d.SoundInstallProgress(ctx); err == nil {
				if status.State == SoundInstallStateDownloading || status.State == SoundInstallStateInstalling {
					continue
				}

				if status.Error != SoundInstallErrorNo {
					return errors.New("return error code " + strconv.FormatUint(status.Error, 10))
				}

				return nil
			}
		}
	}

	return err
}

func (d *Device) Find(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "find_me", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) DoNotDisturb(ctx context.Context) (result DoNotDisturb, err error) {
	type response struct {
		miio.Response

		Result []struct {
			DoNotDisturb

			Enabled uint64 `json:"enabled"`
		} `json:"result"`
	}

	var reply response

	err = d.Client().Send(ctx, "get_dnd_timer", nil, &reply)
	if err == nil {
		r := &reply.Result[0]
		result.Enabled = r.Enabled == 1
		result.StartHour = r.StartHour
		result.StartMinute = r.StartMinute
		result.EndHour = r.EndHour
		result.EndMinute = r.EndMinute
	}

	return result, err
}

func (d *Device) DoNotDisturbDisable(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "close_dnd_timer", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) SetDoNotDisturb(ctx context.Context, startHour, startMinute, endHour, endMinute uint64) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "set_dnd_timer", []uint64{startHour, startMinute, endHour, endMinute}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Timezone(ctx context.Context) (*time.Location, error) {
	type response struct {
		miio.Response

		Result []string `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_timezone", nil, &reply)
	if err != nil {
		return nil, err
	}

	return time.LoadLocation(reply.Result[0])
}

func (d *Device) SetTimezone(ctx context.Context, zone time.Location) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "set_timezone", []string{zone.String()}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) Locale(ctx context.Context) (Locale, error) {
	type response struct {
		miio.Response

		Result []Locale `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "app_get_locale", nil, &reply)
	if err != nil {
		return Locale{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) Gateway(ctx context.Context) (Gateway, error) {
	type response struct {
		miio.Response

		Result []Gateway `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_gateway", nil, &reply)
	if err != nil {
		return Gateway{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) LogUploadStatus(ctx context.Context) (LogUploadStatus, error) {
	type response struct {
		miio.Response

		Result []LogUploadStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_log_upload_status", nil, &reply)
	if err != nil {
		return LogUploadStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SetLabStatus(ctx context.Context, enabled bool) error {
	var (
		reply  miio.ResponseOK
		status int64
	)

	if enabled {
		status = 1
	}

	err := d.Client().Send(ctx, "set_lab_status", []int64{status}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}
