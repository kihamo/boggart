// Code generated by go-swagger; DO NOT EDIT.

package client

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/wiim/httpapi/client/device"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client/eq"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client/networking"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client/playback"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client/presets"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client/track"
)

// Default wiim HTTP client.
var Default = NewHTTPClient(nil)

const (
	// DefaultHost is the default Host
	// found in Meta (info) section of spec file
	DefaultHost string = "localhost"
	// DefaultBasePath is the default BasePath
	// found in Meta (info) section of spec file
	DefaultBasePath string = "/"
)

// DefaultSchemes are the default schemes found in Meta (info) section of spec file
var DefaultSchemes = []string{"https"}

// NewHTTPClient creates a new wiim HTTP client.
func NewHTTPClient(formats strfmt.Registry) *Wiim {
	return NewHTTPClientWithConfig(formats, nil)
}

// NewHTTPClientWithConfig creates a new wiim HTTP client,
// using a customizable transport config.
func NewHTTPClientWithConfig(formats strfmt.Registry, cfg *TransportConfig) *Wiim {
	// ensure nullable parameters have default
	if cfg == nil {
		cfg = DefaultTransportConfig()
	}

	// create transport and client
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	return New(transport, formats)
}

// New creates a new wiim client
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Wiim {
	// ensure nullable parameters have default
	if formats == nil {
		formats = strfmt.Default
	}

	cli := new(Wiim)
	cli.Transport = transport
	cli.Device = device.New(transport, formats)
	cli.Eq = eq.New(transport, formats)
	cli.Networking = networking.New(transport, formats)
	cli.Playback = playback.New(transport, formats)
	cli.Presets = presets.New(transport, formats)
	cli.Track = track.New(transport, formats)
	return cli
}

// DefaultTransportConfig creates a TransportConfig with the
// default settings taken from the meta section of the spec file.
func DefaultTransportConfig() *TransportConfig {
	return &TransportConfig{
		Host:     DefaultHost,
		BasePath: DefaultBasePath,
		Schemes:  DefaultSchemes,
	}
}

// TransportConfig contains the transport related info,
// found in the meta section of the spec file.
type TransportConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

// WithHost overrides the default host,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithHost(host string) *TransportConfig {
	cfg.Host = host
	return cfg
}

// WithBasePath overrides the default basePath,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithBasePath(basePath string) *TransportConfig {
	cfg.BasePath = basePath
	return cfg
}

// WithSchemes overrides the default schemes,
// provided by the meta section of the spec file.
func (cfg *TransportConfig) WithSchemes(schemes []string) *TransportConfig {
	cfg.Schemes = schemes
	return cfg
}

// Wiim is a client for wiim
type Wiim struct {
	Device device.ClientService

	Eq eq.ClientService

	Networking networking.ClientService

	Playback playback.ClientService

	Presets presets.ClientService

	Track track.ClientService

	Transport runtime.ClientTransport
}

// SetTransport changes the transport on the client and all its subresources
func (c *Wiim) SetTransport(transport runtime.ClientTransport) {
	c.Transport = transport
	c.Device.SetTransport(transport)
	c.Eq.SetTransport(transport)
	c.Networking.SetTransport(transport)
	c.Playback.SetTransport(transport)
	c.Presets.SetTransport(transport)
	c.Track.SetTransport(transport)
}
