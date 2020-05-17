package z_stack

import (
	"errors"
	"fmt"
)

const (
	CommandAfIncomingMessage = 0x81
)

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
			AttributeID      uint16
			DataType         uint8
			NumberOfElements uint16
			AttributeData    []struct {
				Type  uint8
				Value interface{}
			}
		}
	}
}

func (c *Client) AfIncomingMessage(frame *Frame) (*AfIncomingMessage, error) {
	if frame.CommandID() != CommandAfIncomingMessage {
		return nil, errors.New("frame isn't a af_incoming_msg")
	}

	data := frame.DataAsBuffer()

	msg := &AfIncomingMessage{
		GroupID:        data.ReadUint16(),
		ClusterID:      data.ReadUint16(),
		SrcAddr:        data.ReadUint16(),
		SrcEndpoint:    data.ReadUint8(),
		DstEndpoint:    data.ReadUint8(),
		WasBroadcast:   data.ReadUint8(),
		LinkQuality:    data.ReadUint8(),
		SecurityUse:    data.ReadUint8(),
		TimeStamp:      data.ReadUint32(),
		TransSeqNumber: data.ReadUint8(),
		Len:            data.ReadUint8(),
	}
	msg.Data = data.Next(int(msg.Len))

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

	// parse payload
	switch msg.Frame.Header.FrameType {
	case 0x0: // global
		// onli xiaomi https://github.com/Koenkk/zigbee-herdsman/blob/5151e5b0922a98abf64adf644def3adfb6970c93/src/zcl/zclFrame.ts#L239
		if msg.Frame.Header.CommandIdentifier == 0xA {
			// https://github.com/Koenkk/zigbee-herdsman/blob/5151e5b0922a98abf64adf644def3adfb6970c93/src/zcl/definition/foundation.ts#L120
			msg.Frame.Payload.AttributeID = payload.ReadUint16()
			msg.Frame.Payload.DataType = payload.ReadUint8()

			// struct
			if msg.Frame.Payload.DataType == 0x4C {
				msg.Frame.Payload.NumberOfElements = payload.ReadUint16()

				for i := uint16(0); i < msg.Frame.Payload.NumberOfElements; i++ {
					msg.Frame.Payload.AttributeData = append(msg.Frame.Payload.AttributeData,
						struct {
							Type  uint8
							Value interface{}
						}{
							Type:  payload.ReadUint8(),
							Value: payload.ReadUint16(),
						})
				}
			} else {
				return nil, fmt.Errorf("unsupported data type %d", msg.Frame.Payload.DataType)
			}
			/*
							    report: {
				        ID: 10,
				        parseStrategy: 'repetitive',
				        parameters: [
				            {name: 'attrId', type: DataType.uint16},
				            {name: 'dataType', type: DataType.uint8},
				            {name: 'attrData', type: BuffaloZclDataType.USE_DATA_TYPE},
				        ],
				    },
			*/
		} else {
			return nil, fmt.Errorf("unsupported command identifier %d", msg.Frame.Header.CommandIdentifier)
		}
	case 0x1: // specific
		// TODO:
	default:
		return nil, fmt.Errorf("unsupported frameType %d", msg.Frame.Header.FrameType)
	}

	fmt.Println("frameControlValue", 24)
	fmt.Println("frameType", msg.Frame.Header.FrameType)
	fmt.Println("manufacturerSpecific", msg.Frame.Header.ManufacturerSpecific)
	fmt.Println("direction", msg.Frame.Header.Direction)
	fmt.Println("disableDefaultResponse", msg.Frame.Header.DisableDefaultResponse)
	fmt.Println("transactionSequenceNumber", msg.Frame.Header.TransactionSequenceNumber)
	fmt.Println("commandIdentifier", msg.Frame.Header.CommandIdentifier)
	fmt.Println("AttributeID", msg.Frame.Payload.AttributeID)
	fmt.Println("DataType", msg.Frame.Payload.DataType)
	fmt.Println("NumberOfElements", msg.Frame.Payload.NumberOfElements)
	fmt.Println("AttributeData", msg.Frame.Payload.AttributeData)

	return msg, nil
}
