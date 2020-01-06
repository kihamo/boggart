// +build !linux

package rpi

type SysFS struct {
}

func NewSysFS() *SysFS {
	return &SysFS{}
}

func (s *SysFS) CPUFrequentie() (map[uint64]uint64, error) {
	return nil, ErrNotImplemented
}

func (s *SysFS) Temperature() (float64, error) {
	return -1, ErrNotImplemented
}

func (s *SysFS) Throttled() (Throttled, error) {
	return 0, ErrNotImplemented
}

func (s *SysFS) Model() (string, error) {
	return "", ErrNotImplemented
}
