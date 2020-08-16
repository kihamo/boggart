package z_stack

import (
	"errors"
	"fmt"
)

type FramePayloadReportRecorder struct {
	AttributeID   uint16
	DataType      uint8
	AttributeData interface{}
}

type AfIncomingMessage struct {
	GroupID        uint16
	ClusterID      uint16
	SrcAddr        uint16
	SrcEndpoint    uint8
	DstEndpoint    uint8
	WasBroadcast   uint8
	LinkQuality    uint8
	SecurityUse    uint8
	TimeStamp      uint32
	TransSeqNumber uint8
	Data           []byte

	Frame struct {
		Header struct {
			FrameType                 uint8
			ManufacturerSpecific      bool
			Direction                 uint8 // CLIENT_TO_SERVER = 0, SERVER_TO_CLIENT = 1
			DisableDefaultResponse    bool
			ManufacturerCode          uint16
			TransactionSequenceNumber uint8
			CommandIdentifier         uint8
		}
		Payload struct {
			Report *[]FramePayloadReportRecorder
		}
	}
}

type Endpoint struct {
	EndPoint          uint8
	AppProfId         uint16
	AppDeviceId       uint16
	AddDevVer         uint8
	LatencyReq        uint8
	AppInClusterList  []uint16
	AppOutClusterList []uint16
}

func (e Endpoint) AsBuffer() *Buffer {
	dataOut := NewBuffer(nil)
	dataOut.WriteUint8(e.EndPoint)
	dataOut.WriteUint16(e.AppProfId)
	dataOut.WriteUint16(e.AppDeviceId)
	dataOut.WriteUint8(e.AddDevVer)
	dataOut.WriteUint8(e.LatencyReq)
	dataOut.WriteUint8(uint8(len(e.AppInClusterList)))

	for _, id := range e.AppInClusterList {
		dataOut.WriteUint16(id)
	}

	dataOut.WriteUint8(uint8(len(e.AppOutClusterList)))

	for _, id := range e.AppOutClusterList {
		dataOut.WriteUint16(id)
	}

	return dataOut
}

func AfIncomingMessageParse(frame *Frame) (*AfIncomingMessage, error) {
	if frame.SubSystem() != SubSystemAFInterface {
		return nil, errors.New("frame isn't a AF interface")
	}

	if frame.CommandID() != CommandAfIncomingMessage {
		return nil, errors.New("frame isn't a af_incoming_msg")
	}

	dataOut := frame.DataAsBuffer()

	msg := &AfIncomingMessage{
		GroupID:        dataOut.ReadUint16(),
		ClusterID:      dataOut.ReadUint16(),
		SrcAddr:        dataOut.ReadUint16(),
		SrcEndpoint:    dataOut.ReadUint8(),
		DstEndpoint:    dataOut.ReadUint8(),
		WasBroadcast:   dataOut.ReadUint8(),
		LinkQuality:    dataOut.ReadUint8(),
		SecurityUse:    dataOut.ReadUint8(),
		TimeStamp:      dataOut.ReadUint32(),
		TransSeqNumber: dataOut.ReadUint8(),
	}
	l := dataOut.ReadUint8()
	msg.Data = dataOut.Next(int(l))

	payload := NewBuffer(msg.Data)

	// parse header
	flowControlValue := payload.ReadUint8()

	msg.Frame.Header.FrameType = flowControlValue & 0x03
	msg.Frame.Header.ManufacturerSpecific = ((flowControlValue >> 2) & 0x01) == 1
	msg.Frame.Header.Direction = (flowControlValue >> 3) & 0x01
	msg.Frame.Header.DisableDefaultResponse = ((flowControlValue >> 4) & 0x01) == 1

	if msg.Frame.Header.ManufacturerSpecific {
		msg.Frame.Header.ManufacturerCode = payload.ReadUint16()
	}

	msg.Frame.Header.TransactionSequenceNumber = payload.ReadUint8()
	msg.Frame.Header.CommandIdentifier = payload.ReadUint8()

	//fmt.Println("Header:")
	//fmt.Println("  frameControlValue", flowControlValue)
	//fmt.Println("  frameType", msg.Frame.Header.FrameType)
	//fmt.Println("  manufacturerSpecific", msg.Frame.Header.ManufacturerSpecific)
	//fmt.Println("  manufacturerCode", msg.Frame.Header.ManufacturerCode)
	//fmt.Println("  direction", msg.Frame.Header.Direction)
	//fmt.Println("  disableDefaultResponse", msg.Frame.Header.DisableDefaultResponse)
	//fmt.Println("  transactionSequenceNumber", msg.Frame.Header.TransactionSequenceNumber)
	//fmt.Println("  commandIdentifier", msg.Frame.Header.CommandIdentifier)
	//fmt.Println("-----")

	// parse payload
	switch msg.Frame.Header.FrameType {
	case 0x0: // global
		switch msg.Frame.Header.CommandIdentifier {
		// only xiaomi https://github.com/Koenkk/zigbee-herdsman/blob/5151e5b0922a98abf64adf644def3adfb6970c93/src/zcl/zclFrame.ts#L239
		case 10:
			reports := make([]FramePayloadReportRecorder, 0)

			for payload.Len() > 0 {
				report := FramePayloadReportRecorder{
					AttributeID: payload.ReadUint16(),
					DataType:    payload.ReadUint8(),
				}
				report.AttributeData = payload.ReadByType(report.DataType)

				reports = append(reports, report)
			}

			msg.Frame.Payload.Report = &reports

		default:
			return nil, fmt.Errorf("unsupported command identifier %d", msg.Frame.Header.CommandIdentifier)
		}

	case 0x1: // specific
		fmt.Println("Cluster", msg.ClusterID)
		fmt.Println("Qua", msg.LinkQuality)

		if msg.Frame.Header.Direction == 0 { // client to server

		} else { // server to client

		}
		fmt.Println("Command", msg.Frame.Header.CommandIdentifier)
		fmt.Println("fieldControl", payload.ReadUint8())
		fmt.Println("manufacturerCode", payload.ReadUint16())
		fmt.Println("imageType", payload.ReadUint16())
		fmt.Println("fileVersion", payload.ReadUint32())

	default:
		return nil, fmt.Errorf("unsupported frameType %d", msg.Frame.Header.FrameType)
	}

	//fmt.Println("Payload:")
	//fmt.Println("  Report", msg.Frame.Payload)
	//fmt.Println("  AttributeID", (*msg.Frame.Payload.Report)[0].AttributeID)
	//fmt.Println("  DataType", (*msg.Frame.Payload.Report)[0].DataType)
	//fmt.Println("  AttributeData", (*msg.Frame.Payload.Report)[0].AttributeData)
	//fmt.Println("-----")

	return msg, nil
}
