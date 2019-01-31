package telegram

const (
	DefaultDebug          = false
	DefaultUpdatesEnabled = false
	DefaultUpdatesBuffer  = 100
	DefaultUpdatesTimeout = 60
)

type Config struct {
	Token          string
	Debug          bool
	UpdatesEnabled bool `mapstructure:"updates_enabled" yaml:"updates_enabled"`
	UpdatesBuffer  int  `mapstructure:"updates_buffer" yaml:"updates_buffer"`
	UpdatesTimeout int  `mapstructure:"updates_timeout" yaml:"updates_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		Debug:          DefaultDebug,
		UpdatesEnabled: DefaultUpdatesEnabled,
		UpdatesBuffer:  DefaultUpdatesBuffer,
		UpdatesTimeout: DefaultUpdatesTimeout,
	}
}
