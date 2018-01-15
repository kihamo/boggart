package main // import "github.com/kihamo/boggart/cmd/boggart"

import (
	"log"

	//_ "github.com/go-sql-driver/mysql"

	boggart "github.com/kihamo/boggart/components/boggart/instance"
	"github.com/kihamo/shadow"
	//alerts "github.com/kihamo/shadow/components/alerts/instance"
	config "github.com/kihamo/shadow/components/config/instance"
	dashboard "github.com/kihamo/shadow/components/dashboard/instance"
	//database "github.com/kihamo/shadow/components/database/instance"
	logger "github.com/kihamo/shadow/components/logger/instance"
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
			//alerts.NewComponent(),
			config.NewComponent(),
			dashboard.NewComponent(),
			//database.NewComponent(),
			logger.NewComponent(),
			metrics.NewComponent(),
			profiling.NewComponent(),
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
