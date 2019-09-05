package internal

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
	if s == "null" || s == "" || s == "2000-00-00 00:00:00" || s == "0000-00-00 00:00:00" {
		t.Time = time.Time{}
	} else {
		// FIXME: так как регистратор не отдает признак часового пояса
		// то все даты приводим к текущей часовой зоне процесса
		// что бы даты отображались корректно. В случае если на регистраторе
		// другой часовой пояс это может стать проблемой
		t.Time, err = time.ParseInLocation("2006-01-02 15:04:05", s, timeLocation)
	}

	return err
}

type Date struct {
	time.Time
}

func (t *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" || s == "2000-00-00 00:00:00" || s == "0000-00-00 00:00:00" {
		t.Time = time.Time{}
	} else {
		// FIXME: так как регистратор не отдает признак часового пояса
		// то все даты приводим к текущей часовой зоне процесса
		// что бы даты отображались корректно. В случае если на регистраторе
		// другой часовой пояс это может стать проблемой
		t.Time, err = time.ParseInLocation("2006-01-02", s, timeLocation)
	}

	return err
}
