package chromecast

import (
	"context"
	"errors"

	"github.com/barnybug/go-cast/controllers"
	"github.com/kihamo/boggart/mime"
)

func (b *Bind) setVolume(ctx context.Context, level *float64, muted *bool) error {
	b.mutex.RLock()
	ctrlReceiver := b.ctrlReceiver
	b.mutex.RUnlock()

	if ctrlReceiver == nil {
		return errors.New("receiver controller isn't init")
	}

	volume := controllers.Volume{
		Level: level,
		Muted: muted,
	}

	_, err := ctrlReceiver.SetVolume(ctx, &volume)

	return err
}

func (b *Bind) Volume() int64 {
	return int64(b.volume.Load())
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

func (b *Bind) Mute() bool {
	return b.mute.Load()
}

func (b *Bind) SetMute(ctx context.Context, mute bool) error {
	return b.setVolume(ctx, nil, &mute)
}

func (b *Bind) PlayFromURL(ctx context.Context, url string) error {
	mimeType, err := mime.TypeFromURL(url)
	if err != nil {
		return err
	}

	// TODO: check support format https://developers.google.com/cast/docs/media

	if mimeType == mime.TypeUnknown {
		return errors.New("unknown audio format")
	}

	ctrlMedia, err := b.Media(ctx)
	if err != nil {
		return err
	}

	item := controllers.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: mimeType.String(),
	}

	_, err = ctrlMedia.LoadMedia(ctx, item, 0, true, map[string]interface{}{})

	return err
}

func (b *Bind) Stop(ctx context.Context) error {
	b.mutex.RLock()
	ctrlReceiver := b.ctrlReceiver
	b.mutex.RUnlock()

	if ctrlReceiver == nil {
		return errors.New("receiver controller isn't init")
	}

	ctrlMedia, err := b.Media(ctx)
	if err != nil {
		return err
	}

	_, err = ctrlMedia.Stop(ctx)
	if err != nil {
		return err
	}

	_, err = ctrlReceiver.QuitApp(ctx)

	return err
}

func (b *Bind) Seek(ctx context.Context, second uint64) error {
	return nil
}

func (b *Bind) Resume(ctx context.Context) error {
	if b.status.Load() != PlayerStatePaused {
		return nil
	}

	ctrlMedia, err := b.Media(ctx)
	if err != nil {
		return err
	}

	_, err = ctrlMedia.Play(ctx)

	return err
}

func (b *Bind) Pause(ctx context.Context) error {
	if b.status.Load() != PlayerStatePlaying {
		return nil
	}

	ctrlMedia, err := b.Media(ctx)
	if err != nil {
		return err
	}

	_, err = ctrlMedia.Pause(ctx)

	return err
}
