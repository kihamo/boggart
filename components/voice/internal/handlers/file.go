package handlers

import (
	"fmt"
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
	message := r.URL().Query().Get(":message")

	if message == "" {
		h.NotFound(w, r)
		return
	}

	// TODO: check exists

	/*
		data, err := base64.StdEncoding.DecodeString(message)
		fmt.Println(data, err)
		if err != nil {
			h.NotFound(w, r)
			return
		}
	*/
	provider := h.Voice.TextToSpeechProvider()
	if provider == nil {
		h.NotFound(w, r)
		return
	}

	config := dashboard.ConfigFromContext(r.Context())

	providerURL := provider.GenerateURL(
		r.Context(),
		message,
		config.String(voice.ConfigYandexSpeechKitCloudLanguage),
		config.String(voice.ConfigYandexSpeechKitCloudSpeaker),
		config.String(voice.ConfigYandexSpeechKitCloudEmotion),
		config.String(voice.ConfigYandexSpeechKitCloudFormat),
		config.String(voice.ConfigYandexSpeechKitCloudQuality),
		config.Float64(voice.ConfigYandexSpeechKitCloudSpeed))

	storageURL, err := h.Storage.SaveURLToFile("voice", providerURL, false)
	if err != nil && err != storage.ErrFileAlreadyExist {
		panic(err.Error())
	}

	fmt.Println(storageURL)

	// h.Redirect(storageURL, http.StatusTemporaryRedirect, w, r)
}
