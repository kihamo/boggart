package mercury

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

func ConvertSerialNumber(serial string) []byte {
	number, _ := strconv.ParseInt(serial[len(serial)-6:], 10, 0)
	h, _ := hex.DecodeString(fmt.Sprintf("%06x", number))

	return h
}
