package xmeye

import (
	"fmt"
	"time"
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

func (c *Client) OPTime() (*time.Time, error) {
	response := &OPTimeQuery{}

	err := c.CmdWithResult(CmdTimeRequest, "OPTimeQuery", response)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timeLayout, response.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) SystemFunctions() (*SystemFunctions, error) {
	result := &SystemFunctions{}

	err := c.CmdWithResult(CmdAbilityGetRequest, "SystemFunction", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) SystemInfo() (*SystemInfo, error) {
	result := &SystemInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "SystemInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) OEMInfo() (*OEMInfo, error) {
	result := &OEMInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "OEMInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) StorageInfo() (*StorageInfo, error) {
	result := &StorageInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "StorageInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) WorkState() (*WorkState, error) {
	result := &WorkState{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "WorkState", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) LogSearch(begin time.Time, end time.Time) (interface{}, error) {
	var result interface{}

	err := c.CallWithResult(CmdLogSearchRequest, map[string]interface{}{
		"Name":      "OPLogQuery",
		"SessionID": c.sessionIDAsString(),
		"OPLogQuery": map[string]interface{}{
			"BeginTime":   begin.Format(timeLayout),
			"EndTime":     end.Format(timeLayout),
			"LogPosition": 0,
			"Type":        "LogAll",
		},
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) Reboot() error {
	_, err := c.Call(CmdSysManagerResponse, map[string]interface{}{
		"Name":      "OPMachine",
		"SessionID": c.sessionIDAsString(),
		"OPMachine": map[string]string{
			"Action": "Reboot",
		},
	})

	return err
}
