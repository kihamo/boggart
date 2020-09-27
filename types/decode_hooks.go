package types

import (
	"fmt"
	"net"
	"net/url"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

var (
	DecodeHookExec               = mapstructure.DecodeHookExec
	ComposeDecodeHookFunc        = mapstructure.ComposeDecodeHookFunc
	StringToSliceHookFunc        = mapstructure.StringToSliceHookFunc
	StringToTimeDurationHookFunc = mapstructure.StringToTimeDurationHookFunc
	StringToIPNetHookFunc        = mapstructure.StringToIPNetHookFunc
	StringToTimeHookFunc         = mapstructure.StringToTimeHookFunc
	WeaklyTypedHook              = mapstructure.WeaklyTypedHook
)

func StringToIPHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(net.IP{}) {
			ip := net.ParseIP(data.(string))

			if ip == nil {
				return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
			}

			return ip, nil
		}

		if t == reflect.TypeOf(IP{}) {
			ip := net.ParseIP(data.(string))

			if ip == nil {
				return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
			}

			return IP{
				IP: ip,
			}, nil
		}

		return data, nil
	}
}

func StringToMACHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(net.HardwareAddr{}) {
			return net.ParseMAC(data.(string))
		}

		if t == reflect.TypeOf(HardwareAddr{}) {
			a, err := net.ParseMAC(data.(string))
			if err != nil {
				return HardwareAddr{}, err
			}

			return HardwareAddr{
				HardwareAddr: a,
			}, nil
		}

		return data, nil
	}
}

func StringToURLHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(url.URL{}) {
			return url.Parse(data.(string))
		}

		if t == reflect.TypeOf(URL{}) {
			u, err := url.Parse(data.(string))
			if err != nil {
				return URL{}, err
			}

			return URL{
				URL: *u,
			}, nil
		}

		return data, nil
	}
}

func StringToFileModeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t == reflect.TypeOf(FileMode{}) {
			mode := &FileMode{}
			err := mode.UnmarshalText([]byte(data.(string)))
			return *mode, err
		}

		return data, nil
	}
}
