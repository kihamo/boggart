package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/kihamo/boggart/components/openhab"
	"github.com/kihamo/boggart/components/openhab/client/things"
	"github.com/kihamo/shadow/components/dashboard"
)

type ProxyHandler struct {
	dashboard.Handler
}

func (h *ProxyHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()

	switch query.Get(":type") {
	case "thing":
		id := query.Get(":id")
		configKey := query.Get(":key")
		if id == "" || configKey == "" {
			h.NotFound(w, r)
			return
		}

		client := dashboard.ComponentFromContext(r.Context()).(openhab.Component).Client()

		// TODO: cache
		status, err := client.Things.GetByUID(&things.GetByUIDParams{
			ThingUID: id,
			Context:  r.Context(),
		})

		if err != nil {
			break
		}

		configVal, ok := status.Payload.Configuration[configKey]
		if !ok {
			break
		}

		response, err := http.Get(configVal.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer response.Body.Close()

		for headerKey, headerValues := range response.Header {
			for _, headerValue := range headerValues {
				w.Header().Add(headerKey, headerValue)
			}
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(response.StatusCode)
		w.Write(body)
	}

	h.NotFound(w, r)
}
