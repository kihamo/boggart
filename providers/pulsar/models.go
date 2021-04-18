package pulsar

type Diagnostic struct {
	bit uint8
}

func NewDiagnostic(bit uint8) *Diagnostic {
	return &Diagnostic{
		bit: bit,
	}
}

func (m *Diagnostic) IsEmptyBattery() bool {
	return m.bit&DiagnosticsEmptyBattery != 0
}

func (m *Diagnostic) IsErrorReadWriteEEPROM() bool {
	return m.bit&DiagnosticsErrorReadWriteEEPROM != 0
}

func (m *Diagnostic) IsResetCounters() bool {
	return m.bit&DiagnosticsResetCounters != 0
}

func (m *Diagnostic) IsThermometerInBroke() bool {
	return m.bit&DiagnosticsThermometerInBroke != 0
}

func (m *Diagnostic) IsThermometerOutBroke() bool {
	return m.bit&DiagnosticsThermometerOutBroke != 0
}

func (m *Diagnostic) IsNegativeTemperaturesDelta() bool {
	return m.bit&DiagnosticsNegativeTemperaturesDelta != 0
}
