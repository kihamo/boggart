package internal

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	connection "github.com/kihamo/boggart/protocols/http"
)

type Auth func(ctx context.Context) error

type Client struct {
	connection *connection.Client

	auth Auth

	baseURL  string
	login    string
	password string

	session         string
	sessionLifeTime time.Time
	sessionLock     sync.RWMutex
}

func New(baseURL, login, password string, auth Auth) *Client {
	return &Client{
		connection: connection.NewClient(),
		auth:       auth,
		baseURL:    baseURL,
		login:      login,
		password:   password,
	}
}

func (c *Client) Login() string {
	return c.login
}

func (c *Client) Password() string {
	return c.password
}

func (c *Client) SetSession(session string) {
	c.sessionLock.Lock()
	c.session = session
	c.sessionLifeTime = time.Now().Add(time.Minute * 5)
	c.sessionLock.Unlock()
}

func (c *Client) Session() string {
	c.sessionLock.RLock()
	lt := c.sessionLifeTime
	c.sessionLock.RUnlock()

	if time.Now().After(lt) {
		c.SetSession("")
	}

	c.sessionLock.RLock()
	defer c.sessionLock.RUnlock()

	return c.session
}

func (c *Client) DoRequestRaw(ctx context.Context, u *url.URL, body io.Reader, writer io.Writer) error {
	values := u.Query()

	// auto login
	if !values.Has("action") || values.Get("action") != "auth" {
		session := c.Session()

		if session == "" {
			if err := c.auth(ctx); err != nil {
				return err
			}

			session = c.Session()
		}
	}

	u.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.connection.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status isn't 200 OK returns " + strconv.Itoa(resp.StatusCode))
	}

	if writer != nil {
		_, err = io.Copy(writer, resp.Body)
	}

	return err
}

func (c *Client) DoRequest(ctx context.Context, action string, query, body map[string]string, data interface{}) error {
	values := url.Values{}

	if action != "" {
		values.Add("action", action)
	}

	for key, val := range query {
		values.Add(key, val)
	}

	// auto login
	if action != "auth" {
		session := c.Session()

		if session == "" {
			if err := c.auth(ctx); err != nil {
				return err
			}

			session = c.Session()
		}

		values.Add("session", session)
	}

	resp, err := c.connection.Post(ctx, c.baseURL+"?"+values.Encode(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status isn't 200 OK returns " + strconv.Itoa(resp.StatusCode))
	}

	if writer, ok := data.(io.Writer); ok {
		_, err = io.Copy(writer, resp.Body)
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r := &response{}
	if data != nil {
		r.Data = data
	}

	err = json.Unmarshal(b, r)
	if err != nil {
		err = multierror.Append(err, errors.New("bad content "+string(b)))
	}

	if err == nil && r.ErrorMessage != "" {
		if strings.Compare(r.ErrorMessage, "Session not found, specified query requires authentication") == 0 {
			c.SetSession("")
		}

		err = errors.New(r.ErrorMessage)
	}

	return err
}
