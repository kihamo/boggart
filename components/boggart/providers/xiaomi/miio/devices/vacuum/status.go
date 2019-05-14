package vacuum

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xiaomi/miio"
)

const (
	StatusStateUnknown StatusState = iota
	StatusStateInitiating
	StatusStateSleeping
	StatusStateWaiting
	StatusStateUnknown4
	StatusStateCleaning
	StatusStateReturningHome
	StatusStateRemoteControl
	StatusStateCharging
	StatusStateChargingError
	StatusStatePause
	StatusStateSpotCleaning
	StatusStateInError
	StatusStateShuttingDown
	StatusStateUpdating
	StatusStateDocking
	StatusStateGoTo
	StatusStateZoneCleaning
	StatusStateFull StatusState = 100
)

type StatusState uint64
type StatusError uint64

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

type Status struct {
	MessageVersion  uint32        `json:"msg_ver"`
	MessageSequence uint32        `json:"msg_seq"`
	State           StatusState   `json:"state"`
	Battery         uint32        `json:"battery"`
	CleanTime       time.Duration `json:"clean_time"`
	CleanArea       uint32        `json:"clean_area"` // mm2
	Error           StatusError   `json:"error_code"`
	MapPresent      bool          `json:"map_present"`
	InCleaning      bool          `json:"in_cleaning"`
	InReturning     bool          `json:"in_returning"`
	InFreshState    bool          `json:"in_fresh_state"`
	LabStatus       bool          `json:"lab_status"`
	FanPower        uint32        `json:"fan_power"`
	DNDEnabled      bool          `json:"dnd_enabled"`
}

func (s StatusState) String() string {
	switch s {
	case StatusStateInitiating:
		return "initiating"
	case StatusStateSleeping:
		return "sleeping"
	case StatusStateWaiting:
		return "waiting"
	case StatusStateCleaning:
		return "cleaning"
	case StatusStateReturningHome:
		return "returning home"
	case StatusStateRemoteControl:
		return "remote control"
	case StatusStateCharging:
		return "charging"
	case StatusStateChargingError:
		return "charging error"
	case StatusStatePause:
		return "pause"
	case StatusStateSpotCleaning:
		return "spot cleaning"
	case StatusStateInError:
		return "in error"
	case StatusStateShuttingDown:
		return "shutting down"
	case StatusStateUpdating:
		return "updating"
	case StatusStateDocking:
		return "docking"
	case StatusStateGoTo:
		return "go to"
	case StatusStateZoneCleaning:
		return "zone cleaning"
	case StatusStateFull:
		return "full"
	}

	return "unknown #" + strconv.FormatUint(uint64(s), 10)
}

func (e StatusError) String() string {
	switch e {
	case 0:
		return "none"
	case 1:
		return "slightly turn the laser (orange) rangefinder, to ensure unobstruction of its motion"
	case 2:
		return "wipe and lightly press the collision sensor"
	case 3:
		return "move the vacuum cleaner to a different location"
	case 4:
		return "wipe the fall sensor and move the vacuum cleaner away from the edge (for example, from a step)"
	case 5:
		return "pull out the main brush. It is necessary to clean the brush and fix of the main axis of the brush"
	case 6:
		return "pull and clean side brushes"
	case 7:
		return "make sure that no foreign objects have entered the main wheel and move the device to a new location"
	case 8:
		return "provide enough space around the vacuum cleaner"
	case 9:
		return "install the dust bag and filter"
	case 10:
		return "make sure the filter is dry or rinse the filter"
	case 11:
		return "a strong magnetic field is detected, move the vacuum cleaner away from the special tape (virtual wall)"
	case 12:
		return "charge level is too low, charge the device"
	case 13:
		return "the problems with charging, make sure the contact between a vacuum cleaner and a docking station"
	case 14:
		return "problems with charging"
	case 15:
		return "wipe the distance to the wall sensor"
	case 16:
		return "install the vacuum cleaner on a flat surface and turn on the device"
	case 17:
		return "problems with the operation of the side brushes, please reset the system settings"
	case 18:
		return "problems with the operation of the suction fan, reset the system"
	case 19:
		return "unpowered charging station"
	case 20:
		return "error unknown #20"
	case 21:
		return "malfunctions in the movement of the laser rangefinder, remove any foreign objects"
	case 22:
		return "it is necessary to wipe the contact areas to charge the device"
	case 23:
		return "it is necessary to wipe the signal area of the docking station"
	}

	return "unknown #" + strconv.FormatUint(uint64(e), 10)
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
		result.Error = r.Error
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
