package pulsar

import (
	"encoding/binary"
	"time"
)

func (p *Pulsar) DaylightSavingTime() (bool, error) {
	value, err := p.ReadSettings(ParamDaylightSavingTime)
	if err != nil {
		return false, err
	}

	return binary.BigEndian.Uint64(value) == 1, nil
}

func (p *Pulsar) Diagnostics() ([]byte, error) {
	value, err := p.ReadSettings(ParamDiagnostics)
	if err != nil {
		return nil, err
	}

	// TODO: split result
	return value, nil
}

func (p *Pulsar) Version() (uint16, error) {
	value, err := p.ReadSettings(ParamVersion)
	if err != nil {
		return 0, err
	}

	return uint16(binary.BigEndian.Uint64(value)), nil
}

func (p *Pulsar) OperatingTime() (time.Duration, error) {
	value, err := p.ReadSettings(ParamOperatingTime)
	if err != nil {
		return -1, err
	}

	return time.Hour * time.Duration(binary.BigEndian.Uint64(value)), nil
}
