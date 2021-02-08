package handlers

import (
	"bytes"
	"net/http"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

type LogsHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewLogsHandler(b boggart.Component) *LogsHandler {
	return &LogsHandler{
		component: b,
	}
}

func (h *LogsHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	id := r.URL().Query().Get(":id")

	if id == "" {
		h.NotFound(w, r)
		return
	}

	bindItem := h.component.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	bindSupport, ok := bindItem.Bind().(di.LoggerContainerSupport)
	if !ok {
		h.NotFound(w, r)
		return
	}

	if r.IsPost() && r.URL().Query().Get("clean") == "1" {
		bindSupport.Clean()

		r.Session().FlashBag().Success("Logs cleaned")

		h.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	type logView struct {
		Level   string
		Time    time.Time
		Message string
		Context string
	}

	records := bindSupport.LastRecords()
	response := make([]logView, len(records))

	for i, record := range bindSupport.LastRecords() {
		response[i].Level = record.Level.String()
		response[i].Time = record.Time
		response[i].Message = record.Message

		val := record.ContextMap()
		if len(val) == 0 {
			continue
		}

		buf := bytes.NewBuffer(nil)
		enc := yaml.NewEncoder(buf)

		err := enc.Encode(val)
		if err != nil {
			enc.Close()

			h.InternalError(w, r, err)

			return
		}

		enc.Close()

		response[i].Context = buf.String()
	}

	h.Render(r.Context(), "logs", map[string]interface{}{
		"bind": bindItem,
		"logs": response,
	})
}
