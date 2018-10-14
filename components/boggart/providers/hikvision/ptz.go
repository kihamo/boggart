package hikvision

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

const (
	proxyPTZPrefixURL = "/PTZCtrl"
)

type PTZStatus struct {
	AbsoluteHigh PTZDataAbsoluteHigh `xml:"AbsoluteHigh"`
}

type PTZChannelList struct {
	Channels []PTZChannel `xml:"PTZChannel"`
}

type PTZChannel struct {
	ID                   uint64 `xml:"id"`
	Enabled              bool   `xml:"enabled"`
	VideoInputID         uint64 `xml:"videoInputID"`
	PanMaxSpeed          uint64 `xml:"panMaxSpeed"`
	TiltMaxSpeed         uint64 `xml:"tiltMaxSpeed"`
	PresetSpeed          uint64 `xml:"presetSpeed"`
	AutoPatrolSpeed      uint64 `xml:"autoPatrolSpeed"`
	KeyBoardControlSpeed string `xml:"keyBoardControlSpeed"`
	ControlProtocol      string `xml:"controlProtocol"`
	DefaultPresetID      uint64 `xml:"defaultPresetID"`
	PanSupport           bool   `xml:"panSupport"`
	TiltSupport          bool   `xml:"tiltSupport"`
	ZoomSupport          bool   `xml:"zoomSupport"`
	ManualControlSpeed   string `xml:"manualControlSpeed"`
}

type PTZData struct {
	XMLName      xml.Name             `xml:"PTZData"`
	Relative     *PTZDataRelative     `xml:"Relative,omitempty"`
	AbsoluteHigh *PTZDataAbsoluteHigh `xml:"AbsoluteHigh,omitempty"`
}

type PTZDataContinuous struct {
	XMLName   xml.Name `xml:"PTZData"`
	Pan       int64    `xml:"pan,omitempty"`
	Tilt      int64    `xml:"tilt,omitempty"`
	Zoom      int64    `xml:"zoom,omitempty"`
	Momentary struct {
		Duration int64 `xml:"duration"`
	} `xml:"Momentary,omitempty"`
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

func (a *ISAPI) PTZChannels(ctx context.Context) (list PTZChannelList, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.channels")
	defer span.Finish()

	u := a.address + proxyPTZPrefixURL + "/channels"

	err = a.DoXML(ctx, http.MethodGet, u, nil, &list)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return list, err
}

func (a *ISAPI) PTZStatus(ctx context.Context, channel uint64) (status PTZStatus, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.status")
	defer span.Finish()

	u := a.address + proxyPTZPrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/status"

	err = a.DoXML(ctx, http.MethodGet, u, nil, &status)
	if err != nil {
		tracing.SpanError(span, err)
	} else {
		span.LogFields(
			log.Uint64("azimuth", status.AbsoluteHigh.Azimuth),
			log.Int64("elevation", status.AbsoluteHigh.Elevation),
			log.Uint64("zoom", status.AbsoluteHigh.AbsoluteZoom),
		)
	}

	return status, err
}

func (a *ISAPI) PTZPresetGoTo(ctx context.Context, channel, preset uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.preset_goto")
	defer span.Finish()

	span.LogFields(
		log.Uint64("preset", preset),
	)

	u := a.address + proxyPTZPrefixURL + "/channels/" +
		strconv.FormatUint(channel, 10) +
		"/presets/" +
		strconv.FormatUint(preset, 10) +
		"/goto"

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, nil, &result)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (a *ISAPI) PTZRelative(ctx context.Context, channel uint64, x, y, zoom int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.relative")
	defer span.Finish()

	if zoom < -100 {
		zoom = -100
	} else if zoom > 100 {
		zoom = 100
	}

	span.LogFields(
		log.Int64("x", x),
		log.Int64("y", y),
		log.Int64("zoom", zoom),
	)

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
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (a *ISAPI) PTZAbsolute(ctx context.Context, channel uint64, elevation int64, azimuth, absoluteZoom uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.absolute")
	defer span.Finish()

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

	span.LogFields(
		log.Int64("elevation", elevation),
		log.Uint64("azimuth", azimuth),
		log.Uint64("absolute_zoom", absoluteZoom),
	)

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
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (a *ISAPI) PTZContinuous(ctx context.Context, channel uint64, pan, tilt, zoom int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.continuous")
	defer span.Finish()

	if pan < -100 {
		pan = -100
	} else if pan > 100 {
		pan = 100
	}

	if tilt < -100 {
		tilt = -100
	} else if tilt > 100 {
		tilt = 100
	}

	if zoom < -100 {
		zoom = -100
	} else if zoom > 100 {
		zoom = 100
	}

	span.LogFields(
		log.Int64("pan", pan),
		log.Int64("tilt", tilt),
		log.Int64("zoom", zoom),
	)

	u := a.address + proxyPTZPrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/continuous"

	data := PTZDataContinuous{
		Pan:  pan,
		Tilt: tilt,
		Zoom: zoom,
	}

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, data, &result)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}

func (a *ISAPI) PTZMomentary(ctx context.Context, channel uint64, pan, tilt, zoom int64, duration time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".ptz.momentary")
	defer span.Finish()

	if pan < -100 {
		pan = -100
	} else if pan > 100 {
		pan = 100
	}

	if tilt < -100 {
		tilt = -100
	} else if tilt > 100 {
		tilt = 100
	}

	if zoom < -100 {
		zoom = -100
	} else if zoom > 100 {
		zoom = 100
	}

	if duration < 0 {
		duration = 0
	}

	span.LogFields(
		log.Int64("pan", pan),
		log.Int64("tilt", tilt),
		log.Int64("zoom", zoom),
	)

	u := a.address + proxyPTZPrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/momentary"

	data := PTZDataContinuous{
		Pan:  pan,
		Tilt: tilt,
		Zoom: zoom,
	}
	data.Momentary.Duration = int64(duration.Seconds() * 1000)

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, data, &result)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}
