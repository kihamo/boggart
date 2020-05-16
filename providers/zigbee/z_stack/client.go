package z_stack

import (
	"bytes"
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

func (c *Client) Close() error {
	c.lock.RLock()
	if c.done != nil {
		c.done <- struct{}{}
	}
	c.lock.RUnlock()

	return c.conn.Close()
}
