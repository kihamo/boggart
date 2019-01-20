package alsa

const (
	DefaultVolume = 50
	DefaultMute   = false
)

type Config struct {
	Volume int64 `valid:"range(0|100)"`
	Mute   bool
}

func (t Type) Config() interface{} {
	return &Config{
		Volume: DefaultVolume,
		Mute:   DefaultMute,
	}
}
