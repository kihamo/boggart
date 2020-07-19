package storage

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"
)

var (
	defaultHTTPClient = &http.Client{
		Transport: http.DefaultTransport,
	}
)

const (
	MIMETypeUnknown     = MIMEType("")
	MIMETypeJPEG        = MIMEType("image/jpeg")
	MIMETypeJPG         = MIMEType("image/jpg")
	MIMETypeGIF         = MIMEType("image/gif")
	MIMETypePNG         = MIMEType("image/png")
	MIMETypeMPEG        = MIMEType("audio/mpeg")
	MIMETypeWAVE        = MIMEType("audio/vnd.wave")
	MIMETypeFLAC        = MIMEType("audio/flac")
	MIMETypeOGG         = MIMEType("application/ogg")
	MIMETypeOctetStream = MIMEType("application/octet-stream")
)

var (
	ErrorUnknownMIMEType = errors.New("unknown mime type")
)

type MIMEType string

func (m MIMEType) String() string {
	return string(m)
}

func MimeTypeFromHTTPHeader(header http.Header) (MIMEType, error) {
	contentType := header.Get("Content-type")
	if contentType == "" {
		return MIMETypeUnknown, ErrorUnknownMIMEType
	}

	for _, v := range strings.Split(contentType, ",") {
		if t, _, err := mime.ParseMediaType(v); err == nil {
			return MIMEType(t), nil
		}
	}

	return MIMETypeUnknown, nil
}

func MimeTypeFromData(data io.Reader) (MIMEType, error) {
	buf := make([]byte, 128)

	if _, err := data.Read(buf); err != nil {
		return MIMETypeUnknown, err
	}

	// Hack for MP3 file
	// https://en.wikipedia.org/wiki/List_of_file_signatures
	if buf[0] == 0xFF && buf[1] == 0xFB {
		return MIMETypeMPEG, nil
	}

	t := http.DetectContentType(buf)

	return MIMEType(t), nil
}

func MimeTypeFromURL(url string) (mimeType MIMEType, err error) {
	// попытка вычитать HEAD
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return MIMETypeUnknown, err
	}

	response, err := defaultHTTPClient.Do(request)
	if err != nil {
		return MIMETypeUnknown, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		mimeType, err = MimeTypeFromHTTPHeader(response.Header)
		if err == nil && mimeType != MIMETypeOctetStream {
			return mimeType, nil
		}
	}

	// если не удалось, то делает GET
	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MIMETypeUnknown, err
	}

	response, err = defaultHTTPClient.Do(request)
	if err != nil {
		return MIMETypeUnknown, err
	}
	defer response.Body.Close()

	mimeType, err = MimeTypeFromHTTPHeader(response.Header)
	if err != nil {
		return MIMETypeUnknown, err
	}

	if mimeType == MIMETypeUnknown || mimeType == MIMETypeOctetStream {
		mimeType, err = MimeTypeFromData(response.Body)
	}

	return mimeType, err
}
