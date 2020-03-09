package tvt

const (
	ErrorWrongUsername int64 = 536870947
	ErrorWrongPassword int64 = 53687094
)

var errorText = map[int64]string{
	ErrorWrongUsername: "wrong username",
	ErrorWrongPassword: "wrong password",
}

func ErrorMessage(code int64) string {
	if text, ok := errorText[code]; ok {
		return text
	}

	return ""
}
