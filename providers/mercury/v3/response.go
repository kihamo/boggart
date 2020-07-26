package v3

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"
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
	address uint8
	payload []byte
	crc     uint16

	lock sync.RWMutex
}

func ParseResponse(data []byte) (*Response, error) {
	if len(data) < 4 {
		return nil, errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	return &Response{
		address: data[0],
		payload: data[1 : len(data)-2],
		crc:     binary.LittleEndian.Uint16(data[len(data)-2:]),
	}, nil
}

func (r *Response) Bytes() []byte {
	packet := append([]byte{r.Address()}, r.Payload()...)

	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, r.CRC())
	packet = append(packet, crc...)

	return packet
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

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
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
