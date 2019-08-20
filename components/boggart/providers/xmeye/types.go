package xmeye

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

var (
	timeLocation = time.Now().Location()
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" || s == "2000-00-00 00:00:00" || s == "0000-00-00 00:00:00" {
		t.Time = time.Time{}
	} else {
		// FIXME: так как регистратор не отдает признак часового пояса
		// то все даты приводим к текущей часовой зоне процесса
		// что бы даты отображались корректно. В случае если на регистраторе
		// другой часовой пояс это может стать проблемой
		t.Time, err = time.ParseInLocation(timeLayout, s, timeLocation)
	}

	return err
}

type Uint32 uint32

func (t *Uint32) UnmarshalJSON(b []byte) (err error) {
	src := bytes.Trim(b, "\"")

	if len(src) < 8 {
		return
	}

	src = src[len(src)-8:]
	dst := make([]byte, 4)

	if _, err = hex.Decode(dst, src); err != nil {
		return err
	}

	*t = Uint32(binary.LittleEndian.Uint32([]byte{dst[3], dst[2], dst[1], dst[0]}))

	return err
}

func (t *Uint32) MarshalJSON() ([]byte, error) {
	src := make([]byte, 4)
	binary.LittleEndian.PutUint32(src, uint32(*t))

	dst := "0x" + strings.ToUpper(hex.EncodeToString([]byte{src[3], src[2], src[1], src[0]}))

	return []byte(dst), nil
}

func (t Uint32) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func (t Uint32) GoString() string {
	return t.String()
}
