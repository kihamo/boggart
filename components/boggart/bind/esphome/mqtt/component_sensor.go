package mqtt

import (
	"encoding/json"
)

type ComponentSensor struct {
	*ComponentBase

	data struct {
		UnitOfMeasurement string `json:"unit_of_measurement"`
	}
}

func NewComponentSensor(id string) *ComponentSensor {
	return &ComponentSensor{
		ComponentBase: NewComponentBase(id, ComponentTypeSensor),
	}
}

func (c *ComponentSensor) State() interface{} {
	s := c.ComponentBase.State().(string)
	if c.data.UnitOfMeasurement != "" {
		s += " " + c.data.UnitOfMeasurement
	}

	return s
}

func (c *ComponentSensor) UnitOfMeasurement() string {
	return c.data.UnitOfMeasurement
}

func (c *ComponentSensor) UnmarshalJSON(b []byte) error {
	if err := c.ComponentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
