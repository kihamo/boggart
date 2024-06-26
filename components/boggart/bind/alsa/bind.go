package alsa

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
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
	di.ConfigBind
	di.MetaBind
	di.MQTTBind
	di.WidgetBind

	playerStatus *atomic.Int64
	volume       *atomic.Int64
	mute         *atomic.Bool
	done         chan struct{}

	mutex   sync.RWMutex
	speaker *speakerWrapper
	stream  *streamWrapper
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	sn, err := machineid.ID()
	if err != nil {
		sn, err = os.Hostname()
		if err != nil {
			return err
		}
	}

	b.Meta().SetSerialNumber(sn)

	cfg := b.config()

	b.done = make(chan struct{}, 1)
	b.playerStatus.Set(StatusStopped.Int64())
	b.volume.Set(cfg.Volume)
	b.mute.Set(cfg.Mute)

	return nil
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
	b.speaker = newSpeakerWrapper(player, bufferSize, numBytes)
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
	b.stream = newStreamWrapper(source, format, b.Volume(), b.Mute())
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

	_ = b.MQTT().PublishAsync(context.Background(), b.config().TopicStateStatus.Format(b.Meta().SerialNumber()), status.String())
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
