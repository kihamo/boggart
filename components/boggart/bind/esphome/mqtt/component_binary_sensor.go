package mqtt

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	DeviceClassConnectivity = "connectivity"
)

type ComponentBinarySensor struct {
	*componentBase

	data struct {
		DeviceClass string `json:"device_class"`
		PayloadOn   string `json:"payload_on"`
		PayloadOff  string `json:"payload_off"`
	}

	state *atomic.BoolNull
}

func NewComponentBinarySensor(id string, discoveryTopic mqtt.Topic) *ComponentBinarySensor {
	return &ComponentBinarySensor{
		componentBase: newComponentBase(id, ComponentTypeBinarySensor, discoveryTopic),
		state:         atomic.NewBoolNull(),
	}
}

func (c *ComponentBinarySensor) State() interface{} {
	if c.state.IsNil() {
		return nil
	}

	return c.state.Load()
}

func (c *ComponentBinarySensor) StateFormat() string {
	if c.state.IsNil() {
		return ""
	}

	if c.state.IsTrue() {
		if c.data.PayloadOn != "" {
			return c.data.PayloadOn
		}

		return "ON"
	}

	if c.data.PayloadOff != "" {
		return c.data.PayloadOff
	}

	return "OFF"
}

func (c *ComponentBinarySensor) SetState(message mqtt.Message) error {
	var value bool
	state := message.String()

	// platform: status
	if c.data.PayloadOn != "" && c.data.PayloadOff != "" {
		switch state {
		case c.data.PayloadOn:
			value = true
		case c.data.PayloadOff:
			value = false
		default:
			return errors.New("unknown status binary sensor state " + state)
		}
	} else {
		value = bytes.Equal([]byte(state), stateON)
	}

	if value {
		c.state.True()
		metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(1)
	} else {
		c.state.False()
		metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(0)
	}

	return nil
}

func (c *ComponentBinarySensor) DeviceClass() string {
	return c.data.DeviceClass
}

func (c *ComponentBinarySensor) PayloadOn() string {
	return c.data.PayloadOn
}

func (c *ComponentBinarySensor) PayloadOff() string {
	return c.data.PayloadOff
}

func (c *ComponentBinarySensor) UnmarshalJSON(b []byte) error {
	if err := c.componentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
