package timelapse

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"time"

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

		to := time.Now()
		from := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

		if queryTime := q.Get("from"); queryTime != "" {
			if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
				from = tm
			} else {
				widget.FlashError(r, "Parse date from failed with error %v", "", err)
			}
		}

		if queryTime := q.Get("to"); queryTime != "" {
			if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
				to = tm
			} else {
				widget.FlashError(r, "Parse date to failed with error %v", "", err)
			}
		}

		if files, err := b.Files(&from, &to); err == nil {
			filesTotal := len(files)
			cfg := b.config()

			offsetLeft := (page - 1) * cfg.FilesOnPage
			if offsetLeft > filesTotal {
				widget.NotFound(w, r)
				return
			}

			offsetRight := page * cfg.FilesOnPage
			if offsetRight > filesTotal {
				offsetRight = filesTotal
			}

			pagesTotal := filesTotal / cfg.FilesOnPage
			if filesTotal%cfg.FilesOnPage > 0 {
				pagesTotal++
			}

			var sizeTotal int64
			for _, f := range files {
				sizeTotal += f.Size()
			}

			vars["files"] = files[offsetLeft:offsetRight]
			vars["files_total"] = filesTotal
			vars["offset_left"] = offsetLeft
			vars["offset_right"] = offsetRight
			vars["size_total"] = sizeTotal
			vars["page"] = page
			vars["pages"] = pagesTotal
		} else {
			widget.FlashError(r, "Get files failed with error %v", "", err)
		}

		vars["date_from"] = from
		vars["date_to"] = to
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
