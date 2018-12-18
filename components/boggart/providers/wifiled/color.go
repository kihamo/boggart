package wifiled

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
		Red:   rgb[0],
		Green: rgb[1],
		Blue:  rgb[2],
	}, nil
}

func ColorFromString(color string) (*Color, error) {
	if matched, _ := regexp.Match("#?[0-9a-f]{6}", []byte(color)); matched {
		if color[0] == '#' {
			color = color[1:]
		}

		return ColorFromHEX(color)
	}

	if matched, _ := regexp.Match("#?[0-9a-f]{3}", []byte(color)); matched {
		if color[0] == '#' {
			color = color[1:]
		}

		expandedName := fmt.Sprintf("%c%c%c%c%c%c", color[0], color[0], color[1], color[1], color[2], color[2])
		return ColorFromHEX(expandedName)
	}

	return nil, errors.New("wrong color format")
}
