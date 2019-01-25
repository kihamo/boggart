package telegram

const (
	DefaultDebug = false
)

type Config struct {
	Token string
	Debug bool
}

func (t Type) Config() interface{} {
	return &Config{
		Debug: DefaultDebug,
	}
}
