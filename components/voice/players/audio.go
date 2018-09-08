package players

import (
	"errors"
	"fmt"
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
	ErrorAlreadyPlaying     = errors.New("Already playing")
	ErrorUnknownAudioFormat = errors.New("Unknown audio format")
	ErrorStreamEmpty        = errors.New("Stream isn't set")
)

type AudioPlayer struct {
	status        int64
	volumePercent int64
	done          chan struct{}

	mutex   sync.RWMutex
	speaker *speakerWrapper
	stream  *streamWrapper
}

func NewAudio() *AudioPlayer {
	p := &AudioPlayer{
		done: make(chan struct{}, 1),
	}
	p.setStatus(StatusStopped)
	p.SetVolume(50)

	return p
}

func (p *AudioPlayer) PlayFromURL(url string) error {
	if p.Status() == StatusPlaying {
		return ErrorAlreadyPlaying
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

	mimeType, err := mimeFromHeader(response.Header)
	if err != nil {
		return err
	}

	if mimeType == "" {
		mimeType, err = mimeFromData(response.Body)
	}

	var format string

	switch mimeType {
	case "audio/mpeg":
		format = AudioFormatMP3

		// TODO: wav, flac

	default:
		return ErrorUnknownAudioFormat
	}

	if err := p.initStream(response.Body, format); err != nil {
		return err
	}

	return p.Play()
}

func (p *AudioPlayer) PlayFromReader(reader io.ReadCloser) error {
	if p.Status() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	mimeType, err := mimeFromData(reader)
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

	if err := p.initStream(reader, format); err != nil {
		return err
	}

	return p.Play()
}

func (p *AudioPlayer) PlayFromFile(file string) error {
	if p.Status() == StatusPlaying {
		return ErrorAlreadyPlaying
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	return p.PlayFromReader(f)
}

func (p *AudioPlayer) Play() error {
	if p.Status() == StatusPlaying {
		return ErrorAlreadyPlaying
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

func (p *AudioPlayer) Stop() error {
	if p.Status() == StatusPlaying {
		p.done <- struct{}{}
	} else {
		p.getSpeaker().Close()
		p.getStream().Close()
	}

	p.setStatus(StatusStopped)

	return nil
}

func (p *AudioPlayer) Pause() error {
	if p.Status() == StatusPlaying {
		p.setStatus(StatusPause)
		p.done <- struct{}{}
	}

	return nil
}

func (p *AudioPlayer) Volume() int64 {
	return atomic.LoadInt64(&p.volumePercent)
}

func (p *AudioPlayer) SetVolume(percent int64) error {
	if percent > 100 {
		percent = 100
	} else if percent < 0 {
		percent = 0
	}

	atomic.StoreInt64(&p.volumePercent, percent)
	p.getStream().Volume(percent)

	return nil
}

func (p *AudioPlayer) Close() {
	p.Stop()
}

func (p *AudioPlayer) initSpeaker() error {
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

	player.SetUnderrunCallback(func() {
		fmt.Println("UNDERRUN, YOUR CODE IS SLOW")
	})

	p.mutex.Lock()
	p.speaker = NewSpeakerWrapper(player, bufferSize, numBytes)
	p.mutex.Unlock()

	return nil
}

func (p *AudioPlayer) initStream(s io.ReadCloser, f string) (err error) {
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
		return ErrorUnknownAudioFormat
	}

	if err != nil {
		return err
	}

	p.mutex.Lock()
	p.stream = NewStreamWrapper(source, format)
	p.stream.Volume(p.Volume())
	p.mutex.Unlock()

	return nil
}

func (p *AudioPlayer) getSpeaker() *speakerWrapper {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.speaker
}

func (p *AudioPlayer) getStream() *streamWrapper {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.stream
}

func (p *AudioPlayer) setStatus(status Status) {
	atomic.StoreInt64(&p.status, status.Int64())
}

func (p *AudioPlayer) Status() Status {
	return Status(atomic.LoadInt64(&p.status))
}

func (p *AudioPlayer) play() {
	p.setStatus(StatusPlaying)

	defer func() {
		p.getSpeaker().Close()

		if p.Status() != StatusPause {
			p.getStream().Close()
			p.setStatus(StatusStopped)
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
