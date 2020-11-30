package pantum

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/pantum/client"
)

func isSeparator(r rune) bool {
	return r == '(' || r == ')' || r == ','
}

type property struct {
	Value  string `json:"value"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Module string `json:"module"`
	Index  int    `json:"index"`
}

type Client struct {
	*client.Pantum
}

func New(address *url.URL, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().
		WithSchemes([]string{address.Scheme}).
		WithHost(address.Host)

	cl := &Client{
		Pantum: client.NewHTTPClientWithConfig(nil, cfg),
	}

	if rt, ok := cl.Pantum.Transport.(*httptransport.Runtime); ok {
		rt.Consumers["text/html"] = omConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func omConsumer() runtime.Consumer {
	jsonConsumer := runtime.JSONConsumer().Consume

	return runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		jsonBody := bytes.NewBuffer(nil)

		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanLines)

		prop := &property{}

		for scanner.Scan() {
			parts := bytes.FieldsFunc(scanner.Bytes(), isSeparator)

			if len(parts) > 4 {
				for i, v := range parts {
					parts[i] = bytes.Trim(v, " '\"")
				}

				valueDecode, _ := base64.StdEncoding.DecodeString(string(parts[1]))

				prop.Value = string(valueDecode)
				prop.Name = string(parts[2])
				prop.Type = strings.TrimPrefix(string(parts[3]), "SN.TYPE.")
				prop.Module = strings.TrimPrefix(string(parts[4]), "MODULE_")
				prop.Index, _ = strconv.Atoi(string(parts[5]))

				if line, err := json.Marshal(prop); err == nil {
					if jsonBody.Len() == 0 {
						jsonBody.WriteRune('[')
					} else {
						jsonBody.WriteByte(',')
					}

					jsonBody.Write(line)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if jsonBody.Len() > 0 {
			jsonBody.WriteByte(']')
		}

		return jsonConsumer(jsonBody, data)
	})
}
