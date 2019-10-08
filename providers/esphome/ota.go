package esphome

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/kihamo/boggart/protocols/connection"
)

const (
	DefaultESP8266Port = 8266
	DefaultESP32Port   = 3232

	Version_1_0 = 1

	CodeOK            = 0
	CodeRequestAuth   = 1
	CodeHeaderOK      = 64
	CodeAuthOK        = 65
	CodeUpdatePrepare = 66
	CodeBinMD5        = 67
	CodeReceiveOK     = 68
	CodeUpdateEnd     = 69

	ErrorMagic                   = 128
	ErrorUpdatePrepare           = 129
	ErrorAuthInvalid             = 130
	ErrorWritingFlash            = 131
	ErrorUpdateEnd               = 132
	ErrorInvalidBootstrapping    = 133
	ErrorWrongCurrentFlashConfig = 134
	ErrorWrongNewFlashConfig     = 135
	ErrorESP8266NotEnoughSpace   = 136
	ErrorESP32NotEnoughSpace     = 137
	ErrorUnknown                 = 255
)

var (
	errorText = map[int64]string{
		ErrorMagic:                   "Invalid magic byte",
		ErrorUpdatePrepare:           "Couldn't prepare flash memory for update. Is the binary too big? Please try restarting the ESP",
		ErrorAuthInvalid:             "Authentication invalid. Is the password correct?",
		ErrorWritingFlash:            "Wring OTA data to flash memory failed. See USB logs for more information",
		ErrorUpdateEnd:               "Finishing update failed. See the MQTT/USB logs for more information",
		ErrorInvalidBootstrapping:    "Please press the reset button on the ESP. A manual reset is required on the first OTA-Update after flashing via USB",
		ErrorWrongCurrentFlashConfig: "ESP has been flashed with wrong flash size. Please choose the correct 'board' option (esp01_1m always works) and then flash over USB",
		ErrorWrongNewFlashConfig:     "ESP does not have the requested flash size (wrong board). Please choose the correct 'board' option (esp01_1m always works) and try uploading again",
		ErrorESP8266NotEnoughSpace:   "ESP does not have enough space to store OTA file. Please try flashing a minimal firmware (remove everything except ota)",
		ErrorESP32NotEnoughSpace:     "The OTA partition on the ESP is too small. ESPHome needs to resize this partition, please flash over USB",
		ErrorUnknown:                 "Unknown error from ESP",
	}

	magicBytes  = []byte{0x6C, 0x26, 0xF7, 0x5C, 0x45}
	headerBytes = []byte{0x0}
)

type OTA struct {
	address  string
	password string
}

func NewOTA(address, password string) *OTA {
	if _, port, err := net.SplitHostPort(address); err != nil || port == "" {
		address += ":" + strconv.FormatInt(DefaultESP8266Port, 10)
	}

	return &OTA{
		address:  address,
		password: password,
	}
}

func (o *OTA) UploadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return o.Upload(f)
}

func (o *OTA) Upload(reader io.Reader) error {
	conn, err := connection.New("tcp4://" + o.address + "?dump=true&once=true")

	//conn, err := net.Dial("tcp4", o.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// >>> magic bytes
	if _, err = conn.Write(magicBytes); err != nil {
		return errors.New("error sending magic bytes: " + err.Error())
	}

	var response []byte

	// <<< version 2
	if response, err = o.receive(conn, 1, []byte{CodeOK}); err != nil {
		return err
	}

	if response[0] != Version_1_0 {
		return fmt.Errorf("unsupported OTA version %x", response[0])
	}

	// >>> features
	if _, err = conn.Write(headerBytes); err != nil {
		return errors.New("error sending features: " + err.Error())
	}

	// <<< features 1
	if response, err = o.receive(conn, 0, []byte{CodeHeaderOK}); err != nil {
		return err
	}

	// <<< auth
	if response, err = o.receive(conn, 1, nil); err != nil {
		return err
	}

	var generator hash.Hash

	if response[0] != CodeAuthOK {
		if o.password == "" {
			return errors.New("ESP requests password, but no password given")
		}

		// <<< authentication nonce 32
		if response, err = o.receive(conn, 32, nil); err != nil {
			return err
		}

		nonce := response

		fmt.Println("Auth: Nonce is", nonce)

		// >>> auth nonce
		token := make([]byte, 32)
		if _, err = rand.Read(token); err != nil {
			return errors.New("error cnonce calculate: " + err.Error())
		}

		generator = md5.New()
		if _, err = generator.Write(token); err != nil {
			return errors.New("error cnonce calculate: " + err.Error())
		}

		cnonce := generator.Sum(nil)

		if _, err = conn.Write(cnonce); err != nil {
			return errors.New("error sending auth cnonce: " + err.Error())
		}

		// >>> auth result
		generator = md5.New()
		if _, err = generator.Write([]byte(o.password)); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}
		if _, err = generator.Write(nonce); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}
		if _, err = generator.Write(cnonce); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}

		authResult := generator.Sum(nil)
		fmt.Println("Auth result", authResult)

		if _, err = conn.Write(authResult); err != nil {
			return errors.New("error sending auth result: " + err.Error())
		}

		// <<< auth result 1
		if response, err = o.receive(conn, 0, []byte{CodeAuthOK}); err != nil {
			return err
		}
	}

	// >>> binary size
	size := bufio.NewReader(reader).Size()
	sizeBytes := []byte{
		byte((size >> 24) & 0xFF),
		byte((size >> 16) & 0xFF0),
		byte((size >> 8) & 0xFF),
		byte((size >> 0) & 0xFF),
	}

	if _, err = conn.Write(sizeBytes); err != nil {
		return errors.New("error sending binary size: " + err.Error())
	}

	// <<< binary size 1
	if response, err = o.receive(conn, 0, []byte{CodeUpdatePrepare}); err != nil {
		return err
	}

	// >>> file checksum
	generator = md5.New()
	if _, err = io.Copy(generator, reader); err != nil {
		return errors.New("error file checksum calculate: " + err.Error())
	}

	if _, err = conn.Write(generator.Sum(nil)); err != nil {
		return errors.New("error sending file checksum: " + err.Error())
	}

	// <<< checksum 1
	if response, err = o.receive(conn, 0, []byte{CodeBinMD5}); err != nil {
		return err
	}

	// >>> byte as byte
	if _, err = io.Copy(conn, reader); err != nil {
		return errors.New("error sending data: " + err.Error())
	}

	fmt.Println("Waiting results")

	// <<< receive OK 1
	if response, err = o.receive(conn, 0, []byte{CodeReceiveOK}); err != nil {
		return err
	}

	// <<< Update end 1
	if response, err = o.receive(conn, 0, []byte{CodeUpdateEnd}); err != nil {
		return err
	}

	// >>> end acknowledgement
	if _, err = conn.Write([]byte{CodeOK}); err != nil {
		return errors.New("error sending end acknowledgement: " + err.Error())
	}

	return nil
}

func (o *OTA) receive(conn connection.Conn, amount int, expect []byte) ([]byte, error) {
	l := len(expect)

	buf := make([]byte, l+amount)
	var shift int

	start := l
	if start < 1 {
		start = 1
	}

	n, err := conn.Read(buf[:start])
	if err != nil {
		return nil, err
	}

	if n != start {
		return nil, fmt.Errorf("unexpected response start bytes %d expected %d", n, start)
	}
	shift += n

	// check first byte
	if text, ok := errorText[int64(buf[0])]; ok {
		return nil, errors.New(text)
	}

	// check expect bytes
	if l > 0 && !bytes.Equal(buf[:l], expect) {
		return nil, fmt.Errorf("unexpected response from ESP: %x expected %x", buf, expect)
	}

	for {
		if shift >= l+amount {
			break
		}

		n, err = conn.Read(buf[shift:])

		if err != nil {
			return nil, err
		}

		shift += n
	}

	// skip expect bytes
	return buf[l:], nil
}
