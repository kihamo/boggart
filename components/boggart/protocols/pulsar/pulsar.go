package pulsar

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jacobsa/go-serial/serial"
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
)

type Pulsar struct {
	config  serial.OpenOptions
	address []byte
	port    string
}

func NewPulsar(port string, address []byte) (*Pulsar, error) {
	return &Pulsar{
		address: address,
		port:    port,
	}, nil
}

func (p *Pulsar) ReadTime() (time.Time, error) {
	response, err := p.Request(FunctionReadTime, nil, 6)
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

func (p *Pulsar) ReadSettings(param int64) (interface{}, error) {
	bs := append(big.NewInt(param).Bytes(), 0)

	response, err := p.Request(FunctionReadSettings, bs, 8)
	if err != nil {
		return -1, err
	}

	fmt.Println(response)

	return response, nil
}

func (p *Pulsar) generateRequestId() []byte {
	id, _ := rand.Int(rand.Reader, big.NewInt(0xFFFF))
	return id.Bytes()
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

func (p *Pulsar) HexString(data []byte) string {
	var s string

	for _, b := range data {
		s += fmt.Sprintf("%02x ", b)
	}

	return s[:len(s)-1]
}

func (p *Pulsar) RequestRaw(request []byte, size uint) ([]byte, error) {
	serial, err := serial.Open(serial.OpenOptions{
		PortName:              p.port,
		BaudRate:              9600,
		DataBits:              8,
		ParityMode:            serial.PARITY_NONE,
		StopBits:              1,
		InterCharacterTimeout: 0,
		MinimumReadSize:       size,
	})

	if err != nil {
		return nil, err
	}
	defer serial.Close()

	if _, err := serial.Write(request); err != nil {
		return nil, err
	}

	buffer := make([]byte, 128)
	n, err := serial.Read(buffer)

	return buffer[:n], err
}

func (p *Pulsar) Request(function byte, data []byte, size uint) ([]byte, error) {
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

	response, err := p.RequestRaw(request, 4+1+1+size+2+2)
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
	if bytes.Compare(response[l-4:l-2], requestId) != 0 {
		return nil, errors.New("Error ID of response packet")
	}

	// check error
	if response[4] == FunctionBadCommand {
		return nil, fmt.Errorf("Device returns error code #%d", response[6])
	}

	return response[6 : l-4], nil
}
