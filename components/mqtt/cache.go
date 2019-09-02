package mqtt

type Cache interface {
	Get(topic string) (payload []byte, ok bool)
	Add(topic string, payload []byte)
	Payloads() map[string][]byte
	Resize(size int) error
	Len() int
	Purge()
}
