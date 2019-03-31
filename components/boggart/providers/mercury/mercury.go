package mercury

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
)

// https://github.com/mrkrasser/MercuryStats
// http://incotex-support.blogspot.ru/2016/05/blog-post.html

const (
	// группа сетевых команд установки
	FunctionWriteDisplayMode = 0x09 // установка режима индикации
	FunctionWriteDisplayTime = 0x0D // установка времени индикации

	// группа сетевых команд чтения
	FunctionReadAddressGroup         = 0x20 // чтение группового адреса счетчика
	FunctionReadDatetime             = 0x21 // чтение установленных даты и времени
	FunctionReadPowerMaximum         = 0x22 // чтение лимита мощности
	FunctionReadEnergyMaximum        = 0x23 // чтение лимита энергии
	FunctionReadDaylightSavingTime   = 0x24 // чтение значения флага сезонного времени
	FunctionReadTimeCorrection       = 0x25 // чтение величины коррекции времени
	FunctionReadPowerCurrent         = 0x26 // чтение текущей мощности в нагрузке
	FunctionReadPowerCounters        = 0x27 // чтение содержимого тарифных аккумуляторов
	FunctionReadVersion              = 0x28 // чтение версии ПО
	FunctionReadBatteryVoltage       = 0x29 // чтение напряжения встроенной батарейки
	FunctionReadDisplayMode          = 0x2A // чтение режима индикации
	FunctionReadLastPowerOffDatetime = 0x2B // чтение времени последнего отключения напряжения
	FunctionReadLastPowerOnDatetime  = 0x2C // чтение времени последнего включения напряжения
	FunctionReadSerialNumber         = 0x2F // чтение серийного номера
	FunctionReadTariffCount          = 0x2E // чтение количества действующих тарифов

	// группа дополнительных сетевых команд чтения
	FunctionReadHolidays      = 0x30 // чтение таблицы праздничных дней
	FunctionReadMonthlyStat   = 0x32 // чтение месячных срезов
	FunctionReadCurrentTariff = 0x60 // чтение тарифа
	FunctionReadLastCloseCap  = 0x62 // чтение времение последнего закрытия крышки счетчика
	FunctionReadParamsCurrent = 0x63 // чтение текущих параметров сети U, I, P
	FunctionReadWordType      = 0x65 // чтение слова исполнения
	FunctionReadMakeDate      = 0x66 // чтение даты изготовления
	FunctionReadDisplayTime   = 0x67 // чтение времени индекации
	FunctionReadWorkingTime   = 0x69 // чтение времени наработки
)

const (
	displayModeTariff1 = 1 << iota
	displayModeTariff2
	displayModeTariff3
	displayModeTariff4
	displayModeAmount
	displayModePower
	displayModeTime
	displayModeDate
)

type ElectricityMeter200 struct {
	address    []byte
	location   *time.Location
	connection *rs485.Connection
}

func NewMercury(address []byte, location *time.Location, connection *rs485.Connection) *ElectricityMeter200 {
	if location == nil {
		location = time.Now().Location()
	}

	return &ElectricityMeter200{
		address:    address,
		location:   location,
		connection: connection,
	}
}

func (d *ElectricityMeter200) Address() []byte {
	return d.address
}

func (d *ElectricityMeter200) AddressGroup() ([]byte, error) {
	return d.Request(FunctionReadAddressGroup, nil)
}

func (d *ElectricityMeter200) Connection() *rs485.Connection {
	return d.connection
}

func (d *ElectricityMeter200) Request(function byte, data []byte) ([]byte, error) {
	request := []byte{0x00}

	// device address
	request = append(request, d.address...)

	// function
	request = append(request, function)

	// data in
	request = append(request, data...)

	// check sum CRC16
	request = append(request, rs485.GenerateCRC16(request)...)

	// fmt.Println("Request: ", request, hex.EncodeToString(request), " with function", strings.ToUpper(hex.EncodeToString([]byte{function})))

	response, err := d.connection.Request(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Response: ", response, hex.EncodeToString(response))

	l := len(response)
	if l < 7 {
		return nil, fmt.Errorf("error length of response packet %d", l)
	}

	// check crc16
	crc16 := rs485.GenerateCRC16(response[:l-2])
	if !bytes.Equal(response[l-2:], crc16) {
		return nil, errors.New(
			"error CRC16 of response packet have" +
				hex.EncodeToString(crc16) + " want " +
				hex.EncodeToString(response[l-2:]))
	}

	// check ADDR
	if !bytes.Equal(response[:4], request[:4]) {
		return nil, errors.New(
			"error ADDR of response packet have " +
				hex.EncodeToString(response[:4]) + " want " +
				hex.EncodeToString(request[:4]))
	}

	return response[5 : l-2], nil
}

func (d *ElectricityMeter200) responseDatetime(function byte, data []byte) (time.Time, error) {
	response, err := d.Request(function, data)
	if err != nil {
		return time.Time{}, err
	}

	// skip day of week

	hour, err := strconv.ParseInt(hex.EncodeToString(response[1:2]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	minute, err := strconv.ParseInt(hex.EncodeToString(response[2:3]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	second, err := strconv.ParseInt(hex.EncodeToString(response[3:4]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.ParseInt(hex.EncodeToString(response[4:5]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.ParseInt(hex.EncodeToString(response[5:6]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	year, err := strconv.ParseInt(hex.EncodeToString(response[6:7]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(2000+int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, d.location), nil
}

func (d *ElectricityMeter200) Datetime() (time.Time, error) {
	return d.responseDatetime(FunctionReadDatetime, nil)
}

func (d *ElectricityMeter200) SerialNumber() (int64, error) {
	response, err := d.Request(FunctionReadSerialNumber, nil)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(hex.EncodeToString(response), 16, 0)
}

func (d *ElectricityMeter200) WordType() error {
	response, err := d.Request(FunctionReadWordType, nil)
	if err != nil {
		return err
	}

	// TODO:
	_ = response

	return nil
}

func (d *ElectricityMeter200) MakeDate() (time.Time, error) {
	response, err := d.Request(FunctionReadMakeDate, nil)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.ParseInt(hex.EncodeToString(response[0:1]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.ParseInt(hex.EncodeToString(response[1:2]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	year, err := strconv.ParseInt(hex.EncodeToString(response[2:3]), 10, 0)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(2000+int(year), time.Month(month), int(day), 0, 0, 0, 0, d.location), nil
}

func (d *ElectricityMeter200) Version() (string, time.Time, error) {
	response, err := d.Request(FunctionReadVersion, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	version1, err := strconv.ParseInt(hex.EncodeToString(response[0:1]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	version2, err := strconv.ParseInt(hex.EncodeToString(response[1:2]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	version3, err := strconv.ParseInt(hex.EncodeToString(response[2:3]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	day, err := strconv.ParseInt(hex.EncodeToString(response[3:4]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	month, err := strconv.ParseInt(hex.EncodeToString(response[4:5]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	year, err := strconv.ParseInt(hex.EncodeToString(response[5:6]), 10, 0)
	if err != nil {
		return "", time.Time{}, err
	}

	return fmt.Sprintf("%d.%d.%d", version1, version2, version3),
		time.Date(2000+int(year), time.Month(month), int(day), 0, 0, 0, 0, d.location),
		nil
}

// PowerMaximum return maximum of power in W
func (d *ElectricityMeter200) PowerMaximum() (int64, error) {
	response, err := d.Request(FunctionReadPowerMaximum, nil)
	if err != nil {
		return -1, err
	}

	value, err := strconv.ParseInt(hex.EncodeToString(response), 10, 0)
	if err != nil {
		return -1, err
	}

	return value * 10, nil
}

// EnergyMaximum return maximum of energy in W/h
func (d *ElectricityMeter200) EnergyMaximum() (int64, error) {
	response, err := d.Request(FunctionReadEnergyMaximum, nil)
	if err != nil {
		return -1, err
	}

	value, err := strconv.ParseInt(hex.EncodeToString(response), 10, 0)
	if err != nil {
		return -1, err
	}

	return value * 1000, nil
}

// BatteryVoltage return voltage of battery in V
func (d *ElectricityMeter200) BatteryVoltage() (float64, error) {
	response, err := d.Request(FunctionReadBatteryVoltage, nil)
	if err != nil {
		return -1, err
	}

	v, err := strconv.ParseUint(hex.EncodeToString(response), 10, 0)
	if err != nil {
		return -1, err
	}

	return float64(v) / 100, nil
}

func (d *ElectricityMeter200) DisplayMode() (bool, bool, bool, bool, bool, bool, bool, bool, error) {
	response, err := d.Request(FunctionReadDisplayMode, nil)
	if err != nil {
		return false, false, false, false, false, false, false, false, err
	}

	bit := int(response[0])

	return bit&displayModeTariff1 != 0,
		bit&displayModeTariff2 != 0,
		bit&displayModeTariff3 != 0,
		bit&displayModeTariff4 != 0,
		bit&displayModeAmount != 0,
		bit&displayModePower != 0,
		bit&displayModeTime != 0,
		bit&displayModeDate != 0,
		nil
}

// true  / разрешает индикацию 1 тарифа
// true  / разрешает индикацию 2 тарифа
// false / разрешает индикацию 3 тарифа
// false / разрешает индикацию 4 тарифа
// true  / разрешает индикацию суммы
// false / разрешает индикацию мощности
// false / разрешает индикацию времени
// false / разрешает индикацию даты
func (d *ElectricityMeter200) SetDisplayMode(t1, t2, t3, t4, amount, power, time, date bool) error {
	bit := 0

	if t1 {
		bit |= displayModeTariff1
	}

	if t2 {
		bit |= displayModeTariff2
	}

	if t3 {
		bit |= displayModeTariff3
	}

	if t4 {
		bit |= displayModeTariff4
	}

	if amount {
		bit |= displayModeAmount
	}

	if power {
		bit |= displayModePower
	}

	if time {
		bit |= displayModeTime
	}

	if date {
		bit |= displayModeDate
	}

	_, err := d.Request(FunctionWriteDisplayMode, []byte{byte(bit)})
	return err
}

// PowerCounters returns value of T1, T2, T3 and T4 in W/h
func (d *ElectricityMeter200) PowerCounters() (uint64, uint64, uint64, uint64, error) {
	response, err := d.Request(FunctionReadPowerCounters, nil)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	values := make([]uint64, 4)
	for i := 0; i < len(response); i += 4 {
		value, err := strconv.ParseUint(hex.EncodeToString(response[i:i+4]), 10, 0)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		values[i/4] = value
	}

	return values[0] * 10, values[1] * 10, values[2] * 10, values[3] * 10, nil
}

// PowerUser return power in W
func (d *ElectricityMeter200) PowerCurrent() (uint64, error) {
	response, err := d.Request(FunctionReadPowerCurrent, nil)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(hex.EncodeToString(response), 10, 0)
}

func (d *ElectricityMeter200) DaylightSavingTime() (bool, error) {
	response, err := d.Request(FunctionReadDaylightSavingTime, nil)
	if err != nil {
		return false, err
	}

	return !bytes.Equal(response, []byte{0}), nil
}

func (d *ElectricityMeter200) TimeCorrection() (uint64, error) {
	response, err := d.Request(FunctionReadTimeCorrection, nil)
	if err != nil {
		return 0, err
	}

	// Коррекция времени в счетчике Меркурий 200 возможна в диапазоне ±30 минут в течение года.
	// Коррекция времени в счетчике Меркурий 230 возможна в диапазоне ±30 секунд в течении суток.

	return uint64(response[0]), nil
}

// ParamsCurrent returns current value of voltage in V, amperage in A, power in W
func (d *ElectricityMeter200) ParamsCurrent() (uint64, float64, uint64, error) {
	response, err := d.Request(FunctionReadParamsCurrent, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	voltage, err := strconv.ParseUint(hex.EncodeToString(response[0:2]), 10, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	amperage, err := strconv.ParseFloat(hex.EncodeToString(response[2:4]), 10)
	if err != nil {
		return 0, 0, 0, err
	}

	power, err := strconv.ParseUint(hex.EncodeToString(response[4:7]), 10, 0)
	if err != nil {
		return 0, 0, 0, err
	}

	return voltage / 10, amperage / 100, power, nil
}

func (d *ElectricityMeter200) LastPowerOffDatetime() (time.Time, error) {
	return d.responseDatetime(FunctionReadLastPowerOffDatetime, nil)
}

func (d *ElectricityMeter200) LastPowerOnDatetime() (time.Time, error) {
	return d.responseDatetime(FunctionReadLastPowerOnDatetime, nil)
}

func (d *ElectricityMeter200) LastCloseCap() (time.Time, error) {
	return d.responseDatetime(FunctionReadLastCloseCap, nil)
}

func (d *ElectricityMeter200) TariffCount() (uint64, error) {
	response, err := d.Request(FunctionReadTariffCount, nil)
	if err != nil {
		return 0, err
	}

	return uint64(response[0]), nil
}

func (d *ElectricityMeter200) Holidays() ([]time.Time, error) {
	response1, err := d.Request(FunctionReadHolidays, []byte{0})
	if err != nil {
		return nil, err
	}

	response2, err := d.Request(FunctionReadHolidays, []byte{1})
	if err != nil {
		return nil, err
	}

	response := append(response1, response2...)
	days := make([]time.Time, 0)
	year := time.Now().Year()

	for i := 0; i < len(response); i += 2 {
		if response[i] < 1 || response[i] > 31 || response[i+1] > 12 || response[i+1] < 1 {
			continue
		}

		days = append(days, time.Date(year, time.Month(response[i+1]), int(response[i]), 0, 0, 0, 0, time.UTC))
	}

	return days, nil
}

func (d *ElectricityMeter200) monthlyStat(month byte) (uint64, uint64, uint64, uint64, error) {
	response, err := d.Request(FunctionReadMonthlyStat, []byte{month})
	if err != nil {
		return 0, 0, 0, 0, err
	}

	values := make([]uint64, 4)
	for i := 0; i < len(response); i += 4 {
		value, err := strconv.ParseUint(hex.EncodeToString(response[i:i+4]), 10, 0)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		values[i/4] = value
	}

	return values[0] * 10, values[1] * 10, values[2] * 10, values[3] * 10, nil
}

func (d *ElectricityMeter200) MonthlyStat() (uint64, uint64, uint64, uint64, error) {
	// 0x0F текущий месяц, но модель 200 возвращает не корректные значения
	// поэтому лучше указывать месяц явно

	return d.monthlyStat(0x0F)
}

// значения счетчика на 1 число месяца
func (d *ElectricityMeter200) MonthlyStatByMonth(month time.Month) (uint64, uint64, uint64, uint64, error) {
	return d.monthlyStat(byte(int(month) - 1))
}

func (d *ElectricityMeter200) CurrentTariff() (uint64, error) {
	response, err := d.Request(FunctionReadCurrentTariff, nil)
	if err != nil {
		return 0, err
	}

	return uint64(response[0]), nil
}

func (d *ElectricityMeter200) DisplayTime() (uint64, uint64, uint64, uint64, error) {
	response, err := d.Request(FunctionReadDisplayTime, nil)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return uint64(response[0]), uint64(response[1]), uint64(response[2]), uint64(response[3]), nil
}

// default 10 10 5 30
// t1 / 10 / время индикации энергии не текущих тарифов и суммы
// t2 / 10 / время индикации энергии текущего тарифа
// t3 /  5 / время индикации мощности, времени и даты
// t4 / 30 / время индикации после нажатия кнопки
func (d *ElectricityMeter200) SetDisplayTime(t1, t2, t3, t4 uint64) error {
	_, err := d.Request(FunctionWriteDisplayTime, []byte{
		byte(t1),
		byte(t2),
		byte(t3),
		byte(t4),
	})
	return err
}

func (d *ElectricityMeter200) WorkingTime() (uint64, uint64, error) {
	response, err := d.Request(FunctionReadWorkingTime, nil)
	if err != nil {
		return 0, 0, err
	}

	under, err := strconv.ParseUint(hex.EncodeToString(response[0:3]), 10, 0)
	if err != nil {
		return 0, 0, err
	}

	without, err := strconv.ParseUint(hex.EncodeToString(response[3:6]), 10, 0)
	if err != nil {
		return under, 0, err
	}

	return under, without, nil
}
