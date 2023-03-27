package main // import "github.com/kihamo/boggart/cmd/server"

import (
	"log"

	_ "github.com/kihamo/boggart/components/boggart/bind/apcupsd"
	_ "github.com/kihamo/boggart/components/boggart/bind/aqicn"
	_ "github.com/kihamo/boggart/components/boggart/bind/astro"
	_ "github.com/kihamo/boggart/components/boggart/bind/broadlink"
	_ "github.com/kihamo/boggart/components/boggart/bind/chromecast"
	_ "github.com/kihamo/boggart/components/boggart/bind/ds18b20"
	_ "github.com/kihamo/boggart/components/boggart/bind/esphome"
	_ "github.com/kihamo/boggart/components/boggart/bind/fcm"
	_ "github.com/kihamo/boggart/components/boggart/bind/google_home"
	_ "github.com/kihamo/boggart/components/boggart/bind/gpio"
	_ "github.com/kihamo/boggart/components/boggart/bind/grafana"
	_ "github.com/kihamo/boggart/components/boggart/bind/herospeed"
	_ "github.com/kihamo/boggart/components/boggart/bind/hikvision"
	_ "github.com/kihamo/boggart/components/boggart/bind/hilink"
	_ "github.com/kihamo/boggart/components/boggart/bind/homie"
	_ "github.com/kihamo/boggart/components/boggart/bind/influxdb"
	_ "github.com/kihamo/boggart/components/boggart/bind/integratorit"
	_ "github.com/kihamo/boggart/components/boggart/bind/led_wifi"
	_ "github.com/kihamo/boggart/components/boggart/bind/lg_webos"
	_ "github.com/kihamo/boggart/components/boggart/bind/mail"
	_ "github.com/kihamo/boggart/components/boggart/bind/mc6"
	_ "github.com/kihamo/boggart/components/boggart/bind/mercury"
	_ "github.com/kihamo/boggart/components/boggart/bind/mikrotik"
	_ "github.com/kihamo/boggart/components/boggart/bind/myheat"
	_ "github.com/kihamo/boggart/components/boggart/bind/neptun"
	_ "github.com/kihamo/boggart/components/boggart/bind/network"
	_ "github.com/kihamo/boggart/components/boggart/bind/nut"
	_ "github.com/kihamo/boggart/components/boggart/bind/octoprint"
	_ "github.com/kihamo/boggart/components/boggart/bind/openhab"
	_ "github.com/kihamo/boggart/components/boggart/bind/openweathermap"
	_ "github.com/kihamo/boggart/components/boggart/bind/owntracks"
	_ "github.com/kihamo/boggart/components/boggart/bind/pantum"
	_ "github.com/kihamo/boggart/components/boggart/bind/pass24online"
	_ "github.com/kihamo/boggart/components/boggart/bind/premiergc"
	_ "github.com/kihamo/boggart/components/boggart/bind/pulsar"
	_ "github.com/kihamo/boggart/components/boggart/bind/rpi"
	_ "github.com/kihamo/boggart/components/boggart/bind/samsung_tizen"
	_ "github.com/kihamo/boggart/components/boggart/bind/serial"
	_ "github.com/kihamo/boggart/components/boggart/bind/smcenter"
	_ "github.com/kihamo/boggart/components/boggart/bind/softvideo"
	_ "github.com/kihamo/boggart/components/boggart/bind/syslog"
	_ "github.com/kihamo/boggart/components/boggart/bind/telegram"
	_ "github.com/kihamo/boggart/components/boggart/bind/text2speech"
	_ "github.com/kihamo/boggart/components/boggart/bind/timelapse"
	_ "github.com/kihamo/boggart/components/boggart/bind/tvt"
	_ "github.com/kihamo/boggart/components/boggart/bind/wol"
	_ "github.com/kihamo/boggart/components/boggart/bind/xiaomi"
	_ "github.com/kihamo/boggart/components/boggart/bind/xmeye"
	_ "github.com/kihamo/boggart/components/boggart/bind/zigbee/z_stack"
	_ "github.com/kihamo/boggart/components/boggart/bind/zigbee/zigbee2mqtt"
	_ "github.com/kihamo/boggart/components/boggart/instance"
	_ "github.com/kihamo/boggart/components/mqtt/instance"
	_ "github.com/kihamo/boggart/components/storage/instance"
	"github.com/kihamo/shadow"
	_ "github.com/kihamo/shadow/components/config/instance"
	_ "github.com/kihamo/shadow/components/dashboard/instance"
	_ "github.com/kihamo/shadow/components/i18n/instance"
	_ "github.com/kihamo/shadow/components/logging/instance"
	_ "github.com/kihamo/shadow/components/metrics/instance"
	_ "github.com/kihamo/shadow/components/profiling/instance"
)

var (
	Name    = "Boggart Server"
	Version = "0.0"
	Build   = "0-0000000"
)

func main() {
	shadow.SetName(Name)
	shadow.SetVersion(Version)
	shadow.SetBuild(Build)

	if err := shadow.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
