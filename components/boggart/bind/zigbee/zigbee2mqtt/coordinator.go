package zigbee2mqtt

import (
	"encoding/json"
	"time"
)

type Log struct {
	Message []interface{} `json:"message"`
	Type    string
}

type LogNewAPI struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Settings struct {
	Coordinator struct {
		Meta struct {
			TransportRevision uint8  `json:"transportrev"`
			Product           uint8  `json:"product"`
			MajorRelease      uint8  `json:"majorrel"`
			MinorRelease      uint8  `json:"minorrel"`
			MainTrel          uint8  `json:"maintrel"`
			HardwareRevision  uint32 `json:"revision"`
		} `json:"meta"`
		Type string `json:"type"`
	} `json:"coordinator"`
	LogLevel string `json:"log_level"`
	Network  struct {
		ExtendedPanID string `json:"extendedPanID"`
		Channel       uint32 `json:"channel"`
		PanID         uint16 `json:"panID"`
	} `json:"network"`
	PermitJoin bool   `json:"permit_join"`
	Commit     string `json:"commit"`
	Version    string `json:"version"`

	// New API
	Config *ConfigNewAPI `json:"config,omitempty"`
}

type ConfigNewAPI struct {
	Serial struct {
		DisableLed bool   `json:"disable_led"`
		Port       string `json:"port"`
	} `json:"serial"`
}

type HealthCheck struct {
	Data struct {
		Healthy bool `json:"healthy"`
	} `json:"data"`
	Status string `json:"status"`
}

type Address struct {
	IEEEAddress    string `json:"ieeeAddr"`
	NetworkAddress uint16 `json:"networkAddress"`
}

type Device struct {
	FriendlyName   string `json:"friendly_name"`
	IEEEAddress    string `json:"ieeeAddr"`
	NetworkAddress uint16 `json:"networkAddress"`
	LastSeen       Time   `json:"lastSeen"`
	Type           string `json:"type"`
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	/*
		TODO:
		Optional: Add a last_seen attribute to MQTT messages, contains date/time of last Zigbee message
		possible values are: disable (default), ISO_8601, ISO_8601_local, epoch (default: disable)
		last_seen: 'disable'
	*/

	var v int64

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	(*t).Time = time.Unix(v/1000, v%1000)
	return nil
}

type NetworkMap struct {
	Links []struct {
		Depth                uint8         `json:"depth"`
		LinkQuality          uint8         `json:"linkquality"`
		LQI                  uint8         `json:"lqi"`
		Relationship         uint8         `json:"relationship"`
		Routes               []interface{} `json:"routes"`
		Source               Address       `json:"source"`
		SourceIEEEAddress    string        `json:"sourceIeeeAddr"`
		SourceNetworkAddress uint16        `json:"sourceNwkAddr"`
		Target               Address       `json:"target"`
		TargetIEEEAddress    string        `json:"targetIeeeAddr"`
	} `json:"links"`
	Nodes []struct {
		Address

		Definition       interface{}   `json:"definition"`
		Failed           []interface{} `json:"failed"`
		FriendlyName     string        `json:"friendlyName"`
		LastSeen         *Time         `json:"lastSeen,omitempty"`
		ManufacturerName *string       `json:"manufacturerName,omitempty"`
		ModelID          *string       `json:"modelID,omitempty"`
		Type             string        `json:"type"`
	} `json:"nodes"`
}
