package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/mergermarket/go-pkcs7"
)

type Cipher struct {
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
}

func NewCipher(key, iv []byte) (*Cipher, error) {
	block, err := aes.NewCipher(key)

	return &Cipher{
		encrypter: cipher.NewCBCEncrypter(block, iv),
		decrypter: cipher.NewCBCDecrypter(block, iv),
	}, err
}

func (c *Cipher) Encrypt(payload []byte) []byte {
	paddedPayload, _ := pkcs7.Pad(payload, aes.BlockSize)
	encryptedPayload := make([]byte, len(paddedPayload))
	c.encrypter.CryptBlocks(encryptedPayload, paddedPayload)

	return encryptedPayload
}

func (c *Cipher) Decrypt(payload []byte) []byte {
	decryptedPayload := make([]byte, len(payload))
	c.decrypter.CryptBlocks(decryptedPayload, payload)

	unpaddedPayload, _ := pkcs7.Unpad(decryptedPayload, aes.BlockSize)

	return unpaddedPayload
}
