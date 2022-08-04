package mqtt

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
)

// https://www.devglan.com/online-tools/aes-encryption-decryption

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind

	key []byte
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.key = []byte(b.config().AES256Key)

	if mac := b.config().MAC.String(); mac != "" {
		b.Meta().SetMACAsString(mac)
	}

	return nil
}

func (b *Bind) Decrypt(crypted []byte) ([]byte, error) {
	if len(crypted) < 1 {
		return nil, errors.New("wrong crypted data")
	}

	block, err := aes.NewCipher(b.key)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()

	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	length := len(out)
	unpadding := int(out[length-1])
	return out[:(length - unpadding)], nil
}

func (b *Bind) Encrypt(uncrypted []byte) ([]byte, error) {
	if len(uncrypted) < 1 {
		return nil, errors.New("wrong uncrypted data")
	}

	block, err := aes.NewCipher(b.key)
	if err != nil {
		return nil, err
	}

	bs := block.BlockSize()

	padding := bs - len(uncrypted)%bs
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	uncrypted = append(uncrypted, padtext...)

	if len(uncrypted)%bs != 0 {
		return nil, errors.New("wrong padding")
	}

	out := make([]byte, len(uncrypted))
	dst := out
	for len(uncrypted) > 0 {
		block.Encrypt(dst, uncrypted[:bs])
		uncrypted = uncrypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

func (b *Bind) ParseUpdate(payload []byte) (*Update, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(payload)))
	n, err := base64.StdEncoding.Decode(buf, payload)
	if err != nil {
		return nil, err
	}

	decrypt, err := b.Decrypt(buf[:n])
	if err != nil {
		return nil, err
	}

	b.Logger().Debugf("receive %s", decrypt)

	update := &Update{}

	err = json.Unmarshal(decrypt, update)
	return update, err
}

func (b *Bind) GenerateUpdate(update *Update) ([]byte, error) {
	payload, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	b.Logger().Debugf("transmit %s", payload)

	out, err := b.Encrypt(payload)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(out)))
	base64.StdEncoding.Encode(buf, out)

	return buf, nil
}
