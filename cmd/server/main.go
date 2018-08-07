package main // import "github.com/kihamo/boggart/cmd/server"

import (
	"log"

	boggart "github.com/kihamo/boggart/components/boggart/instance"
	mqtt "github.com/kihamo/boggart/components/mqtt/instance"
	openhab "github.com/kihamo/boggart/components/openhab/instance"
	syslog "github.com/kihamo/boggart/components/syslog/instance"
	"github.com/kihamo/shadow"
	annotations "github.com/kihamo/shadow/components/annotations/instance"
	config "github.com/kihamo/shadow/components/config/instance"
	dashboard "github.com/kihamo/shadow/components/dashboard/instance"
	i18n "github.com/kihamo/shadow/components/i18n/instance"
	logger "github.com/kihamo/shadow/components/logger/instance"
	messengers "github.com/kihamo/shadow/components/messengers/instance"
	metrics "github.com/kihamo/shadow/components/metrics/instance"
	profiling "github.com/kihamo/shadow/components/profiling/instance"
	workers "github.com/kihamo/shadow/components/workers/instance"
)

var (
	Version = "0.0"
	Build   = "0-0000000"
)

func main() {
	application, err := shadow.NewApp(
		"Boggart Server",
		Version,
		Build,
		[]shadow.Component{
			mqtt.NewComponent(),
			boggart.NewComponent(),
			openhab.NewComponent(),
			annotations.NewComponent(),
			config.NewComponent(),
			dashboard.NewComponent(),
			i18n.NewComponent(),
			logger.NewComponent(),
			messengers.NewComponent(),
			metrics.NewComponent(),
			profiling.NewComponent(),
			syslog.NewComponent(),
			workers.NewComponent(),
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = application.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
