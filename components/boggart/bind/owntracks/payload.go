package owntracks

type Command struct {
	Type   string `json:"_type"`
	Action string `json:"action"`
}

type Location struct {
	Type string   `json:"_type"`
	Lat  *float64 `json:"lat,omitempty"`
	Lon  *float64 `json:"lon,omitempty"`
	Acc  *int64   `json:"acc,omitempty"`
	Alt  *int64   `json:"alt,omitempty"`
	Batt *float64 `json:"batt,omitempty"`
	Vel  *int64   `json:"vel,omitempty"`
	Conn *string  `json:"conn,omitempty"`
}
