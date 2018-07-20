// Code generated by "enumer -type=DeviceType -trimprefix=DeviceType -output=device_type_enumer.go"; DO NOT EDIT.

package boggart

import (
	"fmt"
)

const _DeviceTypeName = "CameraDoorElectricityMeterHeatMeterInternetProviderPhoneRouterTVPCVideoRecorderWaterMeterThermometerBarometerHygrometer"

var _DeviceTypeIndex = [...]uint8{0, 6, 10, 26, 35, 51, 56, 62, 64, 66, 79, 89, 100, 109, 119}

func (i DeviceType) String() string {
	if i < 0 || i >= DeviceType(len(_DeviceTypeIndex)-1) {
		return fmt.Sprintf("DeviceType(%d)", i)
	}
	return _DeviceTypeName[_DeviceTypeIndex[i]:_DeviceTypeIndex[i+1]]
}

var _DeviceTypeValues = []DeviceType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

var _DeviceTypeNameToValueMap = map[string]DeviceType{
	_DeviceTypeName[0:6]:     0,
	_DeviceTypeName[6:10]:    1,
	_DeviceTypeName[10:26]:   2,
	_DeviceTypeName[26:35]:   3,
	_DeviceTypeName[35:51]:   4,
	_DeviceTypeName[51:56]:   5,
	_DeviceTypeName[56:62]:   6,
	_DeviceTypeName[62:64]:   7,
	_DeviceTypeName[64:66]:   8,
	_DeviceTypeName[66:79]:   9,
	_DeviceTypeName[79:89]:   10,
	_DeviceTypeName[89:100]:  11,
	_DeviceTypeName[100:109]: 12,
	_DeviceTypeName[109:119]: 13,
}

// DeviceTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DeviceTypeString(s string) (DeviceType, error) {
	if val, ok := _DeviceTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DeviceType values", s)
}

// DeviceTypeValues returns all values of the enum
func DeviceTypeValues() []DeviceType {
	return _DeviceTypeValues
}

// IsADeviceType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DeviceType) IsADeviceType() bool {
	for _, v := range _DeviceTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
