// Code generated by go-swagger; DO NOT EDIT.

package printer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new printer API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for printer API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	GetBedState(params *GetBedStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetBedStateOK, error)

	GetChamberState(params *GetChamberStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetChamberStateOK, error)

	GetPrinterState(params *GetPrinterStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetPrinterStateOK, error)

	GetSDState(params *GetSDStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetSDStateOK, error)

	GetToolState(params *GetToolStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetToolStateOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetBedState retrieves the current bed state
*/
func (a *Client) GetBedState(params *GetBedStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetBedStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetBedStateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getBedState",
		Method:             "GET",
		PathPattern:        "/api/printer/bed",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetBedStateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetBedStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getBedState: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetChamberState retrieves the current chamber state
*/
func (a *Client) GetChamberState(params *GetChamberStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetChamberStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetChamberStateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getChamberState",
		Method:             "GET",
		PathPattern:        "/api/printer/chamber",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetChamberStateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetChamberStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getChamberState: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetPrinterState retrieves the current printer state
*/
func (a *Client) GetPrinterState(params *GetPrinterStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetPrinterStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPrinterStateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getPrinterState",
		Method:             "GET",
		PathPattern:        "/api/printer",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetPrinterStateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetPrinterStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getPrinterState: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetSDState retrieves the current s d state
*/
func (a *Client) GetSDState(params *GetSDStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetSDStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetSDStateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getSDState",
		Method:             "GET",
		PathPattern:        "/api/printer/sd",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetSDStateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetSDStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getSDState: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  GetToolState retrieves the current tool state
*/
func (a *Client) GetToolState(params *GetToolStateParams, authInfo runtime.ClientAuthInfoWriter) (*GetToolStateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetToolStateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getToolState",
		Method:             "GET",
		PathPattern:        "/api/printer/tool",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetToolStateReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetToolStateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getToolState: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
