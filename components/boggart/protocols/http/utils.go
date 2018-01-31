package http

import (
	"io/ioutil"
	"net/http"
)

func BodyFromResponse(r *http.Response) string {
	defer r.Body.Close()

	if body, err := ioutil.ReadAll(r.Body); err == nil {
		return string(body)
	}

	return ""
}
