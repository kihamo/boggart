package hikvision

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
)

const (
	proxyPTZPrefixURL = "/PTZCtrl"
)

type PTZData struct {
	XMLName      xml.Name             `xml:"PTZData"`
	Relative     *PTZDataRelative     `xml:"Relative,omitempty"`
	AbsoluteHigh *PTZDataAbsoluteHigh `xml:"AbsoluteHigh,omitempty"`
}

type PTZDataRelative struct {
	PositionX    int64 `xml:"positionX,omitempty"`
	PositionY    int64 `xml:"positionY,omitempty"`
	RelativeZoom int64 `xml:"relativeZoom,omitempty"`
}

type PTZDataAbsoluteHigh struct {
	Elevation    int64  `xml:"elevation,omitempty"`
	Azimuth      uint64 `xml:"azimuth,omitempty"`
	AbsoluteZoom uint64 `xml:"absoluteZoom,omitempty"`
}

func (a *ISAPI) PTZPresetGoTo(ctx context.Context, channel, preset uint64) error {
	u := a.address + proxyPTZPrefixURL + "/channels/" +
		strconv.FormatUint(channel, 10) +
		"/presets/" +
		strconv.FormatUint(preset, 10) +
		"/goto"

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, nil, &result)
	if err != nil {
		return err
	}

	if result.StatusCode != 1 {
		return errors.New(result.StatusString)
	}

	return nil
}

func (a *ISAPI) PTZRelative(ctx context.Context, channel uint64, x, y, zoom int64) error {
	if zoom < -100 {
		zoom = -100
	} else if zoom > 100 {
		zoom = 100
	}

	u := a.address + proxyPTZPrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/relative"

	data := PTZData{
		Relative: &PTZDataRelative{
			PositionX:    x,
			PositionY:    y,
			RelativeZoom: zoom,
		},
	}

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, data, &result)
	if err != nil {
		return err
	}

	if result.StatusCode != 1 {
		return errors.New(result.StatusString)
	}

	return nil
}

func (a *ISAPI) PTZAbsolute(ctx context.Context, channel uint64, elevation int64, azimuth, absoluteZoom uint64) error {
	if elevation < -900 {
		elevation = -900
	} else if elevation > 2700 {
		elevation = 2700
	}

	if azimuth < 0 {
		azimuth = 0
	} else if azimuth > 3600 {
		azimuth = 3600
	}

	if absoluteZoom < 0 {
		absoluteZoom = 0
	} else if absoluteZoom > 1000 {
		absoluteZoom = 1000
	}

	u := a.address + proxyPTZPrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/absolute"

	data := PTZData{
		AbsoluteHigh: &PTZDataAbsoluteHigh{
			Elevation:    elevation,
			Azimuth:      azimuth,
			AbsoluteZoom: absoluteZoom,
		},
	}

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, data, &result)
	if err != nil {
		return err
	}

	if result.StatusCode != 1 {
		return errors.New(result.StatusString)
	}

	return nil
}
