package xmeye

import (
	"bytes"
	"encoding/json"
	"errors"
)

const (
	CodeOK                  = 100
	CodeError               = 101
	CodeVersion             = 102
	CodeRequest             = 103
	CodeLogged              = 104
	CodeNotLogged           = 105
	CodeWrongCredentials    = 106
	CodeAccessDenied        = 107
	CodeTimeout             = 108
	CodeFileNotFound        = 109
	CodeCompleteSearch      = 110
	CodePartialSearch       = 111
	CodeUserExists          = 112
	CodeUserNotExists       = 113
	CodeGroupExists         = 114
	CodeGroupNotExists      = 115
	CodeInvalidMessage      = 117
	CodePTZ                 = 118
	CodeSearchNoResults     = 119
	CodeDisabled            = 120
	CodeChannelNotConnected = 121
	CodeReboot              = 150
	Code202                 = 202
	CodePasswordWrong       = 203
	CodeUsernameWrong       = 204
	CodeLockOut             = 205
	CodeBanned              = 206
	CodeConflict            = 207
	CodeInput               = 208
	Code209                 = 209
	Code210                 = 210
	CodeObject              = 211
	CodeAccount             = 212
	CodeSubnet              = 213
	CodePasswordCharacters  = 214
	CodePasswordMatch       = 215
	CodeUsernameReserved    = 216
	CodeCommand             = 502
	CodeIntercomOn          = 503
	CodeIntercomOff         = 504
	CodeUpgradeStarted      = 511
	CodeUpgradeNotStarted   = 512
	CodeUpgradeData         = 513
	CodeUpgradeSuccessful   = 514
	CodeUpgradeFailed       = 515
	CodeResetFailed         = 521
	CodeResetSuccessful     = 522
	CodeResetInvalid        = 523
	CodeImportSuccessful    = 602
	CodeImportReboot        = 603
	CodeWriting             = 604
	CodeFeature             = 605
	CodeReading             = 606
	CodeNoImport            = 607
	CodeSyntax              = 608
)

var codeErrorsText = map[int64]string{
	CodeError:               "unknown error",
	CodeVersion:             "invalid version",
	CodeRequest:             "invalid request",
	CodeLogged:              "already logged in",
	CodeNotLogged:           "not logged in",
	CodeWrongCredentials:    "wrong username or password",
	CodeAccessDenied:        "access denied",
	CodeTimeout:             "timeout",
	CodeFileNotFound:        "file not found",
	CodeUserExists:          "user already exists",
	CodeUserNotExists:       "user does not exist",
	CodeGroupExists:         "group already exists",
	CodeGroupNotExists:      "group does not exist",
	CodeInvalidMessage:      "invalid message",
	CodePTZ:                 "PTZ protocol not set'",
	CodeDisabled:            "disabled",
	CodeChannelNotConnected: "channel not connected",
	Code202:                 "error 202",
	CodePasswordWrong:       "wrong password",
	CodeUsernameWrong:       "wrong username",
	CodeLockOut:             "locked out",
	CodeBanned:              "banned",
	CodeConflict:            "already logged in",
	CodeInput:               "illegal value",
	Code209:                 "error 209",
	Code210:                 "error 210",
	CodeObject:              "object does not exist",
	CodeAccount:             "account in use",
	CodeSubnet:              "subset larger than superset",
	CodePasswordCharacters:  "illegal characters in password",
	CodePasswordMatch:       "passwords do not match",
	CodeUsernameReserved:    "username reserved",
	CodeCommand:             "illegal command",
	CodeUpgradeNotStarted:   "upgrade not started",
	CodeUpgradeData:         "invalid upgrade data",
	CodeUpgradeFailed:       "upgrade failed",
	CodeResetFailed:         "reset failed",
	CodeResetInvalid:        "reset data invalid",
	CodeWriting:             "configuration write failed",
	CodeFeature:             "unsupported feature in configuration",
	CodeReading:             "configuration read failed",
	CodeNoImport:            "configuration not found",
	CodeSyntax:              "illegal configuration syntax",
}

var codeSuccessText = map[int64]string{
	CodeOK:                "ok",
	CodeCompleteSearch:    "complete search results",
	CodePartialSearch:     "partial search results",
	CodeSearchNoResults:   "no search results",
	CodeReboot:            "reboot",
	CodeIntercomOn:        "intercom turned on",
	CodeIntercomOff:       "intercom turned off",
	CodeUpgradeStarted:    "upgrade started",
	CodeUpgradeSuccessful: "upgrade successful",
	CodeResetSuccessful:   "reset successful, reboot required",
	CodeImportSuccessful:  "import successful, restart required",
	CodeImportReboot:      "import successful, reboot required",
}

type Payload struct {
	*bytes.Buffer
}

func NewPayload() *Payload {
	return &Payload{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (p *Payload) JSONUnmarshal(v interface{}) error {
	// обрезаем признак конца строки
	payload := p.Bytes()

	if len(payload) > 2 {
		payload = payload[:len(payload)-len(payloadEOF)]
	}

	return json.Unmarshal(payload, v)
}

func (p *Payload) Error() error {
	var result Response

	if err := p.JSONUnmarshal(&result); err == nil {
		text, ok := codeErrorsText[int64(result.Ret)]
		if ok {
			return errors.New(text)
		}
	}

	return nil
}
