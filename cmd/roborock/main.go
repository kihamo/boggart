package main // import "github.com/kihamo/boggart/cmd/roborock"

import (
	"log"

	_ "github.com/kihamo/boggart/components/mqtt/instance"
	_ "github.com/kihamo/boggart/components/roborock/instance"
	_ "github.com/kihamo/boggart/components/voice/instance"
	"github.com/kihamo/shadow"
	_ "github.com/kihamo/shadow/components/config/instance"
	_ "github.com/kihamo/shadow/components/dashboard/instance"
	_ "github.com/kihamo/shadow/components/i18n/instance"
	_ "github.com/kihamo/shadow/components/logger/instance"
	_ "github.com/kihamo/shadow/components/metrics/instance"
	_ "github.com/kihamo/shadow/components/profiling/instance"
	_ "github.com/kihamo/shadow/components/tracing/instance"
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
