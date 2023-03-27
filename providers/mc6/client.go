package mc6

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/kihamo/boggart/protocols/modbus"
)

const (
	AddressRoomTemperature           uint16 = 0
	AddressFloorTemperature          uint16 = 1
	AddressHumidity                  uint16 = 2
	AddressHeatingValve              uint16 = 3
	AddressCoolingValve              uint16 = 4
	AddressHeatingOutput             uint16 = 9
	AddressWindowsOpen               uint16 = 13
	AddressHoldingFunction           uint16 = 15
	AddressFloorOverheat             uint16 = 17
	AddressDeviceType                uint16 = 18
	AddressFanSpeedNumbers           uint16 = 20
	AddressSystemError               uint16 = 21
	AddressTemperatureFormat         uint16 = 60
	AddressStatus                    uint16 = 61
	AddressSystemMode                uint16 = 62
	AddressFanSpeed                  uint16 = 63
	AddressTargetTemperature         uint16 = 64
	AddressAway                      uint16 = 65
	AddressAwayTemperature           uint16 = 66
	AddressHoldingTime               uint16 = 67
	AddressHoldingTemperatureAndTime uint16 = 68
	AddressHoldingTemperature        uint16 = 69
	AddressPanelLock                 uint16 = 78
	AddressPanelLockPin1             uint16 = 79
	AddressPanelLockPin2             uint16 = 80
	AddressPanelLockPin3             uint16 = 81
	AddressPanelLockPin4             uint16 = 82
	AddressTargetTemperatureMaximum  uint16 = 83
	AddressTargetTemperatureMinimum  uint16 = 84
	AddressFloorTemperatureLimit     uint16 = 85

	TemperatureFormatCelsius    uint16 = 0
	TemperatureFormatFahrenheit uint16 = 1

	FanSpeedMode1 uint16 = 0
	FanSpeedMode3 uint16 = 1

	FanSpeedHigh   uint16 = 0
	FanSpeedMedium uint16 = 1
	FanSpeedLow    uint16 = 2
	FanSpeedAuto   uint16 = 3

	SystemModeHeat       uint16 = 0
	SystemModeCool       uint16 = 1
	SystemModeVent       uint16 = 2
	SystemModeDehumidity uint16 = 3
	SystemModeAuto       uint16 = 4

	writeResponseSuccess uint16 = 2
)

type MC6 struct {
	client *modbus.Client
}

func New(address *url.URL, opts ...modbus.Option) *MC6 {
	address.Scheme = "tcp"

	return &MC6{
		client: modbus.NewClient(address, opts...),
	}
}

func (m *MC6) Close() error {
	return m.client.Close()
}

func (m *MC6) ReadTemperature(address uint16) (float64, error) {
	value, err := m.client.ReadHoldingRegistersUint16(address)
	if err != nil {
		return 0, err
	}

	// по цельсию от 0 до 500
	// по фарингейту от 320 до 1220

	// TODO: выставить определения шкалы, чтобы валидатор корректно срабатывал
	if value > 500 {
		return 0, fmt.Errorf("temperature sensor returned wrong value %d", value)
	}

	return float64(value) / 10, err
}

func (m *MC6) ReadTemperatureUint(address uint16) (uint16, error) {
	value, err := m.ReadTemperature(address)
	if err != nil {
		return 0, err
	}

	return uint16(value), err
}

func (m *MC6) ReadDuration(address uint16) (time.Duration, error) {
	value, err := m.client.ReadHoldingRegistersUint16(address)
	if err != nil {
		return 0, err
	}

	return time.Duration(value) * time.Minute, err
}

func (m *MC6) Write(address, quantity uint16, payload []byte) (err error) {
	response, err := m.client.WriteMultipleRegisters(address, quantity, payload)

	if err == nil {
		code := binary.BigEndian.Uint16(response)

		if code != writeResponseSuccess {
			err = fmt.Errorf("device return not success response %d", code)
		}
	}

	return err
}

func (m *MC6) WriteUint16(address, value uint16) error {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, value)

	return m.Write(address, 2, payload)
}

func (m *MC6) WriteUint32(address uint16, value uint32) error {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, value)

	return m.Write(address, 2, payload)
}

func (m *MC6) WriteBool(address uint16, flag bool) error {
	var value uint16

	if flag {
		value = 1
	}

	return m.WriteUint16(address, value)
}

func (m *MC6) WriteTemperature(address uint16, value float64) error {
	value = m.RoundTemperature(value) * 10

	if value < 50 || value > 350 {
		return errors.New("wrong temperature value 50 >= value <= 350")
	}

	return m.WriteUint16(address, uint16(value))
}

// устанавливаемое значение всегда кратно 0.5 и округляется в меньшую сторону
// даже на устройстве шаг 0.5, поэтому принудительно округляем
func (m *MC6) RoundTemperature(value float64) float64 {
	val := int(value * 10)
	val -= val % 5
	return float64(val) / 10
}
