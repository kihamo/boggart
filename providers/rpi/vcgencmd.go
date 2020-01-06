package rpi

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

const (
	VoltsIDCore   VoltsID = "core"
	VoltsIDSDramC VoltsID = "sdram_c"
	VoltsIDSDramI VoltsID = "sdram_i"
	VoltsIDSDramP VoltsID = "sdram_p"
)

type VoltsID string

func (v VoltsID) String() string {
	return string(v)
}

type VCGenCMD struct {
}

func NewVCGenCMD() *VCGenCMD {
	return &VCGenCMD{}
}

func (v *VCGenCMD) Execute(name string, arg ...string) (*bytes.Buffer, error) {
	cmd := exec.Command("vcgencmd", append([]string{name}, arg...)...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &out, nil
}

func (v *VCGenCMD) SupportCommands() ([]string, error) {
	out, err := v.Execute("commands")
	if err != nil {
		return nil, err
	}

	text := strings.TrimPrefix(out.String(), "commands=")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\"")

	return strings.Fields(text), nil
}

func (v *VCGenCMD) Voltage(id VoltsID) (float64, error) {
	out, err := v.Execute("measure_volts", id.String())
	if err != nil {
		return -1, err
	}

	text := strings.TrimPrefix(out.String(), "volt=")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "V")

	return strconv.ParseFloat(text, 64)
}

func (v *VCGenCMD) Temperature() (float64, error) {
	out, err := v.Execute("measure_temp")
	if err != nil {
		return -1, err
	}

	text := strings.TrimPrefix(out.String(), "temp=")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "'C")

	return strconv.ParseFloat(text, 64)
}

func (v *VCGenCMD) Throttled() (Throttled, error) {
	out, err := v.Execute("get_throttled")
	if err != nil {
		return 0, err
	}

	text := strings.TrimPrefix(out.String(), "Throttled=")
	text = strings.TrimSpace(text)

	value, err := strconv.ParseUint(text, 0, 64)
	if err != nil {
		return 0, err
	}

	return Throttled(value), nil
}
