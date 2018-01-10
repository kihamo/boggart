package softvideo

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	AccountURL = "https://user.softvideo.ru/"
	UserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	randSrc       = rand.NewSource(time.Now().UnixNano())
	balanceRegexp = regexp.MustCompile(`badge-success">([^\s]+)[^<]*</span>`)
)

type Client struct {
	connection *http.Client
	login      string
	password   string
}

func NewClient(login, password string) *Client {
	return &Client{
		connection: &http.Client{},
		login:      login,
		password:   password,
	}
}

func (c *Client) generateSID() string {
	n := 32
	b := make([]byte, n)

	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (c *Client) Balance() (float64, error) {
	data := url.Values{
		"bootstrap[username]": {c.login},
		"bootstrap[password]": {c.password},
		"bootstrap[send]:":    {"<i class=\"icon-ok\"></i>"},
	}

	request, err := http.NewRequest(http.MethodPost, AccountURL, strings.NewReader(data.Encode()))
	if err != nil {
		return -1, err
	}

	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie := http.Cookie{Name: "PHPSESSID", Value: c.generateSID()}
	request.AddCookie(&cookie)

	/*
		if dump, err := httputil.DumpRequestOut(request, true); err != nil {
			return -1, err
		} else {
			fmt.Printf("%q", dump)
		}
	*/

	response, err := c.connection.Do(request)
	if err != nil {
		return -1, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}

	/*
		if dump, err := httputil.DumpResponse(response, true); err != nil {
			return -1, err
		} else {
			fmt.Printf("%q", dump)
		}
	*/

	submatch := balanceRegexp.FindStringSubmatch(string(body))
	if len(submatch) != 2 {
		return -1, errors.New("Balance string not found in page")
	}

	return strconv.ParseFloat(submatch[1], 10)
}
