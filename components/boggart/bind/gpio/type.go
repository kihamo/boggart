package gpio

import (
	"fmt"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type Type struct{}

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

	device := &Bind{
		pin:  g,
		mode: mode,
	}

	device.SetSerialNumber(g.Name())

	if _, ok := g.(gpio.PinIn); ok {
		go func() {
			device.waitForEdge()
		}()
	}

	return device, nil
}
