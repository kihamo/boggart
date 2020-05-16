package z_stack

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const (
	// A POLL command is used to retrieve queued data. This command is only applicable to SPI transport.
	// For a POLL command the subsystem and Id are set to zero and data length is zero.
	TypePoll = 0

	// A synchronous request that requires an immediate response.
	// For example, a function call with a return value would use an SREQ command.
	TypeSREQ = 1

	// An asynchronous request.
	// For example, a callback event or a function call with no return value would use an AREQ command.
	TypeAREQ = 2

	// A synchronous response. This type of command is only sent in response to a SREQ command.
	// For an SRSP command the subsystem and Id are set to the same values as the corresponding SREQ.
	// The length of an SRSP is generally nonzero, so an SRSP with length=0 can be used to indicate an error.
	TypeSRSP = 3

	SubSystemReserved       = 0x00
	SubSystemSysInterface   = 0x01
	SubSystemMACInterface   = 0x02
	SubSystemNWKInterface   = 0x03
	SubSystemAFInterface    = 0x04
	SubSystemZDOInterface   = 0x05
	SubSystemSAPIInterface  = 0x06
	SubSystemUtilInterface  = 0x07
	SubSystemDebugInterface = 0x08
	SubSystemAppInterface   = 0x09

	SOF            = byte(0xFE)
	FrameLengthMin = 5
	FrameLengthMax = 258

	PositionFrameLength = 1
	PositionCommand1    = 2
	PositionCommand2    = 3
	PositionData        = 4
)

type Frame struct {
	length    uint16
	typ       uint16
	subSystem uint16
	commandID uint16
	data      []byte
	fcs       byte
}

func (f *Frame) SetCommand0(value uint16) {
	f.typ = (value & 0xE0) >> 5
	f.subSystem = value & 0x1F
}

func (f *Frame) Type() uint16 {
	return f.typ
}

func (f *Frame) SetType(value uint16) {
	f.typ = value
}

func (f *Frame) SubSystem() uint16 {
	return f.subSystem
}

func (f *Frame) SetSubSystem(value uint16) {
	f.subSystem = value
}

func (f *Frame) CommandID() uint16 {
	return f.commandID
}

func (f *Frame) SetCommandID(value uint16) {
	f.commandID = value
}

func (f *Frame) Data() []byte {
	return f.data
}

func (f *Frame) SetData(value []byte) {
	f.length = uint16(len(value))
	f.data = value
}

func (f Frame) MarshalBinary() ([]byte, error) {
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

	// MT CMD = LEN (1) + CMD (2) + DATA (0-250)
	f.length = uint16(data[PositionFrameLength])

	f.SetCommand0(uint16(data[PositionCommand1]))

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

	buffer.WriteString(strconv.FormatUint(uint64(f.length), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.typ), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.subSystem), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.commandID), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(fmt.Sprint(f.data))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.fcs), 10))

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
