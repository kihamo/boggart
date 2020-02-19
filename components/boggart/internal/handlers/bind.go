package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

type BindYAML struct {
	Type        string
	ID          string
	Description string
	Tags        []string
	Config      map[string]interface{}
}

type BindHandler struct {
	dashboard.Handler

	componentBoggart boggart.Component
	componentMQTT    mqtt.Component
}

func NewBindHandler(b boggart.Component, m mqtt.Component) *BindHandler {
	return &BindHandler{
		componentBoggart: b,
		componentMQTT:    m,
	}
}

func (h *BindHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	var bindItem boggart.BindItem

	id := q.Get(":id")
	if id != "" {
		bindItem = h.componentBoggart.Bind(id)
		if bindItem == nil {
			h.NotFound(w, r)
			return
		}
	}

	switch q.Get(":action") {
	case "unregister":
		h.actionDelete(w, r, bindItem)
		return

	case "readiness":
		h.actionProbe(w, r, bindItem, "readiness")
		return

	case "liveness":
		h.actionProbe(w, r, bindItem, "liveness")
		return

	case "logs":
		h.actionLogs(w, r, bindItem)
		return

	case "mqtt":
		h.actionMQTT(w, r, bindItem)
		return

	case "":
		h.actionCreateOrUpdate(w, r, bindItem)
		return

	default:
		h.NotFound(w, r)
		return
	}
}

func (h *BindHandler) registerByYAML(oldID string, code []byte) (bindItem boggart.BindItem, upgraded bool, err error) {
	bindParsed := &BindYAML{}

	err = yaml.Unmarshal(code, bindParsed)
	if err != nil {
		return nil, false, err
	}

	if bindParsed.Type == "" {
		return nil, false, errors.New("bind type is empty")
	}

	kind, err := boggart.GetBindType(bindParsed.Type)
	if err != nil {
		return nil, false, err
	}

	cfg, err := boggart.ValidateBindConfig(kind, bindParsed.Config)
	if err != nil {
		return nil, false, err
	}

	bind, err := kind.CreateBind(cfg)
	if err != nil {
		return nil, false, err
	}

	removeIDs := make([]string, 0, 2)

	// check new ID
	if bindParsed.ID != "" {
		removeIDs = append(removeIDs, bindParsed.ID)
	}

	// check old ID
	if oldID != "" && oldID != bindParsed.ID {
		removeIDs = append(removeIDs, oldID)
	}

	for _, id := range removeIDs {
		if bindExists := h.componentBoggart.Bind(id); bindExists != nil {
			upgraded = true

			if err := h.componentBoggart.UnregisterBindByID(id); err != nil {
				return nil, false, err
			}
		}
	}

	bindItem, err = h.componentBoggart.RegisterBind(bindParsed.ID, bind, bindParsed.Type, bindParsed.Description, bindParsed.Tags, cfg)

	return bindItem, upgraded, err
}

func (h *BindHandler) actionCreateOrUpdate(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	buf := bytes.NewBuffer(nil)

	var err error

	if r.IsPost() {
		code := r.Original().FormValue("yaml")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		buf.WriteString(code)

		var (
			id       string
			bind     boggart.BindItem
			upgraded bool
		)

		if b != nil {
			id = b.ID()
		}

		if bind, upgraded, err = h.registerByYAML(id, buf.Bytes()); err == nil {
			if upgraded {
				r.Session().FlashBag().Info("Bind " + bind.ID() + " upgraded")
			} else {
				r.Session().FlashBag().Success("Bind register success with id " + bind.ID())
			}

			h.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	} else {
		enc := yaml.NewEncoder(buf)
		defer enc.Close()

		if b == nil {
			err = enc.Encode(&BindYAML{
				Description: "Description of new bind",
				Tags:        []string{"tag_label"},
				Config: map[string]interface{}{
					"config_key": "config_value",
				},
			})
		} else {
			err = enc.Encode(b)
		}
	}

	vars := map[string]interface{}{
		"yaml": buf.String(),
	}

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	if b != nil {
		vars["bindId"] = b.ID()
	}

	h.Render(r.Context(), "bind", vars)
}

func (h *BindHandler) actionDelete(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	if b == nil {
		h.NotFound(w, r)
		return
	}

	err := h.componentBoggart.UnregisterBindByID(b.ID())

	type response struct {
		Result  string `json:"result"`
		Message string `json:"message,omitempty"`
	}

	if err != nil {
		_ = w.SendJSON(response{
			Result:  "failed",
			Message: err.Error(),
		})
		return
	}

	_ = w.SendJSON(response{
		Result: "success",
	})
}

func (h *BindHandler) actionProbe(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem, t string) {
	if b == nil {
		h.NotFound(w, r)
		return
	}

	bindSupport, ok := di.ProbesContainerBind(b.Bind())

	if !ok {
		h.NotFound(w, r)
		return
	}

	var err error

	switch t {
	case "readiness":
		err = bindSupport.ReadinessCheck(r.Context())
	case "liveness":
		err = bindSupport.LivenessCheck(r.Context())
	}

	response := struct {
		Result string `json:"result"`
		Error  string `json:"error,omitempty"`
	}{
		Result: "success",
	}

	if err != nil {
		response.Error = err.Error()
	}

	_ = w.SendJSON(response)
}

func (h *BindHandler) actionLogs(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bindSupport, ok := b.Bind().(di.LoggerContainerSupport)
	if !ok {
		h.NotFound(w, r)
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
		"bind": b,
		"logs": response,
	})
}

func (h *BindHandler) actionMQTT(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bindSupport, ok := b.Bind().(di.MQTTContainerSupport)
	if !ok {
		h.NotFound(w, r)
		return
	}

	type itemView struct {
		Topic      string
		CacheTopic string
		Datetime   time.Time
		Payload    interface{}
	}

	items := make([]itemView, 0)
	publishes := bindSupport.MQTT().Publishes()

	for _, item := range h.componentMQTT.CacheItems() {
		for _, publish := range publishes {
			if publish.IsInclude(item.Topic()) {
				items = append(items, itemView{
					Topic:      publish.String(),
					CacheTopic: item.Topic().String(),
					Datetime:   item.Datetime(),
					Payload:    item.Payload(),
				})

				break
			}
		}
	}

	h.Render(r.Context(), "mqtt", map[string]interface{}{
		"bind":  b,
		"items": items,
	})
}
