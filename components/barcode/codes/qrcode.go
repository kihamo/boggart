package codes

import (
	"image"
	"io"

	"github.com/bieber/barcode"
)

var qrCodeScanner = barcode.NewScanner().
	SetEnabledAll(false).
	SetEnabledSymbology(barcode.QRCode, true)

func parseQRCode(reader io.Reader, debug bool) (code string, result io.Reader, err error) {
	img, format, err := image.Decode(reader)
	if err != nil {
		return "", nil, err
	}

	symbols, err := qrCodeScanner.ScanImage(barcode.NewImage(img))
	if err != nil {
		return "", nil, nil
	}

	for _, s := range symbols {
		code += s.Data

		if debug {
			img = drawDebugBorder(img, s.Boundary)
		}
	}

	if debug {
		result, err = encode(format, img)
	}

	return code, result, err
}

func DecodeQRCode(reader io.Reader) (code string, err error) {
	code, _, err = parseQRCode(reader, false)

	return code, err
}

func DecodeQRCodeDebug(reader io.Reader) (string, io.Reader, error) {
	return parseQRCode(reader, true)
}
