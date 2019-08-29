package main // import "github.com/kihamo/boggart/cmd/agent"

import (
	"log"

	_ "github.com/kihamo/boggart/components/boggart/bind/ds18b20"
	_ "github.com/kihamo/boggart/components/boggart/bind/gpio"
	_ "github.com/kihamo/boggart/components/boggart/bind/hikvision"
	_ "github.com/kihamo/boggart/components/boggart/bind/hilink"
	_ "github.com/kihamo/boggart/components/boggart/bind/homie"
	_ "github.com/kihamo/boggart/components/boggart/bind/mercury"
	_ "github.com/kihamo/boggart/components/boggart/bind/mikrotik"
	_ "github.com/kihamo/boggart/components/boggart/bind/network"
	_ "github.com/kihamo/boggart/components/boggart/bind/nut"
	_ "github.com/kihamo/boggart/components/boggart/bind/serial"
	_ "github.com/kihamo/boggart/components/boggart/bind/wol"
	_ "github.com/kihamo/boggart/components/boggart/bind/xmeye"
	_ "github.com/kihamo/boggart/components/boggart/instance"
	_ "github.com/kihamo/boggart/components/mqtt/instance"
	"github.com/kihamo/shadow"
	_ "github.com/kihamo/shadow/components/config/instance"
	_ "github.com/kihamo/shadow/components/dashboard/instance"
	_ "github.com/kihamo/shadow/components/i18n/instance"
	_ "github.com/kihamo/shadow/components/logging/instance"
	_ "github.com/kihamo/shadow/components/metrics/instance"
	_ "github.com/kihamo/shadow/components/profiling/instance"
	_ "github.com/kihamo/shadow/components/tracing/instance"
	_ "github.com/kihamo/shadow/components/workers/instance"
)

var (
	Version = "0.0"
	Build   = "0-0000000"
)

func main() {
	shadow.SetName("Boggart Agent")
	shadow.SetVersion(Version)
	shadow.SetBuild(Build)

	if err := shadow.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
