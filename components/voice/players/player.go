package players

import (
	"errors"
	"io"
)

var (
	ErrorUnknownAudioFormat = errors.New("unknown audio format")
	ErrorAlreadyPlaying     = errors.New("already playing")
)

type Player interface {
	PlayFromURL(url string) error
	PlayFromReader(reader io.ReadCloser) error
	PlayFromFile(file string) error
	Play() error
	Stop() error
	Pause() error
	Volume() (int64, error)
	SetVolume(percent int64) error
	Mute() (bool, error)
	SetMute(mute bool) error
	Close()
	Status() Status
}
