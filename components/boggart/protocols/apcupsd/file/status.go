package file

type StatusReader struct {
	reader
}

func NewStatusReader(path string) *StatusReader {
	return &StatusReader{
		reader{
			path: path,
		},
	}
}
