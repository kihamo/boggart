package broadlink

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/kihamo/boggart/providers/broadlink/internal"
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

func LocalAddr() (*net.UDPAddr, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return &net.UDPAddr{
					IP: ipnet.IP,
				}, err
			}
		}
	}

	return nil, errors.New("IP not found")
}

// Порядок подключения:
// - Долго зажать кнопку на устройстве (частое моргание - пауза - частое мограние)
// - Подключиться к незащищенной точке BroadlinkProv
// - Запустить Setup
// - Подключиться к сети указанной в Setup и через Discover обнаружить устройство
func SetupWiFi(SSID, password string, securityMode int) (err error) {
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
	binary.LittleEndian.PutUint16(packet[0x20:], internal.Checksum(packet[:]))

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = conn.WriteTo(packet[:], broadCastAddr); err != nil {
		return err
	}

	return nil
}

/*
func DiscoverDevices() (devices []Device, err error) {
		//Offset       Contents
		//0x00-0x07    00
		//0x08-0x0b    Current offset from GMT as a little-endian 32 bit integer
		//0x0c-0x0d    Current year as a little-endian 16 bit integer
		//0x0e         Current number of seconds past the minute
		//0x0f         Current number of minutes past the hour
		//0x10         Current number of hours past midnight
		//0x11         Current day of the week (Monday = 1, Tuesday = 2, etc)
		//0x12         Current day in month
		//0x13         Current month
		//0x14-0x17    00
		//0x18-0x1b    Local IP address
		//0x1c-0x1d    Source port as a little-endian 16 bit integer
		//0x1e-0x1f    00
		//0x20-0x21    Checksum as a little-endian 16 bit integer
		//0x22-0x25    00
		//0x26         06
		//0x27-0x2f    00

	addrInterface, err := LocalAddr()
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", addrInterface)
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
	ip4 := addrInterface.IP.To4()
	packet[0x18], packet[0x19], packet[0x1a], packet[0x1b] = ip4[3], ip4[2], ip4[1], ip4[0]

	// 0x1c-0x1d Source port as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x1c:], uint16(address.Port))

	packet[0x26] = 0x06

	// 0x20-0x21 Checksum as a little-endian 16 bit integer
	binary.LittleEndian.PutUint16(packet[0x20:], internal.Checksum(packet[:]))

	// send request packet
	err = conn.SetDeadline(time.Now().Add(internal.DefaultTimeout))
	if err != nil {
		return nil, err
	}

	if _, err = conn.WriteTo(packet[:], broadCastAddr); err != nil {
		return nil, err
	}

	// read response packet
	devices = make([]Device, 0)
	response := make([]byte, internal.DefaultBufferSize)

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

		devices = append(devices, NewDevice(int(r[0x35])<<8|int(r[0x34]), mac, *addr, *addrInterface))
	}

	return devices, nil
}
*/
