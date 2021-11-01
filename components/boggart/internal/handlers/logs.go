package handlers

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

var (
	crlf = []byte("\r\n")
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
	q := r.URL().Query()
	action := q.Get(":action")
	id := q.Get(":id")

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

	switch action {
	case "download":
		buf := bytes.NewBuffer(nil)

		buf.WriteString("# bind: ")
		buf.WriteString(id)
		buf.Write(crlf)
		buf.WriteString("# date: ")
		buf.WriteString(time.Now().Format(time.RFC3339))
		buf.Write(crlf)
		buf.Write(crlf)

		var total int

		for _, record := range bindSupport.LastRecords() {
			buf.WriteString(record.Time.Format(time.RFC3339))
			buf.WriteString(" [")
			buf.WriteString(record.Level.String())
			buf.WriteString("] ")
			buf.WriteString(record.Message)
			buf.WriteString(" {")

			total = len(record.Context) - 1
			for i, v := range record.Context {
				buf.WriteByte('"')
				buf.WriteString(v.Key)
				buf.WriteString("\":\"")
				buf.WriteString(v.String)
				buf.WriteByte('"')

				if i != total {
					buf.WriteString(", ")
				}
			}

			buf.WriteByte('}')
			buf.Write(crlf)
		}

		w.Header().Set("Content-Disposition", "attachment; filename=\""+time.Now().Format("20060102150405_log.txt")+"\"")
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		if _, err := io.Copy(w, buf); err != nil {
			panic(err.Error())
		}

		return

	case "clean":
		if r.IsPost() {
			bindSupport.Clean()

			r.Session().FlashBag().Success("Logs cleaned")

			h.Redirect("/"+h.component.Name()+"/logs/"+id+"/", http.StatusFound, w, r)
		} else {
			h.MethodNotAllowed(w, r)
		}

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

	for i, record := range records {
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
