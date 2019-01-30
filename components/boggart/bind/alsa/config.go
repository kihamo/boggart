package alsa

const (
	DefaultVolume = 50
	DefaultMute   = false
)

type Config struct {
	Volume        int64 `valid:"range(0|100)"`
	Mute          bool
	WidgetFileURL string `mapstructure:"widget_file_url" yaml:"widget_file_url"`
}

func (t Type) Config() interface{} {
	return &Config{
		Volume: DefaultVolume,
		Mute:   DefaultMute,
	}
}
