package alsa

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	AudioFormatMP3  = "mp3"
	AudioFormatWAV  = "wav"
	AudioFormatFLAC = "flac"
	AudioFormatOGG  = "ogg"
)

var (
	ErrorUnknownAudioFormat = errors.New("unknown audio format")
	ErrorAlreadyPlaying     = errors.New("already playing")
	ErrorStreamEmpty        = errors.New("stream isn't set")
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	playerStatus *atomic.Int64
	volume       *atomic.Int64
	mute         *atomic.Bool
	done         chan struct{}

	mutex   sync.RWMutex
	speaker *speakerWrapper
	stream  *streamWrapper
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

func (b *Bind) Close() error {
	return b.Stop()
}

func (b *Bind) initSpeaker() error {
	b.getSpeaker().Close()

	stream := b.getStream()
	if stream.IsClosed() {
		return ErrorStreamEmpty
	}

	bufferSize := stream.format.SampleRate.N(time.Second / 10)
	numBytes := bufferSize * 4

	player, err := oto.NewPlayer(stream.SampleRate(), stream.ChannelNum(), stream.BytesPerSample(), numBytes)
	if err != nil {
		return err
	}

	b.mutex.Lock()
	b.speaker = NewSpeakerWrapper(player, bufferSize, numBytes)
	b.mutex.Unlock()

	return nil
}

func (b *Bind) initStream(s io.ReadCloser, f string) (err error) {
	b.getStream().Close()

	var (
		source beep.StreamSeekCloser
		format beep.Format
	)

	switch f {
	case AudioFormatMP3:
		source, format, err = mp3.Decode(s)
	case AudioFormatWAV:
		source, format, err = wav.Decode(s)
	case AudioFormatFLAC:
		source, format, err = flac.Decode(s)
	case AudioFormatOGG:
		source, format, err = vorbis.Decode(s)

	default:
		return ErrorUnknownAudioFormat
	}

	if err != nil {
		return err
	}

	b.mutex.Lock()
	b.stream = NewStreamWrapper(source, format, b.Volume(), b.Mute())
	b.mutex.Unlock()

	return nil
}

func (b *Bind) getSpeaker() *speakerWrapper {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.speaker
}

func (b *Bind) getStream() *streamWrapper {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.stream
}

func (b *Bind) setPlayerStatus(status Status) {
	if ok := b.playerStatus.Set(status.Int64()); !ok {
		return
	}

	sn := mqtt.NameReplace(b.SerialNumber())
	ctx := context.Background()

	_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicStateStatus.Format(sn), status.String())
}

func (b *Bind) PlayerStatus() Status {
	return Status(b.playerStatus.Load())
}

func (b *Bind) play() {
	b.setPlayerStatus(StatusPlaying)

	defer func() {
		b.getSpeaker().Close()

		if b.PlayerStatus() != StatusPause {
			b.getStream().Close()
			b.setPlayerStatus(StatusStopped)
		}
	}()

	speaker := b.getSpeaker()
	stream := b.getStream()

	samples := make([][2]float64, speaker.BufferSize())
	buf := make([]byte, speaker.NumBytes())

	for {
		select {
		case <-b.done:
			return

		default:
			if _, ok := stream.Stream(samples); !ok {
				return
			}

			for i := range samples {
				for c := range samples[i] {
					val := samples[i][c]

					if val < -1 {
						val = -1
					}

					if val > +1 {
						val = +1
					}

					valInt16 := int16(val * (1<<15 - 1))
					buf[i*4+c*2+0] = byte(valInt16)
					buf[i*4+c*2+1] = byte(valInt16 >> 8)
				}
			}

			if _, err := speaker.Write(buf); err != nil {
				return
			}
		}
	}
}
