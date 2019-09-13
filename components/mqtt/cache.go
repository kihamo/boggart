package mqtt

type Cache interface {
	Get(topic Topic) (payload []byte, ok bool)
	Add(topic Topic, payload []byte)
	Payloads() map[Topic][]byte
	Resize(size int) error
	Len() int
	Purge()
}
