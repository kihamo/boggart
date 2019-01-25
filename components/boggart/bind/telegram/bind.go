package telegram

import (
	"io"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"gopkg.in/telegram-bot-api.v4"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *tgbotapi.BotAPI
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

func (b *Bind) SendMessage(to, message string) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, message)
	_, err = b.client.Send(msg)

	return err
}

func (b *Bind) SendPhoto(to, name string, file io.Reader) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhotoUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   -1,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) SendAudio(to, name string, file io.Reader) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewAudioUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   -1,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) SendDocument(to, name string, file io.Reader) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   -1,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) chatId(to string) (int64, error) {
	chatId, err := strconv.Atoi(to)
	if err != nil {
		return -1, err
	}

	return int64(chatId), err
}
