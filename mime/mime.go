package mime

import (
	"bytes"
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
	TypeUnknown     = Type("")
	TypeJPEG        = Type("image/jpeg")
	TypeJPG         = Type("image/jpg")
	TypeGIF         = Type("image/gif")
	TypePNG         = Type("image/png")
	TypeMPEG        = Type("audio/mpeg")
	TypeWAVE        = Type("audio/vnd.wave")
	TypeFLAC        = Type("audio/flac")
	TypeOGG         = Type("application/ogg")
	TypeOctetStream = Type("application/octet-stream")
)

var (
	ErrorUnknownMIMEType = errors.New("unknown mime type")
)

type Type string

func (t Type) String() string {
	return string(t)
}

func (t Type) Extension() string {
	switch t {
	case TypeJPEG, TypeJPG:
		return "jpg"
	case TypePNG:
		return "png"
	case TypeGIF:
		return "gif"
	case TypeMPEG:
		return "mp3"
	case TypeWAVE:
		return "wav"
	case TypeFLAC:
		return "flac"
	case TypeOGG:
		return "ogg"
	}

	return ""
}

func (t Type) IsImage() bool {
	return t == TypeJPEG || t == TypeJPG || t == TypePNG || t == TypeGIF
}

func (t Type) IsAudio() bool {
	return t == TypeMPEG || t == TypeWAVE || t == TypeFLAC || t == TypeOGG
}

func TypeFromHTTPHeader(header http.Header) (Type, error) {
	contentType := header.Get("Content-type")
	if contentType == "" {
		return TypeUnknown, ErrorUnknownMIMEType
	}

	for _, v := range strings.Split(contentType, ",") {
		if t, _, err := mime.ParseMediaType(v); err == nil {
			return Type(t), nil
		}
	}

	return TypeUnknown, nil
}

func TypeFromDataRestored(data io.Reader) (Type, io.Reader, error) {
	buf := make([]byte, 128)

	if _, err := data.Read(buf); err != nil {
		return TypeUnknown, nil, err
	}

	t, err := TypeFromData(bytes.NewBuffer(buf))
	if err != nil {
		return TypeUnknown, nil, err
	}

	restored := bytes.NewBuffer(buf)
	if _, err := io.Copy(restored, data); err != nil {
		return TypeMPEG, nil, err
	}

	return t, restored, nil
}

func TypeFromData(data io.Reader) (Type, error) {
	buf := make([]byte, 128)

	if _, err := data.Read(buf); err != nil {
		return TypeUnknown, err
	}

	// Hack for MP3 file
	// https://en.wikipedia.org/wiki/List_of_file_signatures
	if buf[0] == 0xFF && buf[1] == 0xFB {
		return TypeMPEG, nil
	}

	t := http.DetectContentType(buf)

	return Type(t), nil
}

func TypeFromURL(url string) (mimeType Type, err error) {
	// попытка вычитать HEAD
	request, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return TypeUnknown, err
	}

	response, err := defaultHTTPClient.Do(request)
	if err != nil {
		return TypeUnknown, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		mimeType, err = TypeFromHTTPHeader(response.Header)
		if err == nil && mimeType != TypeOctetStream {
			return mimeType, nil
		}
	}

	// если не удалось, то делает GET
	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return TypeUnknown, err
	}

	response, err = defaultHTTPClient.Do(request)
	if err != nil {
		return TypeUnknown, err
	}
	defer response.Body.Close()

	mimeType, err = TypeFromHTTPHeader(response.Header)
	if err != nil {
		return TypeUnknown, err
	}

	if mimeType == TypeUnknown || mimeType == TypeOctetStream {
		mimeType, err = TypeFromData(response.Body)
	}

	return mimeType, err
}
