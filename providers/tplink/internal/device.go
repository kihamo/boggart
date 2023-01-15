package internal

// https://github.com/fishbigger/TapoP100
// https://github.com/rk295/tapo-go

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/performance"
	connection "github.com/kihamo/boggart/protocols/http"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	DefaultScheme        = "http"
	DefaultPort          = 80
	DefaultPath          = "app"
	DefaultQueryKeyToken = "token"

	MethodSecurePassThrough = "securePassthrough"
	MethodHandshake         = "handshake"
	MethodLoginDevice       = "login_device"
	MethodGetDeviceInfo     = "get_device_info"
)

type Device struct {
	connection *connection.Client

	u        *url.URL
	username string
	password string

	cipher      *Cipher
	cipherMutex sync.RWMutex
	token       *atomic.String
}

func NewDevice(addr, username, password string) *Device {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = addr + ":" + strconv.Itoa(DefaultPort)
	}

	// encrypt credentials
	hash := sha1.New()
	hash.Write([]byte(username))
	digest := hex.EncodeToString(hash.Sum(nil))

	encryptedUsername := base64.StdEncoding.EncodeToString([]byte(digest))
	encryptedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	return &Device{
		connection: connection.NewClient(). /*WithTimeout(time.Second * 2).*/ WithDebug(true),

		u: &url.URL{
			Scheme: DefaultScheme,
			Host:   addr,
			Path:   DefaultPath,
		},
		username: encryptedUsername,
		password: encryptedPassword,
		token:    atomic.NewString(),
	}
}

func (d *Device) getCipher(ctx context.Context) (*Cipher, error) {
	d.cipherMutex.RLock()
	if d.cipher == nil {
		d.cipherMutex.RUnlock()

		if err := d.handshake(ctx); err != nil {
			return nil, err
		}

		d.cipherMutex.RLock()
	}

	defer d.cipherMutex.RUnlock()
	return d.cipher, nil
}

func (d *Device) getURL() string {
	if d.token.IsEmpty() {
		return d.u.String()
	}

	cp := &url.URL{}
	*cp = *d.u

	q := cp.Query()
	q.Set(DefaultQueryKeyToken, d.token.String())
	cp.RawQuery = q.Encode()

	return cp.String()
}

func (d *Device) handshake(ctx context.Context) error {
	d.connection.Reset()
	d.token.Set("")

	// generate key pair
	const bits = 1024

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		return fmt.Errorf("generate RSA key failed: %w", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("marshal public key failed: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// send request
	type params struct {
		Key             string `json:"key"`
		RequestTimeMils int    `json:"requestTimeMils"`
	}

	request := NewRequest().
		WithMethod(MethodHandshake).
		WithParams(params{
			Key:             performance.UnsafeBytes2String(publicKeyPEM),
			RequestTimeMils: 0,
		})

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("json marshal request body failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.getURL(), bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("create http request failed: %w", err)
	}

	resp, err := d.connection.Do(req)
	if err != nil {
		return fmt.Errorf("call http request failed: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Key string `json:"key,omitempty"`
		} `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("json decode of response body failed: %w", err)
	}

	if response.ErrorCode != 0 {
		return fmt.Errorf("error code %d", response.ErrorCode)
	}

	decodedEncryptionKey, err := base64.StdEncoding.DecodeString(response.Result.Key)
	if err != nil {
		return fmt.Errorf("decode encrypted key failed: %w", err)
	}

	encryptionKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedEncryptionKey)
	if err != nil {
		return fmt.Errorf("decrypt failed: %w", err)
	}

	cipher, err := NewCipher(encryptionKey[:16], encryptionKey[16:])
	if err != nil {
		return fmt.Errorf("creater cipher failed: %w", err)
	}

	d.cipherMutex.Lock()
	d.cipher = cipher
	d.cipherMutex.Unlock()

	return nil
}

func (d *Device) login(ctx context.Context) error {
	type requestParams struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	request := NewRequest().
		WithMethod(MethodLoginDevice).
		WithParams(requestParams{
			Username: d.username,
			Password: d.password,
		})

	type loginResponse struct {
		Token string `json:"token"`
	}

	response := loginResponse{}

	if err := d.InvokeWithContext(ctx, request, &response); err != nil {
		return err
	}

	if response.Token != "" {
		d.token.Set(response.Token)
	}

	return nil
}

func (d *Device) InvokeWithContext(ctx context.Context, request *Request, response interface{}) error {
	if request == nil {
		return errors.New("request is nil")
	}

	// auto login
	if d.token.IsEmpty() && request.Method != MethodLoginDevice {
		if err := d.login(ctx); err != nil {
			return err
		}
	}

	cpr, err := d.getCipher(ctx)
	if err != nil {
		return fmt.Errorf("get cipher failed: %w", err)
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	fmt.Println("Request body", string(requestBody))

	type requestParams struct {
		Request string `json:"request"`
	}

	securityRequest := NewRequest().
		WithMethod(MethodSecurePassThrough).
		WithParams(requestParams{
			Request: base64.StdEncoding.EncodeToString(cpr.Encrypt(requestBody)),
		})

	securityRequestBody, err := json.Marshal(securityRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.getURL(), bytes.NewBuffer(securityRequestBody))
	if err != nil {
		return err
	}
	req.Close = true

	resp, err := d.connection.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if val := resp.Header.Get("Content-Type"); val != "application/json;charset=UTF-8" {
		return fmt.Errorf("wrong content type response %s", val)
	}

	type responseEncrypted struct {
		ErrorCode int `json:"error_code"`
		Result    struct {
			Response string `json:"Response"`
		} `json:"result"`
	}

	var encryptedResponse responseEncrypted
	if err = json.NewDecoder(resp.Body).Decode(&encryptedResponse); err != nil {
		return err
	}

	if encryptedResponse.ErrorCode != 0 {
		return fmt.Errorf("error code %d", encryptedResponse.ErrorCode)
	}

	// skip decode if target response is empty
	if response == nil {
		return nil
	}

	decodedResponse, err := base64.StdEncoding.DecodeString(encryptedResponse.Result.Response)
	if err != nil {
		return err
	}

	targetResponse := NewResponse().WithResult(response)
	if err := json.Unmarshal(cpr.Decrypt(decodedResponse), &targetResponse); err != nil {
		return err
	}

	return nil
}

func (d *Device) GetDeviceInfo(ctx context.Context) (interface{}, error) {
	request := NewRequest().
		WithMethod(MethodGetDeviceInfo)

	var response interface{}

	err := d.InvokeWithContext(ctx, request, &response)
	if err != nil {
		return nil, err
	}

	return nil, err
}
