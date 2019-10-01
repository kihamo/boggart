// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
// DO NOT EDIT!

package xmeye

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

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\xdd\x6f\xdb\x38\x12\x7f\xef\x5f\x31\xc7\xe6\x4e\x36\x50\x49\x4e\xf6\xe3\x0e\x5a\x3b\x8b\xdb\x4b\x17\x7d\x48\xbb\x41\x93\x1e\x90\x7b\x29\x68\x89\xb2\x99\xa5\x48\x2d\x49\xf9\xe3\x0c\xff\xef\x07\xea\xcb\xb2\x2d\xd9\x92\x9d\x6c\xaf\xd8\xf5\x43\x4b\x48\xc3\x1f\x67\x86\xbf\x99\x21\x47\x59\xad\x20\x20\x21\xe5\x04\x90\x2f\xb8\x26\x5c\x23\x58\xaf\x5f\x0d\x03\x3a\x03\x9f\x61\xa5\x46\x28\xc6\x13\x62\x6b\xaa\x19\x41\xd7\xaf\x00\x00\xaa\x2f\x43\x21\x23\x7b\x22\x45\x12\x43\x9c\x30\x66\x4b\x3a\x99\xea\x5c\x6e\x57\x96\xf2\x38\xd1\x99\x70\x45\x22\x95\x62\x78\x4c\x58\x21\x37\xd6\xbc\x56\x2a\x95\xc4\x30\x95\x24\x1c\xa1\x1f\xb1\xaf\xa9\xe0\x23\x5f\xf0\x90\x4e\x94\x4d\x16\xb1\x90\x1a\x55\x30\xc0\xe0\x50\x5f\x64\x03\x95\xf8\x3e\x51\x2a\x1b\x47\x35\xc8\x29\x3a\x2d\xed\xc2\x0a\x42\x6c\xfb\x62\x82\x20\xb5\x7d\x84\x56\x2b\xa0\x97\xff\xe0\x80\xde\xa6\x4b\x41\xb6\x72\x22\xb1\xd1\x03\x81\x03\xeb\x35\xba\x1e\xba\xb4\x46\x69\x17\xb7\xb0\x84\x89\x17\x33\x63\x4a\x70\xa0\x88\x6e\x34\xc5\x2c\xdd\xcd\x82\xa1\x9b\x6e\x59\x65\xa3\xdd\x80\xce\x72\x7e\x64\xc3\xfc\xbf\x2d\x2e\xf9\x8c\x60\x19\xd2\x85\x59\x66\xff\xad\x14\xf3\x1a\x8a\xf9\x82\xd9\x51\x60\x5f\x5e\x81\x19\xa9\xa8\x18\x2d\x94\x7d\x79\xd5\x40\xb5\xc5\xe7\x18\x73\xc2\x76\x69\xb6\x25\x51\xf0\xbd\xc6\x5a\x23\x27\x85\xf1\x94\xc6\xe3\x3a\xa4\x52\x32\x29\x59\xcb\xf1\x0c\x38\x9e\xd9\x1a\x8f\x15\x8c\xb1\xfc\x6c\x06\x68\x03\xc3\xa8\xaa\x5b\xab\x44\x62\x34\x97\x8d\x25\x51\x84\xeb\x8c\x56\x66\xa7\x42\x20\xbf\x81\x93\xb1\x04\x90\x89\xcf\x62\x4d\xf3\x6c\x46\x8c\x10\xe1\x01\xac\xd7\xd7\x25\xa9\xea\xe7\x3d\xe1\x19\x56\xbe\xa4\xb1\xf6\x66\x82\x06\xbd\x41\xff\x07\x33\x97\x29\x02\xeb\xf5\x6a\x05\xce\x47\xf2\x5b\x42\x94\x76\x3e\x7d\xbc\x75\xee\xb0\x9e\x66\x8f\x33\x70\x74\x5d\xd2\xe6\x7e\xa9\x34\x89\x32\xc6\x18\x6a\x0c\x5d\x56\x43\x9a\xce\xa6\x65\x34\x3c\xcd\xbc\x62\xee\x29\x26\x56\x23\xb0\xce\xdc\xdb\x32\x3c\x9e\xd1\xd8\x90\x32\x72\xb2\xb5\xe5\xe4\x73\xcc\x4d\x41\xea\xec\xfd\x39\x43\x3f\x6e\xf0\xd0\x4d\x58\xc3\x9b\x4a\xa8\x69\x3c\xb6\x9b\x83\x6d\x6b\xc6\x4e\xd0\x55\x11\xcc\x13\x08\x71\x40\x20\x73\x14\x50\x7e\x00\xad\x99\x20\xcd\x0a\x98\x20\x25\x95\x25\x19\x81\xf4\x5f\x5b\x69\x49\x63\x12\x40\x80\x35\xce\x9e\x07\xda\x96\x44\xc5\x82\x2b\xa3\x09\x17\x73\x89\x63\x04\x4a\x2f\x8d\xfa\x73\x1a\xe8\xa9\x77\x39\x18\xfc\xf5\x80\x82\xd9\x8a\x26\x31\x1f\x96\xc9\xe4\xe4\x71\xa1\x1c\x70\xb3\x8b\xaf\x8b\x1d\xd4\xd3\x13\x66\x3f\xd0\x88\x9c\x05\xf0\x9e\x28\x85\x27\xe7\x61\x3c\x2c\xe3\xf3\x00\x3e\x29\x22\xbb\x00\x0c\xdd\x63\xae\x36\x38\x47\x37\x6d\xa8\xc7\x22\x58\x1e\x5f\x6e\xb5\x02\x89\xf9\x84\xc0\x05\x7d\x03\x17\x4c\x4c\xc0\x1b\x81\x63\x98\x7a\x88\xa8\x9b\x55\x5a\xd3\x22\x30\x2e\x31\x0b\x38\x77\x42\xd1\x34\x20\x52\x97\xb4\x20\x5f\x01\x30\xcc\xd2\x0c\xe8\x65\x4c\x46\x08\xc7\x31\xa3\x7e\x9a\xdb\xdc\x4d\x0e\x42\xd7\x81\xf0\x93\x88\x70\xed\xcc\x25\xd5\xa4\x17\x60\x4d\x1e\xc4\xbd\x96\x94\x4f\x7a\x56\xa1\x83\xe1\x96\xf3\xb3\x90\x11\xd6\x80\xae\x06\x83\xef\xed\xc1\xa5\x3d\xb8\x7a\xb8\xfc\xce\x1b\x7c\xeb\x0d\xbe\xfb\xcf\xe0\xef\xde\x60\x60\x82\xd5\xea\xf7\x87\x6e\x86\x7d\xdd\x4d\xdd\x62\xad\x1b\xac\x71\x67\x5b\x4b\x45\x97\x31\x39\x79\xb2\xe1\x5e\xeb\xc9\xc7\x89\x07\x19\x5f\xb2\x64\x7d\x8c\xa2\x87\xe9\x37\x74\xd3\x2c\x76\x30\x7b\xa6\x25\xa4\xa9\xea\xfc\x61\x72\xe8\x07\x7c\x66\x16\xfc\x89\x4c\x28\x07\x7d\x6e\x32\x7d\xcb\x83\xf3\x41\xee\xe9\x7f\xcf\x03\xf8\x67\xca\x03\xf5\xd5\x24\x54\xc3\xd7\x34\xa3\xa6\xc4\x7d\x91\x94\x6a\x90\x1d\x73\x5e\x32\x4c\xf9\x52\x39\x35\x55\x22\x65\xda\xef\x91\x59\x9f\x53\xe9\xb7\x3c\xf8\xbd\x8a\x41\xb9\x53\xb7\x84\x4f\xd2\x73\x70\x27\x8c\x32\xa7\x91\x85\xb6\x7d\xc2\x35\x91\x47\x32\xd3\x16\x40\xe5\x38\x5c\xb6\x37\xa0\x1c\xd9\x0b\xd5\x01\x0c\xaa\xdd\x03\x63\xda\xa1\x43\x7e\x20\xe6\x9c\x09\x1c\xfc\x8d\xe3\x88\x8c\xea\x28\xbb\xd7\x6c\xa8\xf6\x18\xa8\x2f\x38\xba\xde\x6b\x28\x18\x0c\xbb\x80\xde\x6f\x2b\xdc\x94\x6f\x2e\xca\x9e\x42\x6d\x13\xa4\xd1\xbc\x4d\x37\xe1\xb8\xe8\xd7\x5b\x5f\x0f\x55\xd2\x50\xc8\x28\xbf\x0f\x99\x21\x82\x7c\x43\xd1\x8f\x08\x22\xa2\xa7\x22\x18\xa1\x58\x28\x8d\x80\x06\x23\xa4\x88\xd6\x94\x9b\x1b\x8e\x29\xb0\xb6\x16\x93\x89\x99\x39\xc3\x8c\x06\x58\x8b\x63\x5c\xad\xab\xdb\xcf\x53\x9f\xe1\x85\x6a\x74\xa1\x6b\x14\xd8\xbe\x60\xf6\x65\xe5\xfa\xfa\x4b\xbc\xe9\xc8\x9d\x58\xe9\xfe\x8d\x59\xd2\xa9\x56\xb6\xe3\x57\xab\x5a\x07\x1d\xea\x1d\x94\x57\x5d\xc7\x1c\x0f\x3e\xfb\x89\x94\x84\xeb\x36\xa5\x0e\xba\x38\x1c\xf2\x24\x28\xc5\x5c\xc5\x98\x8f\xd0\x15\x6a\xb8\x2a\xb6\xcc\xa8\x39\x60\xc7\xa4\xb7\xd3\x15\x5c\x28\xfb\xdb\x8e\x79\x33\x85\x49\x1b\xd1\x79\xfd\x32\xf9\x1c\x6d\x75\xb3\x7d\xc1\xb5\x14\x0c\x41\x9a\x31\x91\xf1\xab\x9d\xfb\x15\xc1\xcc\x10\x23\xcd\x74\x5b\x0e\xdf\x2a\x60\xce\xe0\xd2\x19\x5c\x41\x5e\xc0\xbe\x47\x69\x8e\x35\x41\xba\x0d\x25\x09\x0e\x04\x67\x4b\x70\x3b\x7a\xa1\x7d\x6e\x84\xd6\xf9\x11\x5a\x73\x18\x4e\xe0\xcd\x97\xd9\xe6\x56\x9f\x1e\x5a\xa3\x75\x26\xcd\x33\x90\x05\x75\xe5\x46\xa9\xad\x89\xd2\x1a\xe3\xed\xb1\x3e\xd4\x35\x3b\x0a\x3b\x4e\xb4\x16\x7c\xef\xeb\x04\x0f\x05\xca\x5d\x93\x49\x54\xf8\xce\xc5\xbc\x92\x2b\x3e\x88\x79\x91\x2a\x32\xc9\x13\xed\x73\x8d\x81\x27\x70\xa2\x5b\xec\x9c\x30\xe5\x25\xc2\xad\xc5\x71\x61\x0b\xb8\x6b\x74\xfe\x5f\x64\xf2\x6f\xfe\x0c\xf1\x6e\xda\xfe\x19\xe2\x0d\xf3\xbe\xd6\x10\xcf\x6e\x1e\x6d\x0e\x91\xc7\x0f\x87\x47\xaf\x20\x43\xd7\x90\xfa\xf0\x15\xe5\xf0\x55\xe8\x80\xd7\x1a\x5e\xd5\x3c\xde\x79\xd4\xf8\x15\xb9\xe2\x9f\xca\x9f\x2c\x98\xc3\x74\xd9\x8f\xac\xff\xe4\xf9\xaa\x62\x90\xd2\x58\x53\xff\xdd\xc3\xfb\x5b\xe8\x65\xe3\x4f\x1f\x6f\x01\xb9\x01\x56\xd3\xb1\xc0\x32\x70\xb1\x52\x44\x2b\x77\x46\x78\x20\xa4\x72\xc7\x42\x68\xa5\x25\x8e\xed\x00\x6b\x92\x76\x96\x62\xea\xff\x4a\xa4\xeb\x2b\xe5\xee\x3c\x73\x22\xca\x1d\x5f\x29\x04\x21\x66\x8a\xf4\x2b\x7a\xe5\xaa\xd7\x1b\xf1\xa4\x5e\xd2\x84\x48\x44\x84\x6b\xf7\xa9\x18\xa5\x5a\x3e\xed\x29\xf9\xfc\x1e\x7a\xaa\x77\xd0\xf6\xd2\x9b\x8d\x6f\xd5\x54\xda\xa2\x4e\x98\xf0\xcc\x49\x8a\x68\x53\xb8\x7a\xba\x0f\xab\x3d\xce\x5d\x38\xf8\x09\x2f\x7a\xfb\x2f\xcc\xcf\xac\xe6\x81\x75\xf7\xcb\xfd\x83\xf5\xa6\x56\x22\x91\xcc\x83\x8b\x9e\xf5\xba\xb8\x60\x5b\x7d\x27\x96\x22\xee\x59\xd9\x0e\x59\xfd\xfa\x79\xe6\x16\xee\xd5\xa8\x53\xae\x4c\x23\xe2\x81\xae\x7d\xbf\xae\x87\xcc\xdb\x32\xde\xc6\xf0\x9e\xac\xb3\xb8\xf8\xd1\x10\x7a\xd2\x91\x44\x25\x4c\xc3\x68\x34\x02\x2b\xc4\x94\x91\xc0\x3a\x34\xc9\xfc\x38\x99\xc3\xdd\x07\xa1\x69\xb8\x6c\xf0\xdb\xb6\x25\x9a\x19\x27\xbe\x95\x52\xc8\x06\x2f\x6e\xc9\x93\x85\xf6\x40\x3a\x51\xf6\x7d\xb0\xc5\x84\x6c\x93\x48\x4b\xfc\x29\x0d\x88\x97\x51\xec\xb8\xb0\xd2\x4b\x46\xf9\xc4\x03\xab\x24\xf2\x37\xd6\xc1\x59\xeb\xfe\x0f\x8d\xef\xd7\xe5\x57\x93\x5e\x69\x1f\xfc\xc5\x78\x3e\xe1\x59\xb4\xbf\x94\xf3\xef\x33\x6e\xbc\xa0\xfb\x55\xeb\x15\xbe\xe4\x06\xd4\xc7\xd3\xde\xd3\x5d\x8c\x9d\xaa\x7b\xd1\x2b\xba\xd7\x7d\xc7\xdc\xd2\x97\xbd\x4d\xc8\xd5\xe6\x98\x9e\xf5\xda\x04\xb4\xd5\x77\x76\x52\x5e\xc3\x16\x2a\xca\x27\x8c\xdc\x60\x4d\xee\x52\x31\x0f\xb4\x4c\x1a\xbc\x65\x72\xe1\xf2\xa7\xf4\x20\xf5\x2f\x73\x26\x23\xca\x38\x6b\xd3\xa6\x6d\xd8\x11\x26\x7c\x6c\xa8\xd1\xcc\xa1\x30\x3d\xa2\x7a\x60\x3d\x3e\x3e\x3e\x3a\xef\xdf\x3b\x37\x37\xf0\xee\x9d\x17\x45\x9e\x52\xf5\x5b\x50\xe3\xc8\x37\x65\x36\x32\xf5\x42\xd6\xa6\xe0\xd4\xe2\x3c\x47\x67\x45\x28\x97\x75\x32\x15\x7a\xb5\x1a\xf4\xfb\xb5\xdb\xd6\xe8\x7c\x73\x0a\xb5\xfa\x8e\xcf\xa8\xff\xeb\x91\xfd\xaa\x2a\x64\xa5\xd3\x6a\x96\x7a\xb5\xbb\x72\x31\x2e\xbf\x41\x34\x56\xf7\xff\x05\x00\x00\xff\xff\x92\x44\xa2\xc9\x55\x29\x00\x00"

func templatesViewsWidgetHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesViewsWidgetHtml,
		"templates/views/widget.html",
	)
}

func templatesViewsWidgetHtml() (*asset, error) {
	bytes, err := templatesViewsWidgetHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/views/widget.html", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x52\xdf\x4b\xe4\x56\x14\xfe\x62\x6c\x6b\xa7\x3f\x68\xa5\x0f\xa5\xf4\xe1\xf6\x41\x69\x1f\x62\x27\xb6\x05\x89\x46\x4b\xad\x42\xa9\xd3\x8a\x1d\x85\x3e\x86\x99\x6b\x4c\x3b\x93\x0c\xb9\x89\x55\xb0\x74\x1c\xa5\xb6\x20\x2d\x08\xcb\xc2\x2e\xec\x32\x3e\xef\xc2\xec\xca\x60\xf0\x47\xfc\x17\xce\x7d\xda\xb7\x85\x7d\xdb\x87\xfd\x23\x96\x9b\x64\x5c\x16\x16\xf6\x42\x72\xbf\x73\xce\x77\xbe\xf3\x9d\x90\xc7\xa3\xc3\x37\x00\xe0\x3d\x00\x9f\x02\xd8\x07\xf0\x21\x80\x67\xc8\x4f\x59\x03\x46\x01\x98\x1a\xf0\x11\x80\xef\x35\xe0\x1d\x00\xbf\x6a\xc0\x30\x00\xae\x01\x6f\x01\xf8\x4d\x03\x74\x00\xa1\x06\xbc\x09\x60\xb3\xa8\xff\x59\xdc\x9d\xe2\x3e\xd0\x80\x37\x00\x1c\x6a\xc0\xfb\x00\x8e\x34\x60\x4e\x03\xba\x1a\x30\x05\xe0\xc9\x10\xf0\x0d\x80\x4f\x74\x60\x0c\xc0\xb2\x0e\x8c\x28\x3d\x3d\xf7\xf1\x97\x5e\xf4\xe9\x40\x09\x40\xb7\xb8\xef\xe9\xf9\xdc\x53\x1d\xf8\x18\x40\xaa\x03\x1f\x00\x78\x54\xe4\x9f\xea\x80\x86\xdc\x83\x3a\x6a\x87\x91\x02\xab\xfa\xbb\x05\x56\x5a\x6a\x1f\xbd\x88\x87\x90\xfb\x55\xdf\xe7\x6d\x95\x98\x0f\xfc\x75\xcf\x65\x82\x47\x4c\xc4\xb5\x1a\x17\x02\x0b\x5b\xad\x20\x8c\x58\x2d\xab\xc4\xa1\x13\x79\x81\x3f\x48\x36\x02\x57\x60\x49\xbd\x2a\x5c\x08\xc7\xe5\xf8\x29\xf8\x03\x3f\xb7\x32\x4e\xd5\x6b\x72\x54\xb7\x5b\x1c\xab\x82\x87\x58\x73\x1a\x31\x07\x75\x29\x91\xbb\xb2\x43\x7d\xba\xa0\x1e\x56\xb8\xd2\x31\x2a\xc2\xf5\xea\xc6\x77\xb1\x2b\x8c\x6a\x60\xb1\x3a\xdf\xfc\xf6\x77\x6f\xc3\x69\x06\x13\x61\x5c\x5a\x72\x44\x64\x54\x43\xc7\x17\x0d\x27\x0a\x42\x8b\xfd\x98\x95\x58\x25\x0e\x9d\x66\x50\x0f\xd8\xcc\x4b\xfc\xd9\xd2\x92\xe3\xbb\xb1\xe3\x72\xa3\xca\x9d\xa6\xc5\xae\x63\x8b\xad\xc4\x42\x78\x8e\x5f\xaa\xfc\x50\x59\x30\xd6\x78\x28\xbc\xc0\xb7\x98\x39\x51\x2e\xcd\x07\x7e\xc4\xfd\xc8\x50\x7e\x2d\x16\xf1\xad\xe8\xcb\x56\xc3\xf1\xfc\x69\x56\xdb\x70\x42\xc1\x23\x7b\xb5\xba\x68\x4c\xbd\xe0\x29\x3f\xeb\x3c\x34\x16\xfc\x5a\x50\xf7\x7c\xd7\x62\xa5\xe5\x46\x1c\x3a\x0d\x63\x31\x08\x9b\xc2\x62\x7e\x2b\x0b\x85\xfd\xd5\x34\xcb\xa1\xed\x8f\x99\x65\xdb\x36\xd9\xf8\x38\x53\xb0\xfc\x99\x6d\x9a\x6c\x8e\x95\x99\x95\xc5\xb3\xf6\xe4\xa0\x34\x63\x7f\xad\xe0\xe7\x19\x6d\xc6\x2c\xb3\x9d\x9d\xbc\x65\xd6\x9e\x2c\x7f\xc1\xe6\x98\xc9\x2c\x36\x39\x0d\xba\x45\x29\x5d\xca\x7d\x4a\xe8\x44\xee\xc9\x36\xf5\xe4\xdf\x94\xc8\xff\x99\xdc\x93\xbb\x74\x45\x7d\xf9\x0f\x5d\x52\xca\x28\xa1\x53\xba\xa0\x3e\x5d\x66\x4f\x0f\x74\x9f\xce\x32\x46\x2a\xdb\xb2\x43\x89\x6c\x53\x4a\x0f\xa9\x27\x3b\xf2\x90\xd1\xd9\x2b\x55\xff\x7b\x4d\xd7\x39\xa5\x74\x42\x09\xe8\xf6\x00\x74\x29\xa5\x94\x1e\xc8\x7f\xb3\xa9\x09\xf5\x41\xc7\xd4\xa7\x33\xb9\x97\xa5\xfa\xa0\x3b\x74\x95\x3b\x06\x1d\xc9\xb6\xfa\x2b\x14\x3c\xa6\x84\xae\x40\x77\x29\xa5\x73\x79\x48\xa7\x83\x29\xd4\x57\x31\xe8\xa6\xda\x41\x1e\x5c\xab\xfe\xb2\x2d\x22\xde\xc4\xf3\x00\x00\x00\xff\xff\xb6\x0e\x85\xf6\xef\x03\x00\x00"

func localesRuLc_messagesWidgetMoBytes() ([]byte, error) {
	return bindataRead(
		_localesRuLc_messagesWidgetMo,
		"locales/ru/LC_MESSAGES/widget.mo",
	)
}

func localesRuLc_messagesWidgetMo() (*asset, error) {
	bytes, err := localesRuLc_messagesWidgetMoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/widget.mo", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"templates/views/widget.html": templatesViewsWidgetHtml,
	"locales/ru/LC_MESSAGES/widget.mo": localesRuLc_messagesWidgetMo,
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
	"locales": &bintree{nil, map[string]*bintree{
		"ru": &bintree{nil, map[string]*bintree{
			"LC_MESSAGES": &bintree{nil, map[string]*bintree{
				"widget.mo": &bintree{localesRuLc_messagesWidgetMo, map[string]*bintree{}},
			}},
		}},
	}},
	"templates": &bintree{nil, map[string]*bintree{
		"views": &bintree{nil, map[string]*bintree{
			"widget.html": &bintree{templatesViewsWidgetHtml, map[string]*bintree{}},
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
