package file

import (
	"bytes"
	"context"
	"io"
	"os"
)

type reader struct {
	path string
}

func (r *reader) Reader(_ context.Context) (io.Reader, error) {
	// TODO: проверить блокровку на файл
	f, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := bytes.NewBuffer(nil)

	if _, err := io.Copy(buffer, f); err != nil {
		return nil, err
	}

	return buffer, nil
}
