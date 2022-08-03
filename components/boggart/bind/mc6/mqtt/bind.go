package mqtt

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"

	"github.com/kihamo/boggart/components/boggart/di"
)

const (
	size = 16
)

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

func (b *Bind) Decrypt(payload []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(b.key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(payload))

	for bs, be := 0, size; bs < len(payload); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], payload[bs:be])
	}

	paddingSize := int(decrypted[len(decrypted)-1])
	return decrypted[0 : len(decrypted)-paddingSize], nil
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
