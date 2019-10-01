package xmeye

import (
	"bytes"
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
		t.Time, err = time.ParseInLocation(TimeLayout, s, timeLocation)
	}

	return err
}

type Uint32 uint32

func (t *Uint32) UnmarshalJSON(b []byte) (err error) {
	src := bytes.Trim(b, "\"")

	if len(src) < 8 {
		return
	}

	source := strings.Replace(string(src), "0x", "", 1)
	val, err := strconv.ParseUint(source, 16, 64)
	if err != nil {
		return err
	}

	*t = Uint32(val)

	return err
}

func (t *Uint32) MarshalJSON() ([]byte, error) {
	dst := "0x" + strconv.FormatUint(uint64(*t), 16)

	return []byte(dst), nil
}

func (t Uint32) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

func (t Uint32) GoString() string {
	return t.String()
}

func (t Uint32) Uint32() uint32 {
	return uint32(t)
}
