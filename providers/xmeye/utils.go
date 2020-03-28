package xmeye

import (
	"crypto/md5"
)

func HashPassword(password []byte) string {
	hashMD5 := md5.Sum(password)

	var n int

	hash := make([]byte, 8)

	for i := 0; i < 8; i++ {
		n = (int(hashMD5[2*i]) + int(hashMD5[2*i+1])) % 0x3e

		if n > 9 {
			if n > 35 {
				n += 61
			} else {
				n += 55
			}
		} else {
			n += 0x30
		}

		hash[i] = uint8(n)
	}

	return string(hash)
}
