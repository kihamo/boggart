package http

import (
	"encoding/json"
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

func JsonUnmarshal(r *http.Response, v interface{}) error {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
