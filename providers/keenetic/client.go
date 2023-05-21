package keenetic

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http/cookiejar"
	"net/textproto"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/keenetic/client"
	"github.com/kihamo/boggart/providers/keenetic/client/user"
)

type Client struct {
	*client.Keenetic
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cl := &Client{
		Keenetic: client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(address)),
	}

	if rt, ok := cl.Keenetic.Transport.(*httptransport.Runtime); ok {
		rt.Jar, _ = cookiejar.New(nil)

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	// 1. инициирующий запрос, чтобы получить признаки шифрования
	_, err := c.User.GetAuth(user.NewGetAuthParamsWithContext(ctx))
	if err == nil {
		return err
	}

	if www, ok := err.(*user.GetAuthUnauthorized); ok {
		realm := www.XNDMRealm
		challenge := www.XNDMChallenge

		//  пытаем найти в заголовке WWW-Authenticate
		if realm == "" || challenge == "" {
			values := ParsePairs(www.WWWAuthenticate)

			realm = values["realm"]
			challenge = values["challenge"]
		}

		if realm == "" || challenge == "" {
			return err
		}

		h := md5.New()
		h.Write([]byte(username))
		h.Write([]byte(":"))
		h.Write([]byte(realm))
		h.Write([]byte(":"))
		h.Write([]byte(password))

		pwd := sha256.New()
		pwd.Write([]byte(challenge))
		pwd.Write([]byte(hex.EncodeToString(h.Sum(nil))))
		pwdHash := hex.EncodeToString(pwd.Sum(nil))

		params := user.NewPostAuthParamsWithContext(ctx)
		params.Request.Login = username
		params.Request.Password = pwdHash

		_, err = c.User.PostAuth(params)
	}

	return err
}

func ParsePairs(line string) map[string]string {
	m := make(map[string]string)

	var s scanner.Scanner

	s.Init(strings.NewReader(textproto.TrimString(line)))
	s.Mode = scanner.ScanIdents | scanner.ScanStrings
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '_' || ch == '-' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}

	var (
		key, part string
		isValue   bool
	)

	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		part = s.TokenText()
		if u, e := strconv.Unquote(part); e == nil {
			part = u
		}

		if part == "=" {
			isValue = true
			continue
		}

		if isValue {
			isValue = false
			m[key] = part
		} else {
			key = part
			m[key] = ""
		}
	}

	return m
}
