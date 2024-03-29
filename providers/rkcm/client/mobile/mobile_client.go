// Code generated by go-swagger; DO NOT EDIT.

package mobile

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new mobile API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for mobile API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetAccountIdents gets account idents
*/
func (a *Client) GetAccountIdents(params *GetAccountIdentsParams) (*GetAccountIdentsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetAccountIdentsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getAccountIdents",
		Method:             "GET",
		PathPattern:        "/MobileAPI/GetAccountIdents.ashx",
		ProducesMediaTypes: []string{""},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetAccountIdentsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetAccountIdentsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getAccountIdents: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetDebt gets debt
*/
func (a *Client) GetDebt(params *GetDebtParams) (*GetDebtOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetDebtParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getDebt",
		Method:             "GET",
		PathPattern:        "/MobileAPI/GetDebt.ashx",
		ProducesMediaTypes: []string{""},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetDebtReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetDebtOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getDebt: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
