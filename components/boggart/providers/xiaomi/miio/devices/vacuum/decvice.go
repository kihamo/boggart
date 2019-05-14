package vacuum

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
)

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

type Locale struct {
	Name     string           `json:"name"`
	Bom      string           `json:"bom"`
	Location string           `json:"location"`
	Language string           `json:"language"`
	WiFiPlan string           `json:"wifiplan"`
	Timezone boggart.Location `json:"timezone"`
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
