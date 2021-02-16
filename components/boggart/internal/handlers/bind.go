package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"sort"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow/components/dashboard"
	"gopkg.in/yaml.v2"
)

type BindYAML struct {
	Type        string
	ID          string
	Description string
	Tags        []string
	Config      interface{}
}

type BindHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewBindHandler(b boggart.Component) *BindHandler {
	return &BindHandler{
		component: b,
	}
}

func (h *BindHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	var bindItem boggart.BindItem

	id := q.Get(":id")
	if id != "" {
		bindItem = h.component.Bind(id)
		if bindItem == nil {
			h.NotFound(w, r)
			return
		}
	}

	switch q.Get(":action") {
	case "unregister":
		if bindItem.Type() == boggart.ComponentName {
			h.NotFound(w, r)
			return
		}

		h.actionDelete(w, r, bindItem)
		return

	case "":
		if bindItem != nil && bindItem.Type() == boggart.ComponentName {
			h.NotFound(w, r)
			return
		}

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

	cfg, md, err := boggart.ValidateBindConfig(kind, bindParsed.Config)
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
		if bindExists := h.component.Bind(id); bindExists != nil {
			upgraded = true

			if err := h.component.UnregisterBindByID(id); err != nil {
				return nil, false, err
			}
		}
	}

	bindItem, err = h.component.RegisterBind(bindParsed.ID, bind, bindParsed.Type, bindParsed.Description, bindParsed.Tags, cfg)

	if len(md.Unused) > 0 {
		if logger, ok := di.LoggerContainerBind(bind); ok {
			for _, field := range md.Unused {
				logger.Warn("Unused config field", "field", field)
			}
		}
	}

	return bindItem, upgraded, err
}

func (h *BindHandler) actionCreateOrUpdate(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	var err error

	buf := bytes.NewBuffer(nil)
	vars := make(map[string]interface{})

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

			redirectURL := &url.URL{}
			*redirectURL = *r.URL()
			redirectURL.Path = "/" + h.component.Name() + "/manager/"
			redirectURL.RawQuery = "search=" + url.QueryEscape(bind.ID())

			h.Redirect(redirectURL.String(), http.StatusFound, w, r)

			return
		}
	} else {
		enc := yaml.NewEncoder(buf)
		defer enc.Close()

		if b == nil {
			bindYAML := &BindYAML{
				Description: "Description of new bind",
				Tags:        []string{"tag_label"},
				Config: map[string]interface{}{
					"config_key": "config_value",
				},
			}
			isAjax := false

			if typeName := r.URL().Query().Get("type"); typeName != "" && r.IsAjax() {
				bindYAML.Type = typeName
				isAjax = true
			} else {
				types := make([]string, 0)
				for typeName := range boggart.GetBindTypes() {
					types = append(types, typeName)
				}
				sort.Strings(types)
				vars["types"] = types

				if len(types) > 0 {
					bindYAML.Type = types[0]
				}
			}

			if bindYAML.Type != "" {
				if t, err := boggart.GetBindType(bindYAML.Type); err == nil {
					bindYAML.Config = t.ConfigDefaults()
				} else {
					bindYAML.Type = ""
				}
			}

			err = enc.Encode(bindYAML)

			if isAjax {
				if err == nil {
					_ = w.SendJSON(buf.String())
					return
				}

				h.InternalError(w, r, err)
				return
			}
		} else {
			err = enc.Encode(b)
		}
	}

	vars["yaml"] = buf.String()

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

	err := h.component.UnregisterBindByID(b.ID())
	if err != nil {
		_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))

		return
	}

	_ = w.SendJSON(boggart.NewResponseJSON().Success(""))
}
