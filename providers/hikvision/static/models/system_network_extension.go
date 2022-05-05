package models

import (
	"context"
	"encoding/xml"

	"github.com/go-openapi/strfmt"
)

type SystemNetworkExtensionMulticastAddress struct {
	IPVersion   string `xml:"ipVersion,omitempty"`
	IPAddress   string `xml:"ipAddress,omitempty"`
	IPv6Address string `xml:"ipv6Address,omitempty"`
}

type SystemNetworkExtension struct {
	XMLName xml.Name `xml:"networkExtension"`

	MulticastAddress  *SystemNetworkExtensionMulticastAddress `xml:"multicastAddress,omitempty"`
	EnableVirtualHost bool                                    `xml:"enVirtualHost"`
}

func (m *SystemNetworkExtension) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *SystemNetworkExtension) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
