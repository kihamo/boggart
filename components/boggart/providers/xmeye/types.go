package xmeye

import (
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
	if s == "null" {
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
