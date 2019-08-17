package xmeye

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
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

type Client struct {
	host     string
	username string
	password []byte
	debug    uint32

	connection   *net.TCPConn
	mutex        sync.RWMutex
	mutexRequest sync.Mutex
	done         chan struct{}

	sessionID      uint32
	sequenceNumber uint32
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

func (c *Client) request(code uint16, payload interface{}) (*internal.Packet, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}

	c.mutexRequest.Lock()
	defer c.mutexRequest.Unlock()

	//
	// --- REQUEST ---
	//
	requestPacket := internal.NewPacket()
	requestPacket.SessionID = atomic.LoadUint32(&c.sessionID)
	requestPacket.SequenceNumber = atomic.LoadUint32(&c.sequenceNumber)
	requestPacket.MessageID = code

	if err := requestPacket.LoadPayload(payload); err != nil {
		return nil, err
	}

	requestPacketBytes := requestPacket.Marshal()

	if _, err = conn.Write(requestPacketBytes); err != nil {
		return nil, err
	}

	debug := atomic.LoadUint32(&c.debug)
	if debug > 0 {
		fmt.Println(">>> request")
		fmt.Println(hex.Dump(requestPacketBytes))
	}

	atomic.AddUint32(&c.sequenceNumber, 1)

	//
	// --- RESPONSE ---
	//
	responsePacketHead := make([]byte, 0x14) // read head
	if _, err = conn.Read(responsePacketHead); err != nil {
		return nil, err
	}

	// save session id
	responsePacket := internal.PacketUnmarshal(responsePacketHead)
	if responsePacket.PayloadLen == 0 {
		return &responsePacket, nil
	}

	bufSize := defaultPayloadBuffer
	if bufSize > responsePacket.PayloadLen {
		bufSize = responsePacket.PayloadLen
	}
	buf := make([]byte, bufSize)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			return &responsePacket, err
		}

		responsePacket.Payload.Write(buf[:n])

		if responsePacket.Payload.Len() >= responsePacket.PayloadLen {
			break
		}
	}

	if debug > 0 {
		fmt.Println("<<< response")
		fmt.Println(hex.Dump(responsePacket.Marshal()))
	}

	return &responsePacket, nil
}

func (c *Client) Call(code uint16, payload interface{}) (*internal.Packet, error) {
	p, err := c.request(code, payload)
	if err != nil {
		return nil, err
	}

	err = c.PayloadError(p.Payload)
	return p, err
}

func (c *Client) CallWithResult(code uint16, payload, result interface{}) error {
	packet, err := c.Call(code, payload)
	if err != nil {
		return err
	}

	return packet.Payload.UnmarshalJSON(result)
}

func (c *Client) Cmd(code uint16, cmd string) (*internal.Packet, error) {
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

func (c *Client) PayloadError(payload *internal.Payload) error {
	result := &internal.Response{}

	if err := payload.UnmarshalJSON(result); err == nil {
		switch result.Ret {
		case CodeOK, CodeUpgradeSuccessful, 0:
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

		default:
			return errors.New("unsupported unknown error")
		}
	}

	return nil
}

func (c *Client) Close() (err error) {
	c.mutex.Lock()

	if c.done != nil {
		close(c.done)
		c.done = nil
	}

	if c.connection != nil {
		err = c.connection.Close()
		c.connection = nil
	}

	c.mutex.Unlock()

	return err
}
