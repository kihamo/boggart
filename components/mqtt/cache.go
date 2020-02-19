package mqtt

import (
	"time"
)

type CacheItem interface {
	Datetime() time.Time
	Topic() Topic
	Payload() []byte
}

type Cache interface {
	Get(topic Topic) (item CacheItem, ok bool)
	Add(topic Topic, payload []byte)
	Payloads() []CacheItem
	Resize(size int) error
	Len() int
	Purge()
}
