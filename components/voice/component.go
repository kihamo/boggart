package voice

import (
	"context"
	"io"

	"github.com/kihamo/boggart/components/voice/players"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Players() map[string]players.Player
	Speech(ctx context.Context, player string, text string) error
	SpeechWithOptions(ctx context.Context, player string, text string, volume int64, speed float64, speaker string) error
	PlayURL(ctx context.Context, player string, url string) error
	PlayReader(ctx context.Context, player string, reader io.ReadCloser) error
	Play(ctx context.Context, player string) error
	Pause(ctx context.Context, player string) error
	Stop(ctx context.Context, player string) error
	Volume(ctx context.Context, player string) (volume int64, err error)
	SetVolume(ctx context.Context, player string, percent int64) error
}
