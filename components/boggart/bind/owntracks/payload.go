package owntracks

import (
	"errors"
)

const (
	LWTType        = "lwt"
	LocationType   = "location"
	TransitionType = "transition"
	WayPointsType  = "waypoints"
	WayPointType   = "waypoint"

	TransitionEventLeave = "leave"
	TransitionEventEnter = "enter"
)

type Payload interface {
	Valid() error
}

type CommandPayload struct {
	Type   string `json:"_type"`
	Action string `json:"action"`
}

type LWTPayload struct {
	Type string `json:"_type"`
}

func (p LWTPayload) Valid() error {
	if p.Type != LWTType {
		return errors.New("payload isn't " + LWTType)
	}

	return nil
}

type LocationPayload struct {
	Type      string    `json:"_type"`
	Lat       *float64  `json:"lat,omitempty"`
	Lon       *float64  `json:"lon,omitempty"`
	Acc       *int64    `json:"acc,omitempty"`
	Alt       *int64    `json:"alt,omitempty"`
	Batt      *float64  `json:"batt,omitempty"`
	Vel       *int64    `json:"vel,omitempty"`
	Conn      *string   `json:"conn,omitempty"`
	InRegions *[]string `json:"inregions,omitempty"`
}

func (p LocationPayload) Valid() error {
	if p.Type != LocationType {
		return errors.New("payload isn't " + LocationType)
	}

	if p.Lat == nil {
		return errors.New("lat not found in payload")
	}

	if p.Lon == nil {
		return errors.New("lon not found in payload")
	}

	return nil
}

type TransitionPayload struct {
	Type  string  `json:"_type"`
	Acc   float64 `json:"acc"`
	Desc  string  `json:"desc"`
	Event string  `json:"event"`
}

func (p TransitionPayload) IsEnter() bool {
	return p.Event == TransitionEventEnter
}

func (p TransitionPayload) Valid() error {
	if p.Type != TransitionType {
		return errors.New("payload isn't " + TransitionType)
	}

	if p.Event != TransitionEventLeave && p.Event != TransitionEventEnter {
		return errors.New("event name " + p.Event + " isn't " + TransitionEventLeave + " or " + TransitionEventEnter)
	}

	return nil
}

type WayPointPayload struct {
	Type string  `json:"_type"`
	Desc string  `json:"desc"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Rad  float64 `json:"rad"`
	Tst  int64   `json:"tst"`
}

func (p WayPointPayload) Valid() error {
	if p.Type != WayPointType {
		return errors.New("payload isn't " + WayPointType)
	}

	return nil
}

type WayPointsPayload struct {
	Type      string            `json:"_type"`
	WayPoints []WayPointPayload `json:"waypoints"`
}

func (p WayPointsPayload) Valid() error {
	if p.Type != WayPointsType {
		return errors.New("payload isn't " + WayPointsType)
	}

	return nil
}

type SetWayPointsPayload struct {
	CommandPayload
	WayPoints WayPointsPayload `json:"waypoints"`
}
