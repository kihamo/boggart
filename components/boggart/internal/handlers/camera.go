package handlers

import (
	"bytes"
	"io"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/bind/hikvision"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/shadow/components/dashboard"
)

type CameraHandler struct {
	dashboard.Handler

	manager *manager.Manager
}

func NewCameraHandler(manager *manager.Manager) *CameraHandler {
	return &CameraHandler{
		manager: manager,
	}
}

func (h *CameraHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	id := q.Get(":id")
	channel := q.Get(":channel")

	if id == "" || channel == "" {
		h.NotFound(w, r)
		return
	}

	ch, err := strconv.ParseUint(channel, 10, 64)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	bindItem := h.manager.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	bind, ok := bindItem.Bind().(*hikvision.Bind)
	if !ok {
		h.NotFound(w, r)
		return
	}

	buf := bytes.NewBuffer(nil)
	if err = bind.Snapshot(r.Context(), ch, buf); err != nil {
		h.NotFound(w, r)
		return
	}

	io.Copy(w, buf)
}
