package xmeye

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/performance"
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

	val, err := strconv.ParseUint(string(src), 0, 64)
	if err != nil {
		return err
	}

	*t = Uint32(val)

	return err
}

func (t *Uint32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("0x%08X", uint64(*t))), nil
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

var replacerAsStringCleanValue = strings.NewReplacer(
	`"`, "",
	`0x`, "",
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	val, err := strconv.ParseUint(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)), 16, 64)
	if err != nil {
		return err
	}

	time.Second.String()

	d.Duration = time.Duration(val) * time.Minute

	return err
}
