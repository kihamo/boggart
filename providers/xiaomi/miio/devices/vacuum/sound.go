package vacuum

import (
	"context"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio"
	"github.com/kihamo/boggart/providers/xiaomi/miio/internal"
)

const (
	soundInstallTickerDuration = time.Second
)

const (
	SoundInstallStateUnknown uint64 = iota
	SoundInstallStateDownloading
	SoundInstallStateInstalling
	SoundInstallStateInstalled
	SoundInstallStateError
)

const (
	SoundInstallErrorNo SoundInstallStatusError = iota
	SoundInstallErrorUnknown1
	SoundInstallErrorFailedDownload
	SoundInstallErrorWrongChecksum
	SoundInstallErrorUnknown4
	SoundInstallErrorUnknown5
)

type Sound struct {
	SIDInUse       uint64 `json:"sid_in_use"`
	SIDVersion     uint64 `json:"sid_version"`
	SIDInProgress  uint64 `json:"sid_in_progress"`
	Location       string `json:"location"`
	Bom            string `json:"bom"`
	Language       string `json:"language"`
	MessageVersion uint64 `json:"msg_ver"`
}

type SoundInstallStatusError uint64

type SoundInstallStatus struct {
	Progress      uint64                  `json:"progress"`
	State         uint64                  `json:"state"`
	Error         SoundInstallStatusError `json:"error"`
	SIDInProgress uint64                  `json:"sid_in_progress"`
}

func (s SoundInstallStatus) IsInstalling() bool {
	return s.State == SoundInstallStateDownloading || s.State == SoundInstallStateInstalling
}

func (s SoundInstallStatus) IsError() bool {
	return s.Error != SoundInstallErrorNo
}

func (e SoundInstallStatusError) String() string {
	switch e {
	case SoundInstallErrorNo:
		return "no"
	case SoundInstallErrorFailedDownload:
		return "download failed"
	case SoundInstallErrorWrongChecksum:
		return "wrong checksum"
	}

	return "unknown #" + strconv.FormatUint(uint64(e), 10)
}

func (d *Device) SoundVolumeTest(ctx context.Context) error {
	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "test_sound_volume", nil, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) SoundVolume(ctx context.Context) (uint32, error) {
	type response struct {
		miio.Response

		Result []uint32 `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_volume", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}

func (d *Device) SetSoundVolume(ctx context.Context, volume uint32) error {
	if volume > 100 {
		volume = 100
	}

	var reply miio.ResponseOK

	err := d.Client().Send(ctx, "change_sound_volume", []uint32{volume}, &reply)
	if err != nil {
		return err
	}

	if !miio.ResponseIsOK(reply) {
		return errors.New("device return not OK response")
	}

	return nil
}

// Мой кастомный {"result":[{"sid_in_use":10000,"sid_version":1,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236719}
// English       {"result":[{"sid_in_use":3,"sid_version":2,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236769}
// По-умолчанию  {"result":[{"sid_in_use":1,"sid_version":2,"sid_in_progress":0,"location":"prc","bom":"A.03.0002","language":"prc","msg_ver":2}],"id":1557236821}
func (d *Device) SoundCurrent(ctx context.Context) (Sound, error) {
	type response struct {
		miio.Response

		Result []Sound `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_current_sound", nil, &reply)
	if err != nil {
		return Sound{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstall(ctx context.Context, url, md5sum string, sid uint64) (SoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []SoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "dnld_install_sound", map[string]interface{}{
		"md5": md5sum,
		"url": url,
		"sid": sid,
		//"sver": 2,
	}, &reply)
	if err != nil {
		return SoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstallProgress(ctx context.Context) (SoundInstallStatus, error) {
	type response struct {
		miio.Response

		Result []SoundInstallStatus `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "get_sound_progress", nil, &reply)
	if err != nil {
		return SoundInstallStatus{}, err
	}

	return reply.Result[0], nil
}

func (d *Device) SoundInstallLocalServer(ctx context.Context, file io.ReadSeeker, sid uint64) error {
	server, err := internal.NewServer(file, d.HostnameForLocalServer())
	if err != nil {
		return err
	}
	defer server.Close()

	var status SoundInstallStatus

	status, err = d.SoundInstall(ctx, server.URL().String(), server.MD5(), sid)

	if err == nil {
		if status.IsError() {
			return errors.New("return error " + status.Error.String())
		}

		ticker := time.NewTicker(soundInstallTickerDuration)
		for range ticker.C {
			if status, err := d.SoundInstallProgress(ctx); err == nil {
				if status.IsInstalling() {
					continue
				}

				if status.IsError() {
					return errors.New("return error " + status.Error.String())
				}

				return nil
			}
		}
	}

	return err
}
