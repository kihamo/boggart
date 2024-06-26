package main // import "github.com/kihamo/boggart/cmd/roborock"

import (
	"log"

	_ "github.com/kihamo/boggart/components/boggart/bind/alsa"
	_ "github.com/kihamo/boggart/components/boggart/bind/xiaomi"
	_ "github.com/kihamo/boggart/components/boggart/instance"
	_ "github.com/kihamo/boggart/components/mqtt/instance"
	"github.com/kihamo/shadow"
	_ "github.com/kihamo/shadow/components/config/instance"
	_ "github.com/kihamo/shadow/components/dashboard/instance"
	_ "github.com/kihamo/shadow/components/i18n/instance"
	_ "github.com/kihamo/shadow/components/logging/instance"
	_ "github.com/kihamo/shadow/components/messengers/instance"
	_ "github.com/kihamo/shadow/components/metrics/instance"
	_ "github.com/kihamo/shadow/components/profiling/instance"
)

var (
	Version = "0.0"
	Build   = "0-0000000"
)

func main() {
	shadow.SetName("Boggart Roborock")
	shadow.SetVersion(Version)
	shadow.SetBuild(Build)

	if err := shadow.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
