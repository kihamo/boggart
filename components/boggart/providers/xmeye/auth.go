package xmeye

import (
	"fmt"
)

func (c *Client) Login() error {
	response := &LoginResponse{}

	err := c.CallWithResult(CmdLoginResponse, map[string]string{
		"EncryptType": "MD5",
		"LoginType":   "DVRIP-Web",
		"PassWord":    HashPassword(c.password),
		"UserName":    c.username,
	}, response)

	if err != nil {
		return err
	}

	if response.Ret != CodeOK {
		return fmt.Errorf("response %d isn't ok", response.Ret)
	}

	if response.AliveInterval > 0 {
		go c.keepAlive(response.AliveInterval)
	}

	return err
}

func (c *Client) Logout() error {
	_, err := c.Call(CmdLogoutResponse, nil)
	return err
}
