package telegram

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"gopkg.in/telegram-bot-api.v4"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
	client *tgbotapi.BotAPI
	done   chan struct{}
}

func (b *Bind) Run() (err error) {
	b.client, err = tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return err
	}

	b.SetSerialNumber(strconv.Itoa(b.client.Self.ID))
	b.client.Debug = b.config.Debug

	if b.config.UpdatesEnabled {
		b.client.Buffer = b.config.UpdatesBuffer

		u := tgbotapi.NewUpdate(0)
		u.Timeout = b.config.UpdatesTimeout

		updates, err := b.client.GetUpdatesChan(u)
		if err != nil {
			return err
		}

		b.listenUpdates(updates)
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	return nil
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

func (b *Bind) SendPhoto(to, name string, file io.Reader, size int64) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhotoUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   size,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) SendAudio(to, name string, file io.Reader, size int64) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewAudioUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   size,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) SendDocument(to, name string, file io.Reader, size int64) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocumentUpload(chatId, tgbotapi.FileReader{
		Name:   name,
		Reader: file,
		Size:   size,
	})
	msg.Caption = name

	_, err = b.client.Send(msg)
	return err
}

func (b *Bind) SendFileAsURL(to, name, u string) error {
	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	if _, err = url.Parse(u); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s [view](%s)", name, u))
	msg.ParseMode = "Markdown"

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

				ctx := context.Background()

				if u.Message.Text != "" {
					b.MQTTPublishAsync(ctx, MQTTPublishTopicReceiveMessage.Format(sn, u.Message.Chat.ID), u.Message.Text)
				}

				var (
					fileID    string
					mqttTopic string
				)

				if u.Message.Voice != nil {
					fileID = u.Message.Voice.FileID
					mqttTopic = MQTTPublishTopicReceiveVoice.Format(sn, u.Message.Chat.ID)
				} else if u.Message.Audio != nil {
					fileID = u.Message.Audio.FileID
					mqttTopic = MQTTPublishTopicReceiveAudio.Format(sn, u.Message.Chat.ID)
				}

				if fileID == "" {
					continue
				}

				link, err := b.client.GetFileDirectURL(fileID)
				if err != nil {
					b.Logger().Error("Get file by direct url failed",
						"error", err.Error(),
						"file", fileID,
					)
					continue
				}

				// TODO: log
				_ = b.MQTTPublishAsync(ctx, mqttTopic, link)

			case <-b.done:
				return
			}
		}
	}()
}

func (b *Bind) Close() (err error) {
	close(b.done)
	b.client.StopReceivingUpdates()
	return nil
}
