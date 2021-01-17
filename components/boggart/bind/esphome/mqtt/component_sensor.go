package mqtt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentSensorData struct {
	UnitOfMeasurement *string `json:"unit_of_measurement,omitempty"`
	ExpireAfter       *uint32 `json:"expire_after,omitempty"`
	ForceUpdate       *bool   `json:"force_update,omitempty"`
}

type ComponentSensor struct {
	*componentBase

	data ComponentSensorData

	state            *atomic.Float32Null
	accuracyDecimals *atomic.Uint64
}

func NewComponentSensor(id string) *ComponentSensor {
	return &ComponentSensor{
		componentBase:    newComponentBase(id, ComponentTypeSensor),
		state:            atomic.NewFloat32Null(),
		accuracyDecimals: atomic.NewUint64(),
	}
}

func (c *ComponentSensor) State() interface{} {
	if c.state.IsNil() {
		return nil
	}

	return float64(c.state.Load())
}

func (c *ComponentSensor) StateFormat() string {
	if c.state.IsNil() {
		return ""
	}

	valueTpl := "%." + strconv.FormatUint(c.accuracyDecimals.Load(), 10) + "f"

	if unit := c.UnitOfMeasurement(); unit != "" {
		return fmt.Sprintf(valueTpl+" %s", c.state.Load(), unit)
	}

	return fmt.Sprintf(valueTpl, c.state.Load())
}

func (c *ComponentSensor) SetState(message mqtt.Message) error {
	state := message.String()

	value, err := strconv.ParseFloat(state, 64)
	if err != nil {
		return err
	}

	if i := strings.Index(state, "."); i > -1 {
		c.accuracyDecimals.Set(uint64(len(state[i+1:])))
	}

	c.state.Set(float32(value))
	metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(value)

	return nil
}

func (c *ComponentSensor) UnitOfMeasurement() string {
	if c.data.UnitOfMeasurement != nil {
		return *c.data.UnitOfMeasurement
	}

	return ""
}

func (c *ComponentSensor) AccuracyDecimals() uint64 {
	return c.accuracyDecimals.Load()
}

func (c *ComponentSensor) UnmarshalJSON(b []byte) error {
	if err := c.componentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
