// Code generated by go-swagger; DO NOT EDIT.

package authorization

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new authorization API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for authorization API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	Login(params *LoginParams, authInfo runtime.ClientAuthInfoWriter) (*LoginOK, error)

	Logout(params *LogoutParams, authInfo runtime.ClientAuthInfoWriter) (*LogoutNoContent, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  Login logins
*/
func (a *Client) Login(params *LoginParams, authInfo runtime.ClientAuthInfoWriter) (*LoginOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewLoginParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "login",
		Method:             "POST",
		PathPattern:        "/api/login",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &LoginReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*LoginOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for login: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  Logout logouts
*/
func (a *Client) Logout(params *LogoutParams, authInfo runtime.ClientAuthInfoWriter) (*LogoutNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewLogoutParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "logout",
		Method:             "POST",
		PathPattern:        "/api/logout",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &LogoutReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*LogoutNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for logout: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
