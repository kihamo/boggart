package barcode

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/llgcode/draw2d/draw2dimg"
)

var (
	debugBorderColor         = color.RGBA{255, 0, 0, 255}
	debugBorderWidth float64 = 2
)

func init() {
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

func encode(format string, img image.Image) (_ io.Reader, err error) {
	result := bytes.NewBuffer(nil)

	switch format {
	case "gif":
		err = gif.Encode(result, img, nil)
	case "jpeg":
		err = jpeg.Encode(result, img, nil)
	case "png":
		err = png.Encode(result, img)
	}

	return result, err
}

func drawDebugBorder(img image.Image, points []image.Point) image.Image {
	if len(points) == 0 {
		return img
	}

	size := img.Bounds().Size()

	drawImage := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
	gc := draw2dimg.NewGraphicContext(drawImage)

	// copy
	gc.DrawImage(img)

	// draw debug
	var currentPoint, nextPoint image.Point

	gc.SetLineWidth(debugBorderWidth)
	gc.SetStrokeColor(debugBorderColor)

	i := 0
	for {
		currentPoint = points[i]
		if i < len(points)-1 {
			nextPoint = points[i+1]
		} else {
			nextPoint = points[0]
		}

		gc.MoveTo(float64(currentPoint.X), float64(currentPoint.Y))
		gc.LineTo(float64(nextPoint.X), float64(nextPoint.Y))
		gc.Stroke()

		i++
		if i == len(points) {
			break
		}
	}

	return drawImage
}
