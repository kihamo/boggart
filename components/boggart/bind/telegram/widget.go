package telegram

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/mime"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	q := r.URL().Query()
	id := q.Get(paramFileName)

	if id == "" {
		t.NotFound(w, r)
		return
	}

	var mimeType string
	fileNameExt := ""

	if mimeType = q.Get(paramMIME); mimeType != "" {
		fileNameExt = mime.Type(mimeType).Extension()
		if fileNameExt == "" {
			mimeType = ""
		}
	}

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}

	fileName := time.Now().Format("20060102150405")
	if fileNameExt != "" {
		fileName += "." + fileNameExt
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	r1 := r.Original()
	r2 := new(http.Request)

	*r2 = *r1
	r2.URL = new(url.URL)
	*r2.URL = *r1.URL
	r2.URL.Path = "/" + id

	bind := b.Bind().(*Bind)
	bind.fileServer.ServeHTTP(w, r2)

	if bind.config.FileAutoClean {
		if strings.Contains(r1.Header.Get("user-agent"), "TelegramBot (like TwitterBot)") {
			bind.RemoveFile(id)
		}
	}
}
