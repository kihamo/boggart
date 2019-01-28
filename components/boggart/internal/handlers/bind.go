package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
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

	manager *manager.Manager
}

func NewBindHandler(manager *manager.Manager) *BindHandler {
	return &BindHandler{
		manager: manager,
	}
}

func (h *BindHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	var bindItem boggart.BindItem

	id := q.Get(":id")
	if id != "" {
		bindItem = h.manager.Bind(id)
		if bindItem == nil {
			h.NotFound(w, r)
			return
		}
	}

	switch q.Get(":action") {
	case "unregister":
		h.actionDelete(w, r, bindItem)
		return

	case "widget":
		h.actionWidget(w, r, bindItem)
		return

	case "":
		h.actionCreateOrUpdate(w, r, bindItem)
		return

	default:
		h.NotFound(w, r)
		return
	}
}

func (h *BindHandler) registerByYAML(code []byte) (bindItem boggart.BindItem, upgraded bool, err error) {
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

	if bindParsed.ID != "" {
		if bindExists := h.manager.Bind(bindParsed.ID); bindExists != nil {
			upgraded = true

			if err := h.manager.Unregister(bindParsed.ID); err != nil {
				return nil, false, err
			}
		}

		bindItem, err = h.manager.RegisterWithID(bindParsed.ID, bind, bindParsed.Type, bindParsed.Description, bindParsed.Tags, cfg)
	} else {
		bindItem, err = h.manager.Register(bind, bindParsed.Type, bindParsed.Description, bindParsed.Tags, cfg)
	}

	return bindItem, upgraded, err
}

func (h *BindHandler) actionCreateOrUpdate(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	buf := bytes.NewBuffer(nil)

	var (
		err     error
		message string
	)

	if r.IsPost() {
		code := r.Original().FormValue("yaml")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		buf.WriteString(code)

		var (
			bind     boggart.BindItem
			upgraded bool
		)

		if bind, upgraded, err = h.registerByYAML(buf.Bytes()); err == nil {
			if upgraded {
				message = "Bind " + bind.ID() + " upgraded"
			} else {
				message = "Bind register success with id " + bind.ID()
			}
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
		"yaml":    buf.String(),
		"error":   err,
		"message": message,
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

	err := h.manager.Unregister(b.ID())

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

func (h *BindHandler) actionWidget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	if b == nil {
		h.NotFound(w, r)
		return
	}

	widget, ok := b.BindType().(boggart.BindTypeHasWidget)
	if !ok {
		h.NotFound(w, r)
		return
	}

	componentName := boggart.ComponentName + "/widget/" + b.Type()

	r = r.WithContext(context.WithValue(r.Context(), dashboard.ComponentContextKey, componentName))

	if render := dashboard.RenderFromContext(r.Context()); render != nil {
		fs := widget.WidgetTemplates()
		if fs != nil {
			if !render.IsRegisterComponent(componentName) {
				err := render.RegisterComponent(componentName, fs)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	widget.Widget(w, r, b)
}
