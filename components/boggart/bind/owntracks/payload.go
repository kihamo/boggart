package owntracks

type CommandPayload struct {
	Type   string `json:"_type"`
	Action string `json:"action"`
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

type WayPointPayload struct {
	Type string  `json:"_type"`
	Desc string  `json:"desc"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Rad  float64 `json:"rad"`
	Tst  int64   `json:"tst"`
}

type WayPointsPayload struct {
	Type      string            `json:"_type"`
	WayPoints []WayPointPayload `json:"waypoints"`
}

type SetWayPointsPayload struct {
	CommandPayload
	WayPoints WayPointsPayload `json:"waypoints"`
}
