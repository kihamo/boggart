package v3

import (
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
	address byte
	payload []byte
	crc     []byte

	lock sync.RWMutex
}

func ParseResponse(data []byte) (*Response, error) {
	if len(data) < 4 {
		return nil, errors.New("bad packet " + hex.EncodeToString(data) + " length")
	}

	return &Response{
		address: data[0],
		payload: data[1 : len(data)-2],
		crc:     data[len(data)-2:],
	}, nil
}

func (r *Response) Bytes() []byte {
	packet := append([]byte{r.Address()}, r.Payload()...)
	packet = append(packet, r.CRC()...)

	return packet
}

func (r *Response) Address() byte {
	return r.address
}

func (r *Response) Payload() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.payload...)
}

func (r *Response) CRC() []byte {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return append([]byte(nil), r.crc...)
}

func (r *Response) PayloadAsBuffer() *Buffer {
	return NewBuffer(r.Payload())
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
