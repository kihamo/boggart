package chromecast

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/barnybug/go-cast/controllers"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/storage"
)

func (b *Bind) setVolume(ctx context.Context, level *float64, muted *bool) error {
	volume := controllers.Volume{
		Level: level,
		Muted: muted,
	}

	_, err := b.receiver.SetVolume(ctx, &volume)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
	}

	return err
}

func (b *Bind) Volume() (int64, error) {
	return atomic.LoadInt64(&b.volume), nil
}

func (b *Bind) SetVolume(ctx context.Context, percent int64) error {
	if percent > 100 {
		percent = 100
	} else if percent < 0 {
		percent = 0
	}

	level := float64(percent) / 100

	return b.setVolume(ctx, &level, nil)
}

func (b *Bind) Mute() (bool, error) {
	return atomic.LoadInt64(&b.mute) == 1, nil
}

func (b *Bind) SetMute(ctx context.Context, mute bool) error {
	return b.setVolume(ctx, nil, &mute)
}

func (b *Bind) PlayFromURL(ctx context.Context, url string) error {
	mimeType, err := storage.MimeTypeFromURL(url)
	if err != nil {
		return err
	}

	// TODO: check support format https://developers.google.com/cast/docs/media

	if mimeType == storage.MIMETypeUnknown {
		return errors.New("unknown audio format")
	}

	ctrl, err := b.Media(ctx)
	if err != nil {
		return err
	}

	item := controllers.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: mimeType.String(),
	}

	_, err = ctrl.LoadMedia(ctx, item, 0, true, map[string]interface{}{})
	return err
}

func (b *Bind) Seek(ctx context.Context, second uint64) error {
	return nil
}

func (b *Bind) Resume(ctx context.Context) error {
	return nil
}

func (b *Bind) Stop(ctx context.Context) error {
	_, err := b.receiver.QuitApp(ctx)
	return err
}

func (b *Bind) Pause(ctx context.Context) error {
	return nil
}
