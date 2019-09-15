package native_api

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/golang/protobuf/proto"
	"net"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	defaultPort = 6053
)

type connection struct {
	sync.Mutex

	address    string
	debug      int32
	connection net.Conn
}

func newConnection(address string) (*connection, error) {
	if _, _, err := net.SplitHostPort(address); err != nil {
		address = address + ":" + strconv.Itoa(defaultPort)
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	c := &connection{
		address:    address,
		connection: conn,
	}

	err = conn.(*net.TCPConn).SetKeepAlive(true)

	return c, nil
}

func (c *connection) WithDebug(debug bool) *connection {
	if debug {
		atomic.StoreInt32(&c.debug, 1)
	} else {
		atomic.StoreInt32(&c.debug, 0)
	}

	return c
}

func (c *connection) Close() error {
	return c.connection.Close()
}

func (c *connection) Write(ctx context.Context, request proto.Message) error {
	requestType, ok := messageTypesByName[proto.MessageName(request)]
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
		err = c.connection.SetWriteDeadline(deadline)
	}

	if err == nil {
		_, err = c.connection.Write(requestPacket)
	}

	if err != nil {
		c.connectionCheck()
		return err
	}

	if atomic.LoadInt32(&c.debug) != 0 {
		println(">>> ")
		println(hex.Dump(requestPacket))
	}

	return nil
}

func (c *connection) Read(ctx context.Context) (proto.Message, error) {
	var err error

	if deadline, ok := ctx.Deadline(); ok {
		err = c.connection.SetReadDeadline(deadline)
	}

	var (
		n                  int
		responsePacketHead []byte
	)

	if err == nil {
		responsePacketHead = make([]byte, 3)
		n, err = c.connection.Read(responsePacketHead)
	}

	if err != nil {
		c.connectionCheck()
		return nil, err
	}

	debug := atomic.LoadInt32(&c.debug) != 0

	if debug {
		println("<<<")
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

	responseReflect := proto.MessageType(responseType)
	if responseReflect == nil {
		return nil, errors.New("unknown type: " + responseType)
	}

	response := reflect.New(responseReflect.Elem()).Interface().(proto.Message)

	// empty payload
	if responsePacketHead[1] == 0 {
		return response, nil
	}

	// parse payload
	responsePacketPayload := make([]byte, responsePacketHead[1])
	n, err = c.connection.Read(responsePacketPayload)

	if err != nil {
		c.connectionCheck()
		return nil, err
	}

	if debug {
		print(hex.Dump(responsePacketPayload))
	}

	if err := proto.Unmarshal(responsePacketPayload, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) connectionCheck() {
	if e := ConnectionCheck(c.connection); e != nil {
		// TODO:
		//c.reset()
	}
}
