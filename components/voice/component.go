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
	PlayURL(ctx context.Context, url string) error
	PlayReader(ctx context.Context, reader io.ReadCloser) error
	Play(ctx context.Context) error
	Pause(ctx context.Context) error
	Stop(ctx context.Context) error
	Volume(ctx context.Context) int64
	SetVolume(ctx context.Context, percent int64) error
}
