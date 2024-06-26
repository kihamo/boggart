package mqtt

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type Subscription struct {
	calls uint64
	qos   uint64
	topic Topic

	mutex       sync.RWMutex
	subscribers []Subscriber
}

func NewSubscription(subscriber Subscriber, subscribers ...Subscriber) *Subscription {
	s := &Subscription{
		subscribers: append([]Subscriber{subscriber}, subscribers...),
	}

	qos := byte(0)

	for _, subscriber := range subscribers {
		if subscriber.QOS() > qos {
			qos = subscriber.QOS()
		}
	}

	s.topic = subscriber.Topic()
	s.qos = uint64(qos)

	return s
}

func (c *Subscription) Calls() uint64 {
	return atomic.LoadUint64(&c.calls)
}

func (c *Subscription) AddSubscriber(subscriber Subscriber) {
	c.mutex.Lock()
	c.subscribers = append(c.subscribers, subscriber)
	c.mutex.Unlock()

	if subscriber.QOS() > c.QOS() {
		atomic.StoreUint64(&c.qos, uint64(subscriber.QOS()))
	}
}

func (c *Subscription) RemoveSubscriber(subscriber Subscriber) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i, s := range c.subscribers {
		if s == subscriber {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			return true
		}
	}

	return false
}

func (c *Subscription) Subscribers() []Subscriber {
	subscribers := make([]Subscriber, 0, c.Len())

	c.mutex.RLock()
	copy(subscribers, c.subscribers)
	c.mutex.RUnlock()

	return subscribers
}

func (c *Subscription) Topic() Topic {
	return c.topic
}

func (c *Subscription) QOS() byte {
	return byte(atomic.LoadUint64(&c.qos))
}

func (c *Subscription) Callback(ctx context.Context, client Component, message Message) (err error) {
	atomic.AddUint64(&c.calls, 1)

	c.mutex.RLock()
	subscribers := append([]Subscriber(nil), c.subscribers...)
	c.mutex.RUnlock()

	if len(subscribers) == 0 {
		return nil
	}

	errIn := make(chan error)
	errOut := make(chan error, 1)

	go func() {
		var err error

		for {
			result := <-errIn
			if result == nil {
				errOut <- err
				return
			}

			if err == nil {
				err = errors.Wrap(result, "Call returned error")
			} else {
				err = errors.Wrap(result, err.Error())
			}
		}
	}()

	var wg sync.WaitGroup

	wg.Add(len(subscribers))

	for _, sub := range subscribers {
		go func(s Subscriber) {
			defer wg.Done()

			if errCall := s.Call(ctx, client, message); errCall != nil {
				errIn <- errCall
			}
		}(sub)
	}

	wg.Wait()
	close(errIn)

	err = <-errOut
	close(errOut)

	return err
}

func (c *Subscription) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.subscribers)
}

func (c *Subscription) Match(topic Topic) bool {
	c.mutex.RLock()
	subscribers := c.subscribers
	c.mutex.RUnlock()

	for _, subscriber := range subscribers {
		if subscriber.Topic().IsInclude(topic) {
			return true
		}
	}

	return false
}
