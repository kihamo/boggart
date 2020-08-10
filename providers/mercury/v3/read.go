package v3

import (
	"errors"
	"fmt"
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

	AuxiliaryPower        uint8 = 0x0
	AuxiliaryVoltage      uint8 = 0x1
	AuxiliaryAmperage     uint8 = 0x2
	AuxiliaryPowerFactors uint8 = 0x3
	AuxiliaryFrequency    uint8 = 0x4
	AuxiliaryPhasesAngle  uint8 = 0x5

	PowerNumberP powerNumber = 0x0
	PowerNumberQ powerNumber = 0x1
	PowerNumberS powerNumber = 0x2

	PhaseNumberAll phase = 0x0
	PhaseNumber1   phase = 0x1
	PhaseNumber2   phase = 0x2
	PhaseNumber3   phase = 0x3
)

func (m *MercuryV3) ReadParameter(param byte) (*Buffer, error) {
	request := NewRequest().
		WithCode(CodeReadParameter).
		WithParameterCode(param)

	response, err := m.Invoke(request)
	if err != nil {
		return nil, err
	}

	return response.PayloadAsBuffer(), err
}

// 2.5.18. ЧТЕНИЕ СЕРИЙНОГО НОМЕРА СЧЕТЧИКА И ДАТЫ ВЫПУСКА.
func (m *MercuryV3) SerialNumberAndBuildDate() (serialNumber string, buildDate time.Time, err error) {
	value, err := m.ReadParameter(ParamCodeSerialNumberAndBuildDate)
	if err != nil {
		return
	}

	serialNumber = value.ReadSerialNumber()
	buildDate = value.ReadBuildDate()

	return
}

// 2.5.19. УСКОРЕННЫЙ РЕЖИМ ЧТЕНИЯ ИНДИВИДУАЛЬНЫХ ПАРАМЕТРОВ ПРИБОРА.
func (m *MercuryV3) ForceReadParameters() (serialNumber string, buildDate time.Time, firmwareVersion string, t *Type, err error) {
	value, err := m.ReadParameter(ParamCodeForceReadParameters)
	if err != nil {
		return
	}

	serialNumber = value.ReadSerialNumber()
	buildDate = value.ReadBuildDate()
	firmwareVersion = value.ReadFirmwareVersion()
	// TODO: t = ParseType(value.Next(6))

	return
}

// 2.5.20. ЧТЕНИЕ КОЭФФИЦИЕНТА ТРАНСФОРМАЦИИ СЧЁТЧИКА.
func (m *MercuryV3) TransformationCoefficient() (uint8, uint8, error) {
	value, err := m.ReadParameter(ParamCodeTransformationCoefficient)
	if err != nil {
		return 0, 0, err
	}

	return value.ReadUint8()*10 + value.ReadUint8(), value.ReadUint8()*10 + value.ReadUint8(), nil
}

// 2.5.21. ЧТЕНИЕ ВЕРСИИ ПО СЧЁТЧИКА.
func (m *MercuryV3) FirmwareVersion() (string, error) {
	value, err := m.ReadParameter(ParamCodeVersion)
	if err != nil {
		return "", err
	}

	return value.ReadFirmwareVersion(), nil
}

// 2.5.23. ЧТЕНИЕ СЕТЕВОГО АДРЕСА.
func (m *MercuryV3) Address() (uint8, error) {
	value, err := m.ReadParameter(ParamCodeAddress)
	if err != nil {
		return 0, err
	}

	if value.ReadUint8() != ResponseCodeOK {
		return 0, errors.New("first byte not equal 0")
	}

	return value.ReadUint8(), nil
}

// 2.5.32. ЧТЕНИЕ ВСПОМОГАТЕЛЬНЫХ ПАРАМЕТРОВ.
func (m *MercuryV3) AuxiliaryParameters(bwri uint8) (*Response, error) {
	request := NewRequest().
		WithCode(CodeReadParameter).
		WithParameterCode(ParamCodeAuxiliaryParameters12).
		WithParameterExtension(bwri)

	response, err := m.Invoke(request)
	if err != nil {
		return nil, err
	}

	if err = response.GetError(); err != nil {
		return nil, err
	}

	return response, nil
}

// 2.5.32.1 Ответ прибора на запрос чтения мощности.
func (m *MercuryV3) Power(number powerNumber) (sum, phase1, phase2, phase3 float64, err error) {
	var response *Response

	response, err = m.AuxiliaryParameters(AuxiliaryPower<<4 | uint8(number)<<2 | uint8(PhaseNumberAll))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		dataOut.ReadByte() // направление мощности
		sum = float64(dataOut.ReadUint16()) / 100

		dataOut.ReadByte() // направление мощности
		phase1 = float64(dataOut.ReadUint16()) / 100

		dataOut.ReadByte() // направление мощности
		phase2 = float64(dataOut.ReadUint16()) / 100

		dataOut.ReadByte() // направление мощности
		phase3 = float64(dataOut.ReadUint16()) / 100
	}

	return
}

// 2.5.32.2 Ответ прибора на запрос чтения фазного и линейного напряжения, тока и углов между фазными напряжениями.
func (m *MercuryV3) Voltage() (phase1, phase2, phase3 float64, err error) {
	var response *Response

	response, err = m.AuxiliaryParameters(AuxiliaryVoltage<<4 | 0x1)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		phase1 = float64(dataOut.ReadUint32By3Byte()) / 100
		phase2 = float64(dataOut.ReadUint32By3Byte()) / 100
		phase3 = float64(dataOut.ReadUint32By3Byte()) / 100
	}

	return
}

func (m *MercuryV3) Amperage() (phase1, phase2, phase3 float64, err error) {
	var response *Response

	response, err = m.AuxiliaryParameters(AuxiliaryAmperage<<4 | 0x1)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		phase1 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase2 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase3 = float64(dataOut.ReadUint32By3Byte()) / 1000
	}

	return
}

func (m *MercuryV3) PhasesAngle() (phase1, phase2, phase3 float64, err error) {
	var response *Response

	response, err = m.AuxiliaryParameters(AuxiliaryPhasesAngle<<4 | 0x1)
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		phase1 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase2 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase3 = float64(dataOut.ReadUint32By3Byte()) / 1000
	}

	return
}

// 2.5.32.3 Ответ прибора на запрос чтения коэффициентов мощности.
func (m *MercuryV3) PowerFactors() (sum, phase1, phase2, phase3 float64, err error) {
	var response *Response

	response, err = m.AuxiliaryParameters(AuxiliaryPowerFactors<<4 | uint8(PhaseNumberAll))
	if err == nil {
		dataOut := response.PayloadAsBuffer()

		sum = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase1 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase2 = float64(dataOut.ReadUint32By3Byte()) / 1000
		phase3 = float64(dataOut.ReadUint32By3Byte()) / 1000
	}

	return
}

// 2.5.32.4 Ответ прибора на запрос чтения частоты (запрос 11h, 14h, 16h).
func (m *MercuryV3) Frequency() (float64, error) {
	response, err := m.AuxiliaryParameters(AuxiliaryFrequency << 4)
	if err != nil {
		return -1, err
	}

	dataOut := response.PayloadAsBuffer()

	return float64(dataOut.ReadUint32By3Byte() / 100), nil
}

// 2.5.33. ЧТЕНИЕ ВАРИАНТА ИСПОЛНЕНИЯ.
func (m *MercuryV3) Type() (*Type, error) {
	_, err := m.ReadParameter(ParamCodeType)
	if err != nil {
		return nil, err
	}

	// TODO:
	return nil, nil
}

// 2.5.16.1 Запросы на чтение массивов в пределах 12 месяцев
func (m *MercuryV3) ReadArray(arr array, mo *month, t tariff) (a1, a2, r3, r4 uint32, err error) {
	code := uint8(arr)

	if mo != nil {
		code |= uint8(*mo)
	}

	var response *Response

	request := NewRequest().
		WithCode(CodeReadArray).
		WithParameterCode(code).
		WithParameterExtension(uint8(t))

	response, err = m.Invoke(request)
	if err != nil {
		return
	}

	dataOut := response.PayloadAsBuffer()

	// check length of response
	switch v := dataOut.Len(); {
	case v == 1:
		if err = response.GetError(); err == nil {
			err = fmt.Errorf("response payload length must 12 or 16 bytes not %d", dataOut.Len())
		}

		return
	case arr == ArrayActiveEnergy && v != 12:
		err = fmt.Errorf("response payload length must 12 bytes not %d", dataOut.Len())
		return
	case v != 16:
		err = fmt.Errorf("response payload length must 16 bytes not %d", dataOut.Len())
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
	a1 = dataOut.ReadUint32()
	a2 = dataOut.ReadUint32()
	r3 = dataOut.ReadUint32()

	if dataOut.Len() > 12 {
		r4 = dataOut.ReadUint32()
	}

	return a1, a2, r3, r4, err
}
