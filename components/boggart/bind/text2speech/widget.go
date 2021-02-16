package text2speech

import (
	"bytes"
	"io"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"

	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	text := q.Get("text")
	widget := b.Widget()

	if text == "" {
		widget.NotFound(w, r)
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
			widget.NotFound(w, r)
			return
		}
	}

	if val := q.Get("force"); val != "" {
		force, err = strconv.ParseBool(val)
		if err != nil {
			widget.NotFound(w, r)
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
		format = b.config().Format
	}

	writer := bytes.NewBuffer(nil)

	err = b.GenerateWriter(r.Context(), text, format, quality, language, speaker, emotion, speed, force, writer)
	if err != nil {
		widget.InternalError(w, r, err)
	}

	switch format {
	case speechkit.FormatMP3:
		w.Header().Set("Content-Type", "audio/mpeg")
	case speechkit.FormatWAV:
		w.Header().Set("Content-Type", "audio/wav")
	case speechkit.FormatOPUS:
		w.Header().Set("Content-Type", "audio/opus")
	}

	io.Copy(w, writer)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return nil
}
