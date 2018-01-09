package handlers

import (
	"fmt"
	"regexp"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/pulsar"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

var hexCleaner = regexp.MustCompile("[^a-zA-Z0-9]+")

type IndexHandler struct {
	dashboard.Handler

	Config config.Component
}

func (h *IndexHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	vars := map[string]interface{}{}

	client, err := pulsar.NewPulsar(
		h.Config.GetString(boggart.ConfigPulsarSerialPath),
		[]byte(hexCleaner.ReplaceAllString(h.Config.GetString(boggart.ConfigPulsarAddress), "")),
	)

	if err != nil {
		vars["error"] = err.Error()
	}

	fmt.Println(client.ReadTime())

	h.Render(r.Context(), boggart.ComponentName, "index", vars)
}
