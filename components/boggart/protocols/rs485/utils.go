package rs485

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

func GenerateRequestId() []byte {
	id, _ := rand.Int(rand.Reader, big.NewInt(0xFFFF))
	return Pad(id.Bytes(), 2)
}

func GenerateCRC16(packet []byte) []byte {
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

	return Pad(Reverse(big.NewInt(int64(result)).Bytes()), 2)
}

func Reverse(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	return data
}

func Pad(data []byte, n int) []byte {
	if len(data) >= n {
		return data
	}

	for i := len(data); i < n; i++ {
		data = append(data, 0x0)
	}

	return data
}

func ToUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}

func ToFloat32(data []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(data))
}
