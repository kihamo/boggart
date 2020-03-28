package telegram

import (
	"context"
	"errors"
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
	di.ProbesBind

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

	chatID, err := b.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)

	return err
}

func (b *Bind) SendPhoto(to, name string, file io.Reader, size int64) error {
	bot := b.bot()
	if bot == nil {
		return errors.New("bot is offline")
	}

	chatID, err := b.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileReader{
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

	chatID, err := b.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewAudioUpload(chatID, tgbotapi.FileReader{
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

	chatID, err := b.chatID(to)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocumentUpload(chatID, tgbotapi.FileReader{
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

	chatID, err := b.chatID(to)
	if err != nil {
		return err
	}

	if _, err = url.Parse(u); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, name+" [view]("+u+")")
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

func (b *Bind) chatID(to string) (int64, error) {
	chatID, err := strconv.Atoi(to)
	if err != nil {
		return -1, err
	}

	return int64(chatID), err
}

func (b *Bind) listenUpdates(ch tgbotapi.UpdatesChannel) {
	go func() {
		sn := b.Meta().SerialNumber()

		for {
			select {
			case u := <-ch:
				if u.Message == nil {
					continue
				}

				// TODO: фильтрация по чату

				// TODO: фильтрация по юзеру

				ctx := context.Background()

				var mqttTopic mqtt.Topic

				if u.Message.Text != "" {
					mqttTopic = b.config.TopicReceiveMessage.Format(sn, u.Message.Chat.ID)

					if err := b.MQTT().PublishAsyncWithoutCache(ctx, b.config.TopicReceiveMessage.Format(sn, u.Message.Chat.ID), u.Message.Text); err != nil {
						b.Logger().Error("Publish message to MQTT failed",
							"topic", mqttTopic,
							"message", u.Message.Text,
							"error", err.Error(),
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

					return
				}

				link, err := bot.GetFileDirectURL(fileID)
				if err != nil {
					b.Logger().Error("Get file by direct url failed",
						"error", err.Error(),
						"file", fileID,
						"error", err.Error(),
					)

					continue
				}

				if err := b.MQTT().PublishAsyncWithoutCache(ctx, mqttTopic, link); err != nil {
					b.Logger().Error("Publish link to MQTT failed",
						"topic", mqttTopic,
						"link", link,
						"error", err.Error(),
					)
				}

			case <-b.done:
				return
			}
		}
	}()
}

func (b *Bind) Close() (err error) {
	close(b.done)

	bot := b.bot()
	if bot != nil {
		bot.StopReceivingUpdates()
	}

	return nil
}
