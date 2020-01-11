package mqtt

import (
	"net/http"
	"net/url"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	q := r.URL().Query()

	switch q.Get("action") {
	case "command":
		t.handleCommand(w, r, bind)

	default:
		t.handleIndex(w, r, bind)
	}
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (t Type) handleCommand(w *dashboard.Response, r *dashboard.Request, bind *Bind) {
	q := r.URL().Query()

	componentID := q.Get("component")
	command := q.Get("cmd")

	if componentID == "" || command == "" {
		t.NotFound(w, r)
		return
	}

	component := bind.Component(componentID)
	if component == nil || component.GetCommandTopic() == "" {
		t.NotFound(w, r)
		return
	}

	ctx := r.Context()
	err := bind.MQTTPublishRawWithoutCache(ctx, component.GetCommandTopic(), 1, false, command)

	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	} else {
		r.Session().FlashBag().Success(t.Translate(ctx, "Success toggle", ""))
	}

	redirectUrl := &url.URL{}
	*redirectUrl = *r.Original().URL
	redirectUrl.RawQuery = ""

	t.Redirect(redirectUrl.String(), http.StatusFound, w, r)
}

func (t Type) handleIndex(w *dashboard.Response, r *dashboard.Request, bind *Bind) {
	ctx := r.Context()

	t.Render(ctx, "index", map[string]interface{}{
		"components": bind.Components(),
	})
}
