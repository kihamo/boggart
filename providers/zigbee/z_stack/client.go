package z_stack

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/protocols/connection"
)

var endpoints = []Endpoint{
	{EndPoint: 1, AppProfId: 0x0104},
	{EndPoint: 2, AppProfId: 0x0101},
	{EndPoint: 3, AppProfId: 0x0105},
	{EndPoint: 4, AppProfId: 0x0107},
	{EndPoint: 5, AppProfId: 0x0108},
	{EndPoint: 6, AppProfId: 0x0109},
	{EndPoint: 8, AppProfId: 0x0104},
	{EndPoint: 11, AppProfId: 0x0104, AppDeviceId: 0x0400, AppOutClusterList: []uint16{1280, 1282}},
	{EndPoint: 12, AppProfId: 0xC05E},
	{EndPoint: 13, AppProfId: 0x0104, AppInClusterList: []uint16{25}},
	{EndPoint: 47, AppProfId: 0x0104},
	{EndPoint: 110, AppProfId: 0x0104},
	{EndPoint: 242, AppProfId: 0xA1E0},
}

type Client struct {
	conn     connection.Looper
	loopOnce sync.Once
	loopLock sync.RWMutex

	closed   uint32
	done     chan struct{}
	watchers []*Watcher
}

func New(conn connection.Conn) *Client {
	return &Client{
		conn:     connection.NewLooper(conn),
		done:     make(chan struct{}),
		watchers: make([]*Watcher, 0),
	}
}

func (c *Client) init() {
	c.loopOnce.Do(func() {
		go c.loop()
	})
}

func (c *Client) loop() {
	chReceiveResponses, chReceiveErrors, chConnKill, chConnDone := c.conn.Loop()

	closer := func() {
		atomic.StoreUint32(&c.closed, 1)

		chConnKill <- struct{}{} // kill connect
		c.unregisterAllWatcher()
	}

	for {
		select {
		case data := <-chReceiveResponses:
			if len(data) == 0 {
				continue
			}

			frames := make([]*Frame, 0)

			for {
				i := bytes.IndexByte(data, SOF)
				if i > 0 {
					data = data[i:]
				}

				if len(data) < FrameLengthMin {
					break
				}

				if data[0] != SOF {
					break
				}

				l := uint16(data[PositionFrameLength]) + FrameLengthMin

				var frame Frame
				if err := frame.UnmarshalBinary(data[:l]); err != nil {
					go func(e error) {
						c.loopLock.RLock()
						defer c.loopLock.RUnlock()

						for _, w := range c.watchers {
							w.notifyError(e)
						}
					}(err)

					data = data[:0]
					continue
				}

				frames = append(frames, &frame)

				// cut data
				data = data[l:]
			}

			if len(frames) > 0 {
				go func(fs []*Frame) {
					c.loopLock.RLock()
					defer c.loopLock.RUnlock()
					/*
						if len(fs) > 0 {
							sort.SliceStable(fs, func(i, j int) bool {
								if fs[i].CommandID() == CommandAfIncomingMessage && fs[i].CommandID() == fs[j].CommandID() {
									return binary.LittleEndian.Uint32(fs[i].Data()[11:15]) < binary.LittleEndian.Uint32(fs[j].Data()[11:15])
								}

								return false
							})
						}
					*/
					for _, w := range c.watchers {
						w.notifyFrames(fs...)
					}
				}(frames)
			}

		case err := <-chReceiveErrors:
			if err == nil || err == io.EOF {
				continue
			}

			go func(e error) {
				c.loopLock.RLock()
				defer c.loopLock.RUnlock()

				for _, w := range c.watchers {
					w.notifyError(e)
				}
			}(err)

			closer()

		case <-c.done:
			closer()

		case <-chConnDone:
			return
		}
	}
}

func (c *Client) Watch() *Watcher {
	watcher := NewWatcher()

	c.loopLock.Lock()
	c.watchers = append(c.watchers, watcher)
	c.loopLock.Unlock()

	c.init()

	return watcher
}

func (c *Client) unregisterWatcher(watcher *Watcher) {
	c.loopLock.Lock()
	defer c.loopLock.Unlock()

	for i := len(c.watchers) - 1; i >= 0; i-- {
		if c.watchers[i] == watcher {
			c.watchers = append(c.watchers[:i], c.watchers[i+1:]...)
			watcher.close()
		}
	}
}

func (c *Client) unregisterAllWatcher() {
	c.loopLock.Lock()
	defer c.loopLock.Unlock()

	for i := len(c.watchers) - 1; i >= 0; i-- {
		c.watchers[i].close()
	}

	c.watchers = c.watchers[:0]
}

func (c *Client) Boot(ctx context.Context) error {
	if err := c.SkipBootLoader(); err != nil {
		return err
	}

	/*
	 * Configuration
	 */

	/*
	 * Start as coordinator
	 */
	device, err := c.UtilGetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	if device.DeviceState != DeviceStateCoordinator {
		ctxWait, cancel := context.WithTimeout(ctx, time.Millisecond*6000)
		defer cancel()

		waitResponse, waitErr := c.WaitAsync(ctxWait, func(response *Frame) bool {
			return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == 0xC0
		})

		status, err := c.ZDOStartupFromApp(ctx, 100)
		if err != nil {
			return err
		}

		if status != 0 {
			return errors.New("startup from app failed")
		}

		select {
		case r := <-waitResponse:
			if r.Data()[0] != DeviceStateCoordinator {
				return errors.New("failed state change after startup from app")
			}

		case e := <-waitErr:
			return e
		}
	}

	/*
	 * Register endpoints
	 */
	ctxWait, cancel := context.WithTimeout(ctx, time.Millisecond*6000)
	defer cancel()

	/*
		ZDO_ACTIVE_EP_RSP

		This callback message is in response to the ZDO Active Endpoint Request.

		Usage:
			AREQ:
				         1         |       1     |       1     |    2    |   1    |    2    |       1       |     0-77
				Length = 0x06-0x53 | Cmd0 = 0x45 | Cmd1 = 0x85 | SrcAddr | Status | NwkAddr | ActiveEPCount | ActiveEPList
			Attributes:
				SrcAddr       2 bytes    The message’s source network address.
				Status        1 bytes    This field indicates either SUCCESS or FAILURE.
				NWKAddr       2 bytes    Device’s short address that this response describes.
				ActiveEPCount 1 byte     Number of active endpoint in the list
				ActiveEPList  0-77 bytes Array of active endpoints on this device.

		Example from zigbee2mqtt:
			zigbee-herdsman:adapter:zStack:znp:AREQ <-- ZDO - activeEpRsp - {"srcaddr":0,"status":0,"nwkaddr":0,"activeepcount":0,"activeeplist":[]} +14ms
	*/
	waitResponse, waitErr := c.WaitAsync(ctxWait, func(response *Frame) bool {
		return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == 0x85
	})

	err = c.ZDOActiveEndpoints(ctx)
	if err != nil {
		return err
	}

	registeredEndpoints := make(map[uint8]bool)

	select {
	case r := <-waitResponse:
		dataOut := r.Data()

		if dataOut[2] != 0 {
			return errors.New("get active endpoints failed")
		}

		for _, id := range dataOut[6:] {
			registeredEndpoints[uint8(id)] = true
		}

	case e := <-waitErr:
		return e
	}

	// register endpoints
	for _, endpoint := range endpoints {
		if !registeredEndpoints[endpoint.EndPoint] {
			if err := c.AfRegister(ctx, endpoint); err != nil {
				return err
			}
		}
	}

	/*
	 * Group green power
	 */
	ep := uint8(242)
	groupID := uint16(0x0B84)

	// check exists
	if group, err := c.ZDOExtFindGroup(ctx, ep, groupID); err != nil {
		return err
	} else if group.Status != 0 {
		// register if not exists
		if err = c.ZDOExtAddToGroup(ctx, ep, groupID, nil); err != nil {
			return err
		}
	}

	/*
	 * Start default watcher
	 */

	go func() {
		watcher := c.Watch()
		defer func() {
			c.unregisterWatcher(watcher)
		}()

		for {
			select {
			case frame := <-watcher.NextFrame():
				switch frame.SubSystem() {
				case SubSystemZDOInterface:
					switch frame.CommandID() {
					case 202: // tcDeviceInd
						fmt.Println("WATCH tcDeviceInd")

						if msg, err := c.ZDODeviceJoinedMessage(frame); err == nil {
							fmt.Println(
								"nwkaddr", msg.NetworkAddress,
								"extaddr", hex.EncodeToString(msg.ExtendAddress),
								"parentaddr", msg.ParentAddress,
							)
						}

					case 193: // endDeviceAnnceInd
						fmt.Println("WATCH endDeviceAnnceInd")

						if msg, err := c.ZDOEndDeviceAnnounceMessage(frame); err == nil {
							fmt.Println(msg)
						}

					case 201: // leaveInd
						fmt.Println("WATCH leaveInd")

						if msg, err := c.ZDODeviceLeaveMessage(frame); err == nil {
							fmt.Println(msg)
						}
					}

				case SubSystemAFInterface:

				}

				//if frame.CommandID() == CommandAfIncomingMessage {
				//	msg, err := c.AfIncomingMessage(frame)
				//	if err != nil {
				//		continue
				//	}
				//
				//	fmt.Println("AF INCOMING", msg)
				//	fmt.Println((*msg.Frame.Payload.Report)[0])
				//} else {
				//	fmt.Println("DEFAULT WATCHER", frame)
				//}

			case <-watcher.NextError():
			}
		}
	}()

	return nil
}

func (c *Client) PermitJoin(ctx context.Context) error {
	ctxWait, cancel := context.WithTimeout(ctx, time.Millisecond*6000)
	defer cancel()

	/*
		ZDO_MGMT_PERMIT_JOIN_RSP

		This callback message is in response to the ZDO Management Permit Join Request.

		Usage:
			AREQ:
				       1      |       1     |       1     |    2    |   1
				Length = 0x03 | Cmd0 = 0x45 | Cmd1 = 0xB6 | SrcAddr | Status
			Attributes:
				SrcAddr       2 bytes    Source address of the message.
				Status        1 bytes    This field indicates either SUCCESS (0) or FAILURE (1).

		Example from zigbee2mqtt:
			zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,3,69,182,0,0,0,240] +9ms
			zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,3,69,182,0,0,0,240] +1ms
			zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 3 - 2 - 5 - 182 - [0,0,0] - 240 +0ms
			zigbee-herdsman:adapter:zStack:znp:AREQ <-- ZDO - mgmtPermitJoinRsp - {"srcaddr":0,"status":0} +47ms
	*/

	waitResponse, waitErr := c.WaitAsync(ctxWait, func(response *Frame) bool {
		return response.Type() == TypeAREQ && response.SubSystem() == SubSystemZDOInterface && response.CommandID() == 0xB6
	})

	/*
		Судя по коду zigbee2mqtt в 3 версии протокола 255 (постоянно включено) установить нельзя, так
		происходит защита сети, поэтому для 3 версии протокола включается специальный механизм который
		устанавливает 254 секундный интервал и переактивирует его по истечению этого времени. В версии
		1.2 протокола можно установить 255 то есть активировать постоянно. Поэтому тут TODO сделать хак
		но пока не актуально так как стик 1.2

		В любых версия по истечению интервала (а при 255 интервале сразу после установки) приходит
		пакет с CmdID=0xCB (permitJoinInd), который оповещает что время истекло. В теле пакета
		{name: 'duration', parameterType: ParameterType.UINT8} который содержит интервал который был
		установлен ранее.
	*/

	if err := c.ZDOPermitJoin(ctx, 255); err != nil {
		return err
	}

	select {
	case r := <-waitResponse:
		if r.Data()[2] != 0 {
			return errors.New("enable permit join failed")
		}

	case e := <-waitErr:
		return e
	}

	return nil
}

func (c *Client) SkipBootLoader() error {
	if c.isClosed() {
		return errors.New("connection is closed")
	}

	_, err := c.conn.Write([]byte{0xEF})
	return err
}

func (c *Client) Call(frame *Frame) error {
	if c.isClosed() {
		return errors.New("connection is closed")
	}

	c.init()

	data, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(data)
	return err
}

func (c *Client) CallWithResult(ctx context.Context, request *Frame, waiter func(frame *Frame) bool) (*Frame, error) {
	if c.isClosed() {
		return nil, errors.New("connection is closed")
	}

	watcher := c.Watch()
	defer func() {
		c.unregisterWatcher(watcher)
	}()

	data, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if _, err = c.conn.Write(data); err != nil {
		return nil, err
	}

	return c.Wait(ctx, waiter)
}

func (c *Client) CallWithResultSREQ(ctx context.Context, request *Frame) (*Frame, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*6000)
	defer cancel()

	return c.CallWithResult(ctx, request, func(response *Frame) bool {
		return response.Type() == TypeSRSP && response.SubSystem() == request.SubSystem() && response.CommandID() == request.CommandID()
	})
}

func (c *Client) Wait(ctx context.Context, waiter func(frame *Frame) bool) (*Frame, error) {
	if c.isClosed() {
		return nil, errors.New("connection is closed")
	}

	watcher := c.Watch()
	defer func() {
		c.unregisterWatcher(watcher)
	}()

	for {
		select {
		case response := <-watcher.NextFrame():
			if waiter(response) {
				return response, nil
			}

		case err := <-watcher.NextError():
			return nil, err

		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (c *Client) WaitAsync(ctx context.Context, waiter func(frame *Frame) bool) (<-chan *Frame, <-chan error) {
	response := make(chan *Frame, 1)
	err := make(chan error, 1)

	go func() {
		if r, e := c.Wait(ctx, waiter); e != nil {
			err <- e
		} else {
			response <- r
		}
	}()

	return response, err
}

func (c *Client) isClosed() bool {
	return atomic.LoadUint32(&c.closed) != 0
}

func (c *Client) Close() error {
	if c.isClosed() {
		return nil
	}

	close(c.done)

	return c.conn.Close()
}
