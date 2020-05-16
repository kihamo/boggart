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

	errors chan error
	frames chan *Frame
	done   chan<- struct{}
}

func New(conn connection.Conn) *Client {
	return &Client{
		conn:   connection.NewLooper(conn),
		errors: make(chan error),
		frames: make(chan *Frame),
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
						c.errors <- e
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
					for _, frame := range f {
						c.frames <- frame
					}
				}(frames)
			}

		case err := <-chErrors:
			if err == io.EOF {
				continue
			}

			go func(e error) {
				c.errors <- e
			}(err)
		}
	}
}

func (c *Client) NextFrame() <-chan *Frame {
	c.init()

	return c.frames
}

func (c *Client) NextError() <-chan error {
	c.init()

	return c.errors
}

func (c *Client) SkipBootLoader() error {
	_, err := c.conn.Write([]byte{0xEF})
	return err
}

func (c *Client) Call(frame *Frame) error {
	data, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.conn.Write(data)
	return err
}

func (c *Client) CallWithResult(ctx context.Context, request *Frame, waiter func(frame *Frame) bool) (*Frame, error) {
	data, err := request.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if _, err = c.conn.Write(data); err != nil {
		return nil, err
	}

	for {
		select {
		case response := <-c.NextFrame():
			if waiter(response) {
				return response, nil
			}

		case err := <-c.NextError():
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

	return c.conn.Close()
}
