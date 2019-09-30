package xmeye

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Payload struct {
	*bytes.Buffer
}

func NewPayload() *Payload {
	return &Payload{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (p *Payload) UnmarshalJSON(v interface{}) error {
	// обрезаем признак конца строки
	payload := p.Bytes()

	if len(payload) > 2 {
		payload = payload[:len(payload)-len(payloadEOF)]
	}

	return json.Unmarshal(payload, v)
}

func (p *Payload) Error() error {
	var result Response

	if err := p.UnmarshalJSON(&result); err == nil {
		switch result.Ret {
		case CodeOK, CodeUpgradeSuccessful, 0:
			return nil

		case CodeUnknownError:
			return errors.New("unknown error")

		case CodeUnsupportedVersion:
			return errors.New("unsupported version")

		case CodeRequestNotPermitted:
			return errors.New("request not permitted")

		case CodeUserAlreadyLoggedIn:
			return errors.New("user already logged in")

		case CodeUserUserIsNotLoggedIn:
			return errors.New("user is not logged in")

		case CodeRequestWrongFormat:
			return errors.New("request wrong format")

		case CodeUsernameOrPasswordIsIncorrect:
			return errors.New("username or password is incorrect")

		case CodeUserDoesNotHaveNecessaryPermissions:
			return errors.New("user does not have necessary permissions")

		case CodePasswordIsIncorrect:
			return errors.New("password is incorrect")

		case CodePasswordLockOut:
			return errors.New("locked out")

		case CodeConfigurationIsNotExists:
			return errors.New("configuration is not exists")

		default:
			return errors.New("unsupported unknown error")
		}
	}

	return nil
}
