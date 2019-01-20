package alsa

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/storage"
)

func (b *Bind) PlayFromURL(url string) error {
	if b.PlayerStatus() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	mimeType, err := storage.MimeTypeFromURL(url)
	if err != nil {
		return err
	}

	var format string

	switch mimeType {
	case storage.MIMETypeMPEG:
		format = AudioFormatMP3

		// TODO: wav, flac

	default:
		return ErrorUnknownAudioFormat
	}

	client := &http.Client{
		Transport: http.DefaultTransport,
		// что бы не рвались стримы
		Timeout: 0,
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if err := b.initStream(response.Body, format); err != nil {
		return err
	}

	return b.Play()
}

func (b *Bind) PlayFromReader(reader io.ReadCloser) error {
	if b.PlayerStatus() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	mimeType, err := storage.MimeTypeFromData(reader)
	if err != nil {
		return err
	}

	var format string

	switch mimeType {
	case "audio/mpeg":
		format = AudioFormatMP3

		// TODO: wav, flac

	default:
		return ErrorUnknownAudioFormat
	}

	if err := b.initStream(reader, format); err != nil {
		return err
	}

	return b.Play()
}

func (b *Bind) PlayFromFile(file string) error {
	if b.PlayerStatus() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	return b.PlayFromReader(f)
}

func (b *Bind) Play() error {
	if b.PlayerStatus() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	if b.getStream().IsClosed() {
		return ErrorStreamEmpty
	}

	if b.getSpeaker().IsClosed() {
		if err := b.initSpeaker(); err != nil {
			return err
		}
	}

	go b.play()

	return nil
}

func (b *Bind) Stop() error {
	if b.PlayerStatus() == StatusPlaying {
		b.done <- struct{}{}
	} else {
		b.getSpeaker().Close()
		b.getStream().Close()
	}

	b.setPlayerStatus(StatusStopped)

	return nil
}

func (b *Bind) Pause() error {
	if b.PlayerStatus() == StatusPlaying {
		b.setPlayerStatus(StatusPause)
		b.done <- struct{}{}
	}

	return nil
}

func (b *Bind) Volume() int64 {
	return atomic.LoadInt64(&b.volume)
}

func (b *Bind) SetVolume(percent int64) error {
	if percent > 100 {
		percent = 100
	} else if percent < 0 {
		percent = 0
	}

	prev := b.Volume()
	if percent == prev {
		return nil
	}

	atomic.StoreInt64(&b.volume, percent)
	b.getStream().SetVolume(percent)

	sn := mqtt.NameReplace(b.SerialNumber())
	ctx := context.Background()
	return b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), 0, true, percent)
}

func (b *Bind) Mute() bool {
	return atomic.LoadInt64(&b.mute) == 1
}

func (b *Bind) SetMute(mute bool) error {
	prev := b.Mute()
	if mute == prev {
		return nil
	}

	if mute {
		atomic.StoreInt64(&b.mute, 1)
	} else {
		atomic.StoreInt64(&b.mute, -1)
	}
	b.getStream().SetMute(mute)

	sn := mqtt.NameReplace(b.SerialNumber())
	ctx := context.Background()
	return b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), 0, true, mute)
}
