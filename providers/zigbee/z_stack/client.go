package z_stack

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/kihamo/boggart/protocols/connection"
)

type Client struct {
	conn connection.Looper
	once sync.Once
	lock sync.RWMutex

	responses <-chan []byte
	errors    <-chan error
	done      chan<- struct{}
}

func New(conn connection.Conn) *Client {
	return &Client{
		conn: connection.NewLooper(conn),
	}
}

func (c *Client) Call() error {
	c.once.Do(func() {
		responses, errors, done := c.conn.Loop()

		c.lock.Lock()
		c.responses = responses
		c.errors = errors
		c.done = done
		c.lock.Unlock()
	})

	for {
		select {
		case data := <-c.responses:
			if len(data) == 0 {
				continue
			}

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
					return err
				}

				fmt.Println(frame)

				// cut data
				data = data[l:]
			}

		case e := <-c.errors:
			if e == io.EOF {
				continue
			}

			fmt.Println("E", e)
		}
	}
}

func (c *Client) Close() error {
	c.lock.RLock()
	if c.done != nil {
		c.done <- struct{}{}
	}
	c.lock.RUnlock()

	return nil
}
