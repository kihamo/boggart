package packet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io"
)

type Crypto struct {
	Base

	token  []byte
	iv     []byte
	key    []byte
	cipher cipher.Block
}

func NewCrypto(deviceID []byte, token string) (*Crypto, error) {
	binToken, err := hex.DecodeString(token)
	if err != nil {
		return nil, err
	}

	p := &Crypto{
		token: binToken,
	}

	p.Base = *NewBase()

	p.Header.DeviceID = deviceID

	hash := md5.New()
	_, err = hash.Write(binToken)
	if err != nil {
		return nil, err
	}
	p.key = hash.Sum(nil)

	hash = md5.New()
	_, err = hash.Write(p.key)
	if err != nil {
		return nil, err
	}
	_, err = hash.Write(binToken)
	if err != nil {
		return nil, err
	}
	p.iv = hash.Sum(nil)

	p.cipher, err = aes.NewCipher(p.key)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Crypto) WriteTo(w io.Writer) (int64, error) {
	err := p.SetBody(p.encrypt(p.Base.Body()))
	if err != nil {
		return -1, err
	}

	return p.Base.WriteTo(w)
}

func (p *Crypto) ReadFrom(r io.Reader) (int64, error) {
	n, err := p.Base.ReadFrom(r)
	if err != nil {
		return -1, err
	}

	p.SetBody(p.decrypt(p.body))
	return n, nil
}

func (p *Crypto) SetBody(body []byte) error {
	p.Header.Checksum = p.token
	return p.Base.SetBody(body)
}

func (p *Crypto) encrypt(body []byte) []byte {
	blockSize := p.cipher.BlockSize()
	padding := blockSize - len(body)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	dataOriginal := append(body, padtext...)

	blockMode := cipher.NewCBCEncrypter(p.cipher, p.iv)
	dataEncrypted := dataOriginal
	blockMode.CryptBlocks(dataEncrypted, dataOriginal)

	return dataEncrypted
}

func (p *Crypto) decrypt(body []byte) []byte {
	blockMode := cipher.NewCBCDecrypter(p.cipher, p.iv)
	dataOriginal := make([]byte, len(body))
	blockMode.CryptBlocks(dataOriginal, body)

	length := len(dataOriginal)
	unpadding := int(dataOriginal[length-1])

	return dataOriginal[:(length-unpadding)-1]
}
