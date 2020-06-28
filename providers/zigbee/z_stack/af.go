package z_stack

import (
	"context"
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
	Len            uint8
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

func (c *Client) AfIncomingMessage(frame *Frame) (*AfIncomingMessage, error) {
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
		Len:            dataOut.ReadUint8(),
	}
	msg.Data = dataOut.Next(int(msg.Len))

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

		//case 0x1: // specific

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

/*
	AF_REGISTER

	This command enables the tester to register an application’s endpoint description.

	Usage:
		SREQ:
			          1        |      1      |       1     |     1    |      2    |      2      |     1     |      1     |        1         |       0-32       |         1         |       0-32
			Length = 0x09-0x49 | Cmd0 = 0x24 | Cmd1 = 0x00 | EndPoint | AppProfId | AppDeviceId | AppDevVer | LatencyReq | AppNumInClusters | AppInClusterList | AppNumOutClusters | AppOutClusterList
		Attributes:
			EndPoint          1 byte   Specifies the endpoint of the device
			AppProfId         2 bytes  Specifies the profile Id of the application
			AppDeviceId       2 bytes  Specifies the device description Id for this endpoint
			AddDevVer         1 byte   Specifies the device version number
			LatencyReq        1 byte   Specifies latency.
			                           0x00-No latency
			                           0x01-fast beacons
			                           0x02-slow beacons
			AppNumInClusters  1 byte   The number of Input cluster Id’s following in the AppInClusterList
			AppInClusterList  32 bytes Specifies the list of Input Cluster Id’s
			AppNumOutClusters 1 byte   Specifies the number of Output cluster Id’s following in the AppOutClusterList
			AppOutClusterList 32 bytes Specifies the list of Output Cluster Id’s

		SRSP:
			       1      |       1     |       1     |    1
			Length = 0x01 | Cmd0 = 0x64 | Cmd1 = 0x00 | Status
		Attributes:
			Status 1 byte Status is either Success (0) or Failure (1).

	Example from zigbee2mqtt:
		zigbee-herdsman:adapter:zStack:znp:SREQ --> AF - register - {"appdeviceid":5,"appdevver":0,"appnuminclusters":0,"appinclusterlist":[],"appnumoutclusters":0,"appoutclusterlist":[],"latencyreq":0,"endpoint":1,"appprofid":260} +14ms
		zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,9,36,0,1,4,1,5,0,0,0,0,0,44] +14ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,100,0,0,101] +6ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,100,0,0,101] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 3 - 4 - 0 - [0] - 101 +0ms
		zigbee-herdsman:adapter:zStack:znp:SRSP <-- AF - register - {"status":0} +12ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [] +1ms
*/
func (c *Client) AfRegister(ctx context.Context, endpoint Endpoint) error {
	response, err := c.CallWithResultSREQ(ctx, endpoint.AsBuffer().Frame(0x24, 0x00))
	if err != nil {
		return err
	}

	if response.Command0() != 0x64 {
		return errors.New("bad response")
	}

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		if dataOut[0] == 0xB8 {
			return errors.New("already registered")
		}

		return errors.New("failure")
	}

	return nil
}
