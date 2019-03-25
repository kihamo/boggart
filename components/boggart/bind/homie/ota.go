package homie

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	otaTopicFirmware = MQTTPrefixImpl + "ota/firmware/+"
	otaTopicStatus   = MQTTPrefixImpl + "ota/status"

	// https://github.com/homieiot/homie-esp8266/blob/develop/docs/others/homie-implementation-specifics.md
	otaStatusSuccessfully = 200 // OTA successfully flashed
	otaStatusAccepted     = 202 // OTA request / checksum accepted
	otaStatusInProcess    = 206 // OTA in progress. The data after the status code corresponds to <bytes written>/<bytes total>
	otaStatusNotModified  = 304 // The current firmware is already up-to-date
	otaStatusBadRequest   = 400 // OTA error from your side. The identifier might be BAD_FIRMWARE, BAD_CHECKSUM, NOT_ENOUGH_SPACE, NOT_REQUESTED
	otaStatusNotEnabled   = 403 // OTA not enabled
	otaStatusFlashError   = 500 // OTA error on the ESP8266. The identifier might be FLASH_ERROR
)

var (
	regexpMagic = regexp.MustCompile(`\x25\x48\x4F\x4D\x49\x45\x5F\x45\x53\x50\x38\x32\x36\x36\x5F\x46\x57\x25`)
	// regexpFirmwareName    = regexp.MustCompile(`\x13\x54(.+)`)
	// regexpFirmwareVersion = regexp.MustCompile(`\x6A\x3F\x3E\x0E.(.+?).\x30\x48.\x1A`)
	// regexpFirmwareBrand   = regexp.MustCompile(`.\x2A.\x68.(.+?)\x6E\x2F\x0F.\x2D`)
)

// Какой-то странный баг в регулярках, нельзя задавать HEX начинающийся с буквы
// \xBF\x84\xE4\x13\x54(.+?)\x93\x44\x6B\xA7\x75
// \x6A\x3F\x3E\x0E\xE1(.+?)\xB0\x30\x48\xD4\x1A
// \xFB\x2A\xF5\x68\xC0(.+?)\x6E\x2F\x0F\xEB\x2D
func (b *Bind) OTA(ctx context.Context, file io.Reader, timeout time.Duration) error {
	if b.otaRun.IsTrue() {
		return errors.New("OTA proccess aleady running")
	}

	firmware := bytes.NewBuffer(nil)
	firmware.ReadFrom(file)

	// check firmware signature
	if !regexpMagic.Match(firmware.Bytes()) {
		return errors.New("firmware isn't homie esp8266")
	}

	/*
		// TODO: разобраться с регулярками, сейчас они работают не адекватно
		// https://github.com/homieiot/homie-esp8266/blob/develop/scripts/firmware_parser/firmware_parser.py

		var name, version, brand string

		if matches := regexpFirmwareName.FindSubmatch(buf.Bytes()); len(matches) > 0 {
			name = string(matches[1])
		}

		if matches := regexpFirmwareVersion.FindSubmatch(buf.Bytes()); len(matches) > 0 {
			version = string(matches[1])
		}

		if matches := regexpFirmwareBrand.FindSubmatch(buf.Bytes()); len(matches) > 0 {
			brand = string(matches[1])
		}
	*/
	if b.Status() != boggart.BindStatusOnline && timeout == 0 {
		return errors.New("device isn't online")
	}

	go b.otaDo(firmware, timeout)

	return nil
}

func (b *Bind) OTAIsRunning() bool {
	return b.otaRun.Load()
}

func (b *Bind) OTAChecksum() string {
	return b.otaChecksum.Load()
}

func (b *Bind) OTAProgress() (uint32, uint32) {
	return b.otaWritten.Load(), b.otaTotal.Load()
}

func (b *Bind) otaStart() {
	if b.otaRun.IsFalse() {
		b.otaFlash = make(chan struct{}, 1)
		b.otaRun.True()
	}
}

func (b *Bind) otaAbort() {
	if b.otaRun.False() {
		close(b.otaFlash)
	}
}

func (b *Bind) otaDo(firmware *bytes.Buffer, timeout time.Duration) {
	if b.Status() != boggart.BindStatusOnline && timeout > 0 {
		if ok := b.otaDelayOnline(timeout); !ok {
			b.Logger().Warn("OTA timeout", "timeout", timeout)
			return
		}

		b.Logger().Info("Device is online. Continue OTA flash")
	}

	if b.Status() != boggart.BindStatusOnline {
		b.Logger().Info("Device is offline")
		return
	}

	// running flash
	b.otaStart()

	/*
		1. During startup of the Homie for ESP8266 device, it reports the current firmware's MD5 to $fw/checksum
		   (in addition to $fw/name and $fw/version). The OTA entity may or may not use this information to
		   automatically schedule OTA updates
	*/
	checkSumHash := md5.New()
	checkSumHash.Write(firmware.Bytes())
	checkSum := hex.EncodeToString(checkSumHash.Sum(nil)[:16])
	b.otaChecksum.Set(checkSum)

	/*
		2. The OTA entity publishes the latest available firmware payload to
		   $implementation/ota/firmware/<md5 checksum>, either as binary or
		   as a Base64 encoded string
	*/
	topic := otaTopicFirmware.Format(b.config.BaseTopic, b.SerialNumber(), checkSum)
	if err := b.MQTTPublishRaw(context.Background(), topic, 1, false, firmware.Bytes()); err != nil {
		b.otaAbort()
	}
}

func (b *Bind) otaDelayOnline(timeout time.Duration) bool {
	timer := time.After(timeout)
	b.Logger().Debug("OTA start delay device online with timeout", "timeout", timeout)

	for {
		select {
		case <-b.otaFlash:
			return true

		case <-timer:
			return false
		}
	}
}

func (b *Bind) otaStatusSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.OTAIsRunning() {
		return nil
	}

	info := strings.Fields(strings.TrimSpace(message.String()))
	if len(info) == 0 {
		return errors.New("payload format is wrong")
	}

	status, err := strconv.ParseInt(info[0], 10, 64)
	if err != nil {
		return errors.New("status format is wrong")
	}

	switch status {
	/*
		3. If OTA is disabled, Homie for ESP8266 reports 403 to $implementation/ota/status
		   and aborts the OTA
	*/
	case otaStatusNotEnabled:
		b.Logger().Warn("OTA not enabled")
		b.otaAbort()

	/*
		4. If OTA is enabled and the latest available checksum is the same as what is currently running,
		   Homie for ESP8266 reports 304 and aborts the OTA
	*/
	case otaStatusNotModified:
		b.Logger().Info("OTA device firmware already up to date")
		b.otaAbort()

	/*
			5. If the checksum is not a valid MD5, Homie for ESP8266 reports 400 BAD_CHECKSUM to
		       $implementation/ota/status and aborts the OTA
	*/
	case otaStatusBadRequest:
		b.Logger().Warn("OTA bad request", "reason", info[1])
		b.otaAbort()

	/*
		6. Homie starts to flash the firmware
	*/
	case otaStatusAccepted:
		b.Logger().Info("OTA flash accepted")

	/*
		7. The firmware is updating. Homie for ESP8266 reports progress with 206 <bytes written>/<bytes total>
	*/
	case otaStatusInProcess:
		progress := strings.Split(strings.TrimSpace(info[1]), "/")
		if len(progress) != 2 {
			b.Logger().Warn("OTA progress format is wrong", "error", info[1])
		} else {
			written, _ := strconv.ParseUint(progress[0], 10, 64)
			b.otaWritten.Set(uint32(written))

			total, _ := strconv.ParseUint(progress[1], 10, 64)
			b.otaTotal.Set(uint32(total))

			b.Logger().Infof("OTA flash progress send %d bytes of %d", written, total)
		}

	/*
			8. When all bytes are flashed, the firmware is verified (including the MD5 if one was set)
			   Homie for ESP8266 either reports 200 on success, 400 if the firmware in invalid
		       or 500 if there's an internal error
	*/
	case otaStatusSuccessfully:
		b.Logger().Info("OTA flash success")
		b.otaAbort()

	case otaStatusFlashError:
		b.Logger().Error("OTA flash error", "error", info[1])
		b.otaAbort()

	default:
		b.Logger().Warn("OTA unknown status", "status", info[0])
	}

	return nil
}
