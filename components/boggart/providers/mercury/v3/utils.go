package v3

import (
	"errors"
	"strconv"
)

// Меркурий 230, 231, 233, 234, 236
func ConvertSerialNumber(serial string) byte {
	// Отделите три последние цифры серийного номера, это будет число N.
	number, _ := strconv.ParseInt(serial[len(serial)-3:], 10, 0)

	// Если N>=240 адресом являются две последние цифры серийного номера.
	if number >= 240 {
		number, _ = strconv.ParseInt(serial[len(serial)-2:], 10, 0)
		return byte(number)
	}

	// Если N<240 адресом являются три последние цифры.
	if number < 240 {
		return byte(number)
	}

	return 1
}

func ResponseError(response *Response) error {
	switch response.Payload[0] {
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
