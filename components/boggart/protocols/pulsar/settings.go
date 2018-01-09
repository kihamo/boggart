package pulsar

func (p *Pulsar) DaylightSavingTime() (bool, error) {
	settings, err := p.ReadSettings(ParamDaylightSavingTime)
	if err != nil {
		return false, err
	}

	return settings[0] == 1, nil
}

func (p *Pulsar) Diagnostics() ([]byte, error) {
	settings, err := p.ReadSettings(ParamDiagnostics)
	if err != nil {
		return nil, err
	}

	// TODO: split result
	return settings, nil
}

func (p *Pulsar) Version() (uint16, error) {
	settings, err := p.ReadSettings(ParamVersion)
	if err != nil {
		return 0, err
	}

	return uint16(settings[0]), nil
}
