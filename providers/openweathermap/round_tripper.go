package openweathermap

import (
	"net/http"
)

type RoundTripper struct {
	original http.RoundTripper
	limiter  *Limiter
}

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := rt.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return rt.original.RoundTrip(req)
}
