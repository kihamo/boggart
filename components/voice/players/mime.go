package players

import (
	"io"
	"mime"
	"net/http"
	"strings"
)

type MIMEType string

var (
	defaultHTTPClient = &http.Client{
		Transport: http.DefaultTransport,
	}
)

const (
	MIMETypeUnknown = MIMEType("")
	MIMETypeMPEG    = MIMEType("audio/mpeg")
	MIMETypeWAVE    = MIMEType("audio/vnd.wave")
	MIMETypeOGG     = MIMEType("audio/ogg")
)

func (m MIMEType) String() string {
	return string(m)
}

func MimeTypeFromHTTPHeader(header http.Header) (MIMEType, error) {
	contentType := header.Get("Content-type")
	if contentType == "" {
		return MIMETypeUnknown, ErrorUnknownAudioFormat
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			return MIMETypeUnknown, err
		}

		return MIMEType(t), nil
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

func MimeTypeFromURL(url string) (MIMEType, error) {
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return MIMETypeUnknown, err
	}

	response, err := defaultHTTPClient.Do(request)
	if err != nil {
		return MIMETypeUnknown, err
	}

	// GET fallback
	if response.StatusCode == http.StatusBadRequest {
		request, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return MIMETypeUnknown, err
		}

		response, err = defaultHTTPClient.Do(request)
		if err != nil {
			return MIMETypeUnknown, err
		}
	}

	mimeType, err := MimeTypeFromHTTPHeader(response.Header)
	if err != nil {
		return MIMETypeUnknown, err
	}

	if mimeType == MIMETypeUnknown {
		mimeType, err = MimeTypeFromData(response.Body)
	}

	return mimeType, err
}
