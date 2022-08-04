package mqtt

type Update struct {
	MessageID          *int64    `json:"msgid,omitempty"`
	MAC                *string   `json:"mac,omitempty"`      // MAC address without :
	Version            *int64    `json:"version,omitempty"`  // Version
	Temperature        *int64    `json:"temp,omitempty"`     // Display temperature (upload after multiplying by 10)
	Humidity           *int64    `json:"humi,omitempty"`     // Display humidity (upload after multiplying by 10)
	SetTemperature     *int64    `json:"settemp,omitempty"`  // Setting temperature (upload after multiplying by 10)
	Mode               *int64    `json:"mode,omitempty"`     // Mode (1 cool, 2 heat, 3 vent, 4 auto)
	OnOff              *int64    `json:"onoff,omitempty"`    // Status 1 on, 2 off)
	FrostTemperature   *int64    `json:"frost,omitempty"`    // Frost temperature (upload after multiplying by 10)
	Delay              *int64    `json:"delay,omitempty"`    // Output delay time
	Diff               *int64    `json:"diff,omitempty"`     // Switch difference (upload after multiplying by 10)
	HoldTemperature    *int64    `json:"holdtemp,omitempty"` // Hold temperature (upload after multiplying by 10)
	HoldTime           *int64    `json:"holdtime,omitempty"` // Hold temperature time
	ScreenLock         *int64    `json:"kb,omitempty"`       // Screen lock (1 no, 2 yes)
	LockPin            *string   `json:"kbkey,omitempty"`    // Lock pin (four numbers)
	TemperatureFormat  *int64    `json:"cf,omitempty"`       // Temperature format (0 celsius, 1 fahrenheit)
	Holiday            *int64    `json:"holiday,omitempty"`  // Holiday (0 no, 1 yes)
	HolidayStartTime   *int64    `json:"holiday_startime,omitempty"`
	HolidayEndTime     *int64    `json:"holiday_endtime,omitempty"`
	StandBy            *int64    `json:"standby,omitempty"`  // Standby (1 no, 2 yes)
	Fan                *int64    `json:"fan,omitempty"`      // Fan speed (1 high, 2 medium, 3 low, 4 auto)
	Timezone           *int64    `json:"timezone,omitempty"` // Time zone
	Program            *int64    `json:"prog,omitempty"`     // Schedule (0 none, 1 weekday/weekend, 2 7 days, 3 24 hours)
	ProgramTemperature *struct{} `json:"temp_prog,omitempty"`
	OnlineType         *int64    `json:"online_type,omitempty"`
}
