package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/snitch"
)

type IndexHandler struct {
	dashboard.Handler
}

type MetricValue struct {
	Value float64
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
		r.Component().(boggart.Component).Metrics().Collect(metricsChan)
		close(metricsChan)
	}()

	for metric := range metricsChan {
		name := metric.Description().Name()

		if value, err := metric.Measure(); err != nil {
			errors = append(errors, "Get metric "+name+" failed: "+err.Error())
		} else {
			v := MetricValue{}

			if value.Value != nil {
				v.Value = *(value.Value)
			}

			vars[name] = v
		}
	}

	vars["errors"] = errors

	h.Render(r.Context(), "index", vars)
}
