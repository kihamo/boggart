package hikvision

import (
	"bytes"
	"io"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	var (
		ch  uint64
		err error
	)

	if channel := r.URL().Query().Get("channel"); channel == "" {
		ch = b.Config().(*Config).WidgetChannel
	} else {
		ch, err = strconv.ParseUint(channel, 10, 64)
		if err != nil {
			t.NotFound(w, r)
			return
		}
	}

	bind, ok := b.Bind().(*Bind)
	if !ok {
		t.NotFound(w, r)
		return
	}

	buf := bytes.NewBuffer(nil)
	if err = bind.Snapshot(r.Context(), ch, buf); err != nil {
		t.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
	_, _ = io.Copy(w, buf)
}
