package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/snitch"
)

type IndexHandler struct {
	dashboard.Handler

	Config    config.Component
	Collector metrics.HasMetrics
}

func (h *IndexHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	// TODO: update metric value
	if r.IsAjax() {

		return
	}

	errors := []string{}
	vars := map[string]interface{}{}

	metricsChan := make(chan snitch.Metric, 10000)

	go func() {
		h.Collector.Metrics().Collect(metricsChan)
		close(metricsChan)
	}()

	for metric := range metricsChan {
		name := metric.Description().Name()

		if value, err := metric.Measure(); err != nil {
			errors = append(errors, "Get metric "+name+" failed: "+err.Error())
		} else {
			vars[name] = value
		}
	}

	vars["errors"] = errors

	h.Render(r.Context(), boggart.ComponentName, "index", vars)
}
