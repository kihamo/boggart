// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Status status
//
// swagger:model Status
type Status struct {

	// code
	// Maximum: 7
	Code uint64 `json:"code,omitempty" xml:"statusCode,omitempty"`

	// string
	String string `json:"string,omitempty" xml:"statusString,omitempty"`

	// sub code
	// Enum: [ok riskPassword noMemory serviceUnavailable upgrading deviceBusy reConnectIpc deviceError badFlash 28181Uninitialized notSupport lowPrivilege badAuthorization methodNotAllowed notSetHdiskRedund invalidOperation notActivated hasActivated badXmlFormat badParameters badHostAddress badXmlContent badIPv4Address badIPv6Address conflictIPv4Address conflictIPv6Address badDomainName connectSreverFail conflictDomainName badPort portError importErrorData badNetMask badVersion badDevType badLanguage incorrentUserNameOrPassword invalidStoragePoolOfCloudServer noFreeSpaceOfStoragePool fileFormatError fileContentError UnSupportCapture unableCalibrate pleaseCalibrate SNMPv3PasswordNone SNMPv3NameDifferent notSupportDeicing notMeetDeicing alarmInputOccupied notSupportWithAPMode rebootRequired]
	SubCode string `json:"subCode,omitempty" xml:"subStatusCode,omitempty"`
}

// Validate validates this status
func (m *Status) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubCode(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Status) validateCode(formats strfmt.Registry) error {

	if swag.IsZero(m.Code) { // not required
		return nil
	}

	if err := validate.MaximumInt("code", "body", int64(m.Code), 7, false); err != nil {
		return err
	}

	return nil
}

var statusTypeSubCodePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ok","riskPassword","noMemory","serviceUnavailable","upgrading","deviceBusy","reConnectIpc","deviceError","badFlash","28181Uninitialized","notSupport","lowPrivilege","badAuthorization","methodNotAllowed","notSetHdiskRedund","invalidOperation","notActivated","hasActivated","badXmlFormat","badParameters","badHostAddress","badXmlContent","badIPv4Address","badIPv6Address","conflictIPv4Address","conflictIPv6Address","badDomainName","connectSreverFail","conflictDomainName","badPort","portError","importErrorData","badNetMask","badVersion","badDevType","badLanguage","incorrentUserNameOrPassword","invalidStoragePoolOfCloudServer","noFreeSpaceOfStoragePool","fileFormatError","fileContentError","UnSupportCapture","unableCalibrate","pleaseCalibrate","SNMPv3PasswordNone","SNMPv3NameDifferent","notSupportDeicing","notMeetDeicing","alarmInputOccupied","notSupportWithAPMode","rebootRequired"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		statusTypeSubCodePropEnum = append(statusTypeSubCodePropEnum, v)
	}
}

const (

	// StatusSubCodeOk captures enum value "ok"
	StatusSubCodeOk string = "ok"

	// StatusSubCodeRiskPassword captures enum value "riskPassword"
	StatusSubCodeRiskPassword string = "riskPassword"

	// StatusSubCodeNoMemory captures enum value "noMemory"
	StatusSubCodeNoMemory string = "noMemory"

	// StatusSubCodeServiceUnavailable captures enum value "serviceUnavailable"
	StatusSubCodeServiceUnavailable string = "serviceUnavailable"

	// StatusSubCodeUpgrading captures enum value "upgrading"
	StatusSubCodeUpgrading string = "upgrading"

	// StatusSubCodeDeviceBusy captures enum value "deviceBusy"
	StatusSubCodeDeviceBusy string = "deviceBusy"

	// StatusSubCodeReConnectIpc captures enum value "reConnectIpc"
	StatusSubCodeReConnectIpc string = "reConnectIpc"

	// StatusSubCodeDeviceError captures enum value "deviceError"
	StatusSubCodeDeviceError string = "deviceError"

	// StatusSubCodeBadFlash captures enum value "badFlash"
	StatusSubCodeBadFlash string = "badFlash"

	// StatusSubCodeNr28181Uninitialized captures enum value "28181Uninitialized"
	StatusSubCodeNr28181Uninitialized string = "28181Uninitialized"

	// StatusSubCodeNotSupport captures enum value "notSupport"
	StatusSubCodeNotSupport string = "notSupport"

	// StatusSubCodeLowPrivilege captures enum value "lowPrivilege"
	StatusSubCodeLowPrivilege string = "lowPrivilege"

	// StatusSubCodeBadAuthorization captures enum value "badAuthorization"
	StatusSubCodeBadAuthorization string = "badAuthorization"

	// StatusSubCodeMethodNotAllowed captures enum value "methodNotAllowed"
	StatusSubCodeMethodNotAllowed string = "methodNotAllowed"

	// StatusSubCodeNotSetHdiskRedund captures enum value "notSetHdiskRedund"
	StatusSubCodeNotSetHdiskRedund string = "notSetHdiskRedund"

	// StatusSubCodeInvalidOperation captures enum value "invalidOperation"
	StatusSubCodeInvalidOperation string = "invalidOperation"

	// StatusSubCodeNotActivated captures enum value "notActivated"
	StatusSubCodeNotActivated string = "notActivated"

	// StatusSubCodeHasActivated captures enum value "hasActivated"
	StatusSubCodeHasActivated string = "hasActivated"

	// StatusSubCodeBadXMLFormat captures enum value "badXmlFormat"
	StatusSubCodeBadXMLFormat string = "badXmlFormat"

	// StatusSubCodeBadParameters captures enum value "badParameters"
	StatusSubCodeBadParameters string = "badParameters"

	// StatusSubCodeBadHostAddress captures enum value "badHostAddress"
	StatusSubCodeBadHostAddress string = "badHostAddress"

	// StatusSubCodeBadXMLContent captures enum value "badXmlContent"
	StatusSubCodeBadXMLContent string = "badXmlContent"

	// StatusSubCodeBadIPV4Address captures enum value "badIPv4Address"
	StatusSubCodeBadIPV4Address string = "badIPv4Address"

	// StatusSubCodeBadIPV6Address captures enum value "badIPv6Address"
	StatusSubCodeBadIPV6Address string = "badIPv6Address"

	// StatusSubCodeConflictIPV4Address captures enum value "conflictIPv4Address"
	StatusSubCodeConflictIPV4Address string = "conflictIPv4Address"

	// StatusSubCodeConflictIPV6Address captures enum value "conflictIPv6Address"
	StatusSubCodeConflictIPV6Address string = "conflictIPv6Address"

	// StatusSubCodeBadDomainName captures enum value "badDomainName"
	StatusSubCodeBadDomainName string = "badDomainName"

	// StatusSubCodeConnectSreverFail captures enum value "connectSreverFail"
	StatusSubCodeConnectSreverFail string = "connectSreverFail"

	// StatusSubCodeConflictDomainName captures enum value "conflictDomainName"
	StatusSubCodeConflictDomainName string = "conflictDomainName"

	// StatusSubCodeBadPort captures enum value "badPort"
	StatusSubCodeBadPort string = "badPort"

	// StatusSubCodePortError captures enum value "portError"
	StatusSubCodePortError string = "portError"

	// StatusSubCodeImportErrorData captures enum value "importErrorData"
	StatusSubCodeImportErrorData string = "importErrorData"

	// StatusSubCodeBadNetMask captures enum value "badNetMask"
	StatusSubCodeBadNetMask string = "badNetMask"

	// StatusSubCodeBadVersion captures enum value "badVersion"
	StatusSubCodeBadVersion string = "badVersion"

	// StatusSubCodeBadDevType captures enum value "badDevType"
	StatusSubCodeBadDevType string = "badDevType"

	// StatusSubCodeBadLanguage captures enum value "badLanguage"
	StatusSubCodeBadLanguage string = "badLanguage"

	// StatusSubCodeIncorrentUserNameOrPassword captures enum value "incorrentUserNameOrPassword"
	StatusSubCodeIncorrentUserNameOrPassword string = "incorrentUserNameOrPassword"

	// StatusSubCodeInvalidStoragePoolOfCloudServer captures enum value "invalidStoragePoolOfCloudServer"
	StatusSubCodeInvalidStoragePoolOfCloudServer string = "invalidStoragePoolOfCloudServer"

	// StatusSubCodeNoFreeSpaceOfStoragePool captures enum value "noFreeSpaceOfStoragePool"
	StatusSubCodeNoFreeSpaceOfStoragePool string = "noFreeSpaceOfStoragePool"

	// StatusSubCodeFileFormatError captures enum value "fileFormatError"
	StatusSubCodeFileFormatError string = "fileFormatError"

	// StatusSubCodeFileContentError captures enum value "fileContentError"
	StatusSubCodeFileContentError string = "fileContentError"

	// StatusSubCodeUnSupportCapture captures enum value "UnSupportCapture"
	StatusSubCodeUnSupportCapture string = "UnSupportCapture"

	// StatusSubCodeUnableCalibrate captures enum value "unableCalibrate"
	StatusSubCodeUnableCalibrate string = "unableCalibrate"

	// StatusSubCodePleaseCalibrate captures enum value "pleaseCalibrate"
	StatusSubCodePleaseCalibrate string = "pleaseCalibrate"

	// StatusSubCodeSNMPv3PasswordNone captures enum value "SNMPv3PasswordNone"
	StatusSubCodeSNMPv3PasswordNone string = "SNMPv3PasswordNone"

	// StatusSubCodeSNMPv3NameDifferent captures enum value "SNMPv3NameDifferent"
	StatusSubCodeSNMPv3NameDifferent string = "SNMPv3NameDifferent"

	// StatusSubCodeNotSupportDeicing captures enum value "notSupportDeicing"
	StatusSubCodeNotSupportDeicing string = "notSupportDeicing"

	// StatusSubCodeNotMeetDeicing captures enum value "notMeetDeicing"
	StatusSubCodeNotMeetDeicing string = "notMeetDeicing"

	// StatusSubCodeAlarmInputOccupied captures enum value "alarmInputOccupied"
	StatusSubCodeAlarmInputOccupied string = "alarmInputOccupied"

	// StatusSubCodeNotSupportWithAPMode captures enum value "notSupportWithAPMode"
	StatusSubCodeNotSupportWithAPMode string = "notSupportWithAPMode"

	// StatusSubCodeRebootRequired captures enum value "rebootRequired"
	StatusSubCodeRebootRequired string = "rebootRequired"
)

// prop value enum
func (m *Status) validateSubCodeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, statusTypeSubCodePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Status) validateSubCode(formats strfmt.Registry) error {

	if swag.IsZero(m.SubCode) { // not required
		return nil
	}

	// value enum
	if err := m.validateSubCodeEnum("subCode", "body", m.SubCode); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Status) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Status) UnmarshalBinary(b []byte) error {
	var res Status
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
