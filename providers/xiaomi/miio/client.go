package miio

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"strconv"
	a "sync/atomic"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/protocols/connection"
	transport "github.com/kihamo/boggart/protocols/connection/transport/net"
)

const (
	DefaultPort = 54321
)

type Client struct {
	io.Closer

	packetsCounter uint32
	stampDiff      int64

	conn     connection.ObserverConnection
	connOnce *atomic.Once

	address     string
	token       string
	tokenBytes  []byte
	helloPacket *Packet

	cryptoOnce   *atomic.Once
	cryptoCipher cipher.Block
	cryptoIV     []byte
}

/*
Заметки:
- Пылесос очень чувствителен к счетчику пакетов, он похоже кэширует соответствие ip и счетчика
  поэтому если с одного ip начинать несколько сессий со сбросом счетчик, то пылесос не отвечает
*/

func NewClient(address, token string) *Client {
	t, _ := hex.DecodeString(token)

	return &Client{
		packetsCounter: uint32(time.Now().Unix()),
		address:        address,
		token:          token,
		tokenBytes:     t,
		connOnce:       new(atomic.Once),
		cryptoOnce:     new(atomic.Once),
	}
}

func (c *Client) lazyConnect() (_ connection.ObserverConnection, err error) {
	c.connOnce.Do(func() {
		t := transport.New(
			transport.WithAddress(net.JoinHostPort(c.address, strconv.Itoa(DefaultPort))),
			transport.WithNetwork("udp"),
			// слишком большое значение ставить нельзя, так как вычитка происходит
			// постоянно, большой риск попадание на блокировку длительностью в весь
			// таймаут, если ответ не успеет прийти в эту вычитку и она свалится по таймауту
			transport.WithReadTimeout(time.Second*1),
			transport.WithWriteTimeout(time.Second*3),
		)

		conn := connection.NewObserverConnection(t,
			connection.WithOnceInit(true),
			connection.WithLocalLock(true),
			//connection.WithDumpRead(func(data []byte) {
			//	fmt.Println(time.Now(), "read dump", data)
			//}),
			//connection.WithDumpWrite(func(data []byte) {
			//	fmt.Println(time.Now(), "write dump", data)
			//}),
		)

		// начинаем сессию hello пакетом
		var response *Packet

		response, err = c.hello(conn)
		if err == nil {
			c.helloPacket = response
		}

		c.conn = conn
	})

	if err != nil {
		c.connOnce.Reset()
		return nil, err
	}

	return c.conn, nil
}

func (c *Client) initCrypto() (err error) {
	defer func() {
		if err != nil {
			c.cryptoOnce.Reset()
		}
	}()

	c.cryptoOnce.Do(func() {
		var token []byte

		token, err = hex.DecodeString(c.token)
		if err != nil {
			return
		}

		hash := md5.New()

		if _, err = hash.Write(token); err != nil {
			return
		}

		key := hash.Sum(nil)

		hash.Reset()

		if _, err = hash.Write(key); err != nil {
			return
		}

		if _, err = hash.Write(token); err != nil {
			return
		}

		c.cryptoIV = hash.Sum(nil)

		if c.cryptoCipher, err = aes.NewCipher(key); err != nil {
			return
		}
	})

	return err
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *Client) LocalAddr() (*net.UDPAddr, error) {
	return nil, errors.New("not supported")
}

func (c *Client) SetPacketsCounter(count uint32) {
	a.StoreUint32(&c.packetsCounter, count)
}

func (c *Client) PacketsCounter() uint32 {
	return a.LoadUint32(&c.packetsCounter)
}

func (c *Client) hello(conn connection.Connection) (*Packet, error) {
	requestPacket := NewPacket(nil)
	requestPacket.SetUnknown(0xFFFFFFFF)
	requestPacket.SetDeviceType(0xFFFF)
	requestPacket.SetDeviceID(0xFFFF)
	requestPacket.SetTimestamp(time.Unix(int64(0xFFFFFFFF), 0))

	requestData, err := requestPacket.MarshalBinary()
	if err != nil {
		return nil, err
	}

	responseData, err := conn.Invoke(requestData)
	if err != nil {
		return nil, err
	}

	responsePacket := NewPacket(nil)
	err = responsePacket.UnmarshalBinary(responseData)

	return responsePacket, err
}

func (c *Client) DeviceID(ctx context.Context) (uint16, error) {
	_, err := c.lazyConnect()
	if err != nil {
		return 0, err
	}

	return c.helloPacket.DeviceID(), nil
}

func (c *Client) CallRPC(ctx context.Context, method string, params, response interface{}) error {
	conn, err := c.lazyConnect()
	if err != nil {
		return err
	}

	requestPacket := NewPacket(c.tokenBytes)
	requestPacket.SetPayloadRPC(0, method, params)

	requestPacket.SetDeviceID(c.helloPacket.DeviceID())
	requestPacket.SetDeviceType(c.helloPacket.DeviceType())

	diff := time.Duration(a.LoadInt64(&c.stampDiff))
	requestPacket.SetTimestamp(time.Now().Add(diff))

	requestID := a.AddUint32(&c.packetsCounter, 1)
	requestPacket.SetPayloadRPC(requestID, method, params)

	if err := c.encrypt(requestPacket); err != nil {
		return err
	}

	requestData, err := requestPacket.MarshalBinary()
	if err != nil {
		return err
	}

	done := make(chan struct{}, 1)

	observer := connection.ObserverFunc(func(responseData []byte, e error) {
		closer := func(e error) {
			err = e
			done <- struct{}{}
		}

		if e != nil {
			closer(e)
			return
		}

		responsePacket := NewPacket(nil)

		if e = responsePacket.UnmarshalBinary(responseData); e != nil {
			closer(e)
			return
		}

		if e = c.decrypt(responsePacket); e != nil {
			closer(e)
			return
		}

		a.StoreInt64(&c.stampDiff, int64(time.Since(responsePacket.Timestamp())))

		var responseError ResponseError
		if e = responsePacket.PayloadJSON(&responseError); e == nil && len(responseError.Error.Message) > 0 {
			closer(errors.New(responseError.Error.Message))
			return
		}

		var responseUnknown ResponseUnknownMethod
		if e = responsePacket.PayloadJSON(&responseUnknown); e == nil && responseUnknown.Result == "unknown_method" {
			closer(errors.New("unknown method"))
			return
		}

		var responseDefault Response
		if e = responsePacket.PayloadJSON(&responseDefault); e == nil && responseDefault.ID != requestID {
			// если requestID > responseDefault.ID то можно еще делать попытки вычитывания, обычно такое бывает
			// из-за того что контекст прерывается быстрее чем девай отвечает и следующий реквест получает предыдущий респонс
			return
		}

		e = responsePacket.PayloadJSON(&response)
		closer(e)
	})
	conn.Attach(observer)

	defer conn.Detach(observer)

	if _, e := conn.Write(requestData); e != nil {
		return e
	}

	select {
	case <-done:

	case <-ctx.Done():
		err = ctx.Err()
	}

	return err
}

func (c *Client) encrypt(packet *Packet) error {
	if err := c.initCrypto(); err != nil {
		return err
	}

	payload := packet.Payload()
	blockSize := c.cryptoCipher.BlockSize()

	paddingCount := blockSize - len(payload)%blockSize
	paddingData := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)

	dataOriginal := append(payload, paddingData...)
	dataEncrypted := make([]byte, len(dataOriginal))

	encrypter := cipher.NewCBCEncrypter(c.cryptoCipher, c.cryptoIV)
	encrypter.CryptBlocks(dataEncrypted, dataOriginal)

	packet.SetPayload(dataEncrypted)

	return nil
}

func (c *Client) decrypt(packet *Packet) error {
	if err := c.initCrypto(); err != nil {
		return err
	}

	payload := packet.Payload()

	dataOriginal := make([]byte, len(payload))

	decrypter := cipher.NewCBCDecrypter(c.cryptoCipher, c.cryptoIV)
	decrypter.CryptBlocks(dataOriginal, payload)

	length := len(dataOriginal)
	unPadding := int(dataOriginal[length-1])

	packet.SetPayload(dataOriginal[:(length-unPadding)-1])

	return nil
}
