package handlers

import (
	"strconv"

	"github.com/kihamo/boggart/components/voice"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

type MessageHandler struct {
	dashboard.Handler

	Component voice.Component
}

func (h *MessageHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	config := r.Config()
	locale := i18n.NewOrNopFromRequest(r)

	vars := map[string]interface{}{
		"volume":  config.Int64(voice.ConfigSpeechVolume),
		"speed":   config.Float64(voice.ConfigYandexSpeechKitCloudSpeed),
		"text":    "",
		"speaker": config.String(voice.ConfigYandexSpeechKitCloudSpeaker),
	}

	if r.IsPost() {
		var (
			err    error
			volume int64
			speed  float64
		)

		volume, err = strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)

		if err == nil {
			vars["volume"] = volume
			speed, err = strconv.ParseFloat(r.Original().FormValue("speed"), 64)
		}

		if err == nil {
			text := r.Original().FormValue("text")
			speaker := r.Original().FormValue("speaker")

			vars["speed"] = speed
			vars["text"] = text
			vars["speaker"] = speaker

			err = h.Component.SpeechWithOptions(text, volume, speed, speaker)
		}

		if err != nil {
			vars["error"] = err.Error()
		} else {
			vars["message"] = locale.Translate(r.Component().Name(), "Message send success", "")
		}
	}

	h.Render(r.Context(), "message", vars)
}
