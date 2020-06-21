package z_stack

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Frame struct {
	length    uint16
	typ       uint16
	subSystem uint16
	commandID uint16
	data      []byte
	fcs       byte

	lock sync.RWMutex
}

func (f *Frame) Command0() uint16 {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return ((f.typ << 5) & 0xE0) | (f.subSystem & 0x1F)
}

func (f *Frame) SetCommand0(value uint16) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.typ = (value & 0xE0) >> 5
	f.subSystem = value & 0x1F
}

func (f *Frame) Command1() uint16 {
	return f.CommandID()
}

func (f *Frame) SetCommand1(value uint16) {
	f.SetCommandID(value)
}

func (f *Frame) Length() uint16 {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.length
}

func (f *Frame) Type() uint16 {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.typ
}

func (f *Frame) SetType(value uint16) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.typ = value
}

func (f *Frame) SubSystem() uint16 {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.subSystem
}

func (f *Frame) SetSubSystem(value uint16) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.subSystem = value
}

func (f *Frame) CommandID() uint16 {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.commandID
}

func (f *Frame) SetCommandID(value uint16) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.commandID = value
}

func (f *Frame) Data() []byte {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return append([]byte(nil), f.data...)
}

func (f *Frame) DataAsBuffer() *Buffer {
	return NewBuffer(f.Data())
}

func (f *Frame) SetData(value []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.length = uint16(len(value))
	f.data = value
}

func (f *Frame) SetDataAsBuffer(buf *Buffer) {
	f.SetData(buf.Bytes())
}

func (f *Frame) FCS() byte {
	f.lock.RLock()
	defer f.lock.RUnlock()

	return f.fcs
}

func (f Frame) MarshalBinary() ([]byte, error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	buffer := make([]byte, 0, FrameLengthMin+f.length)
	buffer = append(buffer,
		byte(f.length),
		byte(((f.typ<<5)&0xE0)|(f.subSystem&0x1F)),
		byte(f.commandID))
	buffer = append(buffer, f.data...)
	buffer = append(buffer, checksum(buffer))
	buffer = append([]byte{SOF}, buffer...)

	return buffer, nil
}

func (f *Frame) UnmarshalBinary(data []byte) error {
	// frame size: SOF + MT CMD + FCS

	// min length: 1 + 3 + 1 = 5
	if len(data) < FrameLengthMin {
		return errors.New("frame length less than 5")
	}

	// max length: 1 + 256 + 1 = 258
	if len(data) > FrameLengthMax {
		return errors.New("frame length greater than 258")
	}

	if data[0] != SOF {
		return errors.New("first byte of frame isn't SOF")
	}

	f.SetCommand0(uint16(data[PositionCommand1]))

	f.lock.Lock()
	defer f.lock.Unlock()

	switch f.typ {
	case TypePoll, TypeSREQ, TypeAREQ, TypeSRSP:
		// skip
	default:
		return fmt.Errorf("unknown type of frame command 0x%X in data 0x%X", f.typ, data)
	}

	switch f.subSystem {
	case SubSystemReserved, SubSystemSysInterface, SubSystemMACInterface,
		SubSystemNWKInterface, SubSystemAFInterface, SubSystemZDOInterface,
		SubSystemSAPIInterface, SubSystemUtilInterface, SubSystemDebugInterface,
		SubSystemAppInterface:
		// skip
	default:
		return fmt.Errorf("unknown sub system of frame command 0x%X in data 0x%X", f.subSystem, data)
	}

	// MT CMD = LEN (1) + CMD (2) + DATA (0-250)
	f.length = uint16(data[PositionFrameLength])
	f.commandID = uint16(data[PositionCommand2])
	f.data = data[PositionData : PositionData+f.length]
	f.fcs = data[f.length+FrameLengthMin-1]

	// checksum validate
	if sum := checksum(data[1 : f.length+FrameLengthMin-1]); sum != f.fcs {
		return fmt.Errorf("FCS isn't valid have 0x%X want 0x%X in data 0x%X", f.fcs, sum, data[:f.length+FrameLengthMin])
	}

	return nil
}

func (f *Frame) String() string {
	buffer := bytes.NewBuffer(nil)

	buffer.WriteString("len: ")
	buffer.WriteString(strconv.FormatUint(uint64(f.Length()), 10))
	buffer.WriteString(" type: ")
	buffer.WriteString(strconv.FormatUint(uint64(f.Type()), 10))
	buffer.WriteString(" sub system: ")
	buffer.WriteString(strconv.FormatUint(uint64(f.SubSystem()), 10))
	buffer.WriteString(" (0x")
	buffer.WriteString(fmt.Sprintf("%X", f.SubSystem()))
	buffer.WriteString(") command id: ")
	buffer.WriteString(strconv.FormatUint(uint64(f.CommandID()), 10))
	buffer.WriteString(" (0x")
	buffer.WriteString(fmt.Sprintf("%X", f.CommandID()))
	buffer.WriteString(") data: ")
	buffer.WriteString(fmt.Sprint(f.Data()))
	buffer.WriteString(" fcs: ")
	buffer.WriteString(strconv.FormatUint(uint64(f.FCS()), 10))

	return buffer.String()
}

/*
     SOF  LENGTH       CMD   DATA                                                                  FCS
<-- [254, 27,    68,   129,  0,0,6,0,233,142,1,1,0,0,0,33,55,9,0,0,7,24,33,10,0,0,16,0,233,142,29, 254] +3s
--- parseNext [254,27,68,129,0,0,6,0,233,142,1,1,0,0,0,33,55,9,0,0,7,24,33,10,0,0,16,0,233,142,29,254] +0ms
--> parsed 27 - 2 - 4 - 129 - [0,0,6,0,233,142,1,1,0,0,0,33,55,9,0,0,7,24,33,10,0,0,16,0,233,142,29] - 254 +1ms
*/

func checksum(buffer []byte) (checksum byte) {
	for _, value := range buffer {
		checksum ^= value
	}

	return checksum
}
