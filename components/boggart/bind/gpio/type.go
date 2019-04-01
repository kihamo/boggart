package gpio

import (
	"fmt"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	g := gpioreg.ByName(fmt.Sprintf("GPIO%d", config.Pin))
	if g == nil {
		return nil, fmt.Errorf("GPIO %d not found", config.Pin)
	}

	var mode Mode
	switch strings.ToLower(config.Mode) {
	case "in":
		mode = ModeIn
	case "out":
		mode = ModeOut
	default:
		mode = ModeDefault
	}

	bind := &Bind{
		pin:  g,
		mode: mode,
	}

	bind.SetSerialNumber(g.Name())

	return bind, nil
}
