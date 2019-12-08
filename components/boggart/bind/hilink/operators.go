package hilink

import (
	"regexp"
)

type operator struct {
	BalanceUSSD           string
	BalanceRegexp         *regexp.Regexp
	SMSLimitTrafficRegexp *regexp.Regexp
	SMSLimitTrafficFactor float64
}

/*
На сегодня 11:44 остатки пакетов: безлимитные звонки на Tele2 РФ; 5120 МБ, неиспользованные остатки с прошлого периода: 1902 МБ. Срок действия пакетов до 21.11.19 23:59. Подробнее на tele2.ru/my или *107#
*/

var (
	operatorTele2 = &operator{
		BalanceUSSD:           "*105#",
		BalanceRegexp:         regexp.MustCompile(`OCTATOK (?P<value>\d+\.\d{2})\sp\..*?`),
		SMSLimitTrafficRegexp: regexp.MustCompile(`остатки пакетов:.*?(?P<value1>\d+)\s*МБ,\s*неиспользованные остатки с прошлого периода:.*?(?P<value2>\d+)\s*МБ*?`),
		SMSLimitTrafficFactor: 1024 * 1024,
	}
)
