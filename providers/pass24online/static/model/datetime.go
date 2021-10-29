package swagger

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

var replacerDateTime = strings.NewReplacer(
	`"`, "",
	`\u0432`, "",
	`\u044f\u043d\u0432\u0430\u0440\u044f`, "January",
	`\u0444\u0435\u0432\u0440\u0430\u043b\u044f`, "February",
	`\u043c\u0430\u0440\u0442\u0430`, "March",
	`\u0430\u043f\u0440\u0435\u043b\u044f`, "April",
	`\u043c\u0430\u044f`, "May",
	`\u0438\u044e\u043d\u044f`, "June",
	`\u0438\u044e\u043b\u044f`, "July",
	`\u0430\u0432\u0433\u0443\u0441\u0442\u0430`, "August",
	`\u0441\u0435\u043d\u0442\u044f\u0431\u0440\u044f`, "September",
	`\u043e\u043a\u0442\u044f\u0431\u0440\u044f`, "October",
	`\u043d\u043e\u044f\u0431\u0440\u044f`, "November",
	`\u0434\u0435\u043a\u0430\u0431\u0440\u044f`, "December",
)

var timeLocation, _ = time.LoadLocation("Europe/Moscow")

type DateTime struct {
	strfmt.DateTime
}

func (m *DateTime) Time() time.Time {
	return time.Time(m.DateTime)
}

func (m *DateTime) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *DateTime) UnmarshalJSON(b []byte) error {
	t, err := time.ParseInLocation("02 January 2006  15:04", replacerDateTime.Replace(string(b)), timeLocation)
	if err != nil {
		return err
	}

	m.DateTime = strfmt.DateTime(t)

	return nil
}
