package handlers

import (
	"bytes"

	"gopkg.in/yaml.v2"

	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/shadow/components/dashboard"
)

type BindHandler struct {
	dashboard.Handler

	manager *manager.Manager
}

func NewBindHandler(manager *manager.Manager) *BindHandler {
	return &BindHandler{
		manager: manager,
	}
}

func (h *BindHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	id := q.Get(":id")
	action := q.Get(":action")

	if id == "" && action == "" {
		h.actionRegister(w, r)
		return
	}

	if action == "" {
		h.NotFound(w, r)
		return
	}

	bind := h.manager.Bind(id)
	if bind == nil {
		h.NotFound(w, r)
		return
	}

	switch action {
	case "unregister":
		h.actionUnregister(id, w, r)
		return
	}

	h.NotFound(w, r)
}

func (h *BindHandler) actionUnregister(id string, w *dashboard.Response, r *dashboard.Request) {
	err := h.manager.Unregister(id)

	type response struct {
		Result  string `json:"result"`
		Message string `json:"message,omitempty"`
	}

	if err != nil {
		w.SendJSON(response{
			Result:  "failed",
			Message: err.Error(),
		})
		return
	}

	w.SendJSON(response{
		Result: "success",
	})
}

func (h *BindHandler) actionRegister(w *dashboard.Response, r *dashboard.Request) {
	var device manager.BindItem

	if r.IsPost() {

	} else {

	}

	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)
	defer enc.Close()

	if err := enc.Encode(device); err != nil {
		panic(err.Error())
	}

	h.Render(r.Context(), "bind", map[string]interface{}{
		"yaml": buf.String(),
	})
}
