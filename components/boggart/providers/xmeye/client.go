package xmeye

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultTimeout = time.Second
	DefaultPort    = 34567

	CmdLoginResponse     uint16 = 1000
	CmdKeepAliveRequest  uint16 = 1006
	CmdTimeRequest       uint16 = 1452
	CmdSystemInfoRequest uint16 = 1020
	CmdPhotoGetRequest   uint16 = 1600

	CodeOK = 100
)

var (
	regularPacketHeader = []byte{0xff, 0x01, 0x00, 0x00}
	payloadEOF          = []byte{0x0a, 0x00}
)

type Client struct {
	host     string
	username string
	password []byte

	connection *net.TCPConn
	mutex      sync.RWMutex

	sessionID      uint32
	sequenceNumber uint32

	keepAliveTicker *time.Ticker
	keepAliceDone   chan struct{}
}

func New(host, username, password string) *Client {
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = host + ":" + strconv.Itoa(DefaultPort)
	}

	return &Client{
		host:     host,
		username: username,
		password: []byte(password),

		sessionID:      1,
		sequenceNumber: 2,

		keepAliveTicker: time.NewTicker(time.Second * 20),
		keepAliceDone:   make(chan struct{}),
	}
}

func (c *Client) connect() (*net.TCPConn, error) {
	c.mutex.RLock()
	conn := c.connection
	c.mutex.RUnlock()

	if conn == nil {
		dial, err := net.Dial("tcp", c.host)
		if err != nil {
			return nil, err
		}

		c.mutex.Lock()
		c.connection = dial.(*net.TCPConn)
		conn = c.connection
		c.mutex.Unlock()
	}

	return conn, nil
}

func (c *Client) Login() error {
	response := &LoginResponse{}

	err := c.CallWithResult(CmdLoginResponse, map[string]string{
		"EncryptType": "MD5",
		"LoginType":   "DVRIP-Web",
		"PassWord":    HashPassword(c.password),
		"UserName":    c.username,
	}, response)

	if err != nil {
		return err
	}

	if response.Ret != CodeOK {
		return fmt.Errorf("response %d isn't ok", response.Ret)
	}

	if response.AliveInterval > 0 {
		go c.keepAlive(response.AliveInterval)
	}

	return err
}

func (c *Client) request(code uint16, payload interface{}) ([]byte, error) {
	requestPacket := make([]byte, 0x14)

	// Head Flag, VERSION, RESERVED01, RESERVED02
	copy(requestPacket, regularPacketHeader)

	// SESSION ID
	binary.LittleEndian.PutUint32(requestPacket[0x04:], atomic.LoadUint32(&c.sessionID))

	// SEQUENCE NUMBER
	binary.LittleEndian.PutUint32(requestPacket[0x08:], atomic.LoadUint32(&c.sequenceNumber))

	// Total Packet
	requestPacket[0x0c] = 0

	// CurPacket
	requestPacket[0x0d] = 0

	// Message Id
	binary.LittleEndian.PutUint16(requestPacket[0x0e:], code)

	// Data Length
	var (
		requestPayload []byte
		err            error
	)

	if payload != nil {
		requestPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	requestPayloadLen := len(requestPayload)
	binary.LittleEndian.PutUint32(requestPacket[0x10:], uint32(requestPayloadLen))

	if requestPayloadLen > 0 {
		// DATA
		requestPacket = append(requestPacket, requestPayload...)
		requestPacket = append(requestPacket, payloadEOF...)
	}

	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>")
	fmt.Println(hex.Dump(requestPacket))

	if _, err = conn.Write(requestPacket); err != nil {
		return nil, err
	}

	responsePacket := make([]byte, 20)
	if _, err = conn.Read(responsePacket); err != nil {
		return nil, err
	}

	atomic.AddUint32(&c.sequenceNumber, 1)

	// save session id
	sessionID := binary.LittleEndian.Uint32(responsePacket[0x08:0x0c])
	atomic.StoreUint32(&c.sessionID, sessionID)

	payloadLen := binary.LittleEndian.Uint16(responsePacket[0x10:0x12])

	responsePacket = make([]byte, payloadLen)
	if _, err = conn.Read(responsePacket); err != nil {
		return nil, err
	}

	fmt.Println("<<<")
	fmt.Println(hex.Dump(responsePacket))

	return responsePacket, nil
}

func (c *Client) Cmd(code uint16, cmd string) ([]byte, error) {
	return c.request(code, map[string]string{
		"Name":      cmd,
		"SessionID": c.sessionIDAsString(),
	})
}

func (c *Client) CmdWithResult(code uint16, cmd string, result interface{}) error {
	return c.CallWithResult(code, map[string]string{
		"Name":      cmd,
		"SessionID": c.sessionIDAsString(),
	}, &result)
}

func (c *Client) Call(code uint16, payload interface{}) ([]byte, error) {
	return c.request(code, payload)
}

func (c *Client) CallWithResult(code uint16, payload, result interface{}) error {
	response, err := c.request(code, payload)
	if err != nil {
		return err
	}

	// обрезаем признак конца строки
	response = response[:len(response)-len(payloadEOF)]

	return json.Unmarshal(response, &result)
}

func (c *Client) keepAlive(interval uint64) {
	for {
		select {
		case <-c.keepAliveTicker.C:
			c.Cmd(CmdKeepAliveRequest, "KeepAlive")

		case <-c.keepAliceDone:
			return
		}
	}
}

func (c *Client) Close() error {
	c.keepAliveTicker.Stop()
	close(c.keepAliceDone)

	return nil
}

func (c *Client) sessionIDAsString() string {
	sessionIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionIDBytes, atomic.LoadUint32(&c.sequenceNumber))

	return "0x" + hex.EncodeToString(sessionIDBytes)
}
