package components

type Sensor struct {
	*Base

	UnitOfMeasurement string `json:"unit_of_measurement"`
}

func (c *Sensor) GetState() string {
	s := c.Base.GetState()
	if c.UnitOfMeasurement != "" {
		s += " " + c.UnitOfMeasurement
	}

	return s
}
