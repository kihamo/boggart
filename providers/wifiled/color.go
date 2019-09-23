package wifiled

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	reColor1 = regexp.MustCompile("^#?([0-9a-f]{2}[0-9a-f]{2}[0-9a-f]{2})$")
	reColor2 = regexp.MustCompile("^#?([0-9a-f])([0-9a-f])([0-9a-f])$")
)

type Color struct {
	Red          uint8
	Green        uint8
	Blue         uint8
	WarmWhite    uint8
	UseRGB       bool
	UseWarmWhite bool
}

func (c Color) Bytes() []byte {
	if c.UseWarmWhite {
		return []byte{c.WarmWhite}
	}

	return []byte{c.Red, c.Green, c.Blue}
}

func (c Color) Uint64() (value uint64) {
	for i, b := range c.Bytes() {
		value = value<<8 + uint64(b)
		if i == 7 {
			return
		}
	}

	return
}

func (c Color) HSV() (h uint32, s, v float64) {
	r := float64(c.Red) / 255
	g := float64(c.Green) / 255
	b := float64(c.Blue) / 255

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	if max == min {
		h = 0
	} else if max == r {
		if g >= b {
			h = uint32(60 * ((g - b) / (max - min)))
		} else {
			h = uint32(60*((g-b)/(max-min)) + 360)
		}
	} else if max == g {
		h = uint32(60*((b-r)/(max-min)) + 120)
	} else if max == b {
		h = uint32(60*((r-g)/(max-min)) + 240)
	}

	h %= 360

	if max == 0 {
		s = 0
	} else {
		s = (1 - min/max) * 100
	}

	v = max * 100
	return
}

func (c Color) String() string {
	if c.UseWarmWhite {
		return strconv.FormatUint(uint64(c.WarmWhite), 10)
	}

	return fmt.Sprintf("#%x", c.Bytes())
}

func ColorFromHEX(color string) (*Color, error) {
	rgb, err := hex.DecodeString(color)
	if err != nil {
		return nil, err
	}

	if len(rgb) < 3 {
		return nil, errors.New("wrong color format")
	}

	return &Color{
		Red:    rgb[0],
		Green:  rgb[1],
		Blue:   rgb[2],
		UseRGB: true,
	}, nil
}

func ColorFromString(color string) (*Color, error) {
	color = strings.ToLower(color)

	// HEX: #aabbcc
	if matches := reColor1.FindStringSubmatch(color); len(matches) > 0 {
		return ColorFromHEX(matches[1])
	}

	// HEX: #abc
	if matches := reColor2.FindStringSubmatch(color); len(matches) > 0 {
		return ColorFromHEX(strings.Repeat(matches[1], 2) + strings.Repeat(matches[2], 2) + strings.Repeat(matches[3], 2))
	}

	// HSB: 358,98,100
	parts := strings.Split(color, ",")
	if len(parts) == 3 {
		// 0...360
		hue, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, err
		}

		// 0...100
		saturation, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, err
		}

		if saturation > 1 {
			saturation /= 100
		}

		// 0...100
		value, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, err
		}

		if value > 1 {
			value /= 100
		}

		var r, g, b float64
		c := value * saturation
		x := c * (1 - math.Abs(math.Mod(hue/60, 2)-1))
		m := value - c

		switch {
		case hue < 60:
			r, g, b = c, x, 0
		case hue < 120:
			r, g, b = x, c, 0
		case hue < 180:
			r, g, b = 0, c, x
		case hue < 240:
			r, g, b = 0, x, c
		case hue < 300:
			r, g, b = x, 0, c
		default:
			r, g, b = c, 0, x
		}

		return &Color{
			Red:    uint8(math.Floor((r + m) * 255)),
			Green:  uint8(math.Floor((g + m) * 255)),
			Blue:   uint8(math.Floor((b + m) * 255)),
			UseRGB: true,
		}, nil
	}

	return nil, errors.New("wrong color format")
}
