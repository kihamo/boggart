package text2speech

import (
	"bytes"
	"io"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	q := r.URL().Query()
	text := q.Get("text")

	if text == "" {
		t.NotFound(w, r)
		return
	}

	var (
		err                                         error
		speed                                       float64
		language, speaker, emotion, format, quality string
		force                                       bool
	)

	if val := q.Get("speed"); val != "" {
		speed, err = strconv.ParseFloat(val, 64)
		if err != nil {
			t.NotFound(w, r)
			return
		}
	}

	if val := q.Get("force"); val != "" {
		force, err = strconv.ParseBool(val)
		if err != nil {
			t.NotFound(w, r)
			return
		}
	}

	if val := q.Get("language"); val != "" {
		language = val
	}

	if val := q.Get("speaker"); val != "" {
		speaker = val
	}

	if val := q.Get("emotion"); val != "" {
		emotion = val
	}

	if val := q.Get("format"); val != "" {
		format = val
	}

	if val := q.Get("quality"); val != "" {
		quality = val
	}

	if format == "" {
		format = bind.config.Format
	}

	writer := bytes.NewBuffer(nil)

	err = bind.GenerateWriter(r.Context(), text, format, quality, language, speaker, emotion, speed, force, writer)
	if err != nil {
		t.InternalError(w, r, err)
	}

	switch format {
	case yandex_speechkit_cloud.FormatMP3:
		w.Header().Set("Content-Type", "audio/mpeg")
	case yandex_speechkit_cloud.FormatWAV:
		w.Header().Set("Content-Type", "audio/wav")
	case yandex_speechkit_cloud.FormatOPUS:
		w.Header().Set("Content-Type", "audio/opus")
	}

	io.Copy(w, writer)
}
