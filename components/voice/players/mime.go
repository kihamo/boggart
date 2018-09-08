package players

import (
	"io"
	"mime"
	"net/http"
	"strings"
)

func mimeIsAllow(mime string) bool {
	return mime == "audio/mpeg" || mime == "audio/vnd.wave" || mime == "audio/ogg"
}

func mimeFromHeader(header http.Header) (string, error) {
	contentType := header.Get("Content-type")
	if contentType == "" {
		return "", ErrorUnknownAudioFormat
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			return "", err
		}

		if mimeIsAllow(t) {
			return t, nil
		}
	}

	return "", nil
}

func mimeFromData(data io.Reader) (string, error) {
	buf := make([]byte, 128)

	if _, err := data.Read(buf); err != nil {
		return "", err
	}

	// Hack for MP3 file
	// https://en.wikipedia.org/wiki/List_of_file_signatures
	if buf[0] == 0xFF && buf[1] == 0xFB {
		return "audio/mpeg", nil
	}

	t := http.DetectContentType(buf)
	if mimeIsAllow(t) {
		return t, nil
	}

	return "", nil
}
