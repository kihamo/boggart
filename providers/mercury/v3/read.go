package v3

import (
	"time"
)

type array int
type month int
type tariff int
type powerNumber int
type phase int

const (
	ArrayReset                 array = 0x00
	ArrayCurrentYear           array = 0x10
	ArrayPreviousYear          array = 0x20
	ArrayMonth                 array = 0x30
	ArrayCurrentDay            array = 0x40
	ArrayPreviousDay           array = 0x50
	ArrayActiveEnergy          array = 0x60
	ArrayBeginningCurrentYear  array = 0x90
	ArrayBeginningPreviousYear array = 0xA0
	ArrayBeginningMonth        array = 0xB0
	ArrayBeginningCurrentDay   array = 0xC0
	ArrayBeginningPreviousDay  array = 0xD0

	MonthJanuary   month = 0x01
	MonthFebruary  month = 0x02
	MonthMarch     month = 0x03
	MonthApril     month = 0x04
	MonthMay       month = 0x05
	MonthJune      month = 0x06
	MonthJuly      month = 0x07
	MonthAugust    month = 0x08
	MonthSeptember month = 0x09
	MonthOctober   month = 0x0A
	MonthNovember  month = 0x0B
	MonthDecember  month = 0x0C

	TariffAll tariff = 0x0
	Tariff1   tariff = 0x1
	Tariff2   tariff = 0x2
	Tariff3   tariff = 0x3
	Tariff4   tariff = 0x4

	AuxiliaryPower        int64 = 0x0
	AuxiliaryVoltage      int64 = 0x1
	AuxiliaryAmperage     int64 = 0x2
	AuxiliaryPowerFactors int64 = 0x3
	AuxiliaryFrequency    int64 = 0x4
	AuxiliaryPhasesAngle  int64 = 0x5

	PowerNumberP powerNumber = 0x0
	PowerNumberQ powerNumber = 0x1
	PowerNumberS powerNumber = 0x2

	PhaseNumberAll phase = 0x0
	PhaseNumber1   phase = 0x1
	PhaseNumber2   phase = 0x2
	PhaseNumber3   phase = 0x3
)

func (m *MercuryV3) ReadParameter(param byte) ([]byte, error) {
	resp, err := m.Request(&Request{
		Address:       m.options.address,
		Code:          RequestCodeReadParameter,
		ParameterCode: &[]byte{param}[0],
	})

	if err != nil {
		return nil, err
	}

	return resp.Payload, err
}

// 2.5.18. ЧТЕНИЕ СЕРИЙНОГО НОМЕРА СЧЕТЧИКА И ДАТЫ ВЫПУСКА.
func (m *MercuryV3) SerialNumberAndBuildDate() (serialNumber string, buildDate time.Time, err error) {
	resp, err := m.ReadParameter(ParamCodeSerialNumberAndBuildDate)

	if err != nil {
		return
	}

	serialNumber = ParseSerialNumber(resp[0:4])
	buildDate = ParseBuildDate(resp[4:7])

	return
}

// 2.5.19. УСКОРЕННЫЙ РЕЖИМ ЧТЕНИЯ ИНДИВИДУАЛЬНЫХ ПАРАМЕТРОВ ПРИБОРА.
func (m *MercuryV3) ForceReadParameters() (serialNumber string, buildDate time.Time, firmwareVersion string, t *Type, err error) {
	resp, err := m.ReadParameter(ParamCodeForceReadParameters)

	if err != nil {
		return
	}

	serialNumber = ParseSerialNumber(resp[0:4])
	buildDate = ParseBuildDate(resp[4:7])
	firmwareVersion = ParseFirmwareVersion(resp[7:10])
	t = ParseType(resp[9:16])

	return
}

// 2.5.20. ЧТЕНИЕ КОЭФФИЦИЕНТА ТРАНСФОРМАЦИИ СЧЁТЧИКА.
func (m *MercuryV3) TransformationCoefficient() (uint8, uint8, error) {
	resp, err := m.ReadParameter(ParamCodeTransformationCoefficient)

	if err != nil {
		return 0, 0, err
	}

	return resp[0]*10 + resp[1], resp[2]*10 + resp[3], nil
}

// 2.5.21. ЧТЕНИЕ ВЕРСИИ ПО СЧЁТЧИКА.
func (m *MercuryV3) FirmwareVersion() (string, error) {
	resp, err := m.ReadParameter(ParamCodeVersion)

	if err != nil {
		return "", err
	}

	return ParseFirmwareVersion(resp), nil
}

// 2.5.23. ЧТЕНИЕ СЕТЕВОГО АДРЕСА.
func (m *MercuryV3) Address() (int64, error) {
	resp, err := m.ReadParameter(ParamCodeAddress)

	if err != nil {
		return 0, err
	}

	if err := PayloadError(resp); err != nil {
		return 0, err
	}

	return int64(resp[1]), nil
}

// 2.5.32. ЧТЕНИЕ ВСПОМОГАТЕЛЬНЫХ ПАРАМЕТРОВ.
func (m *MercuryV3) AuxiliaryParameters(bwri int64) ([]byte, error) {
	resp, err := m.Request(&Request{
		Address:            m.options.address,
		Code:               RequestCodeReadParameter,
		ParameterCode:      &[]byte{ParamCodeAuxiliaryParameters12}[0],
		ParameterExtension: &[]byte{byte(bwri)}[0],
	})

	if err != nil {
		return nil, err
	}

	if err := ResponseError(resp); err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// 2.5.32.1 Ответ прибора на запрос чтения мощности.
func (m *MercuryV3) Power(number powerNumber) (sum, phase1, phase2, phase3 float64, err error) {
	var resp []byte

	bwri := AuxiliaryPower<<4 | int64(number)<<2 | int64(PhaseNumberAll)
	resp, err = m.AuxiliaryParameters(bwri)

	if err == nil {
		sum = float64(ParseValue3Bytes(resp[:3])) / 100
		phase1 = float64(ParseValue3Bytes(resp[3:6])) / 100
		phase2 = float64(ParseValue3Bytes(resp[6:9])) / 100
		phase3 = float64(ParseValue3Bytes(resp[9:12])) / 100
	}

	return
}

// 2.5.32.2 Ответ прибора на запрос чтения фазного и линейного напряжения, тока и углов между фазными напряжениями.
func (m *MercuryV3) Voltage() (phase1, phase2, phase3 float64, err error) {
	var resp []byte

	bwri := AuxiliaryVoltage<<4 | 0x1
	resp, err = m.AuxiliaryParameters(bwri)

	if err == nil {
		phase1 = float64(ParseValue3Bytes(resp[:3])) / 100
		phase2 = float64(ParseValue3Bytes(resp[3:6])) / 100
		phase3 = float64(ParseValue3Bytes(resp[6:9])) / 100
	}

	return
}

func (m *MercuryV3) Amperage() (phase1, phase2, phase3 float64, err error) {
	var resp []byte

	bwri := AuxiliaryAmperage<<4 | 0x1
	resp, err = m.AuxiliaryParameters(bwri)

	if err == nil {
		phase1 = float64(ParseValue3Bytes(resp[:3])) / 1000
		phase2 = float64(ParseValue3Bytes(resp[3:6])) / 1000
		phase3 = float64(ParseValue3Bytes(resp[6:9])) / 1000
	}

	return
}

func (m *MercuryV3) PhasesAngle() (phase1, phase2, phase3 float64, err error) {
	var resp []byte

	bwri := AuxiliaryPhasesAngle<<4 | 0x1
	resp, err = m.AuxiliaryParameters(bwri)

	if err == nil {
		phase1 = float64(ParseValue3Bytes(resp[:3])) / 1000
		phase2 = float64(ParseValue3Bytes(resp[3:6])) / 1000
		phase3 = float64(ParseValue3Bytes(resp[6:9])) / 1000
	}

	return
}

// 2.5.32.3 Ответ прибора на запрос чтения коэффициентов мощности.
func (m *MercuryV3) PowerFactors() (sum, phase1, phase2, phase3 float64, err error) {
	var resp []byte

	bwri := AuxiliaryPowerFactors<<4 | int64(PhaseNumberAll)
	resp, err = m.AuxiliaryParameters(bwri)

	if err == nil {
		sum = float64(ParseValue3Bytes(resp[:3])) / 1000
		phase1 = float64(ParseValue3Bytes(resp[3:6])) / 1000
		phase2 = float64(ParseValue3Bytes(resp[6:9])) / 1000
		phase3 = float64(ParseValue3Bytes(resp[9:12])) / 1000
	}

	return
}

// 2.5.32.4 Ответ прибора на запрос чтения частоты (запрос 11h, 14h, 16h).
func (m *MercuryV3) Frequency() (float64, error) {
	bwri := AuxiliaryFrequency << 4
	resp, err := m.AuxiliaryParameters(bwri)

	if err != nil {
		return -1, err
	}

	return float64(ParseValue3Bytes(resp) / 100), nil
}

// 2.5.33. ЧТЕНИЕ ВАРИАНТА ИСПОЛНЕНИЯ.
func (m *MercuryV3) Type() (*Type, error) {
	resp, err := m.ReadParameter(ParamCodeType)

	if err != nil {
		return nil, err
	}

	if err := PayloadError(resp); err != nil {
		return nil, err
	}

	return ParseType(resp), nil
}

func (m *MercuryV3) ReadArray(arr array, mo *month, t tariff) (a1, a2, r3, r4 uint64, err error) {
	code := int(arr)

	if mo != nil {
		code |= int(*mo)
	}

	var resp *Response

	resp, err = m.Request(&Request{
		Address:            m.options.address,
		Code:               RequestCodeReadArray,
		ParameterCode:      &[]byte{byte(code)}[0],
		ParameterExtension: &[]byte{byte(t)}[0],
	})

	if err != nil {
		return
	}

	if err = ResponseError(resp); err != nil {
		return
	}

	// Если поле данных ответа содержит 12 байт, то отводится по четыре двоичных байта на каждую фазу энергии А+ в последовательности:
	// - активная прямая по 1 фазе
	// - активная прямая по 2 фазе
	// - активная прямая по 3 фазе.
	//
	// Если поле данных ответа содержит 16 байт, то отводится по четыре двоичных байта на каждый вид энергии в последовательности:
	// - для кода запроса 5h:
	//   > активная прямая (А+)
	//   > активная обратная (А-)
	//   > реактивная прямая (R+)
	//   > реактивная обратная (R-)
	// - для кода запроса 15h:
	//   > реактивная R1
	//   > R2
	//   > R3
	//   > R4
	a1 = ParseValue4Bytes(resp.Payload[0:4])
	a2 = ParseValue4Bytes(resp.Payload[4:8])
	r3 = ParseValue4Bytes(resp.Payload[8:12])

	if len(resp.Payload) > 12 {
		r4 = ParseValue4Bytes(resp.Payload[12:16])
	}

	return
}
