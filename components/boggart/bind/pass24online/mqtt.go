package pass24online

import (
	"encoding/json"
	"time"
)

var statusName = map[string]string{
	"0":  "Неизвестно",
	"30": "Гость вне территории",
	"40": "Гость на территории",
	"50": "Срок действия окончен",
	"60": "Пропуск закрыт",
}

// 30 -> 40 -- Гость на территории
// 40 -> 50 -- Срок действия окончен
// 40 -> 60 -- Пропуск закрыт
// 40 -> 30 -- Гость вне территории

type FeedEvent struct {
	ModelName   string    `json:"model_name"`
	PlateNumber string    `json:"plate_number"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Datetime    time.Time `json:"datetime"`
}

func (e FeedEvent) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}
