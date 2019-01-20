package handlers

import (
	"net/http"
	"strconv"

	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/voice"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	StorageFileNameSpace = voice.ComponentName
)

type FileHandler struct {
	dashboard.Handler

	Voice   voice.Component
	Storage storage.Component
}

func (h *FileHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	message := q.Get(":message")
	message = voice.TrimMessage(message)

	if message == "" {
		h.NotFound(w, r)
		return
	}

	provider := h.Voice.TextToSpeechProvider()
	if provider == nil {
		h.NotFound(w, r)
		return
	}

	ctx := r.Context()
	config := dashboard.ConfigFromContext(ctx)

	var (
		err                                         error
		speed                                       float64
		language, speaker, emotion, format, quality string
		force                                       bool
	)

	if qSpeed := q.Get("speed"); qSpeed != "" {
		speed, err = strconv.ParseFloat(qSpeed, 64)
		if err != nil {
			h.NotFound(w, r)
			return
		}
	} else {
		speed = config.Float64(voice.ConfigYandexSpeechKitCloudSpeed)
	}

	if qForce := q.Get("force"); qForce != "" {
		force, err = strconv.ParseBool(qForce)
		if err != nil {
			h.NotFound(w, r)
			return
		}
	}

	if qLanguage := q.Get("language"); qLanguage != "" {
		language = qLanguage
	} else {
		language = config.String(voice.ConfigYandexSpeechKitCloudLanguage)
	}

	if qSpeaker := q.Get("speaker"); qSpeaker != "" {
		speaker = qSpeaker
	} else {
		speaker = config.String(voice.ConfigYandexSpeechKitCloudSpeaker)
	}

	if qEmotion := q.Get("emotion"); qEmotion != "" {
		emotion = qEmotion
	} else {
		emotion = config.String(voice.ConfigYandexSpeechKitCloudEmotion)
	}

	if qFormat := q.Get("format"); qFormat != "" {
		format = qFormat
	} else {
		format = config.String(voice.ConfigYandexSpeechKitCloudFormat)
	}

	if qQuality := q.Get("quality"); qQuality != "" {
		quality = qQuality
	} else {
		quality = config.String(voice.ConfigYandexSpeechKitCloudQuality)
	}

	providerURL, err := provider.GenerateURL(ctx, message, language, speaker, emotion, format, quality, speed)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	storageURL, err := h.Storage.SaveURLToFile(StorageFileNameSpace, providerURL, force)
	if err != nil && err != storage.ErrFileAlreadyExist {
		panic(err.Error())
	}

	h.Redirect(storageURL, http.StatusTemporaryRedirect, w, r)
}
