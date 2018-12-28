package home

type eurekaInfoOption string
type eurekaInfoParam string

const (
	EurekaInfoOptionDetail = eurekaInfoOption("detail")

	EurekaInfoParamVersion         = eurekaInfoParam("version")
	EurekaInfoParamAudio           = eurekaInfoParam("audio")
	EurekaInfoParamName            = eurekaInfoParam("name")
	EurekaInfoParamBuildInfo       = eurekaInfoParam("build_info")
	EurekaInfoParamDetail          = eurekaInfoParam("detail")
	EurekaInfoParamDeviceInfo      = eurekaInfoParam("device_info")
	EurekaInfoParamNet             = eurekaInfoParam("net")
	EurekaInfoParamWifi            = eurekaInfoParam("wifi")
	EurekaInfoParamSetup           = eurekaInfoParam("setup")
	EurekaInfoParamSettings        = eurekaInfoParam("settings")
	EurekaInfoParamOptIn           = eurekaInfoParam("opt_in")
	EurekaInfoParamOpencast        = eurekaInfoParam("opencast")
	EurekaInfoParamMultizone       = eurekaInfoParam("multizone")
	EurekaInfoParamProxy           = eurekaInfoParam("proxy")
	EurekaInfoParamNightModeParams = eurekaInfoParam("night_mode_params")
	EurekaInfoParamUserEq          = eurekaInfoParam("user_eq")
	EurekaInfoParamRoomQualizer    = eurekaInfoParam("room_equalizer")
	EurekaInfoParamAll             = eurekaInfoParam(EurekaInfoParamVersion + "," + EurekaInfoParamAudio +
		"," + EurekaInfoParamName + "," + EurekaInfoParamBuildInfo + "," + EurekaInfoParamDetail +
		"," + EurekaInfoParamDeviceInfo + "," + EurekaInfoParamNet + "," + EurekaInfoParamWifi +
		"," + EurekaInfoParamSetup + "," + EurekaInfoParamSettings + "," + EurekaInfoParamOptIn +
		"," + EurekaInfoParamOpencast + "," + EurekaInfoParamMultizone + "," + EurekaInfoParamProxy +
		"," + EurekaInfoParamNightModeParams + "," + EurekaInfoParamUserEq + "," + EurekaInfoParamRoomQualizer)
)

func (o eurekaInfoOption) Value() *string {
	return &[]string{string(o)}[0]
}

func (v eurekaInfoParam) Value() *string {
	return &[]string{string(v)}[0]
}
