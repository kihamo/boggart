package hilink

import (
	"regexp"
)

type operator struct {
	BalanceUSSD           string
	BalanceRegexp         *regexp.Regexp
	SMSLimitTrafficRegexp *regexp.Regexp
}

var (
	operatorTele2 = &operator{
		BalanceUSSD:           "*105#",
		BalanceRegexp:         regexp.MustCompile(`OCTATOK (?P<value>\d+\.\d{2})\sp\..*?`),
		SMSLimitTrafficRegexp: regexp.MustCompile(`остатки пакетов:.*?(?P<value>\d+) МБ Интернет.*?`),
	}
)
