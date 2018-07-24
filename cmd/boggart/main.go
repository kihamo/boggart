package main // import "github.com/kihamo/boggart/cmd/boggart"

import (
	"log"

	//_ "github.com/go-sql-driver/mysql"

	boggart "github.com/kihamo/boggart/components/boggart/instance"
	openhab "github.com/kihamo/boggart/components/openhab/instance"
	"github.com/kihamo/shadow"
	annotations "github.com/kihamo/shadow/components/annotations/instance"
	config "github.com/kihamo/shadow/components/config/instance"
	dashboard "github.com/kihamo/shadow/components/dashboard/instance"
	//database "github.com/kihamo/shadow/components/database/instance"
	syslog "github.com/kihamo/boggart/components/syslog/instance"
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
		"Boggart",
		Version,
		Build,
		[]shadow.Component{
			boggart.NewComponent(),
			openhab.NewComponent(),
			annotations.NewComponent(),
			config.NewComponent(),
			dashboard.NewComponent(),
			//database.NewComponent(),
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
