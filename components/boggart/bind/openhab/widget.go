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
	"github.com/kihamo/boggart/providers/openhab3"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	widget := b.Widget()
	ctx := r.Context()

	switch q.Get("action") {
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

			//_, err := b.provider.UpdateItemStateWithBody(ctx, id, nil, "text/plain", strings.NewReader(value))
			_, err := b.provider.SendItemCommandWithBody(ctx, id, "text/plain", strings.NewReader(value))
			if err != nil {
				widget.InternalError(w, r, err)
				return
			}

			return
		}

		resp, err := b.provider.GetItemByName(ctx, id, nil)
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		response, err := openhab3.ParseGetItemByNameResponse(resp)
		if err != nil || response.JSON200 == nil {
			widget.NotFound(w, r)
			return
		}

		if *response.JSON200.Type == "DateTime" {
			t, _ := time.Parse("2006-01-02T15:04:05.999Z0700", *response.JSON200.State)
			*response.JSON200.State = t.Format("2006-01-02T15:04:05")
		}

		widget.RenderLayout(ctx, "input", "ui", b.initUI(map[string]interface{}{
			"item_label": *response.JSON200.Label,
			"item_type":  *response.JSON200.Type,
			"item_state": *response.JSON200.State,
			"type":       q.Get("type"),
			"rows":       q.Get("rows"),
			"icon_url":   b.iconURL(r, response.JSON200),
		}, r))

		return

	case "link":
		id := strings.TrimSpace(q.Get("id"))
		if id == "" {
			widget.NotFound(w, r)
			return
		}

		resp, err := b.provider.GetItemByName(ctx, id, nil)
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		response, err := openhab3.ParseGetItemByNameResponse(resp)
		if err != nil || response.JSON200 == nil {
			widget.NotFound(w, r)
			return
		}

		widget.RenderLayout(ctx, "link", "ui", b.initUI(map[string]interface{}{
			"item_label": *response.JSON200.Label,
			"item_name":  *response.JSON200.Name,
			"item_state": *response.JSON200.State,
			"icon_url":   b.iconURL(r, response.JSON200),
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
						b.MQTT().PublishAsync(ctx, mqtt.Topic(topic), strings.Join(value, ";"))
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

		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		defer response.Body.Close()

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

		widget.RenderLayout(ctx, "image", "ui", b.initUI(map[string]interface{}{
			"mime":     mimeType,
			"base64":   base64.StdEncoding.EncodeToString(body),
			"filename": time.Now().Format("20060102150405." + ext),
			"refresh":  refresh,
		}, r))

		return

	case "icon":
		req := r.Original().Clone(ctx)
		req.URL.Path = "icon/" + r.URL().Query().Get("icon")

		b.proxy.ServeHTTP(w, req)

		return
	}

	widget.NotFound(w, r)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) iconURL(r *dashboard.Request, payload *openhab3.EnrichedItemDTO) (icon string) {
	if payload.Category != nil && payload.State != nil && *payload.Category != "" {
		icon = r.URL().Path + "/?action=icon&icon=" + *payload.Category +
			"&state=" + *payload.State + "&format=svg&anyFormat=true"

		if key := r.URL().Query().Get(boggart.AccessKeyName); key != "" {
			icon += "&" + boggart.AccessKeyName + "=" + key
		}
	}

	return icon
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
