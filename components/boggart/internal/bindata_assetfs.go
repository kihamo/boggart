// Code generated by go-bindata.
// sources:
// templates/views/devices.html
// assets/js/devices.js
// locales/ru/LC_MESSAGES/boggart.mo
// locales/ru/LC_MESSAGES/devices.mo
// DO NOT EDIT!

package internal

import (
	"github.com/elazarl/go-bindata-assetfs"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesViewsDevicesHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x56\x4d\x8f\xd3\x30\x10\xbd\xf7\x57\x8c\x7c\x82\x43\x6b\x6d\x4f\x1c\xb2\x95\x90\x16\xb4\x48\x85\x03\x2a\x67\xe4\xc4\x93\xed\x94\xd4\x0e\x1e\x13\xba\x8a\xfa\xdf\x51\xbe\xba\x26\x49\xdb\xed\xd2\xe5\xd2\x26\xe3\x79\xf3\xe6\xe3\x8d\x9c\xb2\x04\x8d\x29\x19\x04\x91\x58\xe3\xd1\x78\x01\xfb\xfd\x64\x12\x69\x2a\x20\xc9\x14\xf3\xad\x70\xf6\xb7\x00\xd2\xb7\x42\x63\x41\x09\xb2\x58\x4c\x00\x00\x42\x97\xdd\xf7\x5c\x19\xcc\xda\x93\xe1\xa9\x27\x9f\x61\x70\x5a\x7b\xac\xe7\x8b\xb2\x04\xba\x79\x67\x40\xdc\xb5\xa1\x61\x06\xfb\x7d\x24\xd7\xf3\x9e\x6f\x10\x2d\xc9\x50\xb9\x94\x76\x62\x11\x49\x4d\x45\x40\xd9\x7b\xfd\x2b\x83\xae\xba\x5e\x5c\xaf\xe2\x0c\x3b\xaf\xe6\xa5\xfe\x9d\xb2\x77\x94\xa3\xee\xf9\x37\x98\x35\x2a\x3d\x66\x77\x43\x63\x0b\x78\x2a\x74\xf5\x98\x63\x57\xa5\x5f\x1f\x05\x74\x29\x6d\xf5\x34\xb1\xd9\x74\x2e\x9e\x22\x7c\xba\x7b\x06\x3e\xec\x2c\x27\x8e\x72\x4f\xd6\x5c\xce\x7b\x13\xf0\xae\x14\xff\xe0\x7f\x0b\xf1\x3e\xa9\xd2\x38\x19\x24\x92\x63\x6d\xac\x7c\x07\x4d\x8f\x64\x3d\xaa\x51\x01\xb4\x8f\xed\xdf\xb8\x9c\x33\x62\x8f\x06\xdd\x35\x05\x1d\x51\xe7\x93\x2a\x48\xd5\x54\xa3\x4a\x2b\xa1\xd2\x02\x0e\x5d\x58\x1e\x78\xff\xb7\xda\x03\x9f\x46\xe5\x0e\x39\xb7\x86\xa9\xe8\xd7\x02\x2f\x5c\x0e\x38\xb1\x20\x70\x6a\x49\xa0\xaf\xdb\xf3\x32\x1f\x40\xbe\xa8\xed\xd9\xdd\x1a\x80\x3e\x14\x68\xfc\x59\x5d\x0f\x60\x1f\xc9\xe1\x8b\x50\xec\x21\x25\x87\xfa\x62\xec\x52\x3d\x17\x3a\xbe\x43\x70\x6c\x8f\x60\x6c\x97\x60\x44\x62\xc7\xd6\xab\x2c\x01\x8d\xae\xef\x8d\xe0\x3e\xa9\x88\xea\xcb\xa4\x02\x94\x25\xb0\x57\x9e\x92\xfb\xd5\xe7\x25\xbc\x69\x9e\xbf\x7d\x5d\x82\x90\x5a\xf1\x3a\xb6\xca\x69\xa9\x98\xd1\xb3\x2c\xd0\x68\xeb\x58\x6a\xe5\x55\x9d\x15\xcf\x0c\xfa\x69\xcc\x32\xe1\xc6\xba\x6a\xac\xb1\xb5\x9e\xbd\x53\xf9\x6c\x4b\x66\x96\x30\x0b\x48\x55\xc6\xf8\xf6\x8a\xac\x29\xed\x50\x57\xa5\xa0\xeb\x32\xa8\x4d\xf7\xb5\xe9\x74\x0a\xe3\x7d\xd9\xf0\x15\xbb\x22\x37\x2c\x37\x3f\x7f\xa1\x7b\x9c\x05\x8d\xa9\x72\xd9\xbc\x46\x37\x62\xae\x08\x8f\x8e\xe0\x55\x38\x83\x09\xf4\xc8\xc3\x41\x5c\x4c\x1f\xdb\x87\x07\xe5\x7c\x47\x5e\x45\x6e\xbe\x41\x7a\x61\x0e\x43\xfc\x13\x00\x00\xff\xff\x65\x67\xde\x0b\x2b\x09\x00\x00")

func templatesViewsDevicesHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesViewsDevicesHtml,
		"templates/views/devices.html",
	)
}

func templatesViewsDevicesHtml() (*asset, error) {
	bytes, err := templatesViewsDevicesHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/views/devices.html", size: 2347, mode: os.FileMode(420), modTime: time.Unix(1531045566, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsJsDevicesJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x58\x51\x6f\xdb\x36\x10\x7e\xef\xaf\xb8\xba\x01\x28\x21\xb1\xbc\x61\x18\x30\x28\x91\xf7\xd0\xf4\xa1\x40\xb1\x16\x68\xde\xda\x22\xa0\xc5\x93\xcc\x86\x26\x0d\x92\xb2\xeb\x15\xfe\xef\x03\x45\x49\xb5\x62\x4b\x72\x66\x77\xeb\x30\x3e\x18\x12\xc5\x3b\xde\xdd\xf7\xf1\xee\xe8\x8b\x80\xa9\xb4\x58\xa0\xb4\x61\xa4\x91\xb2\x4d\x90\x15\x32\xb5\x5c\x49\x08\x42\xf8\xfa\x0c\x00\x60\x45\x35\x58\x3a\x13\x78\x8b\x2b\x9e\xa2\x81\x04\x2e\x02\xf2\x82\x55\x6f\xe5\x27\x12\x96\x4b\xdd\x88\x6e\xa9\xa5\x77\x6e\x32\xf8\xda\x4c\xba\xb1\xa4\x39\xbe\x41\x99\xdb\x79\x0c\xbf\xfe\x74\xd5\xfa\x26\xa8\xcc\x0b\x9a\x63\x0c\x6d\x19\x37\x0a\x2d\x62\x20\x13\x46\xcd\x7c\xa6\xa8\x66\x13\x46\x2d\x2d\x77\x35\x13\xfe\xf3\x6f\x32\xfa\x6c\x94\x24\x2d\xb1\x6d\x5b\x3b\xfd\x4c\xbf\xf4\x68\x9e\xa9\x3c\xa7\xda\x4e\x2a\x8f\x26\xbf\xa3\xb4\xdc\x6e\x92\xea\x9d\x5c\xed\x09\x3a\x0b\xde\xeb\x34\x06\xe2\x9e\x7a\xf7\x4e\x95\x28\x16\xd2\xc4\xf0\x61\x4f\xcb\xbe\x41\xb5\xee\x18\x88\xdd\x2c\x0f\x6e\xed\x86\x46\xc9\x50\xc7\xf0\x0d\xab\x72\x75\xd8\xa1\x11\x2a\x10\x53\x25\x2d\x4a\x0b\x09\x10\x72\xfd\xac\x73\x69\xa6\x34\x04\x6e\x3d\x07\x2e\x61\x50\xb3\x77\xd2\x6b\xbe\x4c\x80\xdc\x98\x25\x95\x90\x0a\x6a\x4c\x32\x12\x74\x86\x02\xca\xdf\xb1\x29\xd2\x14\x8d\x19\x4d\x09\x5c\x7a\xb5\x1f\xf8\x27\xb8\x04\x72\x33\x71\x22\x53\x20\xd7\x9d\x9b\x6c\xbb\xcd\xd5\x68\x0b\x2d\x6b\x13\x0e\xab\xd8\xee\xcd\x6e\xf7\x23\xdb\x0b\x07\x67\x1d\x58\xac\xb8\xe1\x33\x81\x31\x64\x54\x18\x3c\x79\x1f\x86\x26\xd5\x7c\xe9\x50\x25\x27\x2b\xb3\xd4\x3c\x98\xfb\x54\x15\xd2\x9e\xae\x0c\xa5\x3b\x72\x5d\x61\xd8\xa7\xa4\x13\xbb\x2a\x71\xbe\x02\xad\xd6\x47\x92\xb3\x87\x97\x3c\x83\x40\xab\x75\x54\x19\x72\x2c\x27\x1d\x25\x67\x85\xb5\xca\x73\x39\x19\xf9\x97\x51\x4d\xd1\x99\x95\x30\xb3\x72\xcc\xa8\xcc\x51\x97\x8f\x3c\x75\xdf\x9d\x03\x63\xab\xf2\x5c\x60\x32\x5a\x28\x46\x45\x3d\x47\x75\x8e\x36\x19\xbd\xd8\x9d\x2c\x9f\xc7\x96\x5b\xb7\xfa\xa5\x92\x19\xd7\x0b\x60\xdc\x38\x5b\xc1\x67\x12\x78\xe1\x98\xef\x5c\xe0\xcc\xf1\xbe\x25\x99\x52\x21\x66\x34\x7d\x48\x46\x7e\xf1\x5d\xb9\x71\xf0\x91\xd4\x32\x1a\x73\x6e\x2c\xea\x7b\x2f\xfc\x91\x84\xd7\xe5\x51\xea\x0d\x82\x1b\xe4\x86\xd7\xbe\xe6\x62\xb3\x9c\x3b\xef\xa0\x79\x1a\x6b\x5c\xa8\x15\x8e\xa0\x32\xfd\xb6\x65\xf2\x68\x7a\x33\xe1\xd3\x7d\xf2\xd4\x63\x0b\x28\x0c\x9e\x19\x89\x2a\x51\x9c\x17\x0a\xcf\x9a\x1f\x1a\x09\xf5\xd0\xa0\xf0\x4a\x3e\x05\x84\xc1\xe4\x48\x6e\x18\x5f\xed\x84\x79\x9c\x6b\x55\x2c\xa1\x79\x1a\x7f\xa9\xd2\x72\x93\xc8\x5d\x56\xf6\xf0\x0c\x7a\x76\x1c\xa8\x5c\x66\xaa\x41\xb4\x72\x6c\x9c\xce\x31\x7d\x18\x4d\x87\x08\x9a\x69\x34\xf3\x26\x36\x2f\x2b\xa1\x09\x9f\x9e\xd9\xc6\x35\xd5\x92\xcb\x7c\xcf\xcc\x25\x97\xf9\xb0\x95\x86\xff\x89\x63\xb3\xa0\x42\x34\xa6\xbe\xf3\x82\x4f\xb3\x74\xc2\xf8\x6a\xda\x51\x0b\x0f\x14\xb2\xd6\xcc\xa7\xe6\x6d\x1b\x56\xa9\xb4\x69\xdf\xde\x38\xd2\x4a\xd4\x75\x03\x27\x9a\xf7\xa3\x5b\xb8\x1f\xb2\x4d\x6b\xfc\xf8\x77\x1a\x35\xce\x4e\xaf\xad\x92\x2e\xf0\x0c\x15\x7a\x85\xd2\x1e\xdf\x33\x3a\xa9\x73\xb7\x8c\xa5\x09\xaf\x99\x6b\x1c\x87\xd4\xc3\xb1\x7d\xa3\xcb\x1c\x3e\x3b\x39\x95\x1f\xaa\x2d\xfe\x33\x9d\x63\xc6\x35\x9a\xd3\xc1\x75\x6a\xee\x33\xae\x8d\x7d\x0a\xc0\xd8\x87\x80\x6b\xa8\x9e\x0f\x2d\x82\x9d\x2a\x72\x52\x98\xdd\x46\x77\xea\xbd\xd5\x5c\xe6\xde\xb4\xef\x19\xf1\x7b\x41\xff\x8f\x91\x1a\x2a\x06\x7b\x17\x77\xb0\x33\xc5\x36\x24\x8c\x94\x0c\x48\x2a\x78\xfa\x40\xae\x80\xf8\x5a\x15\xed\x56\x69\x72\xb5\x13\xaf\x56\x1c\x30\x5a\xea\xf2\x50\xde\x62\x46\x0b\x61\x83\x1d\x73\x5d\x4e\xf0\x5a\x5e\x33\x48\x5a\x7f\x23\x44\x5a\xad\x83\x8b\xc0\xce\xb9\x09\xa3\x54\x28\x83\xc6\x06\xc4\x6a\x12\x86\x91\x03\x32\x08\x77\x1b\xad\x9d\xe4\x73\x11\xb9\x42\xf1\xa8\x32\xb9\xba\x1e\x03\x79\xf7\xf6\xfd\xdd\x23\xd4\x3b\xaa\x47\x99\x4f\x6a\xcb\x2e\x81\x4c\x2a\x2f\x5b\xb2\x55\x27\xfa\x8d\x2a\xc1\x21\x02\xb4\xbc\x72\xb6\x45\x1a\x85\xa2\x2c\x78\x84\xdb\xb6\x85\xc7\xe9\xb8\xb8\xb6\xe4\x54\x58\x9e\x0c\xca\x77\x04\xe2\x71\x63\x3d\xf1\x0e\x0e\x20\xa2\x0f\x41\x52\xde\x15\x4b\x8b\xe1\x79\x92\x00\x29\x24\xc3\x8c\x4b\x64\xa4\xeb\x04\x97\x47\xdc\x8b\xf4\x1d\x72\x89\x6b\x78\xf7\x87\xb2\x3c\xdb\x04\xfd\xa9\xa0\x6c\x01\x63\x20\x6f\xb3\x4c\x70\x89\x1d\xb9\xa8\x59\x8d\x5f\x6c\x0c\xc4\xc3\x00\x3b\x31\xf1\xa1\x00\x6e\x40\x1d\xa7\xc8\xc7\x1f\xb5\x56\x7a\x60\xa9\xb1\x1b\xc1\x65\x1e\x03\x99\x29\x65\x8d\xd5\x74\xf9\x4b\xcf\x2d\xa3\x2b\x0b\x0d\x5d\x01\xff\x46\xc8\xe4\x99\x22\x76\x94\x1e\x1f\xb0\x8a\x58\xff\x4c\xc8\x06\x12\x77\x4f\xa2\x58\x73\xc9\xd4\x3a\xda\xbd\x99\x42\xf2\xed\x30\xd4\x19\x6d\x97\xc2\xdf\x29\x5f\xfa\xeb\xf8\x29\xc7\x53\xa3\x29\x84\x85\xc4\x1d\xd0\x8c\x72\xd1\x73\x3a\x8f\xe3\x50\xcd\x9f\x57\x03\xe4\xf7\xd4\xd1\xd1\x02\x8d\xa1\x39\xf6\x2c\x3c\xee\x30\xcd\x39\xab\xff\x00\xec\x5e\x74\x3c\x75\xba\x68\xe3\xdb\x83\x03\xdc\x39\x67\x45\x7a\xb6\x0d\xaf\xff\x0a\x00\x00\xff\xff\x48\x31\x9e\x90\x16\x18\x00\x00")

func assetsJsDevicesJsBytes() ([]byte, error) {
	return bindataRead(
		_assetsJsDevicesJs,
		"assets/js/devices.js",
	)
}

func assetsJsDevicesJs() (*asset, error) {
	bytes, err := assetsJsDevicesJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/js/devices.js", size: 6166, mode: os.FileMode(420), modTime: time.Unix(1533730621, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _localesRuLc_messagesBoggartMo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8e\x3f\x6f\x13\x4b\x14\xc5\x4f\xfc\xf2\x5e\xe1\xc7\x1f\x81\x28\x29\x6e\x0a\x22\x10\x9a\xb0\x6b\x13\x29\xda\x64\x13\x44\xfe\x48\x88\x58\x44\xc1\x50\xd1\x0c\xf6\x60\xaf\xf0\xce\x58\xb3\xb3\x11\x91\x22\x14\x82\xa8\x40\xa2\x40\x74\x08\x8a\x34\xa1\x8b\x02\x91\x12\x25\x36\x15\xfd\xdd\x82\x96\x86\x8f\xc0\x17\x40\x5e\x2f\x44\x9c\x66\xce\x4f\xf7\x9c\x3b\xf7\xfb\xf9\xd1\xb7\x00\xf0\x1f\x80\x8b\x00\x96\x01\xfc\x0f\x60\x03\x43\xed\x00\x38\x03\xe0\x23\x80\x53\x00\x0e\x00\x9c\x05\xf0\x15\xc0\x69\x00\xdf\x0a\xfe\x01\xa0\x3d\x02\xfc\x04\x30\x06\xe0\x41\x09\xb8\x00\xe0\x69\xf1\xbe\x29\x0d\xfb\x3b\x25\xe0\x1c\x80\xfd\x12\x30\x52\xfc\xf1\x2f\x80\x7f\x70\xa2\x52\xf1\x8e\x16\x77\xe5\x8a\x95\x4e\x47\x17\x64\xd2\x7e\x68\xa4\x6d\x16\xa8\xd6\xa2\x86\x4a\x86\x50\x33\x3a\x72\xc6\x46\xba\x35\xe4\xbb\xaa\x91\xda\xc8\xad\x17\x14\x4b\xeb\xa8\x6d\x62\x85\x55\xd5\x35\xd6\x89\x5a\xd2\x8a\x9a\xe2\x66\xda\x4a\x44\xdd\x04\xd4\x54\x6b\x37\x1e\x47\x6d\x19\x9b\x09\x9b\x96\x57\xee\xd4\xc5\xbc\x55\xd2\x45\x46\x8b\x05\xe9\x54\x40\x15\xcf\x9f\x12\x5e\x55\x54\xaa\x54\xa9\x06\x93\x93\x57\xbd\xaa\xe7\x95\x97\x65\xe2\x44\xdd\x4a\x9d\x74\xa4\x33\x36\xa0\xdb\xf9\x0e\xaa\xa5\x56\xc6\xa6\x69\x68\xe6\xaf\xc5\xb3\xe5\x65\xa9\x5b\xa9\x6c\x29\x51\x57\x32\x0e\xe8\x0f\x07\xb4\x9a\x26\x49\x24\x75\xb9\x76\xab\xb6\x28\xee\x2b\x9b\x44\x46\x07\xe4\x4f\x78\xe5\x79\xa3\x9d\xd2\x4e\xd4\xd7\xbb\x2a\x20\xa7\x9e\xb8\x6b\xdd\x8e\x8c\xf4\x34\x35\xda\xd2\x26\xca\x85\xf7\xea\x4b\x62\xea\x24\x37\xb8\xe7\x91\xb2\x62\x51\x37\x4c\x33\xd2\xad\x80\xca\x2b\x9d\xd4\xca\x8e\x58\x32\x36\x4e\x02\xd2\xdd\x1c\x93\xb0\x3a\x4d\x43\x1b\xea\x4b\xbe\x17\x86\x3e\x8d\x8f\xd3\xc0\x7a\x63\xa1\xef\xd3\x1c\x79\x14\xe4\x3c\x1b\x56\x7e\x8f\x66\xc2\xeb\x03\x7b\x39\x8f\xcd\xf8\x1e\x6d\x6c\x0c\x2b\xb3\x61\xc5\xbb\x42\x73\xe4\x53\x40\x95\x69\xf0\x07\xde\xe5\x1e\xef\xf3\x51\xf6\x8a\xb2\xe7\xfc\x25\xdb\xe4\x5d\xde\xe3\x23\xde\xe7\x1e\x1f\x64\xaf\xc1\xdb\xd9\xb3\x6c\x2b\xdb\xe4\x3e\x1f\x0e\x1c\xef\xf1\x2e\xf8\x1d\xf7\xf3\xf9\x16\xf7\xb3\x4d\x3e\xe0\x1e\x7f\x02\xbf\xcf\x5e\xe4\xf5\xde\x20\xb1\xcd\xc7\xdc\xcb\x5e\xf2\x21\xf1\x67\xee\xf3\x31\x7e\x05\x00\x00\xff\xff\x45\x49\x1e\x80\xc7\x02\x00\x00")

func localesRuLc_messagesBoggartMoBytes() ([]byte, error) {
	return bindataRead(
		_localesRuLc_messagesBoggartMo,
		"locales/ru/LC_MESSAGES/boggart.mo",
	)
}

func localesRuLc_messagesBoggartMo() (*asset, error) {
	bytes, err := localesRuLc_messagesBoggartMoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/boggart.mo", size: 711, mode: os.FileMode(420), modTime: time.Unix(1533730670, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _localesRuLc_messagesDevicesMo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x52\xcf\x4f\x2b\x55\x14\xfe\xe6\x4d\x9f\xfa\xc6\xe7\x8f\xbc\x68\xdc\x68\x72\x5d\x48\x34\x66\x70\xa6\x95\x84\x0c\x0c\xa0\x14\x12\x22\x55\x24\xd5\xfd\xa5\xbd\x94\x09\xed\x4c\x9d\x3b\x45\x4d\x30\x01\x8c\xc1\x48\x82\x09\xc1\x1f\x89\x51\x82\xec\xdc\x58\x09\x0d\xc5\x4a\x59\xba\x31\xf1\x5c\x37\xee\xfc\x1b\xfc\x13\xcc\x9d\x19\x30\xbc\xbb\x98\x73\xbe\x7b\xbe\xf3\x9d\xef\xdc\xcc\x3f\x0f\x0a\x5f\x01\xc0\x53\x00\x9e\x07\x70\x08\xe0\x39\x00\xf7\x0d\xa4\x67\xcd\x00\x1e\x05\x10\x18\xc0\xe3\x00\x3e\xc8\xf1\x27\x06\xf0\x18\x80\xcf\x72\xbc\x6f\x00\x8f\x00\xf8\xda\x00\xee\x02\xf8\x3e\xe7\x9f\x18\xc0\x1d\x00\x3f\x1b\x80\x05\xe0\xd4\x00\xee\x01\xb8\x34\x80\x02\x80\xdf\x73\xfe\x1f\x39\xfe\xcb\xc8\x66\xfe\x6d\x00\x4f\x03\x98\x31\xb3\xf8\xae\x09\x3c\x03\x60\xc5\x04\x1e\x00\xf8\x30\xbf\xff\xc2\x04\x9e\x04\xf0\x9d\x09\xdc\x07\xf0\x93\x99\xf9\x3f\x37\xb3\xb9\x7f\x9a\xc0\x0b\x5a\x37\xef\xfb\xd7\xcc\x7c\x9a\x85\xcc\x8f\x55\xc8\xf0\xb3\x05\x20\x5f\x39\x3d\x7a\x37\x33\xcf\xb5\xbe\xf6\x76\x2f\x7f\x27\x7d\xf4\x2c\xdd\x67\xe5\x58\xbf\x81\xde\xf7\x4e\x8e\xef\xe6\xf1\x89\x6b\xc1\x37\x6a\x49\x10\x85\x12\x65\x21\x6b\x71\xd0\xd6\x00\x65\xb1\x11\xd4\x84\x44\x39\x90\x7c\xa5\x29\xea\x98\x0b\xf3\xb8\x21\xc2\x44\x62\x3e\x88\x45\xfa\x95\x09\x5b\x0d\x62\x51\xc7\x42\x19\x8b\xfc\x06\x2d\x06\x32\x11\xa1\x88\x25\xde\xe6\x2d\x81\x2a\x97\xeb\x12\xd5\x8f\xdb\x02\xcb\xa2\x1d\xc5\x89\x5d\x91\x8d\xa0\x6e\xbf\xd9\x69\x48\xbb\x1a\x79\xac\x2e\x36\x66\xd6\x83\x35\xde\x8a\x46\xe3\x8e\xb5\xf4\x4e\xd5\x9e\x8d\x05\xd7\x56\xec\x32\x4f\x84\xc7\x8a\x8e\x3b\x6e\x3b\x25\xbb\x58\x62\xc5\x92\x37\x36\xf6\xaa\x53\x72\x1c\x4b\x4f\xb4\xab\x31\x0f\x65\x93\x27\x51\xec\xb1\xb7\x52\x0d\x56\xe9\xc4\xbc\x15\xd5\x23\x36\x79\x4b\x78\xca\x5a\xe4\x61\xa3\xc3\x1b\xc2\xae\x0a\xde\xf2\xd8\x0d\xf6\xd8\x72\x47\xca\x80\x87\x56\x65\xa1\x32\x67\xbf\x2f\x62\x19\x44\xa1\xc7\xdc\x51\xc7\x9a\x8d\xc2\x44\x84\x89\xad\xed\x7b\x2c\x11\x1f\x25\xaf\xb5\x9b\x3c\x08\x27\x58\x6d\x8d\xc7\x52\x24\xfe\x7b\xd5\x79\x7b\xfc\x7f\x9e\xf6\xb3\x2a\x62\x7b\x2e\xac\x45\xf5\x20\x6c\x78\xcc\x5a\x6a\x76\x62\xde\xb4\xe7\xa3\xb8\x25\x3d\x16\xb6\x53\x28\xfd\xd2\x04\xcb\x52\x3f\x7c\xc9\x75\x7c\xdf\x65\x23\x23\x4c\xa7\xce\x8b\xbe\xeb\xb2\x69\xe6\x30\x2f\xc5\x53\x7e\xf1\xba\x34\xe9\xbf\xae\xd3\x97\x53\xda\xa4\xeb\xb0\xcd\xcd\xac\x65\xca\x2f\x3a\xaf\xb0\x69\xe6\x32\x8f\x15\x27\x40\x87\xd4\xa3\x0b\xb5\xad\x76\xe8\x94\xfa\xea\x4b\xd0\x0f\x74\x45\x7d\xb5\x4d\x5d\xba\xa4\x3e\xf5\x40\x27\xba\xaa\xb6\x68\x78\xc3\xeb\x82\x0e\xd4\x1e\xfd\x4a\x03\xb5\xaf\x76\xa9\x47\x97\x34\x04\x1d\x3c\x7c\x71\x4c\x43\xfa\x45\xed\xa9\x9d\x4c\x58\xb7\x9c\xd3\x90\x4e\xd5\x1e\xe8\x88\x7a\x6a\x4b\xa7\x74\xc1\xd2\x90\x56\xf4\xff\x41\x47\x34\x54\xdb\x34\xa0\x1e\x9d\xa5\x16\x6e\x11\xe8\x98\x06\xea\x53\xf5\x39\x75\xd5\x0e\xf5\x68\x40\x7d\xd0\xb7\xf4\x9b\xd6\xff\x86\xba\x74\x46\x5d\xb5\x0b\xfa\x91\xfa\x74\x85\xff\x02\x00\x00\xff\xff\xc7\xff\x7a\x4d\x1c\x04\x00\x00")

func localesRuLc_messagesDevicesMoBytes() ([]byte, error) {
	return bindataRead(
		_localesRuLc_messagesDevicesMo,
		"locales/ru/LC_MESSAGES/devices.mo",
	)
}

func localesRuLc_messagesDevicesMo() (*asset, error) {
	bytes, err := localesRuLc_messagesDevicesMoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/devices.mo", size: 1052, mode: os.FileMode(420), modTime: time.Unix(1533730670, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/views/devices.html": templatesViewsDevicesHtml,
	"assets/js/devices.js": assetsJsDevicesJs,
	"locales/ru/LC_MESSAGES/boggart.mo": localesRuLc_messagesBoggartMo,
	"locales/ru/LC_MESSAGES/devices.mo": localesRuLc_messagesDevicesMo,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"assets": &bintree{nil, map[string]*bintree{
		"js": &bintree{nil, map[string]*bintree{
			"devices.js": &bintree{assetsJsDevicesJs, map[string]*bintree{}},
		}},
	}},
	"locales": &bintree{nil, map[string]*bintree{
		"ru": &bintree{nil, map[string]*bintree{
			"LC_MESSAGES": &bintree{nil, map[string]*bintree{
				"boggart.mo": &bintree{localesRuLc_messagesBoggartMo, map[string]*bintree{}},
				"devices.mo": &bintree{localesRuLc_messagesDevicesMo, map[string]*bintree{}},
			}},
		}},
	}},
	"templates": &bintree{nil, map[string]*bintree{
		"views": &bintree{nil, map[string]*bintree{
			"devices.html": &bintree{templatesViewsDevicesHtml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


func assetFS() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}
