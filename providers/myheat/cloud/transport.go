package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/providers/myheat/cloud/models"
)

const (
	requestKeyAction = "action"
	requestKeyApiKey = "key"
	requestKeyLogin  = "login"
)

type transport struct {
	runtime.ClientTransport

	proxied runtime.ClientTransport

	login string
	key   string
}

// Жутчайший костыль авторизации, но через дефолтный авторизатор действовать не получится, там не возможно подменить тело
// запрос, а только заголовки и урл запроса, поэтому пришлось вылезти на уровень транспорта и действовать тут
func (t *transport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	// делаем URL более чистым, так как тоже используется хак чтобы не попасть на дублирование в рамках одного yaml
	if u, err := url.Parse(operation.PathPattern); err == nil {
		u.RawQuery = ""
		operation.PathPattern = u.String()
	}

	params := operation.Params

	operation.Params = runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) (err error) {
		err = params.WriteToRequest(req, reg)

		if err != nil {
			return err
		}

		bodyParams := map[string]interface{}{
			requestKeyAction: operation.ID, // оооочень не явно, но так проще всего
			requestKeyLogin:  t.login,
			requestKeyApiKey: t.key,
		}

		// если есть параметры при запросе, то необходимо смерджить
		if body := req.GetBodyParam(); body != nil {
			value := reflect.ValueOf(body)
			if value.Kind() == reflect.Ptr {
				value = value.Elem()
			}

			if value.Kind() != reflect.Struct {
				return fmt.Errorf("only accepts structs; got %T", value)
			}

			typ := value.Type()
			for i := 0; i < value.NumField(); i++ {
				field := typ.Field(i)

				if tagValue := field.Tag.Get("json"); tagValue != "" {
					bodyParams[tagValue] = value.Field(i).Interface()
				}
			}
		}

		return req.SetBodyParam(bodyParams)
	})

	resp, err := t.proxied.Submit(operation)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// так как хрен мне, а не oneOf в response swagger, то парсим ошибку в ответе сами эмулируем ошибочный запрос
type roundTripper struct {
	proxied http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.proxied.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// эмулируем http ошибку, так как у hilink всегда 200 OK
	body, err := io.ReadAll(resp.Body)
	if err == nil {
		var (
			errorResponse models.Error
			newBody       []byte
		)

		if err := json.Unmarshal(body, &errorResponse); err == nil {
			if errorResponse.Error > 0 {
				if b, err := json.Marshal(errorResponse); err == nil {
					newBody = b
				}
			}
		}

		if len(newBody) > 0 {
			resp.StatusCode = http.StatusBadRequest
			resp.Status = http.StatusText(http.StatusBadRequest)
			resp.Body = io.NopCloser(bytes.NewReader(newBody))
			resp.ContentLength = int64(len(newBody))
		} else {
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}
	}

	return resp, err
}
