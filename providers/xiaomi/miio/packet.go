package miio

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

/*
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | 0x2131 / Magic bytes                   | 0x0020 / Length                        |
 |-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
 | 0x00000000 / Unknown (0хFFFFFFFF for Hello, 0х00000000 for others               |
 |-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
 | 0x0020 / Device type                   | 0x0020 / Device ID                     |
 |-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
 | Timestamp                                                                       |
 |-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|
 | Token (128-bit)                                                                 |
 | All subsequent encryption is based on this number.                              |
 |-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-|

Magic - "магическое" число, всегда равно 0х2131 (2 байта).
Length - длина пакета в байтах(заголовок+данные) (2 байта).
Unknown - поле неизвестного назначения. Всегда заполнено нулями 0х00000000, а у hello-пакета 0хFFFFFFFF (4 байта).
Device type - тип устройства (2 байта).
Device ID - идентификатор устройства (2 байта).
Time stamp - временная отметка, время работы устройства в секундах (4 байта).
Checksum - контрольная сумма всего пакета по алгоритму MD5. Перед расчетом КС это поле временно заполняется нулями (16 байт).
Data - полезная нагрузка произвольной длины - зашифрованные данные, отправляемые устройству. В hello-пакете это поле отсутствует.
*/

var magicBytes = []byte{0x21, 0x31}

type Packet struct {
	magic      []byte
	length     uint16
	unknown    uint32
	deviceType uint16
	deviceID   uint16
	timestamp  uint32
	checksum   []byte
	payload    []byte
	token      []byte
}

func NewPacket(token []byte) *Packet {
	return &Packet{
		magic: magicBytes,
		token: token,
	}
}

func (p *Packet) SetUnknown(v uint32) {
	p.unknown = v
}

func (p *Packet) DeviceType() uint16 {
	return p.deviceType
}

func (p *Packet) SetDeviceType(t uint16) {
	p.deviceType = t
}

func (p *Packet) DeviceID() uint16 {
	return p.deviceID
}

func (p *Packet) SetDeviceID(id uint16) {
	p.deviceID = id
}

func (p *Packet) Timestamp() time.Time {
	return time.Unix(int64(p.timestamp), 0)
}

func (p *Packet) SetTimestamp(ts time.Time) {
	p.timestamp = uint32(ts.Unix())
}

func (p *Packet) Checksum() []byte {
	return p.checksum
}

func (p *Packet) Payload() []byte {
	return append([]byte(nil), p.payload...)
}

func (p *Packet) PayloadJSON(v interface{}) error {
	return json.Unmarshal(p.Payload(), v)
}

func (p *Packet) SetPayload(payload []byte) {
	p.payload = payload
}

func (p *Packet) SetPayloadRPC(id uint32, method string, params interface{}) {
	if params == nil {
		params = []interface{}{}
	}

	payload, _ := json.Marshal(Request{
		ID:     id,
		Method: method,
		Params: params,
	})

	p.SetPayload(payload)
}

func (p *Packet) MarshalBinary() (_ []byte, err error) {
	request := bytes.NewBuffer(p.magic)
	buf := bytes.NewBuffer(nil)

	// Packet length: 16 bits unsigned int
	p.length = uint16(len(p.Payload())) + 32
	if err = binary.Write(buf, binary.LittleEndian, p.length); err != nil {
		return nil, err
	}

	if _, err = request.Write(reverse(buf.Bytes())); err != nil {
		return nil, err
	}

	// Unknown1: 32 bits
	if err = binary.Write(request, binary.LittleEndian, p.unknown); err != nil {
		return nil, err
	}

	// Device type: 16 bits
	if err = binary.Write(request, binary.LittleEndian, p.deviceType); err != nil {
		return nil, err
	}

	// Device ID: 16 bits
	if err = binary.Write(request, binary.LittleEndian, p.deviceID); err != nil {
		return nil, err
	}

	// Timestamp: 32 bits
	buf.Reset()

	if err = binary.Write(buf, binary.LittleEndian, p.timestamp); err != nil {
		return nil, err
	}

	if _, err = request.Write(reverse(buf.Bytes())); err != nil {
		return nil, err
	}

	payload := p.Payload()

	// MD5 checksum
	if len(p.checksum) == 0 {
		if p.unknown != 0xFFFFFFFF {
			hash := md5.New()

			// Перед расчетом КС это поле временно заполняется нулями (16 байт).
			// на самом деле не нулями, а текущим токеном
			// data := append(request.Bytes(), bytes.Repeat([]byte{0}, 16)...)
			data := append(request.Bytes(), p.token...)
			if len(payload) > 0 {
				data = append(data, payload...)
			}

			if _, err = hash.Write(data); err != nil {
				return nil, err
			}

			p.checksum = hash.Sum(nil)
		} else {
			p.checksum = bytes.Repeat([]byte{0xFF}, 16)
		}
	}

	if _, err = request.Write(p.checksum); err != nil {
		return nil, err
	}

	// Payload
	if len(payload) > 0 {
		_, err = request.Write(payload)
	}

	return request.Bytes(), err
}

func (p *Packet) UnmarshalBinary(data []byte) (err error) {
	if len(data) < 32 {
		return errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	if magic := data[:2]; !bytes.Equal(magicBytes, magic) {
		return errors.New("magic number could not be verified. Expected 0x2131, got 0x" + hex.EncodeToString(magic))
	}

	// Packet length: 16 bits unsigned int
	p.length = binary.LittleEndian.Uint16(reverse(data[2:4]))

	response := bytes.NewBuffer(data[4:])

	// Unknown1: 32 bits
	if err = binary.Read(response, binary.LittleEndian, &p.unknown); err != nil {
		return err
	}

	// Device type: 16 bits
	if err = binary.Read(response, binary.LittleEndian, &p.deviceType); err != nil {
		return err
	}

	// Device ID: 16 bits
	if err = binary.Read(response, binary.LittleEndian, &p.deviceID); err != nil {
		return err
	}

	// Timestamp: 32 bits
	p.timestamp = binary.LittleEndian.Uint32(reverse(response.Next(4)))

	// MD5 checksum
	p.checksum = response.Next(16)

	// Payload
	p.payload = response.Bytes()

	return err
}

func reverse(data []byte) []byte {
	dataNew := make([]byte, len(data))

	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		dataNew[i], dataNew[j] = data[j], data[i]
	}

	return dataNew
}
