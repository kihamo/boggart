package voice

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Speech(text string) error
	SpeechWithOptions(text string, volume int64, speed float64, speaker string) error
	PlayURL(url string) error
	Play() error
	Pause() error
	Stop() error
}
