package repository

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
)

type DirectoryRepository struct {
	*MemoryRepository

	lock     sync.RWMutex
	releases []ota.Release
}

func NewDirectoryRepository() *DirectoryRepository {
	return &DirectoryRepository{
		MemoryRepository: NewMemoryRepository(),
	}
}

func (r *DirectoryRepository) Load(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasPrefix(info.Name(), "release-") {
			return nil
		}

		rl, err := release.NewLocalFile(path, "")
		if err == nil {
			r.Add(rl)
		}

		return err
	})
}
