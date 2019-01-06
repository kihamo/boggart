package handlers

import (
	"sort"
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
		"player":  "",
	}

	players := h.Component.Players()
	viewPlayers := make([]string, 0, len(players))
	for name := range players {
		viewPlayers = append(viewPlayers, name)
	}
	sort.Strings(viewPlayers)
	vars["players"] = viewPlayers

	if len(viewPlayers) > 0 {
		vars["player"] = viewPlayers[0]
	}

	if r.IsPost() {
		original := r.Original()

		var (
			err    error
			volume int64
			speed  float64
		)

		volume, err = strconv.ParseInt(original.FormValue("volume"), 10, 64)

		if err == nil {
			vars["volume"] = volume
			speed, err = strconv.ParseFloat(original.FormValue("speed"), 64)
		}

		if err == nil {
			player := original.FormValue("player")
			text := original.FormValue("text")
			speaker := original.FormValue("speaker")

			vars["speed"] = speed
			vars["player"] = player
			vars["text"] = text
			vars["speaker"] = speaker

			err = h.Component.SpeechWithOptions(r.Context(), player, text, volume, speed, speaker)
		}

		if err != nil {
			vars["error"] = err.Error()
		} else {
			vars["message"] = locale.Translate(r.Component().Name(), "Message send success", "")
		}
	}

	h.Render(r.Context(), "message", vars)
}
