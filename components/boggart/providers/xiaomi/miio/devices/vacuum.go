package devices

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
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
INFO	miIO.info	Get device info
ROUTER	miIO.config_router	Set Wifi settings of the device
OTA	miIO.ota	Update firmware over air
OTA_PROG	miIO.get_ota_progress	Update firmware over air Progress
OTA_STATE	miIO.get_ota_state	Update firmware over air Status
*/

const (
	VacuumStatusUnknown uint64 = iota
	VacuumStatusInitiating
	VacuumStatusSleeping
	VacuumStatusWaiting
	VacuumStatusUnknown4
	VacuumStatusCleaning
	VacuumStatusReturningHome
	VacuumStatusRemoteControl
	VacuumStatusCharging
	VacuumStatusChargingError
	VacuumStatusPause
	VacuumStatusSpotCleaning
	VacuumStatusInError
	VacuumStatusShuttingDown
	VacuumStatusUpdating
	VacuumStatusDocking
	VacuumStatusGoTo
	VacuumStatusZoneCleaning
	VacuumStatusFull uint64 = 100

	VacuumFanPowerQuiet    uint64 = 38
	VacuumFanPowerBalanced uint64 = 60
	VacuumFanPowerTurbo    uint64 = 75
	VacuumFanPowerMax      uint64 = 100
	VacuumFanPowerMob      uint64 = 105

	VacuumConsumableFilter    vacuumConsumable = "filter_work_time"
	VacuumConsumableBrushMain vacuumConsumable = "main_brush_work_time"
	VacuumConsumableBrushSide vacuumConsumable = "side_brush_work_time"
	VacuumConsumableSensor    vacuumConsumable = "sensor_dirty_time"

	VacuumConsumableLifetimeFilter    time.Duration = 150
	VacuumConsumableLifetimeBrushMain time.Duration = 300
	VacuumConsumableLifetimeBrushSide time.Duration = 200
	VacuumConsumableLifetimeSensor    time.Duration = 30

	VacuumSoundInstallStateUnknown uint64 = iota
	VacuumSoundInstallStateDownloading
	VacuumSoundInstallStateInstalling
	VacuumSoundInstallStateInstalled
	VacuumSoundInstallStateError

	VacuumSoundInstallErrorNo uint64 = iota
	VacuumSoundInstallErrorUnknown1
	VacuumSoundInstallErrorFailedDownload
	VacuumSoundInstallErrorWrongChecksum
	VacuumSoundInstallErrorUnknown4
	VacuumSoundInstallErrorUnknown5

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
type VacuumStatus struct {
	MessageVersion  uint32
	MessageSequence uint32
	State           uint64
	Battery         uint64
	CleanTime       time.Duration
	CleanArea       uint64 // mm2
	ErrorCode       uint64
	MapPresent      bool
	InCleaning      bool
	InReturning     bool
	InFreshState    bool
	LabStatus       bool
	FanPower        uint64
	DNDEnabled      bool
}

type VacuumCarpetMode struct {
	Enabled         bool
	CurrentIntegral uint64
	CurrentHigh     uint64
	CurrentLow      uint64
	StallTime       uint64
}

type VacuumCleanSummary struct {
	TotalTime     time.Duration
	TotalArea     uint64 // mm2
	TotalCleanups uint64
	CleanupIDs    []uint64
}

type VacuumCleanDetail struct {
	StartTime        time.Time
	EndTime          time.Time
	CleaningDuration time.Duration
	Area             uint64 // mm2
	Completed        bool
}

type VacuumDoNotDisturb struct {
	Enabled     bool
	StartHour   uint64
	StartMinute uint64
	EndHour     uint64
	EndMinute   uint64
}

type VacuumLocale struct {
	Name     string           `json:"name"`
	Bom      string           `json:"bom"`
	Location string           `json:"location"`
	Language string           `json:"language"`
	WiFiPlan string           `json:"wifiplan"`
	Timezone boggart.Location `json:"timezone"`
}

type VacuumSound struct {
	SIDInUse       uint64 `json:"sid_in_use"`
	SIDVersion     uint64 `json:"sid_version"`
	SIDInProgress  uint64 `json:"sid_in_progress"`
	Location       string `json:"location"`
	Bom            string `json:"bom"`
	Language       string `json:"language"`
	MessageVersion uint64 `json:"msg_ver"`
}

type VacuumSoundInstallStatus struct {
	Progress      uint64 `json:"progress"`
	State         uint64 `json:"state"`
	Error         uint64 `json:"error"`
	SIDInProgress uint64 `json:"sid_in_progress"`
}

type VacuumGateway struct {
	IP  boggart.IP           `json:"gateway_ip"`
	MAC boggart.HardwareAddr `json:"gateway_mac"`
}

type VacuumLogUploadStatus struct {
	Status     uint64 `json:"log_upload_status"`
	Location   string `json:"location"`
	PolicyName uint64 `json:"policy_name"`
}

type vacuumConsumable string

type Vacuum struct {
	miio.Device
}

func NewVacuum(address, token string) *Vacuum {
	d := &Vacuum{
		Device: *miio.NewDevice(address, token),
	}

	return d
}

func (d *Vacuum) SerialNumber(ctx context.Context) (string, error) {
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

func (d *Vacuum) Status(ctx context.Context) (result VacuumStatus, err error) {
	type response struct {
		miio.Response

		Result []struct {
			MessageVersion  uint32        `json:"msg_ver"`
			MessageSequence uint32        `json:"msg_seq"`
			State           uint64        `json:"state"`
			Battery         uint64        `json:"battery"`
			CleanTime       time.Duration `json:"clean_time"`
			CleanArea       uint64        `json:"clean_area"` // mm2
			ErrorCode       uint64        `json:"error_code"`
			MapPresent      uint64        `json:"map_present"`
			InCleaning      uint64        `json:"in_cleaning"`
			InReturning     uint64        `json:"in_returning"`
			InFreshState    uint64        `json:"in_fresh_state"`
			LabStatus       uint64        `json:"lab_status"`
			FanPower        uint64        `json:"fan_power"`
			DNDEnabled      uint64        `json:"dnd_enabled"`
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

func (d *Vacuum) Start(ctx context.Context) error {
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

func (d *Vacuum) Spot(ctx context.Context) error {
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

func (d *Vacuum) Stop(ctx context.Context) error {
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

func (d *Vacuum) Pause(ctx context.Context) error {
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

func (d *Vacuum) Home(ctx context.Context) error {
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

func (d *Vacuum) FanPower(ctx context.Context) (uint64, error) {
	type response struct {
		miio.Response

		Result []uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_custom_mode", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SetFanPower(ctx context.Context, power uint64) error {
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

func (d *Vacuum) CarpetMode(ctx context.Context) (result VacuumCarpetMode, err error) {
	type response struct {
		miio.Response

		Result []struct {
			Enabled         uint64 `json:"enabled"`
			CurrentIntegral uint64 `json:"current_integral"`
			CurrentHigh     uint64 `json:"current_high"`
			CurrentLow      uint64 `json:"current_low"`
			StallTime       uint64 `json:"stall_time"`
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

func (d *Vacuum) SetCarpetMode(ctx context.Context, enabled bool, integral, high, low, stallTime uint64) error {
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

func (d *Vacuum) Consumables(ctx context.Context) (map[vacuumConsumable]time.Duration, error) {
	type response struct {
		miio.Response

		Result []map[string]uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_consumable", nil, &reply)
	if err != nil {
		return nil, err
	}

	consumables := make(map[vacuumConsumable]time.Duration, len(reply.Result[0]))
	for n, v := range reply.Result[0] {
		consumables[vacuumConsumable(n)] = time.Duration(v) * time.Second
	}

	return consumables, nil
}

func (d *Vacuum) ConsumableReset(ctx context.Context, consumable vacuumConsumable) error {
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

func (d *Vacuum) CleanSummary(ctx context.Context) (VacuumCleanSummary, error) {
	type response struct {
		miio.Response

		Result []interface{} `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_summary", nil, &reply)
	if err != nil {
		return VacuumCleanSummary{}, err
	}

	result := VacuumCleanSummary{}

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

func (d *Vacuum) CleanDetails(ctx context.Context, id uint64) (VacuumCleanDetail, error) {
	type response struct {
		miio.Response

		Result [][]int64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_clean_record", []uint64{id}, &reply)
	if err != nil {
		return VacuumCleanDetail{}, err
	}

	result := VacuumCleanDetail{}

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

func (d *Vacuum) SoundVolumeTest(ctx context.Context) error {
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

func (d *Vacuum) SoundVolume(ctx context.Context) (uint64, error) {
	type response struct {
		miio.Response

		Result []uint64 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_volume", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SetSoundVolume(ctx context.Context, volume uint64) error {
	if volume > 100 {
		volume = 100
	}

	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "change_sound_volume", []uint64{volume}, &reply)
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

func (d *Vacuum) SoundCurrent(ctx context.Context) (VacuumSound, error) {
	type response struct {
		miio.Response

		Result []VacuumSound `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_current_sound", nil, &reply)
	if err != nil {
		return VacuumSound{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SoundInstall(ctx context.Context, url, md5sum string, sid uint64) (VacuumSoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []VacuumSoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "dnld_install_sound", []map[string]interface{}{{
		"md5": md5sum,
		"url": url,
		"sid": sid,
	}}, &reply)
	if err != nil {
		return VacuumSoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SoundInstallProgress(ctx context.Context) (VacuumSoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []VacuumSoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_progress", nil, &reply)
	if err != nil {
		return VacuumSoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) Find(ctx context.Context) error {
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

func (d *Vacuum) DoNotDisturb(ctx context.Context) (result VacuumDoNotDisturb, err error) {
	type response struct {
		miio.Response

		Result []struct {
			Enabled     uint64 `json:"enabled"`
			StartHour   uint64 `json:"start_hour"`
			StartMinute uint64 `json:"start_minute"`
			EndHour     uint64 `json:"end_hour"`
			EndMinute   uint64 `json:"end_minute"`
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

func (d *Vacuum) DoNotDisturbDisable(ctx context.Context) error {
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

func (d *Vacuum) SetDoNotDisturb(ctx context.Context, startHour, startMinute, endHour, endMinute uint64) error {
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

func (d *Vacuum) Timezone(ctx context.Context) (*time.Location, error) {
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

func (d *Vacuum) SetTimezone(ctx context.Context, zone time.Location) error {
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

func (d *Vacuum) Locale(ctx context.Context) (VacuumLocale, error) {
	type response struct {
		miio.Response

		Result []VacuumLocale `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "app_get_locale", nil, &reply)
	if err != nil {
		return VacuumLocale{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) Gateway(ctx context.Context) (VacuumGateway, error) {
	type response struct {
		miio.Response

		Result []VacuumGateway `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_gateway", nil, &reply)
	if err != nil {
		return VacuumGateway{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) LogUploadStatus(ctx context.Context) (VacuumLogUploadStatus, error) {
	type response struct {
		miio.Response

		Result []VacuumLogUploadStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_log_upload_status", nil, &reply)
	if err != nil {
		return VacuumLogUploadStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Vacuum) SetLabStatus(ctx context.Context, enabled bool) error {
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
