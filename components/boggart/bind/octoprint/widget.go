package octoprint

import (
	"bytes"

	"github.com/disintegration/imaging"
	"github.com/kihamo/boggart/providers/octoprint/client/util"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	b.handleSnapshot(w, r)
}

func (b *Bind) handleSnapshot(w *dashboard.Response, r *dashboard.Request) {
	settings := b.SystemSettings()
	if settings == nil {
		b.Widget().NotFound(w, r)
		return
	}

	if settings == nil || !settings.Webcam.WebcamEnabled {
		b.Widget().NotFound(w, r)
		return
	}

	params := util.NewUtilTestParamsWithContext(r.Context()).
		WithBody(util.UtilTestBody{
			Command:              "url",
			ContentTypeGuess:     true,
			ContentTypeWhitelist: []string{"image/*"},
			ContentTypeBlacklist: []string{},
			Method:               "GET",
			Response:             "bytes",
			Timeout:              settings.Webcam.SnapshotTimeout,
			URL:                  settings.Webcam.SnapshotURL,
			ValidSsl:             settings.Webcam.SnapshotSslValidation,
		})

	response, err := b.provider.Util.UtilTest(params, nil)
	if err != nil {
		b.Widget().InternalError(w, r, err)
		return
	}

	snapshot := response.GetPayload().Response

	if settings.Webcam.FlipH || settings.Webcam.FlipV || settings.Webcam.Rotate90 {
		img, err := imaging.Decode(bytes.NewReader(snapshot.Content))
		if err != nil {
			b.Widget().InternalError(w, r, err)
			return
		}

		if settings.Webcam.FlipH {
			img = imaging.FlipH(img)
		}

		if settings.Webcam.FlipV {
			img = imaging.FlipV(img)
		}

		if settings.Webcam.Rotate90 {
			img = imaging.Rotate90(img)
		}

		w.Header().Set("Content-Type", snapshot.AssumedContentType)
		imaging.Encode(w, img, imaging.JPEG)
	} else {
		w.Header().Set("Content-Type", snapshot.AssumedContentType)
		w.Write(snapshot.Content)
	}

}
