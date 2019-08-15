package xmeye

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultTimeout = time.Second
	DefaultPort    = 34567

	CmdLoginResponse      uint16 = 1000
	CmdLogoutResponse     uint16 = 1002
	CmdKeepAliveResponse  uint16 = 1006
	CmdTimeRequest        uint16 = 1452
	CmdSystemInfoRequest  uint16 = 1020
	CmdLogSearchRequest   uint16 = 1442
	CmdSysManagerRequest  uint16 = 1450
	CmdSysManagerResponse uint16 = 1451

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

		sessionID:      0,
		sequenceNumber: 0,

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

func (c *Client) request(code uint16, payload interface{}) ([]byte, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	requestPacket, err := c.buildRequestPacket(code, payload)
	if err != nil {
		return nil, err
	}

	if _, err = conn.Write(requestPacket); err != nil {
		return nil, err
	}

	//fmt.Println(">>>")
	//fmt.Println(hex.Dump(requestPacket))

	atomic.AddUint32(&c.sequenceNumber, 1)

	// TODO: check read response needed?
	sessionID, responsePayload, err := c.parseResponsePacker(conn)
	if err != nil {
		return nil, err
	}

	atomic.StoreUint32(&c.sessionID, sessionID)

	return responsePayload, nil
}

func (c *Client) buildRequestPacket(code uint16, payload interface{}) ([]byte, error) {
	packet := make([]byte, 0x14) // build head

	// Head Flag, VERSION, RESERVED01, RESERVED02
	copy(packet, regularPacketHeader)

	// SESSION ID
	binary.LittleEndian.PutUint32(packet[0x04:], atomic.LoadUint32(&c.sessionID))

	// SEQUENCE NUMBER
	binary.LittleEndian.PutUint32(packet[0x08:], atomic.LoadUint32(&c.sequenceNumber))

	// Total Packet
	packet[0x0c] = 0

	// CurPacket
	packet[0x0d] = 0

	// Message Id
	binary.LittleEndian.PutUint16(packet[0x0e:], code)

	// Data Length
	var (
		payloadEncode []byte
		err           error
	)

	if payload != nil {
		payloadEncode, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	payloadLen := len(payloadEncode)
	binary.LittleEndian.PutUint32(packet[0x10:], uint32(payloadLen))

	if payloadLen > 0 {
		// DATA
		packet = append(packet, payloadEncode...)
		packet = append(packet, payloadEOF...)
	}

	return packet, nil
}

func (c *Client) parseResponsePacker(conn io.Reader) (sessionID uint32, payload []byte, err error) {
	// FIXME: после reboot через ручку странное поведение, девайс не перезагружается
	// команды принимает, но не отвечает на них

	packet := make([]byte, 0x14) // read head
	if _, err = conn.Read(packet); err != nil {
		return 0, nil, err
	}

	//fmt.Println("<<<")
	//fmt.Println("total", packet[0x0c], "current", packet[0x0d])
	//fmt.Println(hex.Dump(packet))

	// save session id
	sessionID = binary.LittleEndian.Uint32(packet[0x04:0x08])
	payloadLen := binary.LittleEndian.Uint16(packet[0x10:0x12])

	packet = make([]byte, payloadLen)
	if _, err = conn.Read(packet); err != nil {
		return 0, nil, err
	}

	// fmt.Println(hex.Dump(packet))

	return sessionID, packet, err
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
			c.Cmd(CmdKeepAliveResponse, "KeepAlive")

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
	binary.LittleEndian.PutUint32(sessionIDBytes, atomic.LoadUint32(&c.sessionID))

	return "0x" + hex.EncodeToString(sessionIDBytes)
}
