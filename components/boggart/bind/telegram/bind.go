package telegram

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/mime"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	paramFileName = "file"
	paramMIME     = "mime"
	paramRandom   = "_t"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	mutex  sync.RWMutex
	client *tgbotapi.BotAPI
	done   chan struct{}

	fileServer http.Handler
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	tgbotapi.SetLogger(NewLogger(
		func(message string) {
			b.Logger().Debug(message)
		},
		func(message string) {
			b.Logger().Warn(message)
		}))

	b.client = nil
	b.done = make(chan struct{})

	cfg := b.config()

	if cfg.FileURLPrefix.String() == "" {
		if _, err := b.Widget().URL(nil); err != nil {
			cfg.UseURLForSendFile = false
		}
	}

	b.fileServer = http.FileServer(http.Dir(cfg.FileDirectory))

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

func (b *Bind) SendPhoto(to, name string, file io.Reader) error {
	if b.config().UseURLForSendFile {
		u, err := b.SaveFile(file)
		if err != nil {
			return err
		}

		return b.SendFileAsURL(to, name, u)
	}

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
		Size:   -1,
	})
	msg.Caption = name

	_, err = bot.Send(msg)

	return err
}

func (b *Bind) SendAudio(to, name string, file io.Reader) error {
	if b.config().UseURLForSendFile {
		u, err := b.SaveFile(file)
		if err != nil {
			return err
		}

		return b.SendFileAsURL(to, name, u)
	}

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
		Size:   -1,
	})
	msg.Caption = name

	_, err = bot.Send(msg)

	return err
}

func (b *Bind) SendDocument(to, name string, file io.Reader) error {
	if b.config().UseURLForSendFile {
		u, err := b.SaveFile(file)
		if err != nil {
			return err
		}

		return b.SendFileAsURL(to, name, u)
	}

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
		Size:   -1,
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

func (b *Bind) SaveFile(reader io.Reader) (string, error) {
	hash := md5.New()

	file := bytes.NewBuffer(nil)
	defer file.Reset()

	multi := io.MultiWriter(hash, file)

	if _, err := io.Copy(multi, reader); err != nil {
		return "", err
	}

	id := hex.EncodeToString(hash.Sum(nil))
	filePath := filepath.Join(b.config().FileDirectory, id)

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	mimeType, restored, err := mime.TypeFromDataRestored(file)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, restored)
	if err != nil {
		return "", err
	}

	u, err := b.Widget().URL(map[string]string{
		paramFileName: id,
		paramMIME:     mimeType.String(),
		// телега кэширует урлы и второй раз не прийдет и механизм очистки не сработает,
		// поэтому добавляем рандом
		paramRandom: strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (b *Bind) RemoveFile(id string) error {
	filePath := filepath.Join(b.config().FileDirectory, id)

	_, err := os.Stat(filePath)
	if err == nil {
		return os.Remove(filePath)
	}

	return err
}

func (b *Bind) initBot() (*tgbotapi.BotAPI, error) {
	cfg := b.config()

	client, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	b.Meta().SetSerialNumber(strconv.Itoa(client.Self.ID))

	if client.Self.UserName != "" {
		if link, err := url.Parse("https://t.me/" + client.Self.UserName); err == nil {
			b.Meta().SetLink(link)
		}
	} else {
		b.Meta().SetLink(nil)
	}

	client.Debug = cfg.Debug

	if cfg.UpdatesEnabled {
		client.Buffer = cfg.UpdatesBuffer

		u := tgbotapi.NewUpdate(0)
		u.Timeout = cfg.UpdatesTimeout

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
		cfg := b.config()

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
					mqttTopic = cfg.TopicReceiveMessage.Format(sn, u.Message.Chat.ID)

					if err := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicReceiveMessage.Format(sn, u.Message.Chat.ID), u.Message.Text); err != nil {
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
					mqttTopic = cfg.TopicReceiveVoice.Format(sn, u.Message.Chat.ID)
				} else if u.Message.Audio != nil {
					fileID = u.Message.Audio.FileID
					mqttTopic = cfg.TopicReceiveAudio.Format(sn, u.Message.Chat.ID)
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
