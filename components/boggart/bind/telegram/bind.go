package telegram

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"gopkg.in/telegram-bot-api.v4"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind

	config *Config

	mutex  sync.RWMutex
	client *tgbotapi.BotAPI
	done   chan struct{}
}

func (b *Bind) Run() error {
	b.client = nil
	b.done = make(chan struct{})

	return nil
}

func (b *Bind) SendMessage(to, message string) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, message)
	_, err = bot.Send(msg)

	return err
}

func (b *Bind) SendPhoto(to, name string, file io.Reader, size int64) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

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

	_, err = bot.Send(msg)
	return err
}

func (b *Bind) SendAudio(to, name string, file io.Reader, size int64) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

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

	_, err = bot.Send(msg)
	return err
}

func (b *Bind) SendDocument(to, name string, file io.Reader, size int64) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

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

	_, err = bot.Send(msg)
	return err
}

func (b *Bind) SendFileAsURL(to, name, u string) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

	chatId, err := b.chatId(to)
	if err != nil {
		return err
	}

	if _, err = url.Parse(u); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s [view](%s)", name, u))
	msg.ParseMode = "Markdown"

	_, err = bot.Send(msg)

	return err
}

func (b *Bind) initBot() (*tgbotapi.BotAPI, error) {
	client, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return nil, err
	}

	b.Meta().SetSerialNumber(strconv.Itoa(client.Self.ID))
	client.Debug = b.config.Debug

	if b.config.UpdatesEnabled {
		client.Buffer = b.config.UpdatesBuffer

		u := tgbotapi.NewUpdate(0)
		u.Timeout = b.config.UpdatesTimeout

		updates, err := client.GetUpdatesChan(u)
		if err != nil {
			return nil, err
		}

		b.listenUpdates(updates)
	}

	b.mutex.Lock()
	b.client = client
	b.mutex.Unlock()

	return client, nil
}

func (b *Bind) bot() *tgbotapi.BotAPI {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.client
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
		sn := b.Meta().SerialNumber()

		for {
			select {
			case u := <-ch:
				// TODO: фильтрация по чату

				// TODO: фильтрация по юзеру

				if u.Message == nil {
					continue
				}

				ctx := context.Background()
				var mqttTopic mqtt.Topic

				if u.Message.Text != "" {
					mqttTopic = b.config.TopicReceiveMessage.Format(sn, u.Message.Chat.ID)

					if err := b.MQTT().PublishAsync(ctx, b.config.TopicReceiveMessage.Format(sn, u.Message.Chat.ID), u.Message.Text); err != nil {
						b.Logger().Error("Publish message to MQTT failed",
							"topic", mqttTopic,
							"message", u.Message.Text,
						)
					} else {
						b.Logger().Debug("Receive message",
							"chat", u.Message.Chat.ID,
							"message", u.Message.Text,
						)
					}
				}

				var fileID string

				if u.Message.Voice != nil {
					fileID = u.Message.Voice.FileID
					mqttTopic = b.config.TopicReceiveVoice.Format(sn, u.Message.Chat.ID)
				} else if u.Message.Audio != nil {
					fileID = u.Message.Audio.FileID
					mqttTopic = b.config.TopicReceiveAudio.Format(sn, u.Message.Chat.ID)
				}

				if fileID == "" {
					continue
				}

				bot := b.bot()
				if bot == nil {
					b.Logger().Error("Bot if offline",
						"file", fileID,
					)
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

				if err := b.MQTT().PublishAsync(ctx, mqttTopic, link); err != nil {
					b.Logger().Error("Publish link to MQTT failed",
						"topic", mqttTopic,
						"link", link,
					)
				}

			case <-b.done:
				return
			}
		}
	}()
}

func (b *Bind) Close() (err error) {
	bot := b.bot()
	if bot != nil {
		bot.StopReceivingUpdates()
	}

	close(b.done)

	return nil
}
