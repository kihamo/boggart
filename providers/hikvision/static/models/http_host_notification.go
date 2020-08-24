package models

import (
	"encoding/json"
	"encoding/xml"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

type HTTPHostNotification struct {
	// переопределение ради этого поля, потому что стандартный механизм генерирует
	// аббревиатуру HTTP исключительно в верхнем регистре
	XMLName xml.Name `xml:"HttpHostNotification"`

	AddressingFormatType     string  `json:"addressingFormatType,omitempty" xml:"addressingFormatType,omitempty"`
	HostName                 *string `json:"hostName,omitempty" xml:"hostName,omitempty"`
	HTTPAuthenticationMethod string  `json:"httpAuthenticationMethod,omitempty" xml:"httpAuthenticationMethod,omitempty"`
	ID                       uint64  `json:"id,omitempty" xml:"id,omitempty"`
	IntervalBetweenEvents    *int64  `json:"intervalBetweenEvents,omitempty" xml:"Extensions>intervalBetweenEvents,omitempty"`
	IPAddress                *string `json:"ipAddress,omitempty" xml:"ipAddress,omitempty"`
	IPV6Address              *string `json:"ipv6Address,omitempty" xml:"ipv6Address,omitempty"`
	ParameterFormatType      string  `json:"parameterFormatType,omitempty" xml:"parameterFormatType,omitempty"`
	Password                 *string `json:"password,omitempty" xml:"password,omitempty"`
	PortNo                   uint64  `json:"portNo,omitempty" xml:"portNo,omitempty"`
	ProtocolType             string  `json:"protocolType,omitempty" xml:"protocolType,omitempty"`
	URL                      *string `json:"url,omitempty" xml:"url,omitempty"`
	UserName                 *string `json:"userName,omitempty" xml:"userName,omitempty"`
}

func (m *HTTPHostNotification) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddressingFormatType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHTTPAuthenticationMethod(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParameterFormatType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProtocolType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}

	return nil
}

var httpHostNotificationTypeAddressingFormatTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ipaddress","hostname"]`), &res); err != nil {
		panic(err)
	}

	for _, v := range res {
		httpHostNotificationTypeAddressingFormatTypePropEnum = append(httpHostNotificationTypeAddressingFormatTypePropEnum, v)
	}
}

const (
	// HTTPHostNotificationAddressingFormatTypeIpaddress captures enum value "ipaddress"
	HTTPHostNotificationAddressingFormatTypeIpaddress string = "ipaddress"

	// HTTPHostNotificationAddressingFormatTypeHostname captures enum value "hostname"
	HTTPHostNotificationAddressingFormatTypeHostname string = "hostname"
)

// prop value enum
func (m *HTTPHostNotification) validateAddressingFormatTypeEnum(path, location string, value string) error {
	return validate.EnumCase(path, location, value, httpHostNotificationTypeAddressingFormatTypePropEnum, true)
}

func (m *HTTPHostNotification) validateAddressingFormatType(_ strfmt.Registry) error {
	if swag.IsZero(m.AddressingFormatType) { // not required
		return nil
	}

	// value enum
	if err := m.validateAddressingFormatTypeEnum("addressingFormatType", "body", m.AddressingFormatType); err != nil {
		return err
	}

	return nil
}

var httpHostNotificationTypeHTTPAuthenticationMethodPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["MD5digest","none"]`), &res); err != nil {
		panic(err)
	}

	for _, v := range res {
		httpHostNotificationTypeHTTPAuthenticationMethodPropEnum = append(httpHostNotificationTypeHTTPAuthenticationMethodPropEnum, v)
	}
}

const (
	// HTTPHostNotificationHTTPAuthenticationMethodMD5digest captures enum value "MD5digest"
	HTTPHostNotificationHTTPAuthenticationMethodMD5digest string = "MD5digest"

	// HTTPHostNotificationHTTPAuthenticationMethodNone captures enum value "none"
	HTTPHostNotificationHTTPAuthenticationMethodNone string = "none"
)

// prop value enum
func (m *HTTPHostNotification) validateHTTPAuthenticationMethodEnum(path, location string, value string) error {
	return validate.EnumCase(path, location, value, httpHostNotificationTypeHTTPAuthenticationMethodPropEnum, true)
}

func (m *HTTPHostNotification) validateHTTPAuthenticationMethod(_ strfmt.Registry) error {
	if swag.IsZero(m.HTTPAuthenticationMethod) { // not required
		return nil
	}

	// value enum
	if err := m.validateHTTPAuthenticationMethodEnum("httpAuthenticationMethod", "body", m.HTTPAuthenticationMethod); err != nil {
		return err
	}

	return nil
}

var httpHostNotificationTypeParameterFormatTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["XML","querystring"]`), &res); err != nil {
		panic(err)
	}

	for _, v := range res {
		httpHostNotificationTypeParameterFormatTypePropEnum = append(httpHostNotificationTypeParameterFormatTypePropEnum, v)
	}
}

const (
	// HTTPHostNotificationParameterFormatTypeXML captures enum value "XML"
	HTTPHostNotificationParameterFormatTypeXML string = "XML"

	// HTTPHostNotificationParameterFormatTypeQuerystring captures enum value "querystring"
	HTTPHostNotificationParameterFormatTypeQuerystring string = "querystring"
)

// prop value enum
func (m *HTTPHostNotification) validateParameterFormatTypeEnum(path, location string, value string) error {
	return validate.EnumCase(path, location, value, httpHostNotificationTypeParameterFormatTypePropEnum, true)
}

func (m *HTTPHostNotification) validateParameterFormatType(_ strfmt.Registry) error {
	if swag.IsZero(m.ParameterFormatType) { // not required
		return nil
	}

	// value enum
	if err := m.validateParameterFormatTypeEnum("parameterFormatType", "body", m.ParameterFormatType); err != nil {
		return err
	}

	return nil
}

var httpHostNotificationTypeProtocolTypePropEnum []interface{}

func init() {
	var res []string

	if err := json.Unmarshal([]byte(`["HTTP","HTTPS"]`), &res); err != nil {
		panic(err)
	}

	for _, v := range res {
		httpHostNotificationTypeProtocolTypePropEnum = append(httpHostNotificationTypeProtocolTypePropEnum, v)
	}
}

const (
	// HTTPHostNotificationProtocolTypeHTTP captures enum value "HTTP"
	HTTPHostNotificationProtocolTypeHTTP string = "HTTP"

	// HTTPHostNotificationProtocolTypeHTTPS captures enum value "HTTPS"
	HTTPHostNotificationProtocolTypeHTTPS string = "HTTPS"
)

// prop value enum
func (m *HTTPHostNotification) validateProtocolTypeEnum(path, location string, value string) error {
	return validate.EnumCase(path, location, value, httpHostNotificationTypeProtocolTypePropEnum, true)
}

func (m *HTTPHostNotification) validateProtocolType(_ strfmt.Registry) error {
	if swag.IsZero(m.ProtocolType) { // not required
		return nil
	}

	// value enum
	if err := m.validateProtocolTypeEnum("protocolType", "body", m.ProtocolType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HTTPHostNotification) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPHostNotification) UnmarshalBinary(b []byte) error {
	var res HTTPHostNotification

	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}

	*m = res

	return nil
}
