package broadlink

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"net"
)

func LocalAddr() (*net.UDPAddr, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return &net.UDPAddr{
					IP: ipnet.IP,
				}, err
			}
		}
	}

	return nil, errors.New("IP not found")
}

func checksum(packet []byte) (sum uint16) {
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
