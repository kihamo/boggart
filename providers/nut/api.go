package nut

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (s *Session) ListUPS() ([]ListUPS, error) {
	response, err := s.Request("LIST UPS")
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(response, separator)
	result := make([]ListUPS, 0, len(lines))

	var info ListUPS

	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(" "), 3)
		if len(parts) != 3 {
			continue
		}

		info = ListUPS{
			Name: string(parts[1]),
		}

		if len(parts[2]) > 0 {
			info.Description = string(bytes.Trim(parts[2], "\""))
		}

		result = append(result, info)
	}

	return result, nil
}

func (s *Session) ListVariables(ups string) ([]ListVariable, error) {
	response, err := s.Request("LIST VAR " + ups)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(response, separator)
	result := make([]ListVariable, 0, len(lines))

	var info ListVariable

	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(" "), 4)
		if len(parts) != 4 {
			continue
		}

		info = ListVariable{
			UPS:  string(parts[1]),
			Name: string(parts[2]),
		}

		if len(parts[3]) > 0 {
			info.Value = string(bytes.Trim(parts[3], "\""))
		}

		result = append(result, info)
	}

	return result, nil
}

func (s *Session) ListCommands(ups string) ([]ListCommand, error) {
	response, err := s.Request("LIST CMD " + ups)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(response, separator)
	result := make([]ListCommand, 0, len(lines))

	var info ListCommand

	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(" "), 3)
		if len(parts) != 3 {
			continue
		}

		info = ListCommand{
			UPS:  string(parts[1]),
			Name: string(parts[2]),
		}

		result = append(result, info)
	}

	return result, nil
}

func (s *Session) Type(ups, name string) (t Type, err error) {
	response, err := s.Request("GET TYPE " + ups + " " + name)
	if err != nil {
		return t, err
	}

	parts := bytes.SplitN(response, []byte(" "), 4)
	if len(parts) != 4 {
		return t, errors.New("bad response")
	}

	types := strings.Fields(string(parts[3]))

	for _, tp := range types {
		if tp == "RW" {
			t.Writeable = true
			continue
		}

		if strings.HasPrefix(tp, "STRING:") {
			t.Name = "STRING"
			parts := strings.SplitN(tp, ":", 2)

			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return t, err
			}

			t.MaxLength = v

			continue
		}

		t.Name = tp
	}

	return t, nil
}

func (s *Session) Description(ups, name string) (string, error) {
	response, err := s.Request("GET DESC " + ups + " " + name)
	if err != nil {
		return "", err
	}

	parts := bytes.SplitN(response, []byte(" "), 4)
	if len(parts) != 4 {
		return "", errors.New("bad response")
	}

	return string(bytes.Trim(parts[3], "\"")), nil
}

func (s *Session) CommandDescription(ups, name string) (string, error) {
	response, err := s.Request("GET CMDDESC " + ups + " " + name)
	if err != nil {
		return "", err
	}

	parts := bytes.SplitN(response, []byte(" "), 4)
	if len(parts) != 4 {
		return "", errors.New("bad response")
	}

	return string(bytes.Trim(parts[3], "\"")), nil
}

func (s *Session) Set(ups, name string, value interface{}) error {
	return s.Call(fmt.Sprintf("SET VAR %s %s \"%s\"", ups, name, value))
}

func (s *Session) InstantCommand(ups, cmd string) error {
	return s.Call("INSTCMD " + ups + " " + cmd)
}

func (s *Session) Logout() error {
	return s.Call("LOGOUT")
}

func (s *Session) Password(password string) error {
	return s.Call("PASSWORD " + password)
}

func (s *Session) Username(username string) error {
	return s.Call("USERNAME " + username)
}

func (s *Session) Help() ([]string, error) {
	response, err := s.Request("HELP")
	if err != nil {
		return nil, err
	}

	response = bytes.TrimPrefix(response, []byte("Commands:"))
	response = bytes.TrimSuffix(response, []byte("."))
	response = bytes.TrimSpace(response)

	commands := bytes.Split(response, []byte(" "))
	result := make([]string, len(commands))

	for i, cmd := range commands {
		result[i] = string(cmd)
	}

	return result, nil
}

func (s *Session) Version() (string, error) {
	response, err := s.Request("VER")
	return string(response), err
}

func (s *Session) NetVersion() (string, error) {
	response, err := s.Request("NETVER")
	return string(response), err
}
