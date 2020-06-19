package z_stack

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"

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

	return nil, nil
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
