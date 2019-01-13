package internal

import (
	"fmt"
	"net"
	"net/url"
	"reflect"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/mitchellh/mapstructure"
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

		if t == reflect.TypeOf(boggart.IP{}) {
			ip := net.ParseIP(data.(string))

			if ip == nil {
				return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
			}

			return boggart.IP{
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

		if t == reflect.TypeOf(boggart.HardwareAddr{}) {
			a, err := net.ParseMAC(data.(string))
			if err != nil {
				return boggart.HardwareAddr{}, err
			}

			return boggart.HardwareAddr{
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

		if t == reflect.TypeOf(boggart.URL{}) {
			u, err := url.Parse(data.(string))
			if err != nil {
				return boggart.URL{}, err
			}

			return boggart.URL{
				URL: *u,
			}, nil
		}

		return data, nil
	}
}
