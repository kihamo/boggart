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
	TypePoll = 0x0

	// A synchronous request that requires an immediate response.
	// For example, a function call with a return value would use an SREQ command.
	TypeSREQ = 0x2

	// An asynchronous request.
	// For example, a callback event or a function call with no return value would use an AREQ command.
	TypeAREQ = 0x4

	// A synchronous response. This type of command is only sent in response to a SREQ command.
	// For an SRSP command the subsystem and Id are set to the same values as the corresponding SREQ.
	// The length of an SRSP is generally nonzero, so an SRSP with length=0 can be used to indicate an error.
	TypeSRSP = 0x6

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
	Length    uint16
	Type      uint16
	SubSystem uint16
	CommandID uint16
	Data      []byte
	FCS       byte
}

func (f Frame) MarshalBinary() ([]byte, error) {
	return nil, nil
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
	f.Length = uint16(data[PositionFrameLength])

	cmd := uint16(data[PositionCommand1])
	f.Type = (cmd & 0xE0) >> 5
	f.SubSystem = cmd & 0x1F

	switch f.Type {
	case TypePoll, TypeSREQ, TypeAREQ, TypeSRSP:
		// skip
	default:
		return fmt.Errorf("unknown type of frame command 0x%X in data 0x%X", f.Type, data)
	}

	switch f.SubSystem {
	case SubSystemReserved, SubSystemSysInterface, SubSystemMACInterface,
		SubSystemNWKInterface, SubSystemAFInterface, SubSystemZDOInterface,
		SubSystemSAPIInterface, SubSystemUtilInterface, SubSystemDebugInterface,
		SubSystemAppInterface:
		// skip
	default:
		return fmt.Errorf("unknown sub system of frame command 0x%X in data 0x%X", f.SubSystem, data)
	}

	f.CommandID = uint16(data[PositionCommand2])
	f.Data = data[PositionData : f.Length-1]
	f.FCS = data[f.Length+FrameLengthMin-1]

	// checksum validate
	if sum := checksum(data[1 : f.Length+FrameLengthMin-1]); sum != f.FCS {
		return fmt.Errorf("FCS isn't valid have 0x%X want 0x%X in data 0x%X", f.FCS, sum, data[:f.Length+FrameLengthMin])
	}

	return nil
}

func (f *Frame) String() string {
	buffer := bytes.NewBuffer(nil)

	buffer.WriteString(strconv.FormatUint(uint64(f.Length), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.Type), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.SubSystem), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.CommandID), 10))
	buffer.WriteString(" - ")
	buffer.WriteString(fmt.Sprint(f.Data))
	buffer.WriteString(" - ")
	buffer.WriteString(strconv.FormatUint(uint64(f.FCS), 10))

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
