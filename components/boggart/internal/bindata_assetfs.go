// Code generated by go-bindata.
// sources:
// templates/views/detect.html
// templates/views/devices.html
// templates/views/index.html
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

var _templatesViewsDetectHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func templatesViewsDetectHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesViewsDetectHtml,
		"templates/views/detect.html",
	)
}

func templatesViewsDetectHtml() (*asset, error) {
	bytes, err := templatesViewsDetectHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/views/detect.html", size: 0, mode: os.FileMode(420), modTime: time.Unix(1515662635, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesViewsDevicesHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x56\x4f\x73\xda\x4e\x0c\xbd\xe7\x53\x68\xf6\xf4\xfb\x1d\x60\x27\x39\xf5\xe0\x70\x4a\x3b\xe9\x0c\xed\xa1\xa5\xe7\x8e\xec\x95\xc3\x52\xb3\x76\x57\x8a\x4b\xc6\xc3\x77\xef\xd8\xc6\x64\x6b\x1b\x08\x29\xe9\x05\x6c\xad\x9e\x9e\xfe\x3c\xcd\xba\xaa\xc0\x50\x6a\x1d\x81\x4a\x72\x27\xe4\x44\xc1\x76\x7b\x75\x15\x19\x5b\x42\x92\x21\xf3\xad\xf2\xf9\x2f\x05\xd6\xdc\x2a\x43\xa5\x4d\x88\xd5\xec\x0a\x00\x20\x74\xd9\x7c\x2f\xd0\x51\xb6\x3b\x19\x9e\x8a\x95\x8c\x82\xd3\xc6\x63\x79\x33\xab\x2a\xb0\xd7\xef\x1c\xa8\xbb\x5d\x68\x98\xc2\x76\x1b\xe9\xe5\x4d\xcf\x37\x88\x96\x64\x84\x3e\xb5\x1b\x35\x8b\xb4\xb1\x65\x40\xd9\x7b\xfd\x23\x83\xae\xba\x5e\x5c\xc1\x38\xa3\xce\xab\x7d\x69\x7e\x27\x2c\xde\x16\x64\x7a\xfe\x2d\x66\x49\x68\xc6\xec\x7e\x68\xdc\x01\x9e\x0b\x5d\x3c\x15\xd4\x55\x29\xcb\x83\x80\x2e\xa5\xb5\x99\x24\x79\x36\xb9\x51\xcf\x11\x3e\xde\xbd\x00\x1f\x76\x96\x13\x6f\x0b\xb1\xb9\x3b\x9f\xf7\x3a\xe0\x5d\x20\xff\xe0\xbf\x0b\xf1\x55\x50\x1e\x8f\xc6\x88\xf4\x58\x17\x6b\xdf\x41\xcf\x23\xdd\x4c\x6a\x74\xfe\xbb\xc7\xdd\xdf\xb8\x9a\x33\xcb\x42\x8e\xfc\x25\xf5\x1c\xd9\xce\x27\x45\x48\x71\x62\x08\xd3\x5a\xa7\x76\x06\xfb\x26\xcc\xf7\xbc\xff\x5a\xec\x81\x4f\x2b\x72\x4f\x5c\xe4\x8e\x6d\xd9\xaf\x05\x5e\xb9\x1b\x70\x64\x3f\xe0\xd8\x8e\x40\x5f\xb6\xa7\x55\x3e\x80\x7c\xc6\xf5\xc9\xd5\x1a\x80\xde\x97\xe4\xe4\xa4\xac\x07\xb0\x0f\xd6\xd3\xab\x50\x2c\x90\x5a\x4f\xe6\x6c\xec\x1c\x5f\x0a\x1d\xdf\x21\x38\xb4\x47\x30\xb6\x4b\x30\x22\xb1\x43\xeb\x55\x55\x40\xce\x34\xd7\x46\x70\x9d\xd4\x44\xcd\x5d\x52\x03\xaa\x0a\x58\x50\x6c\x72\xbf\xf8\x34\x87\xff\xda\xe7\x6f\x5f\xe6\xa0\xb4\x41\x5e\xc6\x39\x7a\xa3\x91\x99\x84\x75\x49\xce\xe4\x9e\xb5\x41\xc1\x26\x2b\x9e\x3a\x92\x49\xcc\x3a\xe1\xd6\xba\x68\xad\x71\x9e\x0b\x8b\xc7\x62\xba\xb6\x6e\x9a\x30\x2b\x48\x31\x63\xfa\xff\x82\xac\xa9\xdd\x90\xa9\x4b\x21\xdf\x65\xd0\x98\xee\x1b\xd3\xf1\x14\xc6\xfb\xb2\xe2\x0b\x76\x45\xaf\x58\xaf\x7e\x3e\x92\x7f\x9a\x06\x8d\xa9\x73\x59\xbd\x45\x37\x62\xae\x09\x0f\x8e\xe0\x4d\x38\x83\x09\xf4\xc8\xc3\x41\x9c\x4d\x1f\xe7\x0f\x0f\xe8\xa5\x23\xaf\x23\xb7\x9f\x20\xbd\x30\xfb\x21\xfe\x0e\x00\x00\xff\xff\x47\x92\xb8\xb3\x2a\x09\x00\x00")

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

	info := bindataFileInfo{name: "templates/views/devices.html", size: 2346, mode: os.FileMode(420), modTime: time.Unix(1523310421, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesViewsIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\x5b\x6f\x1b\xc7\x15\x7e\xd7\xaf\x18\x10\x06\x24\x03\x21\x29\x51\x01\x8a\xc8\x92\x81\x38\xbd\x02\x49\x1b\x34\xa9\xf3\xb8\x18\xee\x1e\x72\xc7\xde\xdd\xd9\xcc\xcc\xf2\x02\x42\x00\x4d\x39\x4e\x03\xb7\x70\x90\x16\x68\x51\x14\x49\x51\xa0\x8f\x05\x18\xc5\x6a\x54\x59\x64\xfe\x40\x1f\xce\xfc\xa3\x62\x86\x4b\x8a\xa2\x48\x9a\x37\x3b\x92\xa5\x7d\xe1\xee\xce\xcc\x7e\xe7\x9c\xf9\xe6\x5c\x66\x97\x8d\x06\xf1\xa0\xc4\x22\x20\x19\x97\x47\x0a\x22\x95\x21\xfb\xfb\x6b\xbb\x1e\xab\x10\x37\xa0\x52\xee\x65\x04\xaf\x92\x22\x2f\x97\xa9\x50\x99\xbb\x6b\x84\x10\x32\xdc\xea\xf2\x20\x1b\x7a\xd9\xad\x42\xda\x36\xda\x3e\x74\x7b\xb4\xa9\xe6\x28\xa6\x02\x18\xe9\x61\x7b\xf9\x85\xbb\xf8\x0f\x3c\xd2\x8f\xf4\x81\x6e\xea\x47\xfa\xe9\x6e\xde\x2f\x8c\xe9\x37\x2c\x48\x00\x54\x94\x58\x2d\x73\x77\x37\xef\xb1\xca\x08\xec\x98\x5b\xe7\x24\xe9\x2b\x3f\x1d\x43\xf0\xea\x98\x1e\xe6\x68\x34\x08\x2b\x91\xdc\x7b\x3c\x2a\xb1\x72\xee\x1e\xe7\x01\xc9\xa4\x46\xcb\xc5\x49\x20\xa9\xc8\x41\x44\x8b\x01\x78\xd6\xbe\xe3\x1e\x31\xc6\xaa\xdb\xc4\x9c\xd4\x64\x76\xab\x40\xaa\xcc\x2b\x83\x4a\x7f\x1c\x45\x83\xa0\xee\x14\x79\x6d\x82\x3c\x17\x35\x8c\x69\x04\xc1\x94\xde\xf3\xd8\x64\xec\x38\xe6\xed\x65\x04\x48\x9e\x08\x17\xb2\x3e\x50\x95\x65\x51\x86\x48\x55\x0f\x60\x2f\xe3\x03\x2b\xfb\x6a\x67\x7b\x73\x33\xae\xdd\x99\xe1\x89\xe6\xf8\x55\x34\x53\xb7\x46\x83\xe4\x52\x53\x3b\x1e\x54\x98\x0b\x8e\x41\x77\x42\x50\x20\x9c\x9e\xf1\x1d\x05\x61\x0c\x82\xaa\x44\x80\xc3\x22\xc7\x85\x40\xb2\x44\xe6\xee\xd3\x20\x01\xb2\xbf\x4f\xde\x9b\x09\x6a\xb7\x28\x48\x7e\x36\xe1\x7f\x93\xa8\x57\x23\x3d\x4f\xd4\x6b\x10\xff\xa7\x10\x28\xfa\x6a\x14\xf0\xcc\xa3\x17\x51\x61\xcc\x22\x9e\xb3\xcb\x94\xe6\x69\x4d\x6f\xdc\xc2\xe4\x89\x5a\x6a\x65\xe2\x37\xd8\xc5\x13\x6c\xe3\xf7\xd8\xc6\x0e\x1e\xeb\x67\x44\xff\x11\x3b\x78\xa4\x9b\xf8\x1d\x1e\xe3\xf1\xea\x98\x03\x11\x88\x72\xdd\x29\xb3\x32\x75\x79\xc0\x05\x83\x21\xce\x0c\xdf\xbd\xa1\xcf\x14\xbd\x56\x4b\x1f\xbb\x82\x57\x4e\xa0\x16\x1e\xe1\x89\x3e\xd0\x5f\xe0\x11\x1e\xff\xef\x09\xc1\x53\xec\xea\x2f\xb0\x83\x5d\xfd\x48\xb7\x56\xc9\x29\x97\x47\x32\x09\x63\xc5\x78\xe4\xb8\x49\x91\xb9\xa6\x59\x80\x74\x62\x10\x8e\xcf\x13\x71\xc6\x30\xdb\x4c\x7a\xcd\x24\x06\x41\x4c\xf3\x0d\xd5\xa6\xe8\xb5\x32\xaa\x55\xa9\x02\x91\x75\x79\xe0\x2d\x47\xb5\x7f\x61\x17\x5f\x60\x17\x9f\x63\x07\xdb\xfa\x19\xc1\x43\x7b\xd1\x9e\x9b\x4f\x03\xf6\x04\x9e\x63\x85\x73\x5c\x1a\x53\x97\xa9\xfa\x39\x16\x8d\x27\xcf\x0d\x67\xa6\xe8\xb5\x62\xce\xf8\x7c\xc9\xf0\xf6\x27\xec\xea\xa6\x7e\xa6\x3f\x5f\x09\x61\x7c\xae\xde\x20\xbe\x34\x1a\x04\x22\x6f\x52\x11\x33\xad\x0e\x0a\x41\xb8\x89\xa8\x5f\xaf\x42\x08\x02\x70\x95\x60\x66\xd6\x97\xe3\xe4\xbf\xf1\x85\x0d\x90\x2d\xdd\xc4\x63\xfd\xb9\x2d\x8e\x5b\x86\x9a\x73\x64\xf5\x1f\x6f\xed\x9c\x23\x68\x3a\x23\x8e\xda\x32\x04\x55\x36\xf6\x39\x8a\x2b\x1a\x9c\x91\xf2\x93\xbc\x3f\x0f\x42\x61\x02\x42\x61\x65\x08\xdb\x13\x10\xb6\x57\x86\xf0\xf6\x04\x84\xb7\x57\x85\x70\x9f\x07\x8a\x96\x61\x67\x5c\xba\x32\x44\x98\x34\x6b\xe9\xe3\x17\x36\x37\x9d\x4a\x6f\xa4\xfd\x5d\x14\xfd\x5d\x5b\x8b\x2d\x04\x4f\xd3\xa1\xbd\x13\x38\x93\xe0\xdd\x39\xf0\x3f\xe4\x55\x10\x0b\x80\xc7\x66\x9c\x9d\x82\x21\xcd\xe7\xc0\xbd\x47\x95\x02\x51\x5f\x00\xb9\xd8\x1b\x39\xc1\xfa\xf7\x67\x94\xe1\xf2\x79\xec\x97\xef\x94\x0d\x5d\xa6\xa7\xe9\xcf\xa5\xd8\x22\xfc\x27\x1e\xf5\xbc\x22\x76\xf1\x14\x4f\xf5\x81\xa9\x25\x4c\x5d\xa1\x9f\x98\x2a\xf4\x2a\x6f\x1a\x4a\x5e\x52\x15\xe6\x01\xbf\x16\xe1\xd2\xe5\x61\x98\x44\xcc\xa5\xa6\x26\xcb\xb2\x48\x81\x88\x60\xc9\x24\xee\xaf\xd8\x31\x35\xa5\x6e\xda\x8d\x89\xd6\xa2\xe5\x63\x5f\x18\x27\x16\xdc\x4c\x88\x70\x06\x53\xe3\x14\x69\x40\x23\x17\x1c\x91\x14\x03\x90\x23\x31\x61\x26\x40\xdd\xd4\x07\xf8\x6d\xee\x0a\xfa\x0e\xf2\xb2\x6c\x8f\x17\x59\x00\xd7\x90\xbd\xb1\xcf\x23\x58\x8e\xba\x7f\xc7\x63\xec\xe8\x03\xdd\xd2\x4f\x09\xfe\x80\x5d\xa2\x5b\xd8\xb6\x79\xdf\x63\xec\x18\x5f\xa7\x0f\xcc\xfd\x17\xd8\x36\xdd\x16\x65\xb6\x15\xd4\x09\xa1\x4c\x4b\x3c\x72\x12\x09\x9e\x53\xe1\xa6\x21\x64\x51\xa2\x86\x0a\x93\xd9\x84\x3e\xed\x0b\xbd\xea\xbd\xe8\x8f\x3e\xf8\xe8\xf5\x5a\x41\x86\x73\xea\xae\x7f\xaf\x5b\x2f\x5f\xc3\x64\x3e\xb5\x47\x1d\xd8\xeb\xb5\xc1\xc0\xed\x95\x59\x99\x16\xeb\x73\xd3\xe1\x17\xf7\x56\x6f\x90\x2f\xf1\x50\xb7\xb0\x8b\x3f\xe8\xa6\xad\xc8\x4d\xe8\x37\x01\xff\xe8\x55\xe8\x1d\x0b\x1e\xf0\xa8\x6c\xd7\xf4\xe2\x46\x58\xa5\x0d\xae\x60\x0c\xb8\x46\x15\xff\xf9\x20\x50\x65\x25\xb6\x54\x0c\xf8\x84\xfd\x9c\x11\x3c\xc1\x17\x86\xdf\xc6\x11\xe8\xa7\x8b\xb2\x5c\xf0\xc4\x16\x34\xec\xa1\xe0\x8a\x3d\x74\x8c\x6c\x8e\x1b\x30\x88\xd4\xdc\x39\xcb\xe5\xe3\xe0\x55\xaf\x61\xbe\xc1\x2e\x7e\xd7\xdb\x60\xbc\xec\xf5\xca\x9b\xb0\x4c\xab\x40\x95\x0f\x22\xcb\x13\xe5\x71\x2e\x26\xad\xd1\xe1\x69\x21\xd8\xc1\x36\xd1\x07\x66\x2d\xea\x27\x78\x74\xb3\x53\xbf\xbc\xf9\x59\x34\xbb\xf5\x0f\x8d\x1f\x3c\x34\x69\x8f\x6e\xe1\xb1\x6e\xfe\x38\x33\x70\xd5\xfd\xcc\xd7\xfa\x33\xdd\xb4\x6f\x5b\x2f\xbd\x9f\x99\x96\x52\xb8\x34\x04\x41\x73\x3e\x7b\x58\x61\x92\xf1\x28\x27\x95\x00\x50\xd7\x22\xc5\x90\xe0\x26\x82\xa9\x7a\xb6\x0a\x45\x97\x86\xd9\x9e\xea\xcb\x95\x9a\x7f\xc3\x36\x9e\xda\x22\xe3\xa2\x9f\x9b\x3d\x35\xdf\x65\x61\x99\x48\xe1\xee\x65\xf2\xe9\x44\xe5\x7b\x13\x95\xef\xc9\x98\x8f\x05\x54\x18\x54\x07\xb2\x86\xcc\xe4\x48\x9e\xf2\x77\xb6\x7e\xf2\x4e\x5c\xbb\x13\xd2\x5a\x7a\x5d\xd8\x7e\x67\x36\xf1\x2f\x5f\x2a\x42\xe6\xe5\xae\x4f\x83\xe0\x5a\x32\xd7\x28\x9e\xdd\x5a\x1d\x73\x4d\x15\x7c\xa8\x3f\x4b\xdf\xf2\x77\xf1\xbf\x04\x9f\xe3\xa1\x6d\x3c\x5e\x09\x8d\x8d\xc0\xd7\x87\xc4\x6f\x2c\xe5\x0a\x13\x33\x8e\x73\x7e\xf0\x90\x60\x17\xbf\xb5\x5f\x22\x9d\x9a\xe4\xa3\xab\x9b\x86\x49\xf8\xdc\x9c\xe1\x91\xa5\x9b\x71\x92\x8f\x75\xcb\x44\xd3\x9b\x74\x70\xe9\x49\xb2\xdb\x2e\x2e\x0f\x27\xce\xcf\x9f\xed\x2b\x9e\xae\x7e\x8c\x5d\xec\xfc\x88\xab\x66\xa2\x73\x37\xf9\xac\xbc\x5e\xce\x1c\x22\x25\x68\xe4\x42\x76\x5a\x2e\x3f\x9b\x33\xff\x6a\xe0\xbc\xed\x17\x37\xa9\xf3\xd6\x7f\x58\x9b\x75\xd3\x83\x95\x08\x7c\x7a\x61\xeb\xc3\x08\xe6\x94\x63\xc6\x1d\x01\xe0\x39\xb2\xca\x94\xeb\x3b\x52\x51\x35\xf8\xf4\x78\x33\xb7\x39\xeb\x76\x9e\x95\xf4\x2f\xd8\xc6\x13\xdd\xd4\x4f\xcd\xda\x9f\x55\x3a\x08\xe4\xcc\x9b\x86\x16\xe5\x6b\xdd\x5a\x00\x65\xb2\x47\x1f\x3e\x6e\x62\xce\x4c\x23\x96\x59\x19\xf4\xd3\x84\x66\x5d\x2a\x26\xbe\x07\xc2\x2f\xd3\x9a\xf6\x91\x6e\xd9\xa0\xd2\xc5\xff\xdc\x04\x92\x15\x58\xbe\x08\x41\x30\xd1\xe8\x5f\x9d\x4f\x11\xbf\xb7\x5f\xf6\x99\x8b\x93\x2b\xb6\xa3\x70\xb6\x94\xd6\x86\xfe\xc1\xe5\x03\xed\x45\x9e\x5d\xab\x3f\x51\xf5\x18\xf6\x32\x0a\x6a\x2a\xef\x4a\x99\x9a\xb3\xef\x22\x49\x2e\x9d\x1a\xd2\x18\xe0\xc5\xd4\xf3\x58\x54\xde\x21\x9b\x77\xec\xbd\xfd\xb5\x91\x21\x29\x17\xc6\x8c\x28\xc4\xb5\xfe\x98\xdd\xbc\x85\x9f\x24\xe5\x03\x39\x88\x8e\x8d\x06\x31\xbe\x98\xb9\xbf\xfc\xf8\x83\xf7\xc9\x46\xef\xfc\x77\xbf\x7d\x9f\x64\xf2\x1e\x95\x7e\x91\x53\xe1\xe5\xa9\x94\xa0\x64\xbe\x02\x91\xc7\x85\xcc\x83\xeb\x53\xa1\x64\xfe\xc1\xe0\x34\x17\xb2\x28\x67\x9e\x5a\xa2\x81\x84\xdb\x16\x8f\x9c\x45\x6b\x10\x82\x0b\xd9\x87\xdc\x95\xae\x60\xb1\x1a\xb6\xcd\x03\x5a\xa1\xbd\xbb\x43\x8c\xbb\xb5\xe1\x71\x37\x09\x21\x52\xb7\x73\x02\xa8\x57\xdf\x28\x25\x91\xab\x18\x8f\xc8\xc6\xed\x21\x03\xa4\x40\x82\x46\x65\x20\xb7\xe0\x2d\x72\xcb\xe2\x91\x9d\xbd\x51\xe4\xfe\x11\x41\x95\x7c\xf8\x6b\xae\x58\xa9\xbe\xd1\xb8\x40\x0a\xbb\x87\xb3\x43\xd6\x7f\x66\xc6\xae\xbf\x75\xb1\x1d\x6a\x6a\x87\xac\x37\x1a\x7d\xa0\xfd\xfd\x71\xbd\xea\xb1\x79\x08\x4c\x78\x88\xcf\x3c\xd8\xe9\x99\xeb\x62\xa3\x99\x3c\x3b\xa5\xeb\x45\xce\x95\x54\x82\xc6\xdb\xeb\xe7\x7a\xed\xdf\xbe\x33\xaa\xff\x88\x63\xef\xf7\xd8\xcd\xf7\x0c\x3b\xcc\x85\xb3\xb3\xff\x07\x00\x00\xff\xff\x15\xaa\x6e\xe9\x7d\x38\x00\x00")

func templatesViewsIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesViewsIndexHtml,
		"templates/views/index.html",
	)
}

func templatesViewsIndexHtml() (*asset, error) {
	bytes, err := templatesViewsIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/views/index.html", size: 14461, mode: os.FileMode(420), modTime: time.Unix(1518814225, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsJsDevicesJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x57\x4f\x6f\xe2\x3a\x10\xbf\xf7\x53\xcc\x4b\x2b\x39\xa8\x40\xf4\x6e\x4f\x40\x78\x87\xed\x1e\x2a\xad\xb4\x2b\x95\x5b\x5b\x21\x27\x36\xa9\x5b\x63\x23\xdb\x81\xa2\x2a\xdf\x7d\x35\xf9\x57\xc2\xff\x0a\x56\xbb\xd5\xfa\x80\xec\x89\xe7\xe7\x99\xf9\x8d\xc7\xc3\x95\xcf\x74\x9c\x4e\xb9\x72\xad\xae\xe1\x94\x2d\xfd\x49\xaa\x62\x27\xb4\x02\xbf\x05\x6f\x17\x00\x00\x73\x6a\xc0\xd1\x48\xf2\x1b\x3e\x17\x31\xb7\x10\xc2\x95\x4f\x2e\x59\xb9\xca\x3f\x91\x56\xbe\x15\x47\xf7\x86\x3a\x3a\x42\xa1\xff\x56\x0b\x71\x48\xaa\x92\x94\x26\xbc\x07\x4d\x39\x8e\xd4\xc8\x1e\x90\x80\x51\xfb\x14\x69\x6a\x58\xc0\xa8\xa3\x39\xb2\x0d\xc4\xbf\xff\xa9\xee\xb3\xd5\x8a\x34\xd4\xb2\x76\x63\x49\x9f\xe9\xeb\x1e\xe4\x48\x27\x09\x35\x2e\x28\xad\x0e\xfe\xe7\xca\x09\xb7\x0c\xcb\x35\x69\x6f\x28\xa2\x05\x77\x26\xee\x01\xc1\xd9\xde\xb3\x63\x2d\xd3\xa9\xb2\x3d\xb8\xdf\x40\xd9\x34\xa8\xc2\xee\x01\x71\xcb\xd9\xd6\xa3\x71\x18\xae\x18\x37\x3d\x78\xe7\x23\xdf\xdd\xda\x81\x08\x25\x51\xb1\x56\x8e\x2b\x07\x21\x10\xd2\xbf\xd8\xb9\x75\xa2\x0d\xf8\xb8\x5f\x80\x50\x70\x10\xb9\x70\xb2\x40\xbe\x0e\x81\x0c\xec\x8c\x2a\x88\x25\xb5\x36\xf4\x24\x8d\xb8\x84\xfc\xb7\x63\xd3\x38\xe6\xd6\x7a\x43\x02\xd7\x05\xec\xbd\x78\x84\x6b\x20\x83\x00\x55\x86\x40\xfa\x3b\x0f\xc9\x76\x9b\x6b\xb8\x4b\x8d\xaa\x4c\xd8\x0e\x91\x6d\x48\xb3\xcd\xc8\xee\xa5\x43\x30\x72\x32\x06\xe3\x36\x36\x62\x86\x8c\x9d\x0e\xe6\xa8\x7d\xb1\xe3\x58\xa7\xca\x9d\x0e\xc6\x15\x5e\x27\x76\x74\xba\xa1\x5a\x3b\xe7\xb0\x0d\x46\x2f\x8e\x4c\xbc\x3d\x39\x27\x26\xe0\x1b\xbd\xe8\x96\x86\x1c\x9b\x6f\x98\x6e\x51\xea\x9c\x2e\xf2\x34\xf4\x8a\x85\x57\xa5\x5f\xe4\x14\x44\x4e\x75\x18\x55\x09\x37\xf9\x54\xc4\xf8\x1d\x1d\xe8\x38\x9d\x24\x92\x87\xde\x54\x33\x2a\x2b\x19\x35\x09\x77\xa1\x77\xb9\x2a\xcc\xe7\x1d\x27\x1c\xee\xfe\xa2\xd5\x44\x98\x29\x30\x61\xd1\x56\x28\xaa\x04\x5c\x62\x56\xa3\x0b\x82\x61\x4e\x37\x34\x63\x2a\x65\x44\xe3\x97\xd0\x2b\x36\x8f\xf2\x83\xfd\x07\x52\xe9\x18\x9e\x08\xeb\xb8\x19\x17\xca\x0f\xa4\xd5\xcf\xaf\xc9\xde\x20\xe0\x20\x03\x51\xf9\x9a\xc8\xe5\xec\x09\xbd\x83\x7a\xd6\x31\x7c\xaa\xe7\xdc\x83\xd2\xf4\x9b\x86\xc9\xde\x70\x10\x88\xe1\x66\xf2\x54\x23\x03\x2e\x2d\x3f\x33\x13\x65\x11\x38\x2f\x15\x45\xd6\xfc\xd1\x4c\xe8\x97\x9a\x85\xaf\xea\x23\x24\x1c\x2c\x7c\x64\xc0\xc4\x7c\x25\xcc\x9d\xc4\xe8\x74\x06\xf5\xac\xf3\x5a\x96\xdc\xba\x48\x63\xc5\x2d\xe8\x19\x0e\x02\x26\xe6\xc3\x1d\x85\x77\x4b\xd5\x6c\x48\x1e\xeb\x55\xd6\x2a\xef\x76\xdd\x0f\x7c\xc3\x28\x2a\x6e\xaa\x8e\x40\xd6\xeb\xcf\xdd\x13\xd4\x7e\xfc\x9e\xae\xe0\x1c\xcf\x90\xa2\x53\x7e\x86\x27\x63\xce\x95\x3b\xbe\x41\x41\xad\x73\xf7\x27\xb9\x09\xb7\x0c\xbb\x94\x43\xf0\x70\x6c\x93\x22\xd4\x44\x17\xd7\x05\x21\xef\xcb\x23\x3e\x4d\x9b\x32\x11\x86\xdb\xd3\xc9\x45\x98\xf1\x44\x18\xeb\x3e\x42\x30\xdf\xc7\x00\xbe\xf0\xff\x1c\xda\x04\x2b\x65\xed\xa4\x30\xe3\x41\x23\x7d\xe7\x8c\x50\x49\x61\xda\xaf\x8c\xf8\x58\xd2\xbf\x31\x52\x87\x1e\x83\x85\x50\x4c\x2f\xba\xab\x2f\x2d\x84\x75\x20\xfc\x42\x7e\xdb\x68\xf6\xae\xba\x58\x97\xd7\x1e\x02\x6c\x28\x7a\x40\x7e\x7c\xbf\x1b\xad\x05\x79\x47\xb1\xce\xaf\x6f\x89\x8e\x37\x37\x28\xda\x8b\x35\xe5\xb2\x15\x79\xa7\xc6\xdf\x16\xf0\xd5\x7f\xb7\xb9\x71\x5d\xc3\xa5\xa6\xcc\x5f\x8b\x53\xd6\xf0\xbf\x90\x64\xad\xfe\xcf\x00\x00\x00\xff\xff\x99\xc8\x49\xf8\x45\x0f\x00\x00")

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

	info := bindataFileInfo{name: "assets/js/devices.js", size: 3909, mode: os.FileMode(420), modTime: time.Unix(1523320621, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _localesRuLc_messagesBoggartMo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8e\xcd\x6a\x13\x51\x14\xc7\xff\x8d\xd5\x45\x14\xbf\x70\xe9\xe2\x74\x61\x51\xe4\xd6\x99\xc4\x42\x99\x76\x5a\xb1\x1f\x20\x36\x58\x6a\x74\xe1\xee\x9a\x5c\x93\xc1\xcc\xbd\xe1\xce\x9d\xa2\xd0\x45\xad\xe0\x46\xc5\x8d\xee\x04\xc1\x6e\x74\x57\xaa\xc1\xd6\x36\x71\xe3\x03\x9c\x79\x01\xc1\x47\xf0\x0d\x24\xc9\xd4\xe2\xd9\x9c\xdf\xef\xde\x73\xfe\x9c\x5f\xe7\x47\xdf\x01\xc0\x09\x00\x17\x01\x2c\x03\x38\x09\x60\x1d\xc3\xfa\x04\xe0\x34\x80\xcf\xf9\xfb\x2e\x80\x53\x00\x7e\x02\x38\x03\x20\xcb\xfb\x6f\x00\xcd\x11\xe0\x0f\x80\x31\x00\x0f\x0a\xc0\xd9\x7e\x4e\x01\xb8\x00\xe0\x75\xde\x3f\x16\x80\x73\x00\xbe\x15\x80\x11\x1c\xd5\x31\x00\xa3\x39\x17\xf2\x7e\x3c\xbf\x6b\x50\xb1\xd2\xe9\xe8\x82\x4c\x9a\x0f\x8d\xb4\xf5\x5c\x95\x53\x35\x77\xc8\x6b\x51\x4d\x25\x43\xa9\x18\x1d\x39\x63\x23\xdd\x18\xfa\xdd\x58\x5a\x47\x4d\x13\x2b\xac\xaa\xb6\xb1\x4e\x54\x92\x46\x54\x17\x37\xd3\x46\x22\xaa\x26\xa0\xba\x5a\xbb\xf1\x38\x6a\xca\xd8\x4c\xd8\xb4\xb8\x72\xa7\x2a\xe6\xad\x92\x2e\x32\x5a\x2c\x48\xa7\x02\x2a\x79\xfe\x94\xf0\xca\xa2\x54\xa6\x52\x39\x98\x9c\xbc\xea\x95\x3d\xaf\xb8\x2c\x13\x27\xaa\x56\xea\xa4\x25\x9d\xb1\x01\xdd\x1e\x64\x50\x25\xb5\x32\x36\x75\x43\x33\xff\x05\xcf\x16\x97\xa5\x6e\xa4\xb2\xa1\x44\x55\xc9\x38\xa0\x7f\x1e\xd0\x6a\x9a\x24\x91\xd4\xc5\xca\xad\xca\xa2\xb8\xaf\x6c\x12\x19\x1d\x90\x3f\xe1\x15\xe7\x8d\x76\x4a\x3b\x51\x7d\xda\x56\x01\x39\xf5\xc4\x5d\x6b\xb7\x64\xa4\xa7\xa9\xd6\x94\x36\x51\x2e\xbc\x57\x5d\x12\x53\x47\x73\xfd\x7b\x1e\x29\x2b\x16\x75\xcd\xd4\x23\xdd\x08\xa8\xb8\xd2\x4a\xad\x6c\x89\x25\x63\xe3\x24\x20\xdd\x1e\x68\x12\x96\xa7\x69\x88\xa1\xbe\xe4\x7b\x61\xe8\xd3\xf8\x38\xf5\xd1\x1b\x0b\x7d\x9f\xe6\xc8\xa3\x60\xe0\xb3\x61\xe9\xf0\x6b\x26\xbc\xde\xc7\xcb\x83\xb1\x19\xdf\xa3\xf5\xf5\xe1\xca\x6c\x58\xf2\xae\xd0\x1c\xf9\x14\x50\x69\x1a\xfc\x81\xb7\xb9\xcb\x1d\xde\xcf\x5e\x51\xf6\x9c\x7f\x64\x1b\xbc\xcd\x3b\xbc\xcf\x1d\xee\xf2\x6e\xf6\x06\xfc\x96\x3b\xd9\x26\x77\xf8\x7b\xf6\x62\xf8\xb0\x95\x3d\xcb\x36\xb3\x0d\xee\xf1\x5e\x9f\x78\x87\xb7\xc1\xef\xb9\x37\x58\xd8\xe4\x5e\xb6\xc1\xbb\xdc\xe5\x2f\xe0\x2d\x3e\xe0\x6e\xf6\x92\xf7\x88\xbf\x72\x8f\x0f\xf0\x37\x00\x00\xff\xff\xc8\xe2\x62\x17\xc9\x02\x00\x00")

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

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/boggart.mo", size: 713, mode: os.FileMode(420), modTime: time.Unix(1523796281, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _localesRuLc_messagesDevicesMo = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x52\x41\x6b\x2b\x55\x14\xfe\xe6\x4d\x9e\xbe\x17\xeb\xf3\xf1\x10\xdc\x28\x5c\x17\x3e\x14\x99\xe7\x4c\x62\xa1\x4c\x3b\xad\xd8\xb4\x50\x6c\xb4\xd4\xe8\xfe\x9a\xdc\xa6\x43\x93\x99\x30\x77\xa6\x54\xa8\xd0\xa4\x60\x95\x42\x85\x62\xb1\x20\x58\xaa\x88\xe0\x2a\x96\x06\x53\x63\x53\x17\xee\xdc\x9c\xbb\x71\xd9\xff\xe0\xce\xa5\xdc\x99\xb1\x52\xef\x62\xce\xf7\x9d\x7b\xce\x77\xbe\x33\xdc\xeb\x47\x85\x23\x00\x78\x0e\xc0\x8b\x00\xbe\x04\xf0\x02\x80\x09\x03\xe9\x59\x37\x80\x67\x00\xf8\x06\xf0\x34\x80\xc4\x00\xee\x01\xf8\x24\xe7\x9f\x1a\xc0\x53\x00\x0e\x0c\xe0\x2e\x80\xa3\xbc\xfe\x6b\x03\xb8\x03\xe0\x07\x03\x28\x02\xf8\xd1\x00\xee\x03\x18\x18\x40\x01\xc0\x55\xde\xf7\x7b\xde\xf7\x47\x9e\x57\x46\x36\xf3\x4f\x03\x78\x08\x60\xce\x04\x9e\x07\xb0\x62\x02\x8f\x00\x08\x33\xcb\x6f\x99\xc0\x03\x00\x9f\x9b\xc0\x04\x80\x63\x33\xf3\xfd\xbd\x99\xcd\xbd\x34\x81\x97\x00\xfc\x96\xf7\x5d\x9b\xd9\xbc\xbf\xf2\xfa\xbf\xcd\xcc\xd7\xdd\x42\x96\x7f\x58\x00\xf4\xca\xcf\x66\x6b\xa7\xbb\xdd\xc9\xf1\x83\x3c\xde\xcb\xff\x93\x3e\x7a\x47\xed\xfb\x7e\xce\xb5\x86\xd6\x33\x73\x5e\xc8\xe3\x44\x1e\x51\x11\xb2\x1e\xf9\x9d\xd8\x0f\x03\x54\xc4\xa6\x5f\x17\x12\x15\x5f\xf2\x8f\x5a\xa2\x81\x85\x20\x8f\x9b\x22\x88\x25\x16\xfd\x48\xa4\x5f\x19\xb3\x35\x3f\x12\x0d\x2c\x55\xb0\xcc\x6f\xd8\xb2\x2f\x63\x11\x88\x48\xe2\x5d\xde\x16\x78\x3f\xe6\x71\x22\x51\xe3\x72\x43\xa2\xf6\x71\x47\x60\x55\x74\xc2\x28\xb6\xaa\xb2\xe9\x37\xac\xb7\x93\xa6\xb4\x6a\xa1\xcb\x1a\x62\xf3\xad\x0d\x7f\x9d\xb7\xc3\x27\x51\x52\x5c\x79\xaf\x66\xcd\x47\x82\x6b\x47\x56\x85\xc7\xc2\x65\x25\xdb\x99\xb2\xec\xb2\x55\x2a\xb3\x52\xd9\x9d\x9c\x7c\xdd\x2e\xdb\x76\x51\x0f\xb6\x6a\x11\x0f\x64\x8b\xc7\x61\xe4\xb2\x77\x52\x0d\x56\x4d\x22\xde\x0e\x1b\x21\x9b\xb9\x25\x3c\x5b\x5c\xe6\x41\x33\xe1\x4d\x61\xd5\x04\x6f\xbb\xec\x86\xbb\x6c\x35\x91\xd2\xe7\x41\xb1\xba\x54\x5d\xb0\x3e\x14\x91\xf4\xc3\xc0\x65\xce\x13\xbb\x38\x1f\x06\xb1\x08\x62\x4b\xdb\x77\x59\x2c\xb6\xe2\x37\x3a\x2d\xee\x07\xd3\xac\xbe\xce\x23\x29\x62\xef\x83\xda\xa2\x35\xf5\x5f\x9d\xf6\xb3\x26\x22\x6b\x21\xa8\x87\x0d\x3f\x68\xba\xac\xb8\xd2\x4a\x22\xde\xb2\x16\xc3\xa8\x2d\x5d\x16\x74\x52\x2a\xbd\xf2\x34\xcb\xa0\x17\xbc\xe2\xd8\x9e\xe7\xb0\xc7\x8f\x99\x86\xf6\xcb\x9e\xe3\xb0\x39\x66\x33\x37\xe5\xb3\x5e\xe9\xdf\xab\x19\xef\x4d\x0d\x5f\x4d\xcb\x66\x1c\x9b\x6d\x6f\x67\x2d\xb3\x5e\xc9\x7e\x8d\xcd\x31\x87\xb9\xac\x34\x0d\xfa\x86\xae\x68\xa8\xba\xd4\xa7\x4b\x1a\xd2\x00\xf4\x9d\xea\xaa\x9e\xda\xa1\x31\x5d\x68\x44\x67\xd4\x07\x1d\xaa\x7d\xfa\x85\x46\xea\x40\xed\xd1\x80\x2e\x69\x0c\x3a\xfc\x7f\xe2\x94\xc6\xf4\x93\xda\x57\x3d\x1a\xaa\x2f\xb2\x96\x9f\x69\x4c\x67\x6a\x1f\x74\x42\x03\xb5\xa3\x21\x5d\xb0\x34\xa4\x37\xfa\x55\xd0\x09\x8d\x55\x97\x46\x34\xa0\xf3\xd4\xc2\xad\x02\x3a\xa5\x91\xda\x55\x9f\x51\x5f\xf5\x68\x40\x23\x1a\x82\x8e\xe9\x57\xad\x7f\xaa\x7a\x3a\xab\x76\x55\x17\xf4\x15\xf5\xe9\x9c\xfa\x6a\x0f\xf4\x2d\x0d\xe9\x0a\xff\x04\x00\x00\xff\xff\x1d\xef\x77\xe9\x17\x04\x00\x00")

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

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/devices.mo", size: 1047, mode: os.FileMode(420), modTime: time.Unix(1523796281, 0)}
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
	"templates/views/detect.html": templatesViewsDetectHtml,
	"templates/views/devices.html": templatesViewsDevicesHtml,
	"templates/views/index.html": templatesViewsIndexHtml,
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
			"detect.html": &bintree{templatesViewsDetectHtml, map[string]*bintree{}},
			"devices.html": &bintree{templatesViewsDevicesHtml, map[string]*bintree{}},
			"index.html": &bintree{templatesViewsIndexHtml, map[string]*bintree{}},
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
