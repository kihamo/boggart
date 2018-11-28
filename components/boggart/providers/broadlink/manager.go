package broadlink

import (
	"encoding/binary"
	"net"
	"time"
)

const (
	WifiSecurityNone = iota
	WifiSecurityWEP
	WifiSecurityWPA1
	WifiSecurityWPA2
	WifiSecurityWPACCMP
	WifiSecurityUnknown
	WifiSecurityWPATKIP
)

const (
	DevicePort = 80
)

var broadCastAddr = &net.UDPAddr{
	IP:   net.IPv4bcast,
	Port: DevicePort,
}

// Порядок подключения:
// - Долго зажать кнопку на устройстве (частое моргание - пауза - частое мограние)
// - Подключиться к незащищенной точке BroadlinkProv
// - Запустить Setup
// - Подключиться к сети указанной в Setup и через Discover обнаружить устройство
func SetupWiFi(SSID, password string, securityMode int) (err error) {
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	var packet [0x88]byte

	// 0x26 14 (Always 14)
	packet[0x26] = 0x14

	// 0x44-0x63 SSID Name (zero padding is appended)
	copy(packet[0x44:], SSID)

	// 0x64-0x83 Password (zero padding is appended)
	if len(password) > 0x1f { // WIFI password
		password = password[:0x1f]
	}

	copy(packet[0x64:], password)

	// 0x84	Character length of SSID
	packet[0x84] = byte(len(SSID))

	// 0x85	Character length of password
	packet[0x85] = byte(len(password))

	// 0x86	Wireless security mode (00 - none, 01 = WEP, 02 = WPA1, 03 = WPA2, 04 = WPA1/2)
	packet[0x86] = byte(securityMode)

	// 0x20-0x21 Checksum as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x20:], checksum(packet[:]))

	if _, err = conn.WriteTo(packet[:], broadCastAddr); err != nil {
		return err
	}

	return nil
}

func DiscoverDevices() (devices []*Device, err error) {
	ifaceAddr, err := LocalAddr()
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", ifaceAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	address := conn.LocalAddr().(*net.UDPAddr)

	var packet [0x30]byte

	now := time.Now()

	// 0x08-0x0b Current tz from GMT as a little-endian 32 bit integer
	_, tz := now.Zone()
	tz /= 3600

	binary.LittleEndian.PutUint32(packet[0x08:], uint32(tz))

	// 0x0c-0x0d Current year as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x0c:], uint16(now.Year()))

	// 0x0e Current number of seconds past the minute
	packet[0x0e] = byte(now.Second())

	// 0x0f Current number of minutes past the hour
	packet[0x0f] = byte(now.Minute())

	// 0x10	Current number of hours past midnight
	packet[0x10] = byte(now.Hour())

	// 0x11 Current day of the week (Monday = 1, Tuesday = 2, etc)
	packet[0x11] = byte(now.Weekday())

	// 0x12 Current day in month
	packet[0x12] = byte(now.Day())

	// 0x13 Current month
	packet[0x13] = byte(now.Month())

	// 0x18-0x1b Local IP address
	ip4 := ifaceAddr.IP.To4()
	packet[0x18], packet[0x19], packet[0x1a], packet[0x1b] = ip4[3], ip4[2], ip4[1], ip4[0]

	// 0x1c-0x1d Source port as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x1c:], uint16(address.Port))

	packet[0x26] = 0x06

	// 0x20-0x21 Checksum as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x20:], checksum(packet[:]))

	// send request packet
	err = conn.SetDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return nil, err
	}

	if _, err = conn.WriteTo(packet[:], broadCastAddr); err != nil {
		return nil, err
	}

	// read response packet
	devices = make([]*Device, 0)
	response := make([]byte, 2048)

	for {
		size, addr, err := conn.ReadFromUDP(response)
		if err != nil {
			e, ok := err.(net.Error)
			if ok && e.Timeout() {
				return devices, nil
			}
		}

		if size == 0 {
			continue
		}

		r := response[:size]

		//if !checksumCheck(r) {
		//continue
		//}

		if r[0x26] != 0x07 {
			continue
		}

		//ip := net.IPv4(r[0x39], r[0x38], r[0x37], r[0x36])
		//if bytes.Compare(addr.IP.To4(), ip.To4()) != 0 {
		//	continue
		//}

		// 0x3a-0x3f MAC address of the target device
		var mac net.HardwareAddr
		for i := 0x3f; i >= 0x3a; i-- {
			mac = append(mac, r[i])
		}

		devices = append(devices, NewDevice(int(r[0x35])<<8|int(r[0x34]), mac, *addr, *ifaceAddr))
	}

	return devices, nil
}
