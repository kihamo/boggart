package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultTimeout = time.Second
	DefaultPort    = 80

	CommandAuth = 0x0065
)

var (
	aesKey      = []byte{0x09, 0x76, 0x28, 0x34, 0x3f, 0xe9, 0x9e, 0x23, 0x76, 0x5c, 0x15, 0x13, 0xac, 0xcf, 0x8b, 0x02}
	aesIV       = []byte{0x56, 0x2e, 0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58}
	aesBlock, _ = aes.NewCipher(aesKey)

	regularPacketHeader = []byte{0x5a, 0xa5, 0xaa, 0x55, 0x5a, 0xa5, 0xaa, 0x55}
	readBufferPool      sync.Pool
)

func init() {
	readBufferPool.New = func() interface{} {
		buf := make([]byte, DefaultBufferSize)
		return &buf
	}
}

type Device struct {
	timeout        int64
	kind           int
	id             uint32
	packetsCounter uint64

	mac  net.HardwareAddr
	host string

	aesKey   []byte
	aesIV    []byte
	aesBlock cipher.Block

	requestMutex sync.Mutex
}

func NewDevice(kind int, mac net.HardwareAddr, host string) *Device {
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = host + ":" + strconv.Itoa(DefaultPort)
	}

	d := &Device{
		kind: kind,
		mac:  mac,
		host: host,
	}
	d.SetTimeout(DefaultTimeout)

	return d
}

func (d *Device) SetTimeout(duration time.Duration) {
	atomic.StoreInt64(&d.timeout, int64(duration))
}

func (d *Device) Timeout() time.Duration {
	return time.Duration(atomic.LoadInt64(&d.timeout))
}

func (d *Device) ID() uint32 {
	return atomic.LoadUint32(&d.id)
}

func (d *Device) setID(id uint32) {
	atomic.StoreUint32(&d.id, id)
}

func (d *Device) Kind() int {
	return d.kind
}

func (d *Device) MAC() net.HardwareAddr {
	return d.mac
}

func (d *Device) Host() string {
	return d.host
}

func (d *Device) request(cmd byte, payload []byte, waitResult bool) ([]byte, uint64, error) {
	// auto auth
	if cmd != CommandAuth && d.ID() == 0 {
		if err := d.Auth(nil, ""); err != nil {
			return nil, 0, err
		}
	}

	d.requestMutex.Lock()
	defer d.requestMutex.Unlock()

	conn, err := net.Dial("udp", d.host)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	connUDP, ok := conn.(*net.UDPConn)
	if !ok {
		return nil, 0, errors.New("failed cast connect to *net.UDPConn")
	}

	err = connUDP.SetDeadline(time.Now().Add(d.Timeout()))
	if err != nil {
		return nil, 0, err
	}

	requestPacket, requestID := d.buildCmdPacket(cmd, payload)

	_, err = connUDP.Write(requestPacket)
	if err != nil {
		return nil, requestID, err
	}

	if !waitResult {
		return nil, requestID, nil
	}

	buf := readBufferPool.Get().(*[]byte)
	defer readBufferPool.Put(buf)

	size, _, err := connUDP.ReadFromUDP(*buf)
	if err != nil {
		return nil, requestID, err
	}

	return append([]byte(nil), (*buf)[:size]...), requestID, nil
}

func (d *Device) Cmd(cmd byte, payload []byte) error {
	_, _, err := d.request(cmd, payload, false)
	return err
}

func (d *Device) Call(cmd byte, payload []byte) ([]byte, error) {
	result, pcWant, err := d.request(cmd, payload, true)
	if err != nil {
		return nil, err
	}

	// TODO: Checksum check

	// verify packet counter
	pcHave := binary.LittleEndian.Uint16(result[0x28:])
	if uint16(pcWant) != pcHave {
		return nil, fmt.Errorf("invalid packet counter want %d have %d", pcWant, pcHave)
	}

	// verify MAC address
	for i, part := range d.mac {
		if result[0x2a+i] != part {
			return nil, errors.New("invalid MAC address")
		}
	}

	return result, nil
}

/*
	Offset        Contents
	0x00-0x03     00
	0x04-0x12     A 15-digit value that represents this device. Broadlink's implementation uses the IMEI.
	0x13          01
	0x14-0x2c     00
	0x2d          0x01
	0x30-0x7f     NULL-terminated ASCII string containing the device name
*/
func (d *Device) Auth(id []byte, name string) error {
	if len(id) == 0 {
		id = []byte{0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31}
	}

	if name == "" {
		name = d.mac.String()
	}

	if len(id) != 15 {
		return errors.New("device id size must be 15 bytes long")
	}

	nameLength := len(name)
	if nameLength > 0 {
		nameLength--
	}

	nameLength = (nameLength / 16) * 16

	size := 0x30 + nameLength
	if size < 0x50 {
		size = 0x50
	} else if size > 0x80 {
		size = 0x80
	}

	payload := make([]byte, size)
	copy(payload[0x04:0x13], id)
	payload[0x2d] = 0x01
	copy(payload[0x30:], name)

	response, err := d.Call(CommandAuth, payload)
	if err != nil {
		return err
	}

	data, err := d.DecodePacket(response)
	if err != nil {
		return err
	}

	/*
		Offset	     Contents
		0x00-0x03    Device ID
		0x04-0x13    Device encryption key
	*/
	d.setID(binary.LittleEndian.Uint32(data[:0x04]))
	d.setAESKey(data[0x04:0x14])

	return nil
}

func (d *Device) DecodePacket(packet []byte) ([]byte, error) {
	if len(packet) <= 0x38 {
		return nil, errors.New("blank response")
	}

	return d.decrypt(packet[0x38:]), nil
}

/*
	Offset       Contents
	0x00         0x5a
	0x01         0xa5
	0x02         0xaa
	0x03         0x55
	0x04         0x5a
	0x05         0xa5
	0x06         0xaa
	0x07         0x55
	0x08-0x1f    00
	0x20-0x21    Checksum of full packet as a little-endian 16 bit integer
	0x22-0x23    00
	0x24-0x25    Device type as a little-endian 16 bit integer
	0x26-0x27    Command code as a little-endian 16 bit integer
	0x28-0x29    Packet count as a little-endian 16 bit integer
	0x2a-0x2f    Local MAC address
	0x30-0x33    Local device ID (obtained during authentication, 00 before authentication)
	0x34-0x35    Checksum of unencrypted payload as a little-endian 16 bit integer
	0x36-0x37    00
*/
func (d *Device) buildCmdPacket(cmd byte, payload []byte) (packet []byte, packetID uint64) {
	packet = make([]byte, 0x38)

	// Build header
	copy(packet, regularPacketHeader)
	packet[0x24], packet[0x25] = 0x2a, 0x27
	packet[0x26] = cmd

	d.packetsCounter++
	packetID = d.packetsCounter
	binary.LittleEndian.PutUint16(packet[0x28:], uint16(packetID))

	for i, part := range d.mac {
		packet[0x2a+i] = part
	}

	binary.LittleEndian.PutUint32(packet[0x30:], d.ID())
	binary.LittleEndian.PutUint16(packet[0x34:], Checksum(payload))

	encrypted := d.encrypt(payload)
	packet = append(packet, encrypted...)
	binary.LittleEndian.PutUint16(packet[0x20:], Checksum(packet))

	return
}

func (d *Device) encrypt(data []byte) []byte {
	return blockCipher(cipher.NewCBCEncrypter(d.cipherParam()), data)
}

func (d *Device) decrypt(data []byte) []byte {
	return blockCipher(cipher.NewCBCDecrypter(d.cipherParam()), data)
}

func (d *Device) setAESKey(key []byte) {
	if len(key) != aes.BlockSize || bytes.Equal(key, aesKey) {
		key = nil
	}

	if key == nil {
		d.aesKey, d.aesBlock = nil, nil
		return
	}

	d.aesKey = make([]byte, len(key))
	copy(d.aesKey, key)
	d.aesBlock, _ = aes.NewCipher(d.aesKey)
}

/*
func (d *Device) getAESKey() []byte {
	if d.aesKey == nil {
		k := make([]byte, len(aesKey))
		copy(k, aesKey)
		return k
	}

	return d.aesKey
}
*/
func (d *Device) cipherParam() (block cipher.Block, iv []byte) {
	if d.aesIV != nil {
		iv = d.aesIV
	} else {
		iv = aesIV
	}

	if d.aesBlock != nil {
		block = d.aesBlock
	} else {
		block = aesBlock
	}

	return
}
