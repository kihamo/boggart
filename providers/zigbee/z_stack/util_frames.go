package zstack

import (
	"errors"
)

func UtilGetDeviceInfoParse(frame *Frame) (*UtilDeviceInfo, error) {
	if frame.SubSystem() != SubSystemUtilInterface {
		return nil, errors.New("frame isn't a UTIL interface")
	}

	if frame.CommandID() != CommandGetDeviceInfo {
		return nil, errors.New("frame isn't a device joined command")
	}

	dataOut := frame.DataAsBuffer()
	if dataOut.Len() == 0 {
		return nil, errors.New("failure")
	}

	devicesCount := (dataOut.Len() - 14) / 2

	msg := &UtilDeviceInfo{
		Status: dataOut.ReadCommandStatus(),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	msg.IEEEAddr = dataOut.ReadIEEEAddr()
	msg.ShortAddr = dataOut.ReadUint16()
	msg.DeviceType = dataOut.ReadUint8()
	msg.DeviceState = DeviceState(dataOut.ReadUint8())
	msg.NumAssocDevices = dataOut.ReadUint8()
	msg.AssocDevicesList = make([]uint16, 0, devicesCount)

	for i := 0; i < devicesCount; i++ {
		msg.AssocDevicesList = append(msg.AssocDevicesList, dataOut.ReadUint16())
	}

	return msg, nil
}
