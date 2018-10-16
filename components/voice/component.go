package voice

import (
	"context"
	"io"

	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Speech(ctx context.Context, text string) error
	SpeechWithOptions(ctx context.Context, text string, volume int64, speed float64, speaker string) error
	PlayURL(url string) error
	PlayReader(reader io.ReadCloser) error
	Play() error
	Pause() error
	Stop() error
	Volume() int64
	SetVolume(percent int64) error
}
