package nis

type StatusReader struct {
	reader
}

func NewStatusReader(address string) *StatusReader {
	return &StatusReader{
		reader{
			address: address,
			command: "status",
		},
	}
}
