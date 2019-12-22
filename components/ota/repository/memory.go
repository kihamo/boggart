package repository

import (
	"errors"
	"sync"

	"github.com/kihamo/boggart/components/ota"
)

type MemoryRepository struct {
	lock     sync.RWMutex
	releases []ota.Release
}

func NewMemoryRepository(releases ...ota.Release) *MemoryRepository {
	return &MemoryRepository{
		releases: releases,
	}
}

func (r *MemoryRepository) Add(release ota.Release) {
	r.lock.Lock()
	r.releases = append(r.releases, release)
	r.lock.Unlock()
}

func (r *MemoryRepository) Remove(release ota.Release) {
	r.lock.Lock()

	for i, rl := range r.releases {
		if release == rl {
			r.releases = append(r.releases[:i], r.releases[i+1:]...)
			break
		}
	}

	r.lock.Unlock()
}

func (r *MemoryRepository) Releases(arch string) ([]ota.Release, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	releases := make([]ota.Release, 0, len(r.releases))
	for _, release := range r.releases {
		if arch != "" && release.Architecture() != arch && release.Architecture() != ota.ArchitectureUnknown {
			continue
		}

		releases = append(releases, release)
	}

	return releases, nil
}

func (r *MemoryRepository) ReleaseLatest(arch string) (ota.Release, error) {
	releases, err := r.Releases(arch)
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, errors.New("latest release not found")
	}

	return releases[len(releases)-1], nil
}
