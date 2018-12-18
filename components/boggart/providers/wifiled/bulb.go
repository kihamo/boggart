package wifiled

// https://connect.smartliving.ru/profile/1502/blog61.html

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Bulb struct {
	host string
}

func NewBulb(host string) *Bulb {
	return &Bulb{
		host: net.JoinHostPort(host, strconv.Itoa(PortControlLocal)),
	}
}

func (l *Bulb) request(ctx context.Context, data []byte, length int) ([]byte, error) {
	var d net.Dialer

	conn, err := d.DialContext(ctx, "tcp", l.host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// calculate checksum
	sum := byte(0)
	for _, b := range data {
		sum += b
	}

	data = append(data, sum)

	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	if length <= 0 {
		return nil, nil
	}

	response := make([]byte, length)

	size, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	if size != length {
		return nil, fmt.Errorf("bad response size have %d want %d", size, length)
	}

	return response, nil
}

func (l *Bulb) Host() string {
	return l.host
}

func (l *Bulb) Time(ctx context.Context) (*time.Time, error) {
	/*
		0x11      command of requesting time :
		Send      [0X11][0X1A][0X1B][0xF0 remote,0x0F local][check digit] / 5 bytes
		Return    [0xF0 remote,0x0F local][0X11][0x14][10 digit and single digit of year][month][day][hour][minute][second][week][reserved for future use:0x0][check digit] / 12 bytes
	*/
	response, err := l.request(ctx, []byte{CommandTimeGet, CommandTimeGet2, CommandTimeGet3, CommandLocal}, 12)
	if err != nil {
		return nil, err
	}

	// check first two byte
	if response[0] != CommandLocal {
		return nil, fmt.Errorf("bad first byte of response have %#x want %#x", response[0], CommandLocal)
	}

	if response[1] != CommandTimeGet {
		return nil, fmt.Errorf("bad second byte of response have %#x want %#x", response[0], CommandLocal)
	}

	t := time.Date(int(response[3])+2000, time.Month(response[4]), int(response[5]), int(response[6]), int(response[7]), int(response[8]), 0, time.Local)

	return &t, nil
}

func (l *Bulb) SetTime(ctx context.Context, t time.Time) error {
	/*
		0x10      command of syncing time:
		Send      [0X10][0x14][10 digit and single digit of year][month][day][hour][minute][second][week][reserved for future use:0x0][0xF0 remote,0x0F local][check digit】/ 12 bytes
		Return    [0xF0 remote,0x0F local][0X10][0x00][check digit] / 4 bytes
	*/
	year, week := t.ISOWeek()

	_, err := l.request(ctx, []byte{CommandTimeSet, CommandTimeSet2, byte(year - 2000), byte(t.Month()), byte(t.Day()), byte(t.Hour()), byte(t.Minute()), byte(t.Second()), byte(week), 0x0, CommandLocal}, 4)
	if err != nil {
		return err
	}

	return nil
}

func (l *Bulb) SetColorPersist(ctx context.Context, color Color) error {
	/*
			0x31      command of setting color and color temperature
			Send      [0X31][8bit red data][8bit green data][8bit blue data][8bit warm white data][8bit status sign][0xF0 remote,0x0F local][check digit] / 8 bytes
			Return    Local(0x0F): no return
		              Remote(0xF0): [0xF0 remote][0X31][x00][check digit]
		              Status sign: [0XF0] means changing RGB, [0X0F] means W
		              Note:phone send commands which control static color.
		              Range of static color value is 00-0xff.
		              When value is 0, PWM is 0%; when value is 0XFF, PWM is 100%;
	*/
	var status int
	if color.UseRGB {
		status = True
	} else {
		status = False
	}

	_, err := l.request(ctx, []byte{CommandColorSetPersist, byte(color.Red), byte(color.Green), byte(color.Blue), byte(color.WarmWhite), byte(status), CommandLocal}, 0)
	return err
}

func (l *Bulb) SetColor(ctx context.Context, color Color) error {
	/*
		0x41      command of setting color,color temperature(in music mode,this value will not be saved)
		Send      [0X41][8bit red value][8bit green data][8bit blue data][8bit warm white data][8bit status sign][0x0F local][check digit] / 8 bytes
		Return    No return
		          Status sign:[0XF0] means changing RGB, [0X0F] means changing W
	*/
	var status int
	if color.UseRGB {
		status = True
	} else {
		status = False
	}

	_, err := l.request(ctx, []byte{CommandColorSet, byte(color.Red), byte(color.Green), byte(color.Blue), byte(color.WarmWhite), byte(status), CommandLocal}, 0)
	return err
}

func (l *Bulb) SetMode(ctx context.Context, mode Mode, speed uint8) error {
	/*
		0x61      command of setting builted-in mode
		Send      [0x61][8bit mode value][8bit speed value][0xF0 remote,0x0F local][check digit] / 5 bytes
		Return    If command is local(0x0F):no return
		          If command is remote (0xF0): [0xF0 remote][0X61][0x00][check digit]
		          Note:mode value refers to definition in the end of file,range of speed value is 0x01--0x1F
	*/
	_, err := l.request(ctx, []byte{CommandMode, byte(mode), byte(speed), CommandLocal}, 0)
	return err
}

func (l *Bulb) PowerOn(ctx context.Context) error {
	/*
		0x71      command of setting key's value(switcher command) command
		Send      [0X71][8bit value][0xF0remote,0x0F local][check digit] / 4 bytes
		Return    [0xF0remote,0x0F local][0X71][switcher status value][check digit]
		          Note:key value0x23 means "turn on",0x24 means "turn off"
	*/
	_, err := l.request(ctx, []byte{CommandPower, PowerOn, CommandLocal}, 0)
	return err
}

func (l *Bulb) PowerOff(ctx context.Context) error {
	/*
		0x71      command of setting key's value(switcher command) command
		Send      [0X71][8bit value][0xF0remote,0x0F local][check digit] / 4 bytes
		Return    [0xF0remote,0x0F local][0X71][switcher status value][check digit]
		          Note:key value0x23 means "turn on",0x24 means "turn off"
	*/
	_, err := l.request(ctx, []byte{CommandPower, PowerOff, CommandLocal}, 0)
	return err
}

func (l *Bulb) State(ctx context.Context) (*State, error) {
	/*
		0x81      command of requesting devices'status
		Send      [0X81][0X8A][0X8B][check digit] / 4 bytes
		Return    [0X81][8bit device name][8bit turn on/off][8bit mode value][8bit run/pause][8bit speed value][8bit red value][8bit green data][8bit blue data][8bit warm white data][version NO][8bit cool white data][8bit status sign][check digit] / 14 bytes
		          Note:when module received command of checking devices's status, module will reply,
		          [8bit turn on/off]
		              0x23 means turn on
		              0x24 means turn off
		          [8bit run/pause status]
		              0x20 means status in present
		              0x21 means pause status, it is unuseful in this item
		          [8bit speed value] means speed value of dynamic model
		              range:0x01-0x1f,
		              0x01 is the fast
		          [0XF0] Status sign means RGB
		          [0X0F] means W
	*/

	response, err := l.request(ctx, []byte{CommandState, CommandState2, CommandState3}, 14)
	if err != nil {
		return nil, err
	}

	// check first byte
	if response[0] != CommandState {
		return nil, fmt.Errorf("bad first byte of response have %#x want %#x", response[0], CommandState)
	}

	result := &State{
		DeviceName: uint8(response[1]),
		Mode:       Mode(response[3]),
		Speed:      uint8(response[5]),
		Color: Color{
			Red:       uint8(response[6]),
			Green:     uint8(response[7]),
			Blue:      uint8(response[8]),
			WarmWhite: uint8(response[9]),
		},
	}

	switch response[2] {
	case 0x23:
		result.Power = true
	case 0x24:
		result.Power = false
	default:
		return nil, fmt.Errorf("unknown power value %#x", response[2])
	}

	if !((result.Mode >= ModePreset1 && result.Mode <= ModePreset21) || (result.Mode >= ModeCustom && result.Mode <= ModeTesting)) {
		return nil, fmt.Errorf("unknown mode value %#x", response[3])
	}

	if response[12] == False {
		result.Color.UseWarmWhite = true
	} else {
		result.Color.UseRGB = true
	}

	return result, nil
}
