package openhab

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/storage"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/openhab/client/items"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	q := r.URL().Query()

	action := q.Get("action")
	switch action {
	case "input":
		id := strings.TrimSpace(q.Get("id"))
		if id == "" {
			t.NotFound(w, r)
			return
		}

		if r.IsPost() {
			value := r.Original().FormValue("value")
			if value == "" {
				value = "NULL"
			}

			paramsPut := items.NewPutItemStateParams().
				WithItemname(id).
				WithBody(value)

			_, err := bind.provider.Items.PutItemState(paramsPut)
			if err != nil {
				t.InternalError(w, r, err)
				return
			}

			return
		}

		paramsGet := items.NewGetItemDataParams().
			WithItemname(id)

		response, err := bind.provider.Items.GetItemData(paramsGet)
		if err != nil {
			t.NotFound(w, r)
			return
		}

		if response.Payload.Type == "DateTime" {
			t, _ := time.Parse("2006-01-02T15:04:05.999Z0700", response.Payload.State)
			response.Payload.State = t.Format("2006-01-02T15:04:05")
		}

		var iconURL string

		if response.Payload.Category != "" {
			iconURL = r.URL().Path + "/?action=icon&icon=" + response.Payload.Category +
				"&state=" + response.Payload.State + "&format=svg&anyFormat=true"

			if key := q.Get(boggart.AccessKeyName); key != "" {
				iconURL += "&" + boggart.AccessKeyName + "=" + key
			}
		}

		t.RenderLayout(r.Context(), "input", "ui", t.initUI(map[string]interface{}{
			"item":     response.Payload,
			"type":     q.Get("type"),
			"rows":     q.Get("rows"),
			"icon_url": iconURL,
		}, r))

		return

	case "image":
		u := q.Get("url")
		if u == "" {
			t.NotFound(w, r)
			return
		}

		response, err := http.Get(u)
		if err != nil {
			t.InternalError(w, r, err)
			return
		}

		var mime storage.MIMEType

		mime, err = storage.MimeTypeFromHTTPHeader(response.Header)
		if err != nil {
			copyBody := &bytes.Buffer{}
			if _, err := io.CopyN(copyBody, response.Body, 128); err != nil {
				t.InternalError(w, r, err)
				return
			}

			mime, err = storage.MimeTypeFromData(bytes.NewBuffer(copyBody.Bytes()))
			if err != nil {
				t.InternalError(w, r, err)
				return
			}

			// довычитываем все остальное, так как body уже порвался на две части до 128 и послке
			if _, err := io.Copy(copyBody, response.Body); err != nil {
				t.InternalError(w, r, err)
				return
			}

			// присваиваем обратно, чтобы с этим можно было проджолжать работать
			response.Body = ioutil.NopCloser(copyBody)
			defer copyBody.Reset()
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.InternalError(w, r, err)
			return
		}

		var ext string

		switch mime {
		case storage.MIMETypeJPEG, storage.MIMETypeJPG:
			ext = "jpg"
		case storage.MIMETypePNG:
			ext = "png"
		case storage.MIMETypeGIF:
			ext = "gif"
		default:
			ext = "image"
		}

		var refresh int64
		if v := q.Get("refresh"); v != "" {
			refresh, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				t.InternalError(w, r, err)
				return
			}
		}

		t.RenderLayout(r.Context(), "image", "ui", t.initUI(map[string]interface{}{
			"mime":     mime,
			"base64":   base64.StdEncoding.EncodeToString(body),
			"filename": time.Now().Format("20060102150405." + ext),
			"refresh":  refresh,
		}, r))

		return

	case "icon":
		req := r.Original().Clone(r.Context())
		req.URL.Path = "icon/" + r.URL().Query().Get("icon")

		bind.proxy.ServeHTTP(w, req)

		return
	}

	t.NotFound(w, r)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (t Type) initUI(vars map[string]interface{}, r *dashboard.Request) map[string]interface{} {
	q := r.URL().Query()

	style := q.Get("style")
	switch style {
	case "android", "basicui":
		// skip
	default:
		if strings.Contains(strings.ToLower(r.Original().Header.Get("User-Agent")), "android") {
			style = "android"
		} else {
			style = "basicui"
		}
	}
	vars["style"] = style
	vars["theme"] = q.Get("theme")

	return vars
}
