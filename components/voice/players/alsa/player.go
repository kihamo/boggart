package alsa

import (
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/boggart/components/voice/players"
)

/*
	TODO: wav, flac
*/

const (
	AudioFormatMP3  = "mp3"
	AudioFormatWAV  = "wav"
	AudioFormatFLAC = "flac"
)

var (
	ErrorStreamEmpty = errors.New("stream isn't set")
)

type Player struct {
	status        int64
	volumePercent int64
	done          chan struct{}

	mutex   sync.RWMutex
	speaker *speakerWrapper
	stream  *streamWrapper
}

func New() *Player {
	p := &Player{
		done: make(chan struct{}, 1),
	}
	p.setStatus(players.StatusStopped)
	p.SetVolume(50)

	return p
}

func (p *Player) PlayFromURL(url string) error {
	if p.Status() == players.StatusPlaying {
		return players.ErrorAlreadyPlaying
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
		return players.ErrorUnknownAudioFormat
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

	if err := p.initStream(response.Body, format); err != nil {
		return err
	}

	return p.Play()
}

func (p *Player) PlayFromReader(reader io.ReadCloser) error {
	if p.Status() == players.StatusPlaying {
		return players.ErrorAlreadyPlaying
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
		return players.ErrorUnknownAudioFormat
	}

	if err := p.initStream(reader, format); err != nil {
		return err
	}

	return p.Play()
}

func (p *Player) PlayFromFile(file string) error {
	if p.Status() == players.StatusPlaying {
		return players.ErrorAlreadyPlaying
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	return p.PlayFromReader(f)
}

func (p *Player) Play() error {
	if p.Status() == players.StatusPlaying {
		return players.ErrorAlreadyPlaying
	}

	if p.getStream().IsClosed() {
		return ErrorStreamEmpty
	}

	if p.getSpeaker().IsClosed() {
		if err := p.initSpeaker(); err != nil {
			return err
		}
	}

	go p.play()

	return nil
}

func (p *Player) Stop() error {
	if p.Status() == players.StatusPlaying {
		p.done <- struct{}{}
	} else {
		p.getSpeaker().Close()
		p.getStream().Close()
	}

	p.setStatus(players.StatusStopped)

	return nil
}

func (p *Player) Pause() error {
	if p.Status() == players.StatusPlaying {
		p.setStatus(players.StatusPause)
		p.done <- struct{}{}
	}

	return nil
}

func (p *Player) Volume() (int64, error) {
	return atomic.LoadInt64(&p.volumePercent), nil
}

func (p *Player) SetVolume(percent int64) error {
	if percent > 100 {
		percent = 100
	} else if percent < 0 {
		percent = 0
	}

	atomic.StoreInt64(&p.volumePercent, percent)
	p.getStream().SetVolume(percent)

	return nil
}

func (p *Player) Mute() (bool, error) {
	return p.getStream().Mute(), nil
}

func (p *Player) SetMute(mute bool) error {
	p.getStream().SetMute(mute)
	return nil
}

func (p *Player) Close() {
	p.Stop()
}

func (p *Player) initSpeaker() error {
	p.getSpeaker().Close()

	stream := p.getStream()
	if stream.IsClosed() {
		return ErrorStreamEmpty
	}

	bufferSize := stream.format.SampleRate.N(time.Second / 10)
	numBytes := bufferSize * 4

	player, err := oto.NewPlayer(stream.SampleRate(), stream.ChannelNum(), stream.BytesPerSample(), numBytes)
	if err != nil {
		return err
	}

	p.mutex.Lock()
	p.speaker = NewSpeakerWrapper(player, bufferSize, numBytes)
	p.mutex.Unlock()

	return nil
}

func (p *Player) initStream(s io.ReadCloser, f string) (err error) {
	p.getStream().Close()

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

	default:
		return players.ErrorUnknownAudioFormat
	}

	if err != nil {
		return err
	}

	p.mutex.Lock()
	p.stream = NewStreamWrapper(source, format)
	v, _ := p.Volume()
	p.stream.SetVolume(v)
	p.mutex.Unlock()

	return nil
}

func (p *Player) getSpeaker() *speakerWrapper {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.speaker
}

func (p *Player) getStream() *streamWrapper {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.stream
}

func (p *Player) setStatus(status players.Status) {
	atomic.StoreInt64(&p.status, status.Int64())
}

func (p *Player) Status() players.Status {
	return players.Status(atomic.LoadInt64(&p.status))
}

func (p *Player) play() {
	p.setStatus(players.StatusPlaying)

	defer func() {
		p.getSpeaker().Close()

		if p.Status() != players.StatusPause {
			p.getStream().Close()
			p.setStatus(players.StatusStopped)
		}
	}()

	speaker := p.getSpeaker()
	stream := p.getStream()

	samples := make([][2]float64, speaker.BufferSize())
	buf := make([]byte, speaker.NumBytes())

	for {
		select {
		case <-p.done:
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
