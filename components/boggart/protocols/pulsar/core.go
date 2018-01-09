package pulsar

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
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

	ParamDaylightSavingTime = 0x0001
	ParamVersion            = 0x0005
	ParamDiagnostics        = 0x0006

	Channel1  = 0x00000001
	Channel2  = 0x00000002
	Channel3  = 0x00000004 // температура подачи [°C]
	Channel4  = 0x00000008 // температура обратки [°C]
	Channel5  = 0x00000010 // перепад температур [°C]
	Channel6  = 0x00000020 // мощность [Гкал/ч]
	Channel7  = 0x00000040 // энергия [Гкал]
	Channel8  = 0x00000080 // объем [м^3]
	Channel9  = 0x00000100 // расход [м^3/ч]
	Channel10 = 0x00000200 // импульсный вход 1 [м^3]
	Channel11 = 0x00000400 // импульсный вход 2 [м^3]
	Channel12 = 0x00000800 // импульсный вход 3 [м^3]
	Channel13 = 0x00001000 // импульсный вход 4 [м^3]
	Channel14 = 0x00002000 // расход по энергии [м3/ч]
	Channel20 = 0x00080000 // Время нормальной работы [ч]
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

func (p *Pulsar) ReadMetrics(channel int64) ([]byte, error) {
	bs := p.pad(big.NewInt(channel).Bytes(), 4)
	response, err := p.Request(FunctionReadMetrics, bs)
	if err != nil {
		return nil, err
	}

	value := p.reverse(response)

	//fmt.Println(math.Float32frombits(binary.BigEndian.Uint32(value)))

	return value, nil
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
	return p.Request(FunctionReadSettings, bs)
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

func (p *Pulsar) HexString(data []byte) string {
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

	fmt.Println("Request: ", request, p.HexString(request))

	response, err := p.RequestRaw(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response: ", response, p.HexString(response))

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
