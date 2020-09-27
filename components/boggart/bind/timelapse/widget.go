package timelapse

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/mime"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	vars := map[string]interface{}{}
	ctx := r.Context()
	q := r.URL().Query()
	action := q.Get("action")

	switch q.Get("action") {
	case "thumbnail", "download":
		filename := q.Get("file")
		m, buf, err := b.loadFile(ctx, filename)

		if err != nil {
			widget.InternalError(w, r, err)
			return
		}

		if m != mime.TypeUnknown {
			w.Header().Set("Content-Type", m.String())
		}

		if action == "download" {
			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
		}

		_, _ = io.Copy(w, buf)

		return

	default:
		page := 1
		if v, err := strconv.Atoi(q.Get("page")); err == nil && v > 1 {
			page = v
		}

		if files, err := b.Files(); err == nil {
			offsetLeft := (page - 1) * b.config.FilesOnPage
			if offsetLeft > len(files) {
				widget.NotFound(w, r)
				return
			}

			offsetRight := page * b.config.FilesOnPage
			if offsetRight > len(files) {
				offsetRight = len(files)
			}

			totalPages := len(files) / b.config.FilesOnPage
			if len(files)%b.config.FilesOnPage > 0 {
				totalPages++
			}

			vars["files"] = files[offsetLeft:offsetRight]
			vars["total"] = len(files)
			vars["page"] = page
			vars["pages"] = totalPages
		} else {
			widget.FlashError(r, "Get files failed with error %v", "", err)
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) loadFile(ctx context.Context, filename string) (m mime.Type, _ io.Reader, err error) {
	buf := bytes.NewBuffer(nil)
	m = mime.TypeUnknown

	if filename == "" {
		err = b.Capture(ctx, buf)
	} else {
		err = b.Load(filename, buf)
	}

	if err != nil {
		return m, nil, err
	}

	return mime.TypeFromDataRestored(buf)
}
