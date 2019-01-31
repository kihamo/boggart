package telegram

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"gopkg.in/telegram-bot-api.v4"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *tgbotapi.BotAPI
	done   chan struct{}
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

func (b *Bind) listenUpdates(ch tgbotapi.UpdatesChannel) {
	go func() {
		sn := mqtt.NameReplace(b.SerialNumber())

		for {
			select {
			case u := <-ch:
				// TODO: фильтрация по чату

				// TODO: фильтрация по юзеру

				if u.Message == nil {
					continue
				}

				fmt.Println(u.Message.Text)

				var (
					fileID    string
					mqttTopic string
				)

				if u.Message.Voice != nil {
					fileID = u.Message.Voice.FileID
					mqttTopic = MQTTPublishTopicFileVoice.Format(sn, u.Message.Chat.ID)
				} else if u.Message.Audio != nil {
					fileID = u.Message.Audio.FileID
					mqttTopic = MQTTPublishTopicFileAudio.Format(sn, u.Message.Chat.ID)
				}

				if fileID == "" {
					continue
				}

				link, err := b.client.GetFileDirectURL(fileID)
				if err != nil {
					// TODO: log
					continue
				}

				b.MQTTPublishAsync(context.Background(), mqttTopic, 1, false, link)

			case <-b.done:
				return
			}
		}
	}()
}

func (b *Bind) Close() (err error) {
	close(b.done)
	return nil
}
