package pulsar

import (
	"time"
)

func TimeToBytes(t time.Time) []byte {
	return []byte{
		byte(t.Year() - 2000),
		byte(t.Month()),
		byte(t.Day()),
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
	}
}

func BytesToTime(data []byte, location *time.Location) time.Time {
	return time.Date(
		2000+int(data[0]),
		time.Month(data[1]),
		int(data[2]),
		int(data[3]),
		int(data[4]),
		int(data[5]),
		0,
		location)
}
