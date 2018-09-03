package main // import "github.com/kihamo/boggart/cmd/roborock"

import (
	"log"

	mqtt "github.com/kihamo/boggart/components/mqtt/instance"
	roborock "github.com/kihamo/boggart/components/roborock/instance"
	voice "github.com/kihamo/boggart/components/voice/instance"
	"github.com/kihamo/shadow"
	config "github.com/kihamo/shadow/components/config/instance"
	dashboard "github.com/kihamo/shadow/components/dashboard/instance"
	i18n "github.com/kihamo/shadow/components/i18n/instance"
	logger "github.com/kihamo/shadow/components/logger/instance"
	metrics "github.com/kihamo/shadow/components/metrics/instance"
	profiling "github.com/kihamo/shadow/components/profiling/instance"
)

var (
	Version = "0.0"
	Build   = "0-0000000"
)

func main() {
	application, err := shadow.NewApp(
		"Boggart Roborock",
		Version,
		Build,
		[]shadow.Component{
			mqtt.NewComponent(),
			config.NewComponent(),
			dashboard.NewComponent(),
			roborock.NewComponent(),
			voice.NewComponent(),
			i18n.NewComponent(),
			logger.NewComponent(),
			metrics.NewComponent(),
			profiling.NewComponent(),
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
