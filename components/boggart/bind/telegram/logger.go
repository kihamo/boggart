package telegram

import (
	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	cutSize = 1024
)

type logger struct {
	tgbotapi.BotLogger

	f1 func(string)
	f2 func(string)
}

func NewLogger(printf func(string), println func(string)) tgbotapi.BotLogger {
	return &logger{
		f1: printf,
		f2: println,
	}
}

func (l logger) Printf(format string, args ...interface{}) {
	record := fmt.Sprintf(format, args...)

	cut := len(record)
	if cut > cutSize {
		cut = cutSize
	}

	l.f1(record[:cut])
}

func (l logger) Println(args ...interface{}) {
	record := fmt.Sprintln(args...)

	cut := len(record)
	if cut > cutSize {
		cut = cutSize
	}

	l.f2(record[:cut])
}
