package vacuum

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
	"github.com/kihamo/boggart/types"
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
// TIMER_SET	set_timer	Add a timer
// TIMER_UPDATE	upd_timer	Activate/deactivate a timer
// TIMER_GET	get_timer	Get Timers
// TIMER_DEL	del_timer	Remove a timer
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
// OTA	miIO.ota	Update firmware over air
// OTA_PROG	miIO.get_ota_progress	Update firmware over air Progress
// OTA_STATE	miIO.get_ota_state	Update firmware over air Status
*/

type Locale struct {
	Name     string         `json:"name"`
	Bom      string         `json:"bom"`
	Location string         `json:"location"`
	Language string         `json:"language"`
	WiFiPlan string         `json:"wifiplan"`
	Timezone types.Location `json:"timezone"`
}

type Gateway struct {
	IP  types.IP           `json:"gateway_ip"`
	MAC types.HardwareAddr `json:"gateway_mac"`
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
