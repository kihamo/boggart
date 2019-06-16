// Code generated by go-swagger; DO NOT EDIT.

package image

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new image API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for image API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetImageChannels get image channels API
*/
func (a *Client) GetImageChannels(params *GetImageChannelsParams, authInfo runtime.ClientAuthInfoWriter) (*GetImageChannelsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetImageChannelsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getImageChannels",
		Method:             "GET",
		PathPattern:        "/Image/channels",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetImageChannelsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetImageChannelsOK), nil

}

/*
SetImageIrCutFilter set image ir cut filter API
*/
func (a *Client) SetImageIrCutFilter(params *SetImageIrCutFilterParams, authInfo runtime.ClientAuthInfoWriter) (*SetImageIrCutFilterOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetImageIrCutFilterParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setImageIrCutFilter",
		Method:             "PUT",
		PathPattern:        "/Image/channels/{channel}/IrcutFilter",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SetImageIrCutFilterReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetImageIrCutFilterOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}