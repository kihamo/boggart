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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\xdd\x6f\xdb\x38\x12\x7f\xef\x5f\x31\xc7\xe6\x4e\x36\x50\x49\x4e\xf6\xe3\x0e\x5a\x3b\x8b\xdb\x4b\x17\x7d\x48\xbb\x41\x93\x1e\xd0\x7b\x09\x68\x89\xb2\x99\xa5\x48\x2d\x49\xf9\xe3\x0c\xff\xef\x07\xea\xcb\xb2\x2d\xd9\x92\x9d\x6c\xaf\xd8\xf5\x43\x4b\x48\xc3\x1f\x67\x86\xbf\x99\x21\x47\x59\xad\x20\x20\x21\xe5\x04\x90\x2f\xb8\x26\x5c\x23\x58\xaf\x5f\x0d\x03\x3a\x03\x9f\x61\xa5\x46\x28\xc6\x13\x62\x6b\xaa\x19\x41\xd7\xaf\x00\x00\xaa\x2f\x43\x21\x23\x7b\x22\x45\x12\x43\x9c\x30\x66\x4b\x3a\x99\xea\x5c\x6e\x57\x96\xf2\x38\xd1\x99\x70\x45\x22\x95\x62\x78\x4c\x58\x21\x37\xd6\xbc\x56\x2a\x95\xc4\x30\x95\x24\x1c\xa1\x1f\xb1\xaf\xa9\xe0\x23\x5f\xf0\x90\x4e\x94\x4d\x16\xb1\x90\x1a\x55\x30\xc0\xe0\x50\x5f\x64\x03\x95\xf8\x3e\x51\x2a\x1b\x47\x35\xc8\x29\x3a\x2d\xed\xc2\x0a\x42\x6c\xfb\x62\x82\x20\xb5\x7d\x84\x56\x2b\xa0\x97\xff\xe0\x80\xde\xa6\x4b\x41\xb6\x72\x22\xb1\xd1\x03\x81\x03\xeb\x35\xba\x1e\xba\xb4\x46\x69\x17\xb7\xb0\x84\x89\x17\x33\x63\x4a\x70\xa0\x88\x6e\x34\xc5\x2c\xdd\xcd\x82\xa1\x9b\x6e\x59\x65\xa3\xdd\x80\xce\x72\x7e\x64\xc3\xfc\xbf\x2d\x2e\xf9\x8c\x60\x19\xd2\x85\x59\x66\xff\xad\x14\xf3\x1a\x8a\xf9\x82\xd9\x51\x60\x5f\x5e\x81\x19\xa9\xa8\x18\x2d\x94\x7d\x79\xd5\x40\xb5\xc5\x63\x8c\x39\x61\xbb\x34\xdb\x92\x28\xf8\x5e\x63\xad\x91\x93\xc2\x78\x4a\xe3\x71\x1d\x52\x29\x99\x94\xac\xe5\x78\x06\x1c\xcf\x6c\x8d\xc7\x0a\xc6\x58\x3e\x9a\x01\xda\xc0\x30\xaa\xea\xd6\x2a\x91\x18\xcd\x65\x63\x49\x14\xe1\x3a\xa3\x95\xd9\xa9\x10\xc8\x6f\xe0\x64\x2c\x01\x64\xe2\xb3\x58\xd3\x3c\x9b\x11\x23\x44\x78\x00\xeb\xf5\x75\x49\xaa\xfa\x79\x4f\x78\x86\x95\x2f\x69\xac\xbd\x99\xa0\x41\x6f\xd0\xff\xc1\xcc\x65\x8a\xc0\x7a\xbd\x5a\x81\xf3\x91\xfc\x96\x10\xa5\x9d\x4f\x1f\x6f\x9d\x3b\xac\xa7\xd9\xe3\x0c\x1c\x5d\x97\xb4\xb9\x5f\x2a\x4d\xa2\x8c\x31\x86\x1a\x43\x97\xd5\x90\xa6\xb3\x69\x19\x0d\x4f\x33\xaf\x98\x7b\x8a\x89\xd5\x08\xac\x33\xf7\xb6\x0c\x8f\x67\x34\x36\xa4\x8c\x9c\x6c\x6d\x39\xf9\x1c\x73\x53\x90\x3a\x7b\x7f\xce\xd0\x8f\x1b\x3c\x74\x13\xd6\xf0\xa6\x12\x6a\x1a\x8f\xed\xe6\x60\xdb\x9a\xb1\x13\x74\x55\x04\xf3\x04\x42\x1c\x10\xc8\x1c\x05\x94\x1f\x40\x6b\x26\x48\xb3\x02\x26\x48\x49\x65\x49\x46\x20\xfd\xd7\x56\x5a\xd2\x98\x04\x10\x60\x8d\xb3\xe7\x81\xb6\x25\x51\xb1\xe0\xca\x68\xc2\xc5\x5c\xe2\x18\x81\xd2\x4b\xa3\xfe\x9c\x06\x7a\xea\x5d\x0e\x06\x7f\x3d\xa0\x60\xb6\xa2\x49\xcc\x87\x65\x32\x39\x79\x5c\x28\x07\xdc\xec\xe2\xeb\x62\x07\xf5\xf4\x84\xd9\x0f\x34\x22\x67\x01\xbc\x27\x4a\xe1\xc9\x79\x18\x0f\xcb\xf8\x3c\x80\x4f\x8a\xc8\x2e\x00\x43\xf7\x98\xab\x0d\xce\xd1\x4d\x1b\xea\xb1\x08\x96\xc7\x97\x5b\xad\x40\x62\x3e\x21\x70\x41\xdf\xc0\x05\x13\x13\xf0\x46\xe0\x18\xa6\x1e\x22\xea\x66\x95\xd6\xb4\x08\x8c\x4b\xcc\x02\xce\x9d\x50\x34\x0d\x88\xd4\x25\x2d\xc8\x57\x00\x0c\xb3\x34\x03\x7a\x19\x93\x11\xc2\x71\xcc\xa8\x9f\xe6\x36\x77\x93\x83\xd0\x75\x20\xfc\x24\x22\x5c\x3b\x73\x49\x35\xe9\x05\x58\x93\x07\x71\xaf\x25\xe5\x93\x9e\x55\xe8\x60\xb8\xe5\xfc\x2c\x64\x84\x35\xa0\xab\xc1\xe0\x7b\x7b\x70\x69\x0f\xae\x1e\x2e\xbf\xf3\x06\xdf\x7a\x83\xef\xfe\x33\xf8\xbb\x37\x18\x98\x60\xb5\xfa\xfd\xa1\x9b\x61\x5f\x77\x53\xb7\x58\xeb\x06\x6b\xdc\xd9\xd6\x52\xd1\x65\x4c\x4e\x9e\x6c\xb8\xd7\x7a\xf2\x71\xe2\x41\xc6\x97\x2c\x59\x1f\xa3\xe8\x61\xfa\x0d\xdd\x34\x8b\x1d\xcc\x9e\x69\x09\x69\xaa\x3a\x7f\x98\x1c\xfa\x01\x9f\x99\x05\x7f\x22\x13\xca\x41\x9f\x9b\x4c\xdf\xf2\xe0\x7c\x90\x7b\xfa\xdf\xf3\x00\xfe\x99\xf2\x40\x7d\x35\x09\xd5\xf0\x35\xcd\xa8\x29\x71\x5f\x24\xa5\x1a\x64\xc7\x9c\x97\x0c\x53\xbe\x54\x4e\x4d\x95\x48\x99\xf6\x7b\x64\xd6\xe7\x54\xfa\x2d\x0f\x7e\xa7\x62\x30\x4d\x22\xcc\x1f\xc7\x4b\x4d\x54\x65\xd7\x6e\x09\x9f\xa4\x67\xe2\x4e\x78\x65\x7e\x23\x0b\x6d\xfb\x84\x6b\x22\x8f\x64\xa9\x2d\x80\xca\xd1\xb8\x6c\x75\x40\x39\xb2\x17\xaa\x03\x18\x54\x3b\x09\xc6\xad\x87\x0e\xfc\x81\x98\x73\x26\x70\xf0\x37\x8e\x23\x32\xaa\xa3\xef\x5e\xe3\xa1\xda\x6f\xa0\xbe\xe0\xe8\x7a\xaf\xb9\x60\x30\xec\x02\x7a\xbf\xc5\x70\x53\xbe\xb9\x28\xfb\x0b\xb5\x0d\x91\x46\xf3\x36\x9d\x85\xe3\xa2\x5f\x6f\xad\x3d\x54\x55\x43\x21\xa3\xfc\x6e\x64\x86\x08\xf2\x0d\x45\x3f\x22\x88\x88\x9e\x8a\x60\x84\x62\xa1\x34\x02\x1a\x8c\x90\x22\x5a\x53\x6e\x6e\x3b\xa6\xd8\xda\x5a\x4c\x26\x66\xe6\x0c\x33\x1a\x60\x2d\x8e\x71\xb5\xae\x86\x3f\x4f\xad\x86\x17\xaa\xd7\x85\xae\x51\x60\xfb\x82\xd9\x97\x95\xab\xec\x2f\xf1\xa6\x3b\x77\x62\xd5\xfb\x37\x66\x49\xa7\xba\xd9\x8e\x5f\xad\xea\x1e\x74\xa8\x7d\x50\x5e\x7b\x1d\x73\x54\x78\xf4\x13\x29\x09\xd7\x6d\xca\x1e\x74\x71\x38\xe4\x49\x50\x8a\xb9\x8a\x31\x1f\xa1\x2b\xd4\x70\x6d\x6c\x99\x51\x73\xc0\x8e\x49\x6f\xa7\x43\xb8\x50\xf6\xb7\x1d\xf3\x66\x0a\x93\x36\xa5\xf3\x5a\x66\xf2\x39\xda\xea\x6c\xfb\x82\x6b\x29\x18\x82\x34\x63\x22\xe3\x57\x3b\xf7\x2b\x82\x99\x21\x46\x9a\xe9\xb6\x1c\xbe\x55\xcc\x9c\xc1\xa5\x33\xb8\x82\xbc\x98\x7d\x8f\xd2\x1c\x6b\x82\x74\x1b\x4a\x12\x1c\x08\xce\x96\xe0\x76\xf4\x42\xfb\xdc\x08\xad\xf3\x23\xb4\xe6\x30\x9c\xc0\x9b\x2f\xb3\xcd\xad\x3e\x43\xb4\x46\xeb\x4c\x9a\x67\x20\x0b\xea\xca\x8d\x52\x5b\x13\xa5\x35\xc6\xdb\x63\x7d\xa8\x83\x76\x14\x76\x9c\x68\x2d\xf8\xde\x97\x0a\x1e\x0a\x94\xbb\x26\x93\xa8\xf0\x9d\x8b\x79\x25\x57\x7c\x10\xf3\x22\x55\x64\x92\x27\xda\xe7\x1a\x03\x4f\xe0\x44\xb7\xd8\x39\x61\xca\x4b\x84\x5b\x8b\xe3\xc2\x16\x70\xd7\xe8\xfc\xbf\xc8\xe4\xdf\xfc\x19\xe2\xdd\xb4\xfd\x33\xc4\x1b\xe6\x7d\xad\x21\x9e\xdd\x3c\xda\x1c\x22\x8f\x1f\x0e\x8f\x5e\x41\x86\xae\x21\xf5\xe1\x2b\xca\xe1\xab\xd0\x01\xaf\x35\xbc\xaa\x79\xbc\xf3\xa8\xf1\x8b\x72\xc5\x3f\x95\x3f\x5f\x30\x87\xe9\xb2\x37\x59\xff\xf9\xf3\x55\xc5\x20\xa5\xb1\xa6\xfe\xbb\x87\xf7\xb7\xd0\xcb\xc6\x9f\x3e\xde\x02\x72\x03\xac\xa6\x63\x81\x65\xe0\x62\xa5\x88\x56\xee\x8c\xf0\x40\x48\xe5\x8e\x85\xd0\x4a\x4b\x1c\xdb\x01\xd6\x24\xed\x32\xc5\xd4\xff\x95\x48\xd7\x57\xca\xdd\x79\xe6\x44\x94\x3b\xbe\x52\x08\x42\xcc\x14\xe9\x57\xf4\xca\x55\xaf\x37\xe2\x49\xbd\xa4\x09\x91\x88\x08\xd7\xee\x53\x31\x4a\xb5\x7c\xda\x53\xf2\xf9\x3d\xf4\x54\xef\xa0\xed\xa5\x37\x1b\xdf\xaa\xc1\xb4\x45\x9d\x30\xe1\x99\x93\x14\xd1\xa6\x70\xf5\x74\x1f\x56\x7b\x9c\xbb\x70\xf0\x13\x5e\xf4\xf6\x5f\x98\x9f\x59\xcd\x03\xeb\xee\x97\xfb\x07\xeb\x4d\xad\x44\x22\x99\x07\x17\x3d\xeb\x75\x71\xc1\xb6\xfa\x4e\x2c\x45\xdc\xb3\xb2\x1d\xb2\xfa\xf5\xf3\xcc\x2d\xdc\xab\x51\xa7\x5c\x99\x46\xc4\x03\x5d\xfb\x7e\x5d\x0f\x99\xb7\x65\xbc\x8d\xe1\x3d\x59\x67\x71\xf1\xa3\x21\xf4\xa4\x23\x89\x4a\x98\x86\xd1\x68\x04\x56\x88\x29\x23\x81\x75\x68\x92\xf9\x71\x32\x87\xbb\x0f\x42\xd3\x70\xd9\xe0\xb7\x6d\x4b\x34\x33\x4e\x7c\x2b\xa5\x90\x0d\x5e\xdc\x92\x27\x0b\xed\x81\x74\xa2\xec\x5b\x61\x8b\x09\xd9\x26\x91\x96\xf8\x53\x1a\x10\x2f\xa3\xd8\x71\x61\xa5\x97\x8c\xf2\x89\x07\x56\x49\xe4\x6f\xac\x83\xb3\xd6\xfd\x1f\x1a\xdf\xaf\xcb\x2f\x28\xbd\xd2\x3e\xf8\x8b\xf1\x7c\xc2\xb3\x68\x7f\x29\xe7\xdf\x67\xdc\x78\x41\xf7\xab\xd6\x2b\x7c\xc9\x0d\xa8\x8f\xa7\xbd\xa7\xbb\x18\x3b\x55\xf7\xa2\x57\x74\xb2\xfb\x8e\xb9\xa5\x2f\x7b\x9b\x90\xab\xcd\x31\x3d\xeb\xb5\x09\x68\xab\xef\xec\xa4\xbc\x86\x2d\x54\x94\x4f\x18\xb9\xc1\x9a\xdc\xa5\x62\x1e\x68\x99\x34\x78\xcb\xe4\xc2\xe5\x4f\xe9\x41\xea\x5f\xe6\x4c\x46\x94\x71\xd6\xa6\x4d\xdb\xb0\x23\x4c\xf8\xd8\x50\xa3\x99\x43\x61\x7a\x44\xf5\xc0\xfa\xfc\xf9\xf3\x67\xe7\xfd\x7b\xe7\xe6\x06\xde\xbd\xf3\xa2\xc8\x53\xaa\x7e\x0b\x6a\x1c\xf9\xa6\xcc\x46\xa6\x5e\xc8\xda\x14\x9c\x5a\x9c\xe7\xe8\xac\x08\xe5\xb2\x4e\xa6\x42\xaf\x56\x83\x7e\xbf\x76\xdb\x1a\x9d\x6f\x4e\xa1\x56\xdf\xf1\x19\xf5\x7f\x3d\xb2\x5f\x55\x85\xac\x74\x5a\xcd\x52\xaf\x76\x57\x2e\xc6\xe5\xf7\x88\xc6\xea\xfe\xbf\x00\x00\x00\xff\xff\xd9\x95\x6e\x40\x61\x29\x00\x00"

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
