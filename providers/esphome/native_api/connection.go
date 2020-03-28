package nativeapi

import (
	"context"
	"encoding/hex"
	"errors"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	defaultReadTimeout   = time.Second * 10
	defaultWriteTimeout  = time.Second * 10
	defaultCloseDeadline = time.Second
)

var (
	ErrConnectionNotInit = errors.New("connection must initialization")
)

func (c *Client) connectionInit() error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	tcpConn := conn.(*net.TCPConn)

	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}

	if err = tcpConn.SetKeepAlivePeriod(keepAliveInterval); err != nil {
		return err
	}

	c.connectionSet(tcpConn)
	return nil
}

func (c *Client) connectionRun() {
	if !c.isRestart() {
		go c.loop()
		go c.keepalive()
	}
}

func (c *Client) connectionClose() (err error) {
	conn := c.connectionGet()

	if conn != nil {
		err = conn.Close()

		if err == nil {
			c.connectionSet(nil)

			if c.isRestart() {
				atomic.StoreInt64(&c.closeDeadline, time.Now().Add(defaultCloseDeadline).UnixNano())
			}
		}
	}

	return err
}

func (c *Client) connectionSet(conn net.Conn) {
	c.mutex.Lock()
	c.conn = conn
	c.mutex.Unlock()
}

func (c *Client) connectionGet() net.Conn {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.conn
}

func (c *Client) loop() {
	c.mutex.RLock()
	done := c.done
	c.mutex.RUnlock()

	for {
		select {
		case <-done:
			return

		default:
			ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
			message, err := c.read(ctx)
			cancel()

			c.handle(message, err)
		}
	}
}

func (c *Client) write(ctx context.Context, request proto.Message) (err error) {
	err = c.connectionCheck()
	if err != nil {
		return err
	}

	conn := c.connectionGet()

	name := proto.MessageName(request)
	requestType, ok := messageTypesByName[name]
	if !ok {
		return errors.New("unknown request message type")
	}

	requestPayload, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	requestPacket := make([]byte, 3, len(requestPayload)+3)
	requestPacket[0] = packetMagicByte
	requestPacket[1] = byte(len(requestPayload))
	requestPacket[2] = requestType
	requestPacket = append(requestPacket, requestPayload...)

	if deadline, ok := ctx.Deadline(); ok {
		err = conn.SetWriteDeadline(deadline)
	}

	if err == nil {
		_, err = conn.Write(requestPacket)
	}

	if err != nil {
		c.connectionCheckBroken()
		return err
	}

	if c.Debug() {
		println(">>> Model " + name)
		println(">>> Packet")
		print(hex.Dump(requestPacket))
	}

	return nil
}

func (c *Client) read(ctx context.Context) (message proto.Message, err error) {
	err = c.connectionCheck()
	if err != nil {
		return nil, err
	}

	conn := c.connectionGet()

	if deadline, ok := ctx.Deadline(); ok {
		err = conn.SetReadDeadline(deadline)
	}

	var (
		n                  int
		responsePacketHead []byte
	)

	if err == nil {
		responsePacketHead = make([]byte, 3)
		n, err = conn.Read(responsePacketHead)
	}

	if err != nil {
		c.connectionCheckBroken()
		return nil, err
	}

	debug := c.Debug()

	if debug {
		println("<<<  Packet header")
		print(hex.Dump(responsePacketHead))
	}

	if n < 3 {
		return nil, errors.New("header of response packet failed")
	}

	if responsePacketHead[0] != packetMagicByte {
		return nil, errors.New("magic byte of response packet failed")
	}

	// type
	responseType, ok := messageTypesByID[responsePacketHead[2]]
	if !ok {
		return nil, errors.New("unknown response message type")
	}

	if debug {
		println("<<< Model " + responseType)
	}

	responseReflect := proto.MessageType(responseType)
	if responseReflect == nil {
		return nil, errors.New("unknown type: " + responseType)
	}

	message = reflect.New(responseReflect.Elem()).Interface().(proto.Message)

	// empty payload
	if responsePacketHead[1] == 0 {
		return message, nil
	}

	// parse payload
	responsePacketPayload := make([]byte, responsePacketHead[1])
	_, err = conn.Read(responsePacketPayload)

	if err != nil {
		c.connectionCheckBroken()
		return nil, err
	}

	if debug {
		println("<<< Packet payload")
		print(hex.Dump(responsePacketPayload))
	}

	if err := proto.Unmarshal(responsePacketPayload, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (c *Client) connectionCheck() error {
	if conn := c.connectionGet(); conn == nil {
		return ErrConnectionNotInit
	}

	return nil
}

func (c *Client) connectionCheckBroken() {
	conn := c.connectionGet()
	if conn == nil {
		return
	}

	if err := ConnectionCheck(conn); err != nil {
		c.restart()
	}
}
