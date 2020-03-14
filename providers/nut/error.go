package nut

import (
	"errors"
)

const (
	ErrorCodeAccessDenied         = "ACCESS-DENIED"
	ErrorCodeUnknownUPS           = "UNKNOWN-UPS"
	ErrorCodeVariableNonSupported = "VAR-NOT-SUPPORTED"
	ErrorCodeCommandNotSupported  = "CMD-NOT-SUPPORTED"
	ErrorCodeInvalidArgument      = "INVALID-ARGUMENT"
	ErrorCodeInstantCommandFailed = "INSTCMD-FAILED"
	ErrorCodeSetFailed            = "SET-FAILED"
	ErrorCodeReadonly             = "READONLY"
	ErrorCodeTooLong              = "TOO-LONG"
	ErrorCodeFeatureNotSupported  = "FEATURE-NOT-SUPPORTED"
	ErrorCodeFeatureNotConfigured = "FEATURE-NOT-CONFIGURED"
	ErrorCodeAlreadySSLMode       = "ALREADY-SSL-MODE"
	ErrorCodeDriverNotConnected   = "DRIVER-NOT-CONNECTED"
	ErrorCodeDataStale            = "DATA-STALE"
	ErrorCodeAlreadyLoggedIn      = "ALREADY-LOGGED-IN"
	ErrorCodeInvalidPassword      = "INVALID-PASSWORD"
	ErrorCodeAlreadySetPassword   = "ALREADY-SET-PASSWORD"
	ErrorCodeInvalidUsername      = "INVALID-USERNAME"
	ErrorCodeAlreadySetUsername   = "ALREADY-SET-USERNAME"
	ErrorCodeUsernameRequired     = "USERNAME-REQUIRED"
	ErrorCodePasswordRequired     = "PASSWORD-REQUIRED"
	ErrorCodeUnknownCommand       = "UNKNOWN-COMMAND"
	ErrorCodeInvalidValue         = "INVALID-VALUE"
)

var (
	ErrAccessDenied         = errors.New("the client's host and/or authentication details (username, password) are not sufficient to execute the requested command")
	ErrUnknownUPS           = errors.New("the UPS specified in the request is not known to upsd")
	ErrVariableNonSupported = errors.New("the specified UPS doesn’t support the variable in the request")
	ErrCommandNotSupported  = errors.New("the specified UPS doesn’t support the instant command in the request")
	ErrInvalidArgument      = errors.New("the client sent an argument to a command which is not recognized or is otherwise invalid in this context")
	ErrInstantCommandFailed = errors.New("upsd failed to deliver the instant command request to the driver")
	ErrSetFailed            = errors.New("upsd failed to deliver the set request to the driver")
	ErrReadonly             = errors.New("the requested variable in a SET command is not writable")
	ErrTooLong              = errors.New("the requested value in a SET command is too long")
	ErrFeatureNotSupported  = errors.New("this instance of upsd does not support the requested feature")
	ErrFeatureNotConfigured = errors.New("this instance of upsd hasn't been configured properly to allow the requested feature to operate")
	ErrAlreadySSLMode       = errors.New("TLS/SSL mode is already enabled on this connection, so upsd can't start it again")
	ErrDriverNotConnected   = errors.New("upsd can't perform the requested command, since the driver for that UPS is not connected")
	ErrDataStale            = errors.New("upsd is connected to the driver for the UPS, but that driver isn't providing regular updates or has specifically marked the data as stale")
	ErrAlreadyLoggedIn      = errors.New("the client already sent LOGIN for a UPS and can't do it again")
	ErrInvalidPassword      = errors.New("the client sent an invalid PASSWORD - perhaps an empty one")
	ErrAlreadySetPassword   = errors.New("the client already set a PASSWORD and can't set another")
	ErrInvalidUsername      = errors.New("the client sent an invalid USERNAME")
	ErrAlreadySetUsername   = errors.New("the client has already set a USERNAME, and can't set another")
	ErrUsernameRequired     = errors.New("the requested command requires a username for authentication, but the client hasn't set one")
	ErrPasswordRequired     = errors.New("the requested command requires a password for authentication, but the client hasn't set one")
	ErrUnknownCommand       = errors.New("upsd doesn't recognize the requested command")
	ErrInvalidValue         = errors.New("the value specified in the request is not valid")
)

func ErrorByCode(code string) error {
	switch code {
	case ErrorCodeAccessDenied:
		return ErrAccessDenied
	case ErrorCodeUnknownUPS:
		return ErrUnknownUPS
	case ErrorCodeVariableNonSupported:
		return ErrVariableNonSupported
	case ErrorCodeCommandNotSupported:
		return ErrCommandNotSupported
	case ErrorCodeInvalidArgument:
		return ErrInvalidArgument
	case ErrorCodeInstantCommandFailed:
		return ErrInstantCommandFailed
	case ErrorCodeSetFailed:
		return ErrSetFailed
	case ErrorCodeReadonly:
		return ErrReadonly
	case ErrorCodeTooLong:
		return ErrTooLong
	case ErrorCodeFeatureNotSupported:
		return ErrFeatureNotSupported
	case ErrorCodeFeatureNotConfigured:
		return ErrFeatureNotConfigured
	case ErrorCodeAlreadySSLMode:
		return ErrAlreadySSLMode
	case ErrorCodeDriverNotConnected:
		return ErrDriverNotConnected
	case ErrorCodeDataStale:
		return ErrDataStale
	case ErrorCodeAlreadyLoggedIn:
		return ErrAlreadyLoggedIn
	case ErrorCodeInvalidPassword:
		return ErrInvalidPassword
	case ErrorCodeAlreadySetPassword:
		return ErrAlreadySetPassword
	case ErrorCodeInvalidUsername:
		return ErrInvalidUsername
	case ErrorCodeAlreadySetUsername:
		return ErrAlreadySetUsername
	case ErrorCodeUsernameRequired:
		return ErrUsernameRequired
	case ErrorCodePasswordRequired:
		return ErrPasswordRequired
	case ErrorCodeUnknownCommand:
		return ErrUnknownCommand
	case ErrorCodeInvalidValue:
		return ErrInvalidValue
	}

	return errors.New("unknown error " + code)
}
