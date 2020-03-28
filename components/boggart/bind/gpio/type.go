package gpio

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	g := gpioreg.ByName("GPIO" + strconv.FormatUint(config.Pin, 10))
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

	config.TopicPinState = config.TopicPinState.Format(g.Number())
	config.TopicPinSet = config.TopicPinSet.Format(g.Number())

	bind := &Bind{
		config: config,
		pin:    g,
		mode:   mode,
		out:    atomic.NewBool(),
	}

	return bind, nil
}
