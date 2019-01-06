package chromecast

import (
	"context"
	"errors"
	"io"
	"math"
	"net"
	"sync"
	"sync/atomic"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	"github.com/kihamo/boggart/components/voice/players"
)

const (
	DefaultPort = 8009
)

type Player struct {
	mutex  sync.RWMutex
	client *cast.Client
	done   chan struct{}

	host net.IP
	port int

	connected int64
	status    int64
	volume    int64
	mute      int64
}

func New(ip string, port int64) *Player {
	p := &Player{
		host: net.ParseIP(ip),
		port: int(port),
		done: make(chan struct{}, 1),
	}
	p.setStatus(players.StatusStopped)
	p.SetVolume(50)

	return p
}

func (p *Player) PlayFromURL(url string) error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	// TODO: возвращает готовность к проигрыванию, а не реальное проигрывание
	//if p.cast.IsPlaying(ctx) {
	//	return players.ErrorAlreadyPlaying
	//}

	mimeType, err := players.MimeTypeFromURL(url)
	if err != nil {
		return err
	}

	if mimeType == players.MIMETypeUnknown {
		return players.ErrorUnknownAudioFormat
	}

	media, err := client.Media(ctx)
	if err != nil {
		return err
	}

	item := controllers.MediaItem{
		ContentId:   url,
		StreamType:  "BUFFERED",
		ContentType: mimeType.String(),
	}

	_, err = media.LoadMedia(ctx, item, 0, true, map[string]interface{}{})
	return err
}

func (p *Player) PlayFromReader(_ io.ReadCloser) error {
	return errors.New("not support play from reader")
}

func (p *Player) PlayFromFile(_ string) error {
	return errors.New("not support play from reader")
}

func (p *Player) Play() error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	//if p.cast.IsPlaying(ctx) {
	//	return players.ErrorAlreadyPlaying
	//}

	media, err := client.Media(ctx)
	if err != nil {
		return err
	}

	_, err = media.Play(ctx)
	return err
}

func (p *Player) Stop() error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	//if !client.IsPlaying(ctx) {
	//	return nil
	//}

	media, err := client.Media(ctx)
	if err != nil {
		return err
	}

	_, err = media.Stop(ctx)
	return err
}

func (p *Player) Pause() error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	media, err := client.Media(ctx)
	if err != nil {
		return err
	}

	_, err = media.Pause(ctx)
	return err
}

func (p *Player) Volume() (int64, error) {
	return atomic.LoadInt64(&p.volume), nil
}

func (p *Player) SetVolume(percent int64) error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	r := client.Receiver()
	if r == nil {
		return nil
	}

	if percent > 100 {
		percent = 100
	} else if percent < 0 {
		percent = 0
	}

	level := float64(percent) / 100

	return p.setVolume(&level, nil)
}

func (p *Player) Mute() (bool, error) {
	return atomic.LoadInt64(&p.mute) == 1, nil
}

func (p *Player) SetMute(mute bool) error {
	return p.setVolume(nil, &mute)
}

// FIXME: библиотека криво обрабатывает закрытие коннекта, сам коннект скидывает
// а проинициализированный media у клиента не скидывает, хотя там уже хранится закрытый коннект
func (p *Player) Close() {
	var client *cast.Client

	p.mutex.RLock()
	client = p.client
	p.mutex.RUnlock()

	if client != nil {
		r := client.Receiver()
		if r != nil {
			r.QuitApp(context.Background())
		}

		p.client.Close()
	}

	p.setStatus(players.StatusStopped)
	atomic.StoreInt64(&p.connected, 0)

	p.done <- struct{}{}
}

func (p *Player) setVolume(level *float64, muted *bool) error {
	ctx := context.Background()

	client, err := p.connect(ctx)
	if err != nil {
		return err
	}

	r := client.Receiver()
	if r == nil {
		return nil
	}

	volume := controllers.Volume{
		Level: level,
		Muted: muted,
	}
	_, err = r.SetVolume(ctx, &volume)

	return err
}

func (p *Player) setStatus(status players.Status) {
	atomic.StoreInt64(&p.status, status.Int64())
}

func (p *Player) Status() players.Status {
	return players.Status(atomic.LoadInt64(&p.status))
}

func (p *Player) connect(ctx context.Context) (*cast.Client, error) {
	if atomic.LoadInt64(&p.connected) != 1 {
		client := cast.NewClient(p.host, p.port)

		err := client.Connect(ctx)
		if err != nil {
			return nil, err
		}

		atomic.StoreInt64(&p.connected, 1)
		go func() {
			p.doEvents(client.Events)
		}()

		p.mutex.Lock()
		p.client = client
		p.mutex.Unlock()

		return client, nil
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.client, nil
}

func (p *Player) doEvents(ch chan events.Event) {
	for {
		select {
		case event := <-ch:
			switch t := event.(type) {
			case events.Connected:
				//fmt.Println("events.Connected")

			case events.Disconnected:
				p.Close()

				//fmt.Println("events.Disconnected")

			case events.AppStarted:
				p.setStatus(players.StatusStopped)

				//fmt.Println("events.AppStarted")

			case events.AppStopped:
				p.Close()

				//fmt.Println("events.AppStopped")

			case events.StatusUpdated:
				atomic.StoreInt64(&p.volume, int64(math.Round(t.Level*100)))

				if t.Muted {
					atomic.StoreInt64(&p.mute, 1)
				} else {
					atomic.StoreInt64(&p.mute, 0)
				}

				//fmt.Println("events.StatusUpdated", t.Level, t.Muted, t)

			case controllers.MediaStatus:
				switch t.PlayerState {
				case "PLAYING":
					p.setStatus(players.StatusPlaying)

				case "FINISHED", "IDLE":
					p.setStatus(players.StatusStopped)

				case "PAUSED":
					p.setStatus(players.StatusPause)
				}

				// fmt.Println("events.MediaStatus", t.PlayerState)

			default:
				//fmt.Printf("Unknown event: %#v\n", t)
			}

		case <-p.done:
			return
		}
	}
}
