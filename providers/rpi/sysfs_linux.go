package rpi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	cpuCountPath     = "/sys/devices/system/cpu/present"
	cpuFrequentie    = "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_cur_freq"
	temperaturePath  = "/sys/class/thermal/thermal_zone0/temp"
	throttledPath    = "/sys/devices/platform/soc/soc:firmware/get_throttled"
	modelPath        = "/proc/device-tree/model"
	serialNumberPath = "/proc/device-tree/serial-number"
)

var cpuCount []uint64

func init() {
	out, err := ioutil.ReadFile(cpuCountPath)
	if err != nil {
		return
	}

	parts := bytes.Split(bytes.TrimSpace(out), []byte("-"))

	rangeStart, err := strconv.ParseUint(string(parts[0]), 10, 64)
	if err != nil {
		return
	}

	rangeEnd, err := strconv.ParseUint(string(parts[1]), 10, 64)
	if err != nil {
		return
	}

	cpuCount = make([]uint64, 0, rangeEnd)
	for i := rangeStart; i <= rangeEnd; i++ {
		cpuCount = append(cpuCount, i)
	}
}

type SysFS struct {
}

func NewSysFS() *SysFS {
	return &SysFS{}
}

func (s *SysFS) CPUFrequentie() (map[uint64]uint64, error) {
	result := make(map[uint64]uint64, len(cpuCount))

	for _, i := range cpuCount {
		out, err := ioutil.ReadFile(fmt.Sprintf(cpuFrequentie, i))
		if err != nil {
			return nil, err
		}

		value, err := strconv.ParseUint(string(bytes.TrimSpace(out)), 10, 64)
		if err != nil {
			return nil, err
		}

		result[i] = value
	}

	return result, nil
}

func (s *SysFS) Temperature() (float64, error) {
	out, err := ioutil.ReadFile(temperaturePath)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(string(bytes.TrimSpace(out)), 10, 64)
	if err != nil {
		return 0, err
	}

	return float64(value) / 1000, nil
}

func (s *SysFS) Throttled() (Throttled, error) {
	out, err := ioutil.ReadFile(throttledPath)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseUint(string(bytes.TrimSpace(out)), 16, 64)
	if err != nil {
		return 0, err
	}

	return Throttled(value), nil
}

func (s *SysFS) Model() (string, error) {
	out, err := ioutil.ReadFile(modelPath)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(out)), nil
}

func (s *SysFS) SerialNumber() (string, error) {
	out, err := ioutil.ReadFile(serialNumberPath)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(bytes.TrimRight(out, "\x00"))), nil
}
