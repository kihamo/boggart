package nut

import (
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	variables, err := b.Bind().(*Bind).Variables()
	variablesByName := make(map[string]interface{}, len(variables))
	charged := false

	for _, v := range variables {
		variablesByName[v.Name] = v

		switch v.Name {
		case "ups.status":
			charged = strings.HasSuffix(v.Value.(string), "CHRG")
		}
	}

	vars := map[string]interface{}{
		"variables":         variables,
		"variables_by_name": variablesByName,
		"charged":           charged,
		"error":             err,
	}

	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get variables failed with error %s", "", err.Error()))
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
