package v3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/kihamo/boggart/protocols/serial"
)

const (
	ResponseCodeOK               = 0x0 // Всё нормально
	ResponseCodeBadRequest       = 0x1 // Недопустимая команда или параметр
	ResponseCodeInternalError    = 0x2 // Внутренняя ошибка счётчика
	ResponseCodeAccessDenied     = 0x3 // Недостаточен уровень для удовлетворения запроса
	ResponseCodeTimeCorrectFiled = 0x4 // Внутренние часы счётчика уже корректировались в течение текущих суток
	ResponseCodeChannelClosed    = 0x5 // Не открыт канал связи
)

type Response struct {
	payload []byte
	lock    sync.RWMutex
	crc     uint16
	address uint8
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Address() uint8 {
	return r.address
}

func (r *Response) Payload() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.payload...)
}

func (r *Response) PayloadAsBuffer() *Buffer {
	return NewBuffer(r.Payload())
}

func (r *Response) CRC() uint16 {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.crc
}

func (r *Response) MarshalBinary() (_ []byte, err error) {
	packet := append([]byte{r.Address()}, r.Payload()...)

	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, r.CRC())
	packet = append(packet, crc...)

	return packet, nil
}

func (r *Response) UnmarshalBinary(data []byte) (err error) {
	l := len(data)

	if l < 4 {
		return errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	r.address = data[0]
	r.payload = data[1 : l-2]
	r.crc = binary.LittleEndian.Uint16(data[l-2:])

	// check crc
	crc := serial.GenerateCRC16(data[:l-2])
	if !bytes.Equal(crc, data[l-2:]) {
		return errors.New("wrong CRC16 of response packet " +
			hex.EncodeToString(data) + " have " +
			hex.EncodeToString(data[l-2:]) + " want " +
			hex.EncodeToString(crc))
	}

	return nil
}

func (r *Response) String() string {
	data, _ := r.MarshalBinary()
	return hex.EncodeToString(data)
}

func (r *Response) GetError() error {
	switch r.Payload()[0] {
	case ResponseCodeBadRequest:
		return errors.New("bad request")
	case ResponseCodeInternalError:
		return errors.New("internal error")
	case ResponseCodeAccessDenied:
		return errors.New("access denied")
	case ResponseCodeTimeCorrectFiled:
		return errors.New("correct time failed")
	case ResponseCodeChannelClosed:
		return errors.New("channel is closed")
	}

	return nil
}
