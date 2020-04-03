package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

type Client struct {
	mutex sync.RWMutex

	debug      uint64
	connection *http.Client
	cookies    []*http.Cookie
	userAgent  string
}

func NewClient() *Client {
	client := &Client{
		connection: &http.Client{
			Transport: &nethttp.Transport{},
		},
		userAgent: DefaultUserAgent,
	}

	client.Reset()

	return client
}

func (c *Client) WithDebug(debug bool) *Client {
	if debug {
		atomic.StoreUint64(&c.debug, 1)
	} else {
		atomic.StoreUint64(&c.debug, 0)
	}

	return c
}

func (c *Client) WithUserAgent(agent string) *Client {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.userAgent = agent

	return c
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.connection.Timeout = timeout

	return c
}

func (c *Client) Do(request *http.Request) (*http.Response, error) {
	if agent := request.Header.Get("User-Agent"); agent == "" {
		c.mutex.RLock()
		request.Header.Set("User-Agent", c.userAgent)
		c.mutex.RUnlock()
	}

	c.mutex.RLock()
	for _, cookie := range c.cookies {
		request.AddCookie(cookie)
	}
	c.mutex.RUnlock()

	debug := atomic.LoadUint64(&c.debug)

	if debug == 1 {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}

		fmt.Printf("\n\n%q", dump)
	}

	//request, closer := tracing.TraceRequest(opentracing.GlobalTracer(), request)
	//defer closer.Finish()

	response, err := c.connection.Do(request)
	if err != nil {
		return nil, err
	}

	if debug == 1 {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, err
		}

		fmt.Printf("\n\n%q", dump)
	}

	c.mutex.Lock()
	c.cookies = append(c.cookies, response.Cookies()...)
	c.mutex.Unlock()

	return response, nil
}

func (c *Client) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connection.Jar, _ = cookiejar.New(nil)
	c.cookies = []*http.Cookie{}
}

func (c *Client) Get(ctx context.Context, u string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, u, nil)

	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	return c.Do(request)
}

func (c *Client) GetAjax(ctx context.Context, u string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	return c.Do(request)
}

func (c *Client) Post(ctx context.Context, u string, data map[string]string) (*http.Response, error) {
	values := url.Values{}
	for key, val := range data {
		values[key] = []string{val}
	}

	request, err := http.NewRequest(http.MethodPost, u, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.Do(request)
}
