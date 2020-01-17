package mikrotik

import (
	"context"
)

type SystemHealth struct {
	Voltage     float64 `json:"voltage"`
	Temperature uint64  `json:"temperature"`
}

type SystemRouterboard struct {
	Routerboard     bool   `json:"routerboard"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serial-number"`
	FirmwareType    string `json:"firmware-type"`
	FactoryFirmware string `json:"factory-firmware"`
	CurrentFirmware string `json:"current-firmware"`
	UpgradeFirmware string `json:"upgrade-firmware"`
}

type SystemResource struct {
	Uptime               string  `json:"uptime"`
	Version              string  `json:"version"`
	BuildTime            string  `json:"build-time"`
	FactorySoftware      string  `json:"factory-software"`
	FreeMemory           uint64  `json:"free-memory"`
	TotalMemory          uint64  `json:"total-memory"`
	CPU                  string  `json:"cpu"`
	CPUCount             uint64  `json:"cpu-count"`
	CPUFrequency         uint64  `json:"cpu-frequency"`
	CPULoad              uint64  `json:"cpu-load"`
	FreeHDDSpace         uint64  `json:"free-hdd-space"`
	TotalHDDSpace        uint64  `json:"total-hdd-space"`
	WriteSectSinceReboot float64 `json:"write-sect-since-reboot"`
	WriteSectTotal       float64 `json:"write-sect-total"`
	BadBlocks            float64 `json:"bad-blocks"`
	ArchitectureName     string  `json:"architecture-name"`
	BoardName            string  `json:"board-name"`
	Platform             string  `json:"platform"`
}

type SystemPackageUpdate struct {
	Channel          string `json:"channel"`
	InstalledVersion string `json:"installed-version"`
	LatestVersion    string `json:"latest-version,omitempty"`
	Status           string `json:"status,omitempty"`
}

type SystemDisk struct {
	Id    string `json:".id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Disk  string `json:"disk"`
	Free  uint64 `json:"free"`
	Size  uint64 `json:"size"`
}

func (c *Client) SystemHealth(ctx context.Context) (result SystemHealth, err error) {
	var list []SystemHealth

	err = c.doConvert(ctx, []string{"/system/health/print"}, &list)
	if err != nil {
		return result, err
	}

	if len(list) == 0 {
		return result, ErrEmptyResponse
	}

	return list[0], nil
}

func (c *Client) SystemRouterboard(ctx context.Context) (result SystemRouterboard, err error) {
	var list []SystemRouterboard

	err = c.doConvert(ctx, []string{"/system/routerboard/print"}, &list)
	if err != nil {
		return result, err
	}

	if len(list) == 0 {
		return result, ErrEmptyResponse
	}

	return list[0], nil
}

func (c *Client) SystemResource(ctx context.Context) (result SystemResource, err error) {
	var list []SystemResource

	err = c.doConvert(ctx, []string{"/system/resource/print"}, &list)
	if err != nil {
		return result, err
	}

	if len(list) == 0 {
		return result, ErrEmptyResponse
	}

	return list[0], nil
}

// need policies: write, policy
func (c *Client) SystemPackageUpdateCheck(ctx context.Context) (result SystemPackageUpdate, err error) {
	var list []SystemPackageUpdate

	err = c.doConvert(ctx, []string{"/system/package/update/print"}, &list)
	if err != nil {
		return result, err
	}

	if len(list) == 0 {
		return result, ErrEmptyResponse
	}

	return list[len(list)-1], nil
}

func (c *Client) SystemDisk(ctx context.Context) (result []SystemDisk, err error) {
	err = c.doConvert(ctx, []string{"/disk/print"}, &result)
	return result, err
}
