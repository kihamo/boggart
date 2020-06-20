package z_stack

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/protocols/connection"
)

type Client struct {
	conn connection.Looper
	once sync.Once
	lock sync.RWMutex

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
	c.once.Do(func() {
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
						c.lock.RLock()
						defer c.lock.RUnlock()

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
					c.lock.RLock()
					defer c.lock.RUnlock()
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
				c.lock.RLock()
				defer c.lock.RUnlock()

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

	c.lock.Lock()
	c.watchers = append(c.watchers, watcher)
	c.lock.Unlock()

	c.init()

	return watcher
}

func (c *Client) unregisterWatcher(watcher *Watcher) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i := len(c.watchers) - 1; i >= 0; i-- {
		if c.watchers[i] == watcher {
			c.watchers = append(c.watchers[:i], c.watchers[i+1:]...)
			watcher.close()
		}
	}
}

func (c *Client) unregisterAllWatcher() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i := len(c.watchers) - 1; i >= 0; i-- {
		c.watchers[i].close()
	}

	c.watchers = c.watchers[:0]
}

func (c *Client) Boot(ctx context.Context) error {
	/*
	 * Start as coordinator
	 */
	device, err := c.UtilGetDeviceInfo(ctx)
	if err != nil {
		return err
	}

	if device.DeviceState == DeviceStateCoordinator {
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
		for _, id := range r.Data()[6:] {
			registeredEndpoints[uint8(id)] = true
		}

	case e := <-waitErr:
		return e
	}

	// register endpoints
	endpoints := []Endpoint{{
		EndPoint:  1,
		AppProfId: 0x0104,
	}, {
		EndPoint:  2,
		AppProfId: 0x0101,
	}, {
		EndPoint:  3,
		AppProfId: 0x0105,
	}, {
		EndPoint:  4,
		AppProfId: 0x0107,
	}, {
		EndPoint:  5,
		AppProfId: 0x0108,
	}, {
		EndPoint:  6,
		AppProfId: 0x0109,
	}, {
		EndPoint:  8,
		AppProfId: 0x0104,
	}}

	for _, endpoint := range endpoints {
		if !registeredEndpoints[endpoint.EndPoint] {
			if err := c.AfRegister(ctx, endpoint); err != nil {
				return err
			}
		}
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
