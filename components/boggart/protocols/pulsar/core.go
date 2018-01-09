package pulsar

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/goburrow/serial"
)

const (
	FunctionBadCommand    = 0x00
	FunctionReadMetrics   = 0x01
	FunctionReadTime      = 0x04
	FunctionWriteTime     = 0x05
	FunctionReadArchive   = 0x06
	FunctionReadSettings  = 0x0A
	FunctionWriteSettings = 0x0B

	ParamDaylightSavingTime = 0x0001 // uint16  | RW |    | признак автоперехода на летнее время 0 - выкл, 1 - вкл
	ParamPulseDuration      = 0x0003 // float32 | RW | мс | длительность импульса
	ParamPauseDuration      = 0x0004 // float32 | RW | мс | длительность паузы
	ParamVersion            = 0x0005 // uint16  | R  |    | версия прошивки
	ParamDiagnostics        = 0x0006
	ParamOperatingTime      = 0x000C // uint32  | RW | ч  | время наработки

	Channel3  = 0x00000004 // float32 | °C     | температура подачи
	Channel4  = 0x00000008 // float32 | °C     | температура обратки
	Channel5  = 0x00000010 // float32 | °C     | перепад температур
	Channel6  = 0x00000020 // float32 | гкал/ч | мощность
	Channel7  = 0x00000040 // float32 | гкал   | энергия
	Channel8  = 0x00000080 // float32 | м^3    | объем
	Channel9  = 0x00000100 // float32 | м^3/ч  | расход
	Channel10 = 0x00000200 // float32 | м^3    | импульсный вход 1
	Channel11 = 0x00000400 // float32 | м^3    | импульсный вход 2
	Channel12 = 0x00000800 // float32 | м^3    | импульсный вход 3
	Channel13 = 0x00001000 // float32 | м^3    | импульсный вход 4
	Channel14 = 0x00002000 // float32 | м3/ч   | расход по энергии === Channel6
	Channel20 = 0x00080000 // uint32  | ч      | время нормальной работы
)

type Pulsar struct {
	config  *serial.Config
	address []byte
	port    string
}

func NewPulsar(port string, address []byte) (*Pulsar, error) {
	return &Pulsar{
		address: address,
		config: &serial.Config{
			Address:  port,
			BaudRate: 9600,
			Parity:   "N",
			Timeout:  time.Second,
		},
		port: port,
	}, nil
}

func (p *Pulsar) DeviceAddress() ([]byte, error) {
	response, err := p.RequestRaw([]byte{0xF0, 0x0F, 0x0F, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA5, 0x44})
	if err != nil {
		return nil, err
	}

	return response[4:8], nil
}

func (p *Pulsar) ReadMetrics(channel int64) ([][]byte, error) {
	bs := p.pad(big.NewInt(channel).Bytes(), 4)
	response, err := p.Request(FunctionReadMetrics, bs)
	if err != nil {
		return nil, err
	}

	metrics := math.Ceil(float64(len(response) / 4))
	result := make([][]byte, 0, int64(metrics))
	value := p.reverse(response)
	for i := 0; i < len(value); i += 4 {
		result = append(result, value[i:i+4])
	}

	return result, nil
}

func (p *Pulsar) ReadTime() (time.Time, error) {
	response, err := p.Request(FunctionReadTime, nil)
	if err != nil {
		return time.Now(), err
	}

	return time.Date(
		2000+int(response[0]),
		time.Month(response[1]),
		int(response[2]),
		int(response[3]),
		int(response[4]),
		int(response[5]),
		0,
		time.Now().Location()), nil
}

func (p *Pulsar) ReadSettings(param int64) ([]byte, error) {
	bs := p.pad(big.NewInt(param).Bytes(), 2)
	response, err := p.Request(FunctionReadSettings, bs)
	if err != nil {
		return nil, err
	}

	return p.reverse(response), nil
}

func (p *Pulsar) generateRequestId() []byte {
	id, _ := rand.Int(rand.Reader, big.NewInt(0xFFFF))
	return p.pad(id.Bytes(), 2)
}

func (p *Pulsar) generateCRC16(packet []byte) []byte {
	result := 0xFFFF

	for i := 0; i < len(packet); i++ {
		result = ((result << 8) >> 8) ^ int(packet[i])
		for j := 0; j < 8; j++ {
			flag := result & 0x0001
			result >>= 1
			if flag == 1 {
				result ^= 0xA001
			}
		}
	}

	return p.reverse(big.NewInt(int64(result)).Bytes())
}

func (p *Pulsar) reverse(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data
}

func (p *Pulsar) pad(data []byte, n int) []byte {
	if len(data) >= n {
		return data
	}

	for i := len(data); i < n; i++ {
		data = append(data, 0x0)
	}

	return data
}

func (p *Pulsar) ToFloat32(data []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(data))
}

func (p *Pulsar) ToString(data []byte) string {
	var s string

	for _, b := range data {
		s += fmt.Sprintf("%02x ", b)
	}

	if len(s) > 0 {
		return s[:len(s)-1]
	}

	return ""
}

func (p *Pulsar) RequestRaw(request []byte) ([]byte, error) {
	port, err := serial.Open(p.config)

	if err != nil {
		return nil, err
	}
	defer port.Close()

	if _, err := port.Write(request); err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	for {
		b := make([]byte, 512)
		n, err := port.Read(b)
		if err != nil {
			break
		}

		if n != 0 {
			buffer.Write(b[:n])
		}
	}

	return buffer.Bytes(), err
}

func (p *Pulsar) Request(function byte, data []byte) ([]byte, error) {
	var request []byte

	// device address
	request = append(request, p.address...)

	// function
	request = append(request, function)

	// length of packet
	l := len(request) + 1 + len(data) + 2 + 2
	request = append(request, byte(l))

	// data in
	request = append(request, data...)

	// request id
	requestId := p.generateRequestId()
	request = append(request, requestId...)

	// check sum CRC16
	request = append(request, p.generateCRC16(request)...)

	// fmt.Println("Request: ", request, p.ToString(request))

	response, err := p.RequestRaw(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Response: ", response, p.ToString(response))

	l = len(response)
	if l < 10 {
		return nil, errors.New("Error length of response packet")
	}

	// check crc16
	crc16 := p.generateCRC16(response[:l-2])
	if bytes.Compare(response[l-2:], crc16) != 0 {
		return nil, errors.New("Error CRC16 of response packet")
	}

	// check id
	if bytes.Compare(response[l-(2+len(requestId)):l-2], requestId) != 0 {
		return nil, errors.New("Error ID of response packet")
	}

	// check error
	if response[4] == FunctionBadCommand {
		return nil, fmt.Errorf("Device returns error code #%d", response[6])
	}

	return response[6 : l-4], nil
}
