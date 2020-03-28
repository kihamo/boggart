package home

import (
	"strconv"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/google/home/client"
)

type eurekaInfoOption string
type eurekaInfoParam string

const (
	DefaultPort = 8008

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
	EurekaInfoParamAll             = EurekaInfoParamVersion + "," + EurekaInfoParamAudio +
		"," + EurekaInfoParamName + "," + EurekaInfoParamBuildInfo + "," + EurekaInfoParamDetail +
		"," + EurekaInfoParamDeviceInfo + "," + EurekaInfoParamNet + "," + EurekaInfoParamWifi +
		"," + EurekaInfoParamSetup + "," + EurekaInfoParamSettings + "," + EurekaInfoParamOptIn +
		"," + EurekaInfoParamOpencast + "," + EurekaInfoParamMultizone + "," + EurekaInfoParamProxy +
		"," + EurekaInfoParamNightModeParams + "," + EurekaInfoParamUserEq + "," + EurekaInfoParamRoomQualizer
)

func (o eurekaInfoOption) Value() *string {
	return &[]string{string(o)}[0]
}

func (v eurekaInfoParam) Value() *string {
	return &[]string{string(v)}[0]
}

type Client struct {
	*client.GoogleHome
}

func New(address string, debug bool, logger logger.Logger) *Client {
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		address = strings.Join(append(parts, strconv.Itoa(DefaultPort)), ":")
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		GoogleHome: cl,
	}
}
