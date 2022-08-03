package mqtt

type Update struct {
	MessageID          *int64    `json:"msgid"`
	MAC                *string   `json:"mac"`      // MAC address without :
	Version            *int64    `json:"version"`  // Version
	Temperature        *int64    `json:"temp"`     // Display temperature (upload after multiplying by 10)
	Humidity           *int64    `json:"humi"`     // Display humidity (upload after multiplying by 10)
	SetTemperature     *int64    `json:"settemp"`  // Setting temperature (upload after multiplying by 10)
	Mode               *int64    `json:"mode"`     // Mode (1 cool, 2 heat, 3 vent, 4 auto)
	OnOff              *int64    `json:"onoff"`    // Status 1 on, 2 off)
	FrostTemperature   *int64    `json:"frost"`    // Frost temperature (upload after multiplying by 10)
	Delay              *int64    `json:"delay"`    // Output delay time
	Diff               *int64    `json:"diff"`     // Switch difference (upload after multiplying by 10)
	HoldTemperature    *int64    `json:"holdtemp"` // Hold temperature (upload after multiplying by 10)
	HoldTime           *int64    `json:"holdtime"` // Hold temperature time
	ScreenLock         *int64    `json:"kb"`       // Screen lock (1 no, 2 yes)
	LockPin            *string   `json:"kbkey"`    // Lock pin (four numbers)
	TemperatureFormat  *int64    `json:"cf"`       // Temperature format (0 celsius, 1 fahrenheit)
	Holiday            *int64    `json:"holiday"`  // Holiday (0 no, 1 yes)
	HolidayStartTime   *int64    `json:"holiday_startime"`
	HolidayEndTime     *int64    `json:"holiday_endtime"`
	StandBy            *int64    `json:"standby"`  // Standby (1 no, 2 yes)
	Fan                *int64    `json:"fan"`      // Fan speed (1 high, 2 medium, 3 low, 4 auto)
	Timezone           *int64    `json:"timezone"` // Time zone
	Program            *int64    `json:"prog"`     // Schedule (0 none, 1 weekday/weekend, 2 7 days, 3 24 hours)
	ProgramTemperature *struct{} `json:"temp_prog"`
	OnlineType         *int64    `json:"online_type"`
}
