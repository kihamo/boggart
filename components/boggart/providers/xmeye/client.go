package xmeye

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultTimeout       = time.Second
	DefaultPort          = 34567
	defaultPayloadBuffer = 2048

	timeLayout = "2006-01-02 15:04:05"

	CmdLoginResponse      uint16 = 1000
	CmdLogoutResponse     uint16 = 1002
	CmdKeepAliveResponse  uint16 = 1006
	CmdTimeRequest        uint16 = 1452
	CmdSystemInfoRequest  uint16 = 1020
	CmdAbilityGetRequest  uint16 = 1360
	CmdLogSearchRequest   uint16 = 1442
	CmdGuardRequest       uint16 = 1500
	CmdUnGuardRequest     uint16 = 1502
	CmdAlarmRequest       uint16 = 1504
	CmdSysManagerResponse uint16 = 1451

	CodeOK                                  = 100
	CodeUnknownError                        = 101
	CodeUnsupportedVersion                  = 102
	CodeRequestNotPermitted                 = 103
	CodeUserAlreadyLoggedIn                 = 104
	CodeUserUserIsNotLoggedIn               = 105
	CodeUsernameOrPasswordIsIncorrect       = 106
	CodeUserDoesNotHaveNecessaryPermissions = 107
	CodePasswordIsIncorrect                 = 203
	CodeUpgradeSuccessful                   = 515
)

var (
	regularPacketHeader = []byte{0xff, 0x01, 0x00, 0x00}
	payloadEOF          = []byte{0x0a, 0x00}
)

type Client struct {
	host     string
	username string
	password []byte
	debug    uint32

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

func (c *Client) WithDebug(flag bool) *Client {
	var debug uint32
	if flag {
		debug = 1
	}

	atomic.StoreUint32(&c.debug, debug)

	return c
}

func (c *Client) sessionIDAsString() string {
	sessionIDBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sessionIDBytes, atomic.LoadUint32(&c.sessionID))

	return "0x" + hex.EncodeToString(sessionIDBytes)
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

	if atomic.LoadUint32(&c.debug) > 0 {
		fmt.Println(">>> request")
		fmt.Println(hex.Dump(requestPacket))
	}

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
	packetHead := make([]byte, 0x14) // read head
	if _, err = conn.Read(packetHead); err != nil {
		return 0, nil, err
	}

	debug := atomic.LoadUint32(&c.debug)

	if debug > 0 {
		fmt.Println("<<< response")
		fmt.Println(hex.Dump(packetHead))
	}

	// save session id
	sessionID = binary.LittleEndian.Uint32(packetHead[0x04:0x08])
	payloadLen := int(binary.LittleEndian.Uint16(packetHead[0x10:0x12]))

	packetPayload := bytes.NewBuffer(nil)

	bufSize := defaultPayloadBuffer
	if bufSize > payloadLen {
		bufSize = payloadLen
	}
	buf := make([]byte, bufSize)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			return 0, nil, err
		}

		packetPayload.Write(buf[:n])

		if packetPayload.Len() >= payloadLen {
			break
		}
	}

	if debug > 0 {
		fmt.Println(hex.Dump(packetPayload.Bytes()))
		fmt.Println(packetPayload.String())
	}

	return sessionID, packetPayload.Bytes(), err
}

func (c *Client) Cmd(code uint16, cmd string) ([]byte, error) {
	return c.Call(code, map[string]string{
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
	response, err := c.request(code, payload)

	if err != nil {
		return nil, err
	}

	err = c.PayloadError(response)

	return response, err
}

func (c *Client) CallWithResult(code uint16, payload, result interface{}) error {
	response, err := c.Call(code, payload)
	if err != nil {
		return err
	}

	// обрезаем признак конца строки
	response = response[:len(response)-len(payloadEOF)]

	return json.Unmarshal(response, &result)
}

func (c *Client) PayloadError(payload []byte) error {
	// обрезаем признак конца строки
	payload = payload[:len(payload)-len(payloadEOF)]

	result := &Response{}
	if err := json.Unmarshal(payload, &result); err == nil {
		switch result.Ret {
		case CodeOK:
			return nil

		case CodeUnknownError:
			return errors.New("unknown error")

		case CodeUnsupportedVersion:
			return errors.New("unsupported version")

		case CodeRequestNotPermitted:
			return errors.New("request not permitted")

		case CodeUserAlreadyLoggedIn:
			return errors.New("user already logged in")

		case CodeUserUserIsNotLoggedIn:
			return errors.New("user is not logged in")

		case CodeUsernameOrPasswordIsIncorrect:
			return errors.New("username or password is incorrect")

		case CodeUserDoesNotHaveNecessaryPermissions:
			return errors.New("user does not have necessary permissions")

		case CodePasswordIsIncorrect:
			return errors.New("password is incorrect")

		case CodeUpgradeSuccessful:
			return nil

		default:
			return errors.New("unsupported unknown error")
		}
	}

	return nil
}

func (c *Client) Close() error {
	c.keepAliveTicker.Stop()
	close(c.keepAliceDone)

	return nil
}
