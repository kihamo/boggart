package handlers

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kihamo/boggart/components/barcode/codes"
	"github.com/kihamo/shadow/components/dashboard"
)

type DecodeHandler struct {
	dashboard.Handler
}

type response struct {
	Result string `json:"result"`
	Code   string `json:"code,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (h *DecodeHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	byURL := q.Get("url")
	byBase64 := q.Get("base64")

	if byURL == "" && byBase64 == "" {
		h.NotFound(w, r)
		return
	}

	var (
		err    error
		reader io.Reader
		debug  bool
	)

	if q.Get("debug") != "" {
		debug = true
	}

	if byURL != "" {
		var r *http.Response

		r, err = http.Get(byURL)
		if err == nil {
			defer r.Body.Close()
			reader = r.Body
		}
	} else if byBase64 != "" {
		reader, err = codes.FromBase64([]byte(byBase64))
	}

	var (
		debugImage io.Reader
		code       string
	)

	if err == nil {
		if debug {
			code, debugImage, err = codes.DecodeQRCodeDebug(reader)
		} else {
			code, err = codes.DecodeQRCode(reader)
		}
	}

	if err != nil {
		w.SendJSON(response{
			Result: "failed",
			Error:  err.Error(),
		})

		return
	}

	if debugImage != nil {
		body, _ := ioutil.ReadAll(debugImage)
		w.Write(body)

		return
	}

	w.SendJSON(response{
		Result: "success",
		Code:   code,
	})
}
