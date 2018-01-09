package pulsar

func (p *Pulsar) readMetricFloat32(channel int64) (float32, error) {
	value, err := p.ReadMetrics(channel)
	if err != nil {
		return -1, err
	}

	return p.ToFloat32(value[0]), nil
}

func (p *Pulsar) TemperatureIn() (float32, error) {
	return p.readMetricFloat32(Channel3)
}

func (p *Pulsar) TemperatureOut() (float32, error) {
	return p.readMetricFloat32(Channel4)
}

func (p *Pulsar) TemperatureDelta() (float32, error) {
	return p.readMetricFloat32(Channel5)
}

func (p *Pulsar) Power() (float32, error) {
	return p.readMetricFloat32(Channel6)
}

func (p *Pulsar) Energy() (float32, error) {
	return p.readMetricFloat32(Channel7)
}

func (p *Pulsar) Capacity() (float32, error) {
	return p.readMetricFloat32(Channel8)
}

func (p *Pulsar) Consumption() (float32, error) {
	return p.readMetricFloat32(Channel9)
}

func (p *Pulsar) PulseInput1() (float32, error) {
	return p.readMetricFloat32(Channel10)
}

func (p *Pulsar) PulseInput2() (float32, error) {
	return p.readMetricFloat32(Channel11)
}

func (p *Pulsar) PulseInput3() (float32, error) {
	return p.readMetricFloat32(Channel12)
}

func (p *Pulsar) PulseInput4() (float32, error) {
	return p.readMetricFloat32(Channel13)
}
