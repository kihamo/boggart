package v3

import (
	"time"
)

func (d *MercuryV3) ReadParameter(param byte) ([]byte, error) {
	resp, err := d.Request(&Request{
		Address:       d.address,
		Code:          RequestCodeReadParameter,
		ParameterCode: &[]byte{param}[0],
	})

	return resp.Payload, err
}

// 2.5.18. ЧТЕНИЕ СЕРИЙНОГО НОМЕРА СЧЕТЧИКА И ДАТЫ ВЫПУСКА.
func (d *MercuryV3) SerialNumberAndBuildDate() (serialNumber string, buildDate time.Time, err error) {
	resp, err := d.ReadParameter(ParamCodeSerialNumberAndBuildDate)

	if err != nil {
		return
	}

	serialNumber = ParseSerialNumber(resp[0:4])
	buildDate = ParseBuildDate(resp[4:7])

	return
}

// 2.5.19. УСКОРЕННЫЙ РЕЖИМ ЧТЕНИЯ ИНДИВИДУАЛЬНЫХ ПАРАМЕТРОВ ПРИБОРА.
func (d *MercuryV3) ForceReadParameters() (serialNumber string, buildDate time.Time, firmwareVersion string, t *Type, err error) {
	resp, err := d.ReadParameter(ParamCodeForceReadParameters)

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
func (d *MercuryV3) TransformationCoefficient() (uint8, uint8, error) {
	resp, err := d.ReadParameter(ParamCodeTransformationCoefficient)

	if err != nil {
		return 0, 0, err
	}

	return resp[0]*10 + resp[1], resp[2]*10 + resp[3], nil
}

// 2.5.21. ЧТЕНИЕ ВЕРСИИ ПО СЧЁТЧИКА.
func (d *MercuryV3) FirmwareVersion() (string, error) {
	resp, err := d.ReadParameter(ParamCodeVersion)

	if err != nil {
		return "", err
	}

	return ParseFirmwareVersion(resp), nil
}

// 2.5.23. ЧТЕНИЕ СЕТЕВОГО АДРЕСА.
func (d *MercuryV3) Address() (byte, error) {
	resp, err := d.ReadParameter(ParamCodeAddress)

	if err != nil {
		return 0, err
	}

	if err := PayloadError(resp); err != nil {
		return 0, err
	}

	return resp[1], nil
}

// 2.5.32. ЧТЕНИЕ ВСПОМОГАТЕЛЬНЫХ ПАРАМЕТРОВ.
func (d *MercuryV3) AuxiliaryParameters(bwri byte) ([]byte, error) {
	resp, err := d.Request(&Request{
		Address:            d.address,
		Code:               RequestCodeReadParameter,
		ParameterCode:      &[]byte{ParamCodeAuxiliaryParameters}[0],
		ParameterExtension: &[]byte{bwri}[0],
	})

	if err != nil {
		return nil, err
	}

	if err := ResponseError(resp); err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// 2.5.33. ЧТЕНИЕ ВАРИАНТА ИСПОЛНЕНИЯ.
func (d *MercuryV3) Type() (*Type, error) {
	resp, err := d.ReadParameter(ParamCodeType)

	if err != nil {
		return nil, err
	}

	if err := PayloadError(resp); err != nil {
		return nil, err
	}

	return ParseType(resp), nil
}

func (d *MercuryV3) ReadArray(a array, m *month, t tariff) (a1, a2, r3, r4 uint64, err error) {
	code := int(a)

	if m != nil {
		code |= int(*m)
	}

	var resp *Response

	resp, err = d.Request(&Request{
		Address:            d.address,
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
	a1 = ParseArrayValue(resp.Payload[0:4])
	a2 = ParseArrayValue(resp.Payload[4:8])
	r3 = ParseArrayValue(resp.Payload[8:12])

	if len(resp.Payload) > 12 {
		r4 = ParseArrayValue(resp.Payload[12:16])
	}

	return
}
