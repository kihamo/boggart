package swagger

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

var replacerDateTime = strings.NewReplacer(
	`"`, "",
	`\u0432`, "",
	`\u044f\u043d\u0432\u0430\u0440\u044f`, "01",
	`\u0444\u0435\u0432\u0440\u0430\u043b\u044f`, "02",
	`\u043c\u0430\u0440\u0442\u0430`, "03",
	`\u0430\u043f\u0440\u0435\u043b\u044f`, "04",
	`\u043c\u0430\u044f`, "05",
	`\u0438\u044e\u043d\u044f`, "06",
	`\u0438\u044e\u043b\u044f`, "07",
	`\u0430\u0432\u0433\u0443\u0441\u0442\u0430`, "08",
	`\u0441\u0435\u043d\u0442\u044f\u0431\u0440\u044f`, "09",
	`\u043e\u043a\u0442\u044f\u0431\u0440\u044f`, "10",
	`\u043d\u043e\u044f\u0431\u0440\u044f`, "11",
	`\u0434\u0435\u043a\u0430\u0431\u0440\u044f`, "12",
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
	t, err := time.ParseInLocation("02 01 2006  15:04", replacerDateTime.Replace(string(b)), timeLocation)
	if err != nil {
		return err
	}

	m.DateTime = strfmt.DateTime(t)

	return nil
}
