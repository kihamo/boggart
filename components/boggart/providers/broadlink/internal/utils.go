package internal

import (
	"crypto/aes"
	"crypto/cipher"
)

// https://github.com/mjg59/python-broadlink/blob/5cfb92ee12f1afc2e5a997a9f3121f33b5b2ece0/protocol.md
// https://github.com/mixcode/broadlink/blob/59c23b4e7b1ab5188aafc54ca46ff794a5aec6a2/remotecontrol.go

const (
	DefaultBufferSize = 2048
)

func Checksum(packet []byte) (sum uint16) {
	sum = 0xbeaf
	for _, b := range packet {
		sum += uint16(b)
	}

	return sum
}

func blockCipher(mode cipher.BlockMode, data []byte) []byte {
	size := len(data)
	if size == 0 {
		return nil
	}

	// start position of leftover data
	split := (size / aes.BlockSize) * aes.BlockSize
	// calculate padding data size
	pad := 0
	if split < size {
		pad = split + aes.BlockSize - size
	}

	result := make([]byte, size+pad)
	mode.CryptBlocks(result[:split], data[:split])
	if split < size {
		// add padding to leftover data and feed to cipher function
		tmp := make([]byte, aes.BlockSize)
		copy(tmp, data[split:])
		mode.CryptBlocks(result[split:], tmp)
	}

	return result
}
