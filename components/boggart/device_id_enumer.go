// Code generated by "enumer -type=DeviceId -trimprefix=DeviceId -output=device_id_enumer.go -transform=snake"; DO NOT EDIT

package boggart

import (
	"fmt"
)

const _DeviceIdName = "electricity_metercamera_hallcamera_streetheat_meterentrance_doorphoneroutervideo_recorderwater_meter_coldwater_meter_hot"

var _DeviceIdIndex = [...]uint8{0, 17, 28, 41, 51, 64, 69, 75, 89, 105, 120}

func (i DeviceId) String() string {
	if i < 0 || i >= DeviceId(len(_DeviceIdIndex)-1) {
		return fmt.Sprintf("DeviceId(%d)", i)
	}
	return _DeviceIdName[_DeviceIdIndex[i]:_DeviceIdIndex[i+1]]
}

var _DeviceIdNameToValueMap = map[string]DeviceId{
	_DeviceIdName[0:17]:    0,
	_DeviceIdName[17:28]:   1,
	_DeviceIdName[28:41]:   2,
	_DeviceIdName[41:51]:   3,
	_DeviceIdName[51:64]:   4,
	_DeviceIdName[64:69]:   5,
	_DeviceIdName[69:75]:   6,
	_DeviceIdName[75:89]:   7,
	_DeviceIdName[89:105]:  8,
	_DeviceIdName[105:120]: 9,
}

// DeviceIdString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DeviceIdString(s string) (DeviceId, error) {
	if val, ok := _DeviceIdNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DeviceId values", s)
}
