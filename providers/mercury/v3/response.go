package v3

import (
	"encoding/hex"
	"errors"
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
	Address byte
	Payload []byte
	CRC     []byte
}

func ParseResponse(data []byte) (*Response, error) {
	if len(data) < 4 {
		return nil, errors.New("bad packet length")
	}

	return &Response{
		Address: data[0],
		Payload: data[1 : len(data)-2],
		CRC:     data[len(data)-2:],
	}, nil
}

func (r *Response) Bytes() []byte {
	packet := append([]byte{r.Address}, r.Payload...)
	packet = append(packet, r.CRC...)

	return packet
}

func (r *Response) String() string {
	return hex.EncodeToString(r.Bytes())
}
