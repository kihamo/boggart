package esphome

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
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

	connChunkSize     = 1024
	connWriteBuffer   = 8192
	connTimeout       = time.Second * 20
	sleepAfterSuccess = time.Second * 10
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

	running   uint32
	writen    uint64
	total     uint64
	checksum  atomic.Value
	lastError atomic.Value
}

type hasher struct {
	hash.Hash
}

func (h hasher) Bytes() []byte {
	return []byte(hex.EncodeToString(h.Sum(nil)))
}

func (h hasher) String() string {
	return hex.EncodeToString(h.Sum(nil))
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

func (o *OTA) IsRunning() bool {
	return atomic.LoadUint32(&o.running) == 1
}

func (o *OTA) Checksum() string {
	if cs := o.checksum.Load(); cs != nil {
		return cs.(string)
	}

	return ""
}

func (o *OTA) Progress() (uint64, uint64) {
	return atomic.LoadUint64(&o.writen), atomic.LoadUint64(&o.total)
}

func (o *OTA) LastError() error {
	if err := o.lastError.Load(); err != nil {
		if text := err.(string); text != "" {
			return errors.New(text)
		}
	}

	return nil
}

func (o *OTA) UploadAsync(reader io.Reader) error {
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, reader); err != nil {
		return errors.New("error get file data: " + err.Error())
	}

	return o.doAsync(buf, buf.Len())
}

func (o *OTA) UploadAsyncFromFile(filename string) error {
	reader, size, err := o.readerFromFile(filename)
	if err != nil {
		return err
	}

	return o.doAsync(reader, size)
}

func (o *OTA) Upload(reader io.Reader) error {
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, reader); err != nil {
		return errors.New("error get file data: " + err.Error())
	}

	return o.doUpload(buf, buf.Len())
}

func (o *OTA) UploadFromFile(filename string) error {
	reader, size, err := o.readerFromFile(filename)
	if err != nil {
		return err
	}

	return o.doUpload(reader, size)
}

func (o *OTA) readerFromFile(filename string) (io.Reader, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, -1, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, -1, err
	}

	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, -1, err
	}

	f.Close()
	return buf, int(stat.Size()), nil
}

func (o *OTA) doAsync(reader io.Reader, size int) error {
	if o.IsRunning() {
		return errors.New("OTA proccess aleady running")
	}

	go func() {
		atomic.StoreUint32(&o.running, 1)
		atomic.StoreUint64(&o.total, uint64(size))
		atomic.StoreUint64(&o.writen, 0)

		defer func() {
			atomic.StoreUint32(&o.running, 0)
			atomic.StoreUint64(&o.total, 0)
			atomic.StoreUint64(&o.writen, 0)
			o.checksum.Store("")
		}()

		err := o.doUpload(reader, size)

		if err != nil {
			o.lastError.Store(err.Error())
		} else {
			o.lastError.Store("")
		}
	}()

	return nil
}

func (o *OTA) doUpload(reader io.Reader, size int) error {
	if size == 0 {
		return errors.New("error sending data: file is empty")
	}

	// check connect
	conn, err := net.Dial("tcp4", o.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)

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

	hashMD5 := hasher{
		Hash: md5.New(),
	}

	if response[0] != CodeAuthOK {
		if o.password == "" {
			return errors.New("ESP requests password, but no password given")
		}

		// <<< authentication nonce 32
		if response, err = o.receive(conn, 32, nil); err != nil {
			return err
		}

		nonce := response

		// >>> auth nonce
		token := make([]byte, 32)
		if _, err = rand.Read(token); err != nil {
			return errors.New("error cnonce calculate: " + err.Error())
		}

		if _, err = hashMD5.Write(token); err != nil {
			return errors.New("error cnonce calculate: " + err.Error())
		}

		cnonce := hashMD5.Bytes()
		if _, err = conn.Write(cnonce); err != nil {
			return errors.New("error sending auth cnonce: " + err.Error())
		}

		// >>> auth result
		hashMD5.Reset()
		if _, err = hashMD5.Write([]byte(o.password)); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}
		if _, err = hashMD5.Write(nonce); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}
		if _, err = hashMD5.Write(cnonce); err != nil {
			return errors.New("error auth result calculate: " + err.Error())
		}

		if _, err = conn.Write(hashMD5.Bytes()); err != nil {
			return errors.New("error sending auth result: " + err.Error())
		}

		// <<< auth result 1
		if response, err = o.receive(conn, 0, []byte{CodeAuthOK}); err != nil {
			return err
		}
	}

	// >>> binary size
	sizeBytes := []byte{
		byte((size >> 24) & 0xFF),
		byte((size >> 16) & 0xFF),
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
	body := bytes.NewBuffer(nil)

	hashMD5.Reset()
	if _, err = io.Copy(hashMD5, io.TeeReader(reader, body)); err != nil {
		return errors.New("error file checksum calculate: " + err.Error())
	}

	o.checksum.Store(hashMD5.String())

	if _, err = conn.Write(hashMD5.Bytes()); err != nil {
		return errors.New("error sending file checksum: " + err.Error())
	}

	if body.Len() == 0 {
		return errors.New("error sending data: file is empty")
	}

	// <<< checksum 1
	if response, err = o.receive(conn, 0, []byte{CodeBinMD5}); err != nil {
		return err
	}

	// >>> byte as byte
	if err = tcpConn.SetNoDelay(false); err != nil {
		return errors.New("set connection option failed: " + err.Error())
	}

	if err = tcpConn.SetWriteBuffer(connWriteBuffer); err != nil {
		return errors.New("set connection option failed: " + err.Error())
	}

	if err = tcpConn.SetDeadline(time.Now().Add(connTimeout)); err != nil {
		return errors.New("set connection option failed: " + err.Error())
	}

	var offset, n int

	for {
		chunk := body.Next(connChunkSize)
		if len(chunk) == 0 {
			break
		}

		if n, err = conn.Write(chunk); err != nil {
			return errors.New("error sending data: " + err.Error())
		}

		offset += n
		atomic.StoreUint64(&o.writen, uint64(offset))
	}

	if err = tcpConn.SetNoDelay(true); err != nil {
		return errors.New("set connection option failed: " + err.Error())
	}

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

	// wait reload
	time.Sleep(sleepAfterSuccess)

	return nil
}

func (o *OTA) receive(conn net.Conn, amount int, expect []byte) ([]byte, error) {
	l := len(expect)

	buf := make([]byte, l+amount)
	var offset int

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
	offset += n

	// check first byte
	if text, ok := errorText[int64(buf[0])]; ok {
		return nil, errors.New(text)
	}

	// check expect bytes
	if l > 0 && !bytes.Equal(buf[:l], expect) {
		return nil, fmt.Errorf("unexpected response from ESP: %x expected %x", buf, expect)
	}

	for {
		if offset >= l+amount {
			break
		}

		n, err = conn.Read(buf[offset:])

		if err != nil {
			return nil, err
		}

		offset += n
	}

	// skip expect bytes
	return buf[l:], nil
}
