// Code generated by go-swagger; DO NOT EDIT.

package ptz

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new ptz API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for ptz API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetPtzChannels its is used to get the list of p t z channels for the device
*/
func (a *Client) GetPtzChannels(params *GetPtzChannelsParams, authInfo runtime.ClientAuthInfoWriter) (*GetPtzChannelsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPtzChannelsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getPtzChannels",
		Method:             "GET",
		PathPattern:        "/PTZCtrl/channels",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetPtzChannelsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetPtzChannelsOK), nil

}

/*
GetPtzStatus its is used to get currently p t z coordinate position for the device
*/
func (a *Client) GetPtzStatus(params *GetPtzStatusParams, authInfo runtime.ClientAuthInfoWriter) (*GetPtzStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPtzStatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getPtzStatus",
		Method:             "GET",
		PathPattern:        "/PTZCtrl/channels/{channel}/status",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetPtzStatusReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetPtzStatusOK), nil

}

/*
GotoPtzPreset its is used to move a particular p t z channel to a ID preset position for the device
*/
func (a *Client) GotoPtzPreset(params *GotoPtzPresetParams, authInfo runtime.ClientAuthInfoWriter) (*GotoPtzPresetOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGotoPtzPresetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "gotoPtzPreset",
		Method:             "PUT",
		PathPattern:        "/PTZCtrl/channels/{channel}/presets/{preset}/goto",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GotoPtzPresetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GotoPtzPresetOK), nil

}

/*
SetPtzContinuous its is used to control p t z move around and zoom for the device
*/
func (a *Client) SetPtzContinuous(params *SetPtzContinuousParams, authInfo runtime.ClientAuthInfoWriter) (*SetPtzContinuousOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPtzContinuousParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setPtzContinuous",
		Method:             "PUT",
		PathPattern:        "/PTZCtrl/channels/{channel}/continuous",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SetPtzContinuousReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetPtzContinuousOK), nil

}

/*
SetPtzMomentary its is used to control p t z move around and zoom in a period of time for the device
*/
func (a *Client) SetPtzMomentary(params *SetPtzMomentaryParams, authInfo runtime.ClientAuthInfoWriter) (*SetPtzMomentaryOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPtzMomentaryParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setPtzMomentary",
		Method:             "PUT",
		PathPattern:        "/PTZCtrl/channels/{channel}/momentary",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SetPtzMomentaryReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetPtzMomentaryOK), nil

}

/*
SetPtzPositionAbsolute its is used to move a particular p t z channel to a absolute position which is defined by absolute for the device
*/
func (a *Client) SetPtzPositionAbsolute(params *SetPtzPositionAbsoluteParams, authInfo runtime.ClientAuthInfoWriter) (*SetPtzPositionAbsoluteOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPtzPositionAbsoluteParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setPtzPositionAbsolute",
		Method:             "PUT",
		PathPattern:        "/PTZCtrl/channels/{channel}/absolute",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SetPtzPositionAbsoluteReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetPtzPositionAbsoluteOK), nil

}

/*
SetPtzPositionRelative its is used to move the position which is defined by position x position y to the screen center and relative zoom for the device
*/
func (a *Client) SetPtzPositionRelative(params *SetPtzPositionRelativeParams, authInfo runtime.ClientAuthInfoWriter) (*SetPtzPositionRelativeOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPtzPositionRelativeParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "setPtzPositionRelative",
		Method:             "PUT",
		PathPattern:        "/PTZCtrl/channels/{channel}/relative",
		ProducesMediaTypes: []string{"application/xml"},
		ConsumesMediaTypes: []string{"application/xml"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SetPtzPositionRelativeReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*SetPtzPositionRelativeOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
