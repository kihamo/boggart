package serial_network

const (
	maxBufferSize = 1024
)

type Client interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Invoke(request []byte) (response []byte, err error)
}
