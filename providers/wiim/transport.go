package wiim

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
)

type transport struct {
	runtime.ClientTransport

	proxied runtime.ClientTransport
}

// Очередной костыль так как openapi не позволяет обрабатывать {variable} в path но после ?
// добавляем эту логику сами. Строго говоря все что ? opeapi справедливо считает query, но
// без такого механизма не возможно различать запросы у которых одинаковый path как в wiim
func (t *transport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	fake := &reguestFake{}
	operation.Params.WriteToRequest(fake, nil)

	for k, v := range fake.pathParams {
		operation.PathPattern = strings.Replace(operation.PathPattern, "{"+k+"}", url.PathEscape(v), -1)
	}

	return t.proxied.Submit(operation)
}

type reguestFake struct {
	runtime.ClientRequest

	pathParams map[string]string
}

func (r *reguestFake) SetHeaderParam(string, ...string) error {
	return nil
}

func (r *reguestFake) GetHeaderParams() http.Header {
	return nil
}

func (r *reguestFake) SetQueryParam(string, ...string) error {
	return nil
}

func (r *reguestFake) SetFormParam(string, ...string) error {
	return nil
}

func (r *reguestFake) SetPathParam(name string, value string) error {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}

	r.pathParams[name] = value
	return nil
}

func (r *reguestFake) GetQueryParams() url.Values {
	return nil
}

func (r *reguestFake) SetFileParam(string, ...runtime.NamedReadCloser) error {
	return nil
}

func (r *reguestFake) SetBodyParam(interface{}) error {
	return nil
}

func (r *reguestFake) SetTimeout(time.Duration) error {
	return nil
}

func (r *reguestFake) GetMethod() string {
	return ""
}

func (r *reguestFake) GetPath() string {
	return ""
}

func (r *reguestFake) GetBody() []byte {
	return nil
}

func (r *reguestFake) GetBodyParam() interface{} {
	return nil
}

func (r *reguestFake) GetFileParam() map[string][]runtime.NamedReadCloser {
	return nil
}
