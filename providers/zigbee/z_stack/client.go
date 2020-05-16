package z_stack

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/kihamo/boggart/protocols/connection"
)

type Client struct {
	conn connection.Looper
	once sync.Once
	lock sync.RWMutex

	done     chan<- struct{}
	watchers []*Watcher
}

func New(conn connection.Conn) *Client {
	return &Client{
		conn:     connection.NewLooper(conn),
		watchers: make([]*Watcher, 0),
	}
}

func (c *Client) init() {
	c.once.Do(func() {
		go c.loop()
	})
}

func (c *Client) loop() {
	chResponses, chErrors, chDone := c.conn.Loop()

	c.lock.Lock()
	c.done = chDone
	c.lock.Unlock()

	for {
		select {
		case data := <-chResponses:
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
				go func(f []*Frame) {
					c.lock.RLock()
					defer c.lock.RUnlock()

					for _, w := range c.watchers {
						for _, frame := range f {
							w.notifyFrame(frame)
						}
					}
				}(frames)
			}

		case err := <-chErrors:
			if err == io.EOF {
				continue
			}

			go func(e error) {
				c.lock.RLock()
				defer c.lock.RUnlock()

				for _, w := range c.watchers {
					w.notifyError(e)
				}
			}(err)
		}
	}
}

func (c *Client) Watch() *Watcher {
	c.init()

	watcher := NewWatcher()

	c.lock.Lock()
	c.watchers = append(c.watchers, watcher)
	c.lock.Unlock()

	return watcher
}

func (c *Client) unregisterWatcher(watcher *Watcher) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i := len(c.watchers) - 1; i >= 0; i-- {
		if c.watchers[i] == watcher {
			c.watchers = append(c.watchers[:i], c.watchers[i+1:]...)
		}
	}
}

func (c *Client) SkipBootLoader() error {
	_, err := c.conn.Write([]byte{0xEF})
	return err
}

func (c *Client) Call(frame *Frame) error {
	c.init()

	data, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(data)
	return err
}

func (c *Client) CallWithResult(ctx context.Context, request *Frame, waiter func(frame *Frame) bool) (*Frame, error) {
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

func (c *Client) Close() error {
	c.lock.RLock()
	if c.done != nil {
		c.done <- struct{}{}
	}
	c.lock.RUnlock()

	c.lock.Lock()
	c.watchers = c.watchers[:0]
	c.lock.Unlock()

	return c.conn.Close()
}
