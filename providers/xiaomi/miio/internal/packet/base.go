package packet

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

const (
	MaxBufferSize = 1472
)

var magicNumber = []byte{0x21, 0x31}

// https://github.com/OpenMiHome/mihome-binary-protocol/blob/master/doc/PROTOCOL.md
// https://github.com/skysilver-lab/php-miio

type Packet interface {
	io.WriterTo
	io.ReaderFrom

	DeviceID() []byte
	Stamp() time.Time
	Body() []byte
	Dump() string
}

type Base struct {
	Header Header
	body   []byte
}

type Header struct {
	Length   uint16
	Unknown1 []byte
	DeviceID []byte
	Stamp    uint32
	Checksum []byte
}

func NewBase() *Base {
	b := &Base{
		Header: Header{
			Length:   0,
			Unknown1: bytes.Repeat([]byte{0x0}, 4),
		},
	}
	b.SetStamp(time.Now())

	return b
}

func (p *Base) DeviceID() []byte {
	return p.Header.DeviceID
}

func (p *Base) Stamp() time.Time {
	return time.Unix(int64(p.Header.Stamp), 0)
}

func (p *Base) Body() []byte {
	return p.body
}

func (p *Base) Dump() string {
	buf := bytes.NewBuffer(nil)
	p.WriteTo(buf)

	return hex.EncodeToString(buf.Bytes())
}

func (p *Base) build(checksum []byte) []byte {
	req := bytes.NewBuffer(nil)

	// Magic number: 16 bits
	req.Write(magicNumber)

	buf := bytes.NewBuffer(nil)

	// Packet length: 16 bits unsigned int
	binary.Write(buf, binary.LittleEndian, p.Header.Length)
	req.Write(reverse(buf.Bytes()))

	// Unknown1: 32 bits
	req.Write(p.Header.Unknown1)

	// Device ID: 32 bits
	req.Write(p.Header.DeviceID)

	// Stamp: 32 bit unsigned int
	buf.Reset()
	binary.Write(buf, binary.LittleEndian, p.Header.Stamp)
	req.Write(reverse(buf.Bytes()))

	// MD5 checksum
	buf.Reset()
	binary.Write(buf, binary.LittleEndian, checksum)
	req.Write(buf.Bytes())

	if len(p.body) > 0 {
		req.Write(p.body)
	}

	return req.Bytes()
}

func (p *Base) WriteTo(w io.Writer) (int64, error) {
	if p.Header.Checksum == nil {
		hash := md5.New()

		_, err := hash.Write(p.build(bytes.Repeat([]byte{0x00}, 16)))
		if err != nil {
			return -1, err
		}

		p.Header.Checksum = hash.Sum(nil)
	}

	n, err := w.Write(p.build(p.Header.Checksum))

	return int64(n), err
}

func (p *Base) ReadFrom(r io.Reader) (int64, error) {
	var result int

	buf := make([]byte, 2)

	// Magic number: 16 bits
	n, err := r.Read(buf)
	if err != nil {
		return -1, err
	}

	if !bytes.Equal(buf, magicNumber) {
		return -1, errors.New("magic number could not be verified. Expected 0x2131, got 0x" + hex.EncodeToString(buf))
	}

	result += n

	p.Header = Header{}

	// Packet length: 16 bits unsigned int
	n, err = r.Read(buf)
	if err != nil {
		return -1, err
	}

	result += n
	p.Header.Length = binary.LittleEndian.Uint16(reverse(buf))

	// Unknown1: 32 bits
	buf = make([]byte, 4)

	n, err = r.Read(buf)
	if err != nil {
		return -1, err
	}

	p.Header.Unknown1 = reverse(buf)
	result += n

	// Device ID: 32 bits
	n, err = r.Read(buf)
	if err != nil {
		return -1, err
	}

	p.Header.DeviceID = append([]byte(nil), buf...)
	result += n

	// Stamp: 32 bit unsigned int
	n, err = r.Read(buf)
	if err != nil {
		return -1, err
	}

	result += n
	p.Header.Stamp = binary.LittleEndian.Uint32(reverse(buf))

	// MD5 checksum
	buf = make([]byte, 16)

	n, err = r.Read(buf)
	if err != nil {
		return -1, err
	}

	p.Header.Checksum = append([]byte(nil), buf...)
	result += n

	buf = make([]byte, MaxBufferSize)
	n, _ = r.Read(buf)

	p.body = buf[:n]
	result += n

	return int64(result), nil
}

func (p *Base) SetBody(body []byte) error {
	p.Header.Length = uint16(len(body) + 32)
	p.body = body

	checksum := md5.New()

	_, err := p.WriteTo(checksum)
	if err != nil {
		return err
	}

	p.Header.Checksum = checksum.Sum(nil)

	return nil
}

func (p *Base) SetStamp(stamp time.Time) {
	p.Header.Stamp = uint32(stamp.Unix())
}

func reverse(data []byte) []byte {
	dataNew := make([]byte, len(data))

	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		dataNew[i], dataNew[j] = data[j], data[i]
	}

	return dataNew
}
