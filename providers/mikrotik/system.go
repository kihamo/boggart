package mikrotik

import (
	"context"
)

type SystemHealth struct {
	Voltage     float64 `mapstructure:"voltage"`
	Temperature uint64  `mapstructure:"temperature"`
}

type SystemRouterBoard struct {
	RouterBoard     bool   `mapstructure:"routerboard"`
	Model           string `mapstructure:"model"`
	SerialNumber    string `mapstructure:"serial-number"`
	FirmwareType    string `mapstructure:"firmware-type"`
	FactoryFirmware string `mapstructure:"factory-firmware"`
	CurrentFirmware string `mapstructure:"current-firmware"`
	UpgradeFirmware string `mapstructure:"upgrade-firmware"`
}

type SystemResource struct {
	Uptime               string  `mapstructure:"uptime"`
	Version              string  `mapstructure:"version"`
	BuildTime            string  `mapstructure:"build-time"`
	FactorySoftware      string  `mapstructure:"factory-software"`
	FreeMemory           uint64  `mapstructure:"free-memory"`
	TotalMemory          uint64  `mapstructure:"total-memory"`
	CPU                  string  `mapstructure:"cpu"`
	CPUCount             uint64  `mapstructure:"cpu-count"`
	CPUFrequency         uint64  `mapstructure:"cpu-frequency"`
	CPULoad              uint64  `mapstructure:"cpu-load"`
	FreeHDDSpace         uint64  `mapstructure:"free-hdd-space"`
	TotalHDDSpace        uint64  `mapstructure:"total-hdd-space"`
	WriteSectSinceReboot float64 `mapstructure:"write-sect-since-reboot"`
	WriteSectTotal       float64 `mapstructure:"write-sect-total"`
	BadBlocks            float64 `mapstructure:"bad-blocks"`
	ArchitectureName     string  `mapstructure:"architecture-name"`
	BoardName            string  `mapstructure:"board-name"`
	Platform             string  `mapstructure:"platform"`
}

type SystemPackageUpdate struct {
	Channel          string `mapstructure:"channel"`
	InstalledVersion string `mapstructure:"installed-version"`
	LatestVersion    string `mapstructure:"latest-version,omitempty"`
	Status           string `mapstructure:"status,omitempty"`
}

type SystemDisk struct {
	Id    string `mapstructure:".id"`
	Name  string `mapstructure:"name"`
	Label string `mapstructure:"label"`
	Type  string `mapstructure:"type"`
	Disk  string `mapstructure:"disk"`
	Free  uint64 `mapstructure:"free"`
	Size  uint64 `mapstructure:"size"`
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

func (c *Client) SystemRouterBoard(ctx context.Context) (result SystemRouterBoard, err error) {
	var list []SystemRouterBoard

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
