package openhab

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/mime"
	"github.com/kihamo/boggart/providers/openhab/client/items"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()

	action := q.Get("action")
	switch action {
	case "input":
		id := strings.TrimSpace(q.Get("id"))
		if id == "" {
			widget.NotFound(w, r)
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

			_, err := b.provider.Items.PutItemState(paramsPut)
			if err != nil {
				widget.InternalError(w, r, err)
				return
			}

			return
		}

		paramsGet := items.NewGetItemDataParams().
			WithItemname(id)

		response, err := b.provider.Items.GetItemData(paramsGet)
		if err != nil {
			widget.NotFound(w, r)
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

		widget.RenderLayout(r.Context(), "input", "ui", b.initUI(map[string]interface{}{
			"item":     response.Payload,
			"type":     q.Get("type"),
			"rows":     q.Get("rows"),
			"icon_url": iconURL,
		}, r))

		return

	case "image":
		// send to mqtt
		if r.IsPost() {
			topic := q.Get("topic")
			if topic == "" {
				widget.NotFound(w, r)
				return
			}

			if err := r.Original().ParseForm(); err == nil {
				for key, value := range r.Original().PostForm {
					if key == "payload" {
						b.MQTT().PublishAsync(r.Context(), mqtt.Topic(topic), strings.Join(value, ";"))
						break
					}
				}
			}

			return
		}

		u := q.Get("url")
		if u == "" {
			widget.NotFound(w, r)
			return
		}

		response, err := http.Get(u)
		defer response.Body.Close()

		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		var mimeType mime.Type

		mimeType, err = mime.TypeFromHTTPHeader(response.Header)
		if err != nil {
			var restored io.Reader

			mimeType, restored, err = mime.TypeFromDataRestored(response.Body)
			if err != nil {
				widget.InternalError(w, r, err)
				return
			}

			// присваиваем обратно, чтобы с этим можно было проджолжать работать
			response.Body = ioutil.NopCloser(restored)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		ext := mimeType.Extension()
		if ext == "" {
			ext = "image"
		}

		var refresh int64
		if v := q.Get("refresh"); v != "" {
			refresh, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				widget.InternalError(w, r, err)
				return
			}
		}

		widget.RenderLayout(r.Context(), "image", "ui", b.initUI(map[string]interface{}{
			"mime":     mimeType,
			"base64":   base64.StdEncoding.EncodeToString(body),
			"filename": time.Now().Format("20060102150405." + ext),
			"refresh":  refresh,
		}, r))

		return

	case "icon":
		req := r.Original().Clone(r.Context())
		req.URL.Path = "icon/" + r.URL().Query().Get("icon")

		b.proxy.ServeHTTP(w, req)

		return
	}

	widget.NotFound(w, r)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) initUI(vars map[string]interface{}, r *dashboard.Request) map[string]interface{} {
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
