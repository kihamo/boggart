package z_stack

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	a "github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/protocols/connection"
)

type Client struct {
	conn    connection.Connection
	options options

	loopOnce sync.Once
	loopLock sync.RWMutex

	closed   uint32
	done     chan struct{}
	watchers []*Watcher

	versionOnce *a.Once
	version     *SysVersion

	permitJoinState uint32
	ledState        uint32

	devices sync.Map
}

func New(conn connection.Connection, opts ...Option) *Client {
	c := &Client{
		conn:        conn,
		options:     DefaultOptions(),
		done:        make(chan struct{}),
		watchers:    make([]*Watcher, 0),
		versionOnce: new(a.Once),
	}

	for _, opt := range opts {
		opt.apply(&c.options)
	}

	return c
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
			//fmt.Printf("\033[35m Receive <<< %v %X \033[0m\n", data, data)

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

				//fmt.Println("\033[32m Receive and parse <<< ", frame.String(), "\033[0m")

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

func (c *Client) Write(data []byte) (int, error) {
	if c.isClosed() {
		return -1, errors.New("connection is closed")
	}

	c.init()

	//fmt.Printf("\033[35m >>> %v %X \033[0m\n", data, data)

	return c.conn.Write(data)
}

func (c *Client) Call(frame *Frame) error {
	data, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	//fmt.Printf("\033[34m >>> %v %X \033[0m\n", frame.String(), data)

	_, err = c.Write(data)
	return err
}

func (c *Client) CallWithResult(ctx context.Context, frame *Frame, waiter func(frame *Frame) bool) (*Frame, error) {
	if c.isClosed() {
		return nil, errors.New("connection is closed")
	}

	watcher := c.Watch()
	defer func() {
		c.unregisterWatcher(watcher)
	}()

	err := c.Call(frame)
	if err != nil {
		return nil, err
	}

	return c.Wait(ctx, waiter)
}

func (c *Client) CallWithResultSREQ(ctx context.Context, request *Frame) (*Frame, error) {
	return c.CallWithResult(ctx, request, func(response *Frame) bool {
		return response.Type() == TypeSRSP && response.SubSystem() == request.SubSystem() && response.CommandID() == request.CommandID()
	})
}

func (c *Client) Wait(ctx context.Context, waiter func(frame *Frame) bool) (*Frame, error) {
	return c.WaitWithTimeout(ctx, waiter, DefaultWaitTimeout)
}

func (c *Client) WaitWithTimeout(ctx context.Context, waiter func(frame *Frame) bool, timeout time.Duration) (*Frame, error) {
	if c.isClosed() {
		return nil, errors.New("connection is closed")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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
	return c.WaitAsyncWithTimeout(ctx, waiter, DefaultWaitTimeout)
}

func (c *Client) WaitAsyncWithTimeout(ctx context.Context, waiter func(frame *Frame) bool, timeout time.Duration) (<-chan *Frame, <-chan error) {
	response := make(chan *Frame, 1)
	err := make(chan error, 1)

	go func() {
		if r, e := c.WaitWithTimeout(ctx, waiter, timeout); e != nil {
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
