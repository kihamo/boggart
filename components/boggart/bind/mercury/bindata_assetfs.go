// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
// DO NOT EDIT!

package mercury

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x9b\x7d\x6f\xe3\xb6\x19\xc0\xff\xef\xa7\x78\x40\xe4\xd0\xa4\x57\xf9\x45\x52\x6f\x5d\x60\x67\x7f\xec\x36\x60\xc0\x0e\x18\xd6\x60\xc3\xd6\x76\x01\x2d\xd2\x31\x7b\x12\xa9\x23\x29\x3b\x37\xc3\xdf\xbd\x90\x28\xc9\xb4\x2d\xd9\xb2\x23\xf5\x9a\xe8\x80\x0b\xcc\xb7\x87\xcf\xcb\xef\xa1\x6c\x8a\x5c\xaf\x81\xd0\x39\xe3\x14\x50\x20\xb8\xa6\x5c\x23\xd8\x6c\xbe\x02\x00\x58\xaf\x81\xcd\x81\x7e\x82\x01\x0e\x34\x13\x1c\x50\x24\xb8\x5e\x84\x9f\xcb\x1e\x13\xc2\x96\x10\x84\x58\xa9\x29\x92\x62\x85\xee\xb2\xda\xfd\x96\x40\x84\x4e\x44\x9c\xb1\x0b\xe9\x27\x15\x15\x9f\x9e\x94\x33\x76\xad\x31\xfb\xe3\x9e\x1e\x62\xcc\x69\xb8\xd7\xe3\xb0\x97\x66\x3a\xa4\x15\xbd\xb2\x9e\x0b\xf7\x2e\xb5\x63\xfc\x3d\x07\xf4\xc1\xa8\x0f\x4a\x63\xcd\x94\x66\x81\x42\x30\x80\xcd\x66\x32\x5c\xb8\x35\xc3\x6d\x33\x42\x8a\xe5\x9c\x3d\xa1\xbb\xc9\x90\xb0\x65\x85\x56\x35\xd5\x3b\xca\x16\x4e\xae\x99\x2f\x2e\x7a\x6a\xfa\xa4\x9d\x28\xd1\x94\xc0\x5c\x70\xed\x8c\x3d\x88\x9c\x99\xe3\x8d\x6a\x46\xe6\x01\x93\x78\x05\xd7\xc6\xda\x3f\x27\x52\x52\xae\x81\x60\x4d\x41\xb3\x88\x82\xe0\x40\xe8\x92\x05\x14\xde\x64\x96\x73\x16\x96\xff\xaf\x63\xc9\xb8\x06\x34\x51\x81\x64\xb1\x06\xfd\x39\xa6\xd3\x9f\x10\x8e\xe3\x90\x05\x38\x0d\xff\xf0\x17\xbc\xc4\xa6\xf5\x27\x74\x47\x44\x90\x44\x94\xeb\xc1\x4a\x32\x4d\xaf\xd3\x49\xee\xc5\x0f\x5a\x32\xfe\x78\xfd\x35\x82\xeb\x41\x5a\x33\xf8\xab\x90\x11\xd6\x80\xdc\xd1\xe8\x9d\x33\x1a\x3b\x23\xf7\x7e\xfc\xdd\xed\xc8\xbf\x1d\x7d\xf7\xdf\xd1\x1f\x6e\x47\x23\x74\x03\xe8\xeb\x9b\x9b\xc9\xd0\x08\xbe\x43\x37\x37\x05\x5c\x87\xee\x8d\xef\xbe\xaa\x6e\xd1\x78\x16\xd2\xd2\x73\x59\x21\xfb\xeb\x28\x2d\x59\x4c\x49\xea\x03\x6c\xea\x89\x76\x24\x55\xb1\xe0\x8a\x2d\x29\x70\xb1\x92\x38\x46\xa0\xf4\xe7\x90\x4e\xd1\x8a\x11\xbd\xb8\x1d\x8f\x46\x6f\x8e\x78\x79\xa2\x17\x14\x93\x63\xed\xb2\xbe\x31\x17\xb0\xc7\x64\x81\xa1\x5e\x9c\x31\xf2\x7e\x0c\x43\xf8\xf7\x85\x63\xdd\x67\x8c\xf5\x9e\x31\xd6\x6f\x3c\x76\x32\x3c\xe6\xc7\x74\xec\x89\x28\xcc\x04\xf9\x7c\x22\x57\xf8\x23\x85\x2b\xf6\x2d\x5c\xa5\x2b\x82\x82\xdb\x29\x0c\xcc\xa7\x1a\x04\x8d\x60\x69\x16\x46\xcc\x09\x5c\xd3\x4f\xf9\xe0\x41\x16\x48\xb8\x32\xdc\x67\x85\x1b\xbb\xf9\x3f\x14\xcb\xa2\x35\xfd\x9c\x62\x5e\x10\xab\x92\x20\xa0\x4a\xa1\xf5\x1a\x28\x27\xb0\xd9\x9c\xf2\x28\x29\x3d\x6a\x4f\x3e\x30\xe9\x07\x57\xa9\xe8\xf5\x7a\x67\xe6\xcc\xdb\x47\xdc\x55\xc8\x3d\xda\x21\xf7\x5b\x2e\xf7\x7e\x7c\xcc\x4d\x56\x7f\x36\x87\x47\xbd\x1d\x75\x2f\x53\x23\x47\x4d\x06\x67\x5a\xa9\x18\xf3\xc2\x53\x21\x9e\xd1\x10\xb2\xbf\x0e\x49\xc3\x27\xd1\x9d\xad\xd1\x7b\x1a\x6a\x9c\x9a\x3f\x84\xb7\x76\xbd\x99\x33\x75\x42\x2a\xad\x91\x95\x34\x54\x34\x55\x3d\xec\x40\xf5\x22\xe0\x75\xba\xb7\xa4\xfa\x96\xbe\xf6\x54\x67\x7c\x2e\xba\xd6\xdb\xa6\xc5\x48\x6f\x81\x16\x3a\xc7\x49\xa8\xab\x55\x3f\x47\xc7\xcc\xa8\xe3\x79\xd4\x7a\xa6\xb9\x17\x65\x9a\xdb\x61\xa6\xb9\x35\x99\xe6\xb6\x93\x69\xad\xa9\x5e\x95\x69\x6e\x35\xb1\xcf\x51\xdd\xca\xb4\xd6\x54\x3f\xc8\xb4\x0e\xf4\xb6\x69\xe9\x32\xd3\xdc\x97\x92\x69\xde\x45\x99\xe6\x75\x98\x69\x5e\x4d\xa6\x79\xed\x64\x5a\x6b\xaa\x57\x65\x9a\x57\x4d\xec\x73\x54\xb7\x32\xad\x35\xd5\x0f\x32\xad\x03\xbd\x6d\x5a\xba\xcc\x34\xef\xa5\x64\x9a\x7f\x51\xa6\xf9\x1d\x66\x9a\x5f\x93\x69\x7e\x3b\x99\xd6\x9a\xea\x55\x99\xe6\x57\x13\xfb\x1c\xd5\xad\x4c\x6b\x4d\xf5\x83\x4c\xeb\x40\x6f\x9b\x96\x2e\x33\xcd\xff\x12\x99\x76\xfc\x57\xf1\xe9\x19\x26\xc3\x23\xbf\x8b\x27\xc3\x6c\x8f\xa4\xd1\xce\xd6\x5e\x95\x55\xb4\x3e\xee\xd2\x54\x6e\x23\x12\xa6\xe2\x10\xbf\xd8\x6d\xc4\xf7\x85\xfa\xbf\xb3\xbd\xc3\xb9\x90\x11\x30\x32\xdd\xfa\x37\x1f\x99\x36\x38\x0b\x21\xd9\xff\x05\xd7\x38\x84\xac\x6c\xc0\x0e\xe9\x5c\xa3\x6c\x7b\xcc\xd1\xe2\xf1\x31\xa4\x53\xb4\xc4\x21\x23\x58\x0b\x89\x20\xa2\x7a\x21\xc8\x14\xc5\x42\xd5\xcd\xba\xaf\x65\x26\xfb\x51\x8a\x24\x3e\x32\x20\x1b\x64\xf2\xab\x0c\x2e\xd7\x52\x84\x4e\x5e\x69\x42\xed\x17\x91\xf6\xad\x40\xa7\xda\x4f\x51\x24\x08\x7d\xd0\x63\xb4\x8d\xca\x0f\x0b\xb1\x82\xfb\xb1\x89\xca\x6e\x26\x4b\xfa\x29\x61\x92\x12\x74\xf7\x4d\x93\x5c\x9d\x0c\x33\x2d\x4e\x74\x3a\x24\xf3\x5d\xa1\xee\xbb\x5a\x2e\x2b\x25\x31\x1e\x27\xf9\x86\x2b\x0a\x16\x34\xf8\x38\x13\x4f\x65\xec\x7e\x51\x8e\x5a\x31\x1d\x2c\x10\x70\x1c\xd1\xad\xe5\x59\xa4\xcb\x42\xb9\x57\x9f\xef\xe3\x0e\xf2\x96\x1b\x40\x5a\x26\x14\x65\xbb\x4d\xa9\x6c\x4a\xca\x45\x02\x86\xa7\xfc\x50\xc9\x64\xd3\xe6\x2f\x05\x85\x7b\x00\x85\xdb\x13\x28\x5c\x1b\x0a\xb7\x16\x0a\xb7\x87\x50\x78\x07\x50\x78\x3d\x81\xc2\xb3\xa1\xf0\x6a\xa1\xf0\x7a\x08\x85\x7f\x00\x85\xdf\x13\x28\x7c\x1b\x0a\xbf\x16\x0a\xbf\x7f\x50\xe0\x48\x24\x5c\xef\x83\x91\xd7\xf6\x02\x8e\xc2\xd6\x12\x90\xa2\xa2\x1a\x12\xd3\xda\x3f\x50\x62\xb1\xca\xf7\x0f\x2c\x4e\x4c\x65\x2f\x30\xc9\x4d\x2d\x29\xc9\xcb\xd5\x90\x64\x8d\xfd\x63\x44\xb3\x88\xee\x23\x92\xd5\xf5\x82\x10\x63\xe9\xf6\x39\x93\x15\x6b\x9e\x34\x2c\xa2\xfd\xc3\x83\x60\x7d\x80\x47\x56\xd7\x0b\x3c\x8c\xa5\x25\x1e\xa6\x78\x80\x47\x5a\xdd\x2b\x32\xd2\x54\xd8\xdd\xd9\xf8\x1b\xd7\x54\x2e\xcd\xbe\x0d\x70\xa1\x21\xc8\x4f\x76\x69\x2c\xd9\x7c\xfe\xc2\x69\xe1\x49\x34\x4b\x1f\x1c\x86\x8c\xc2\xfa\x8c\x8b\xb2\xb0\xc4\x61\x42\xa7\x68\xbd\x86\x41\x5e\x07\x9b\x0d\x82\x18\x6b\x4d\x25\x9f\xa2\xff\xfd\xe8\xbc\xfd\xf9\x4f\x3f\x8e\x9c\x3f\xfe\xfc\xcd\x15\x7a\x9d\x44\xb8\x75\x44\xbc\x7a\x1a\x5c\x9b\x06\xb7\x82\x06\xb7\x7f\x34\x78\x75\x34\x64\xdf\xb4\xbe\x35\xe7\x3d\x31\x27\xaf\xe1\x79\x52\xc5\x84\x67\x33\xe1\x55\x30\xe1\xf5\x8f\x09\xbf\x8e\x09\x3c\xd7\x54\x42\x2c\xa9\x52\x30\xc3\xc1\x47\x98\x25\x5a\x0b\xfe\x0a\xb1\xf0\x6d\x2c\xfc\x0a\x2c\xfc\xce\xb1\x68\xc4\x45\xc8\x1f\x94\x08\x19\xa9\x7d\x1f\x54\x35\xa8\x39\x4c\xcd\x1c\x5e\x80\x25\xe6\x73\x45\xb5\xe3\x35\x09\x80\x21\x27\x8f\x80\xa4\x8a\xea\xf2\x5b\xde\x4c\x73\x98\x69\x6e\xde\xf1\x16\x9e\x37\x5d\xb6\x5c\xfe\xd3\x0c\x31\x6f\xce\x8c\xb0\x73\x67\x55\xc9\x2c\x62\x87\xd3\x16\x6f\xc5\x8b\x99\x15\x5e\xee\x7c\xb9\x4e\x8b\x67\xcc\x7b\xf9\x1a\x30\x19\xa6\x81\x6a\xff\x25\xea\xcb\x7c\x59\x5a\xdc\x42\x88\x28\x56\x89\xa4\x11\xe5\xfa\xf7\x77\xeb\xe2\xc5\xdd\x1f\x28\x74\x8d\x88\x93\x46\xd5\xfe\x6e\xf8\x61\xeb\xe8\x33\x8e\xe8\x1f\x11\xf8\xaf\x34\x9f\x2e\x3a\xed\xff\x9e\x9a\x5b\x1d\xac\x78\xd6\x7c\xd1\x13\xff\xd0\xc4\xb7\xa6\x13\xb1\x6c\xc0\xda\x7c\x91\xda\xee\xdb\x9c\x3e\x0b\x05\xe5\xd1\xa5\xec\xf7\x6b\x3a\x74\x40\xa5\x14\xb2\xf1\x41\x14\x4d\x76\x2e\x02\x59\xa7\x95\xf6\x24\x0e\xfe\x92\xcb\x6d\xac\x95\xbd\x98\x34\x71\x45\xa3\x8e\x59\x67\xfb\xea\x50\xcd\xcd\xa1\x13\x17\x87\x76\xec\xcb\x16\xf2\x06\x57\x88\x60\xb3\xb1\xaf\x10\x35\xd6\xd7\x44\x88\x53\x6b\x4a\x72\xde\x79\xa1\xad\xe5\x5a\x0a\xfe\x58\x19\xb2\xb3\x04\x95\xe7\xdd\x9e\xad\x52\x21\x6d\xfc\x3d\xff\x47\x98\x48\x1c\x02\x7a\xfb\x86\x80\xa2\x81\xe0\x04\xd9\x85\x74\x39\xde\x9b\xed\x6a\xe7\x6a\xd8\x7e\xeb\x99\x9a\x9c\x89\xdc\x11\x03\x6c\xfd\x7f\x53\xf5\x4f\x9f\xda\xda\xff\x37\x19\x1a\x24\xce\xa2\xf1\x8c\x79\x9a\x67\x7b\x33\x99\x3b\x6b\xde\xce\xcd\xc1\x62\xed\x03\xc1\x21\xa2\xba\x78\xbd\xd1\xe0\x40\xe8\xd1\xe5\x1c\x2e\x5a\x8d\x3f\xe0\x8f\xd4\xfa\x49\x7b\xd6\x4a\x1c\xe1\x8f\x66\x97\xb1\xbd\xa5\x78\x4f\x64\xf7\x6b\xf1\xee\xa4\x75\x0b\x24\x3a\x4b\x85\x0b\x00\x39\x3b\x0c\x5d\xb0\xf0\x77\xac\xb4\xd9\xee\x00\x51\xec\x7b\x9d\x05\x44\x88\x95\x36\x6f\xa6\x1e\xc4\x7c\xfe\xd0\xfe\x93\xfa\xe8\x04\xdd\xc3\xd2\xd2\xf3\xb8\xce\x8a\x0b\x1f\xcf\x5d\x80\xf9\xe5\x01\xe4\xcf\xe4\x8f\x77\x8c\xdf\xbe\xfc\x17\x49\x1f\x7f\x95\xf0\xfd\xf6\x47\xb2\x8d\x6d\xe5\xa7\x5f\x03\x00\x00\xff\xff\xaa\x6b\x35\x62\x03\x42\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x56\x5b\x6f\x5c\x57\x15\xfe\x42\xec\x71\x3d\xb6\xc7\x69\x29\xe1\x56\xe8\x2e\x90\x90\x92\x4c\x32\x17\x17\xca\x24\x4e\x5a\x72\xa9\x1a\x62\x1a\x25\x43\x78\x28\x48\x1c\xcf\x6c\xdb\xa3\xcc\x9c\x33\x3a\xe7\x4c\x52\x43\x25\xec\x71\xdb\xa4\x72\x88\xa1\x02\x91\xa2\x96\x0a\x28\x3c\x00\x12\x53\x27\xd3\x4c\x7d\x19\x0b\x21\x78\x41\xa0\xb5\x91\x10\x0f\x95\xb8\x48\x48\x08\xf8\x01\x95\xe0\x05\xad\xbd\x77\xe6\xe2\x8c\x21\x2f\x39\x52\xbc\xd6\x5e\x7b\x9d\x6f\x7f\x6b\xad\x6f\x9f\xcc\x9f\xee\xeb\xfb\x36\x00\x8c\x01\x78\x00\xc0\xda\x36\xe0\x09\x00\x5f\xdc\x0e\xfd\x3c\xd3\x07\xdc\x0f\x60\xb6\x0f\xd8\x09\xe0\x6a\x1f\xd0\x0f\xe0\x47\x7d\x40\x04\xc0\x4f\xfb\x80\x3d\x00\x96\xfb\x80\x0f\x02\x78\xdb\xae\xdf\xb1\x6b\xd1\x0f\x7c\x18\xc0\x21\x6b\x9f\xee\x37\x78\x5f\xe9\x07\x46\x00\x2c\xf5\x03\xf7\x00\x78\xb5\x1f\x18\x02\xf0\xe3\x7e\x60\x80\xf1\x6c\x7c\xa5\x1f\xd8\x05\xe0\x97\xd6\xfe\xa3\x1f\xf8\x18\x80\xd1\x88\xc1\x1f\x8b\x00\x0f\x02\x38\x1d\x01\x1e\x02\x50\x8e\x98\xf7\x5f\x8a\x00\x7d\x00\xbe\x6b\xed\x6b\x11\x20\x06\xe0\xf5\x88\x39\xf7\x7a\x04\x18\x04\xf0\x8b\x88\xa9\xe7\x37\x11\x80\x4b\x56\x11\xc3\xe3\x0f\x36\xfe\xd7\x08\x70\x1f\x80\xbf\x47\x0c\x9f\x7f\x5b\xfc\xed\x03\x66\x7f\x68\xc0\xe0\xdf\x3b\x60\xf0\x76\x0e\x98\x7d\x61\xed\x1e\x6b\x13\xd6\x3e\x3a\x60\xf0\x1f\xb3\xf9\xa7\x06\x80\x28\x80\x73\x76\xed\xd8\xbc\xf3\xd6\x06\xd6\x7e\xd5\xda\x05\x7b\xee\xe2\x00\x30\xb3\x0d\xf8\xe6\x00\x90\x66\x5e\xf7\x00\x9f\x04\xb0\x7f\x10\x18\x06\x20\xad\xf5\x07\x81\x2f\x03\xf8\xda\x20\x70\x04\xc0\x5f\x06\x81\x29\xe6\x17\x05\x8e\x03\xb8\x14\x05\x0e\x02\xf8\x6d\xd4\x68\x60\x70\x08\xf8\x28\x80\xc3\x43\x66\xde\x93\x43\xa6\x6f\xb3\x43\xc0\x0e\x00\x97\x87\x4c\x3f\x5e\xb6\xf1\x9f\x0d\x19\xfc\xd5\x21\x83\xfb\xc0\x30\x70\x92\x79\x0e\x03\xa7\x01\xfc\x6e\x18\xc8\x02\x78\x70\xc4\xac\xc3\x11\xc3\xeb\x57\x23\xa6\x9f\xbf\xb7\xf6\xed\x11\xe0\xe3\x00\xfe\x36\x62\xe6\x3c\x1c\x33\xf3\x4c\xc5\xcc\x7e\x36\x66\xf4\xf6\xa5\x98\xd1\x4f\x3e\x66\xfa\x36\x1b\x33\xf9\xcf\xc5\x0c\xee\x4f\x62\x86\xd7\xcd\x98\xe1\xfb\xeb\x98\xe1\xfb\x47\xbb\xfe\x67\x0c\x78\x2f\x80\xff\x58\xfb\xee\x51\x63\xf7\x5a\x7b\x74\xd4\xe8\xe9\xe9\x51\xe0\x43\xdc\xf7\x51\x83\xff\xb2\x8d\x6f\x8c\x9a\xf9\xfd\xd9\xda\x7f\x59\xfb\x8e\xb5\xdb\x77\x98\x73\x46\x77\x00\xdb\x60\x34\xff\x1e\x18\x6c\x7e\x98\x33\xd7\xc5\xda\xe5\x7b\xc2\x5a\x7c\x17\xda\x0f\xf7\x80\xcf\xd1\x58\x30\xb5\xf3\xc3\x1a\xe3\xfa\xb9\xb6\x94\x8d\x31\xbf\x41\xeb\xb3\x26\xe2\x1d\x38\x1f\xe1\xda\x00\x3c\x0c\xe0\x13\x5c\x1f\x6b\x10\x66\xb6\x69\x9b\x33\xda\x91\x7f\xa0\xc3\xe7\x7a\xf9\x8e\xed\x07\xb0\x0f\xa6\x7f\xfc\xec\xb4\x76\x87\xfd\x5e\xdc\x6f\xd7\xbb\x61\xb4\xc9\xfd\xff\x00\x80\x24\x80\xf7\xd9\x3d\xfe\xa6\x08\xeb\xdf\x0b\xe0\xfd\xad\x43\xf2\x22\x90\x39\xcf\xcd\xb7\xbd\x00\x7b\xdb\xd1\xbd\x1d\xe1\xc7\xcb\x7e\xa1\x88\xc7\x2b\xd3\x95\x20\xc4\xd1\x19\xc7\x9d\x96\x22\x5f\x08\xca\x45\x67\x56\x94\xbc\xbc\x14\x53\x4e\xa1\x28\xf3\xe2\x62\x21\x9c\x11\xd2\xf7\x3d\x5f\xec\x0a\x7a\x26\x06\x95\x5c\x4e\x06\xb7\xed\x85\x85\xd2\x9d\x81\xe8\xc4\x16\x48\xc5\xf7\xa5\x1b\x8a\xbc\x13\x4a\xe1\xb8\x79\xb3\xeb\xb9\xa2\x24\x43\xe9\x77\x6f\xdf\xda\xca\xcb\x0b\x85\x9c\xd4\xc8\x76\xbb\x24\x9d\xa0\xe2\xcb\x92\x74\xc3\x00\xc7\x3a\xa1\x70\x4c\xe6\x64\x69\x52\xfa\x38\x26\x83\x9c\x5f\x28\x87\x05\xcf\xc5\x31\x43\x05\x27\xe4\xa4\x5f\x71\xfc\x59\x3c\x21\xc3\x3b\xe8\x46\x67\xd6\xd6\xe5\x3e\xe9\x86\xd2\xbf\xe0\x14\xc5\x94\xe7\x0b\x67\x2a\x94\xbe\x28\xfb\x32\x08\xc4\xa4\x93\x3b\x2f\x26\x2b\x61\xe8\xb9\xdd\x49\x39\x5b\x46\xe8\xf8\x85\xa9\xa9\xee\x3d\xd7\x0b\xff\xe7\x7e\xd9\xbb\x28\xfd\x7d\x86\x0e\x17\xcd\x9d\xc2\x49\xc7\xd5\x65\x9d\xac\x14\xf9\x8f\x2b\x71\xca\x09\x42\x93\x2b\xbc\xa9\xa9\xae\xa5\x8b\x09\xe7\xbc\x34\x2f\x4e\x38\x7e\x6e\x06\x13\xce\x2c\x26\xda\x2d\xc5\x84\xe7\x86\x33\xe6\x6f\x71\x56\x04\xa1\x13\x16\x82\xb0\x90\x0b\xf0\x39\xef\x82\x69\xee\x53\xb9\xd0\x63\x7b\x46\x06\x32\xc4\x59\xe7\x82\xc4\x59\x59\x0e\xcd\xe6\xd9\x19\xef\xa2\xc8\x26\xad\x4d\x59\x9b\xb6\x76\xcc\x58\xa7\xe4\x55\xdc\xd0\xf8\x9a\x8a\xf6\x34\x45\xe3\xea\x69\x66\x93\xe2\x80\xf8\xc2\x0c\xb2\x29\x6b\xd3\xd6\x8e\x19\x7b\xce\x29\x56\x24\xce\xc8\xb2\xe7\x87\xf1\x89\x60\xba\x90\x8f\x7f\xa6\x32\x1d\xc4\xb3\x5e\x86\x65\xf3\xd8\xf9\xc2\x8c\x53\xf2\xf6\xfb\x95\xe8\xe9\xa7\xb2\xf1\xa3\xbe\x74\x58\x0f\x71\x96\x4c\x46\xa4\x12\xc9\x4f\xc7\x13\xe9\x78\xea\x53\x22\x95\xce\x3c\xf2\xc8\xde\x44\x3a\x91\x88\x72\xab\xe2\x59\xdf\x71\x83\xa2\x13\x7a\x7e\x46\x7c\x56\x63\x88\x89\x8a\xef\x94\xbc\xbc\x27\x0e\x75\x01\x1f\x8e\x9e\x72\xdc\xe9\x8a\x33\x2d\xe3\x59\xe9\x94\x32\xa2\xb5\xce\x88\x33\x95\x20\x28\x38\x6e\x74\xe2\xc9\x89\xe3\xf1\x73\xd2\x0f\x0a\x9e\x9b\x11\xc9\xfd\x89\xe8\x51\xcf\x0d\xa5\x1b\xc6\xb3\xb3\x65\x99\x11\xa1\x7c\x26\x3c\x50\x2e\x3a\x05\xf7\xa0\xc8\xcd\x38\x7e\x20\xc3\xf1\xcf\x67\x4f\xc4\x1f\x6d\xe7\x31\x9f\x29\xe9\xc7\x8f\xbb\x39\x2f\x5f\x70\xa7\x33\x22\x7a\xba\x58\xf1\x9d\x62\xfc\x84\xe7\x97\x82\x8c\x70\xcb\x7a\x19\x8c\xa7\x0f\x0a\xe3\x8e\xbb\xbb\x92\x89\xf1\xf1\xa4\xd8\xbd\x5b\xb0\x9b\x78\x68\x3c\x99\x14\x47\x44\x42\x64\xf4\xfa\xf0\x78\xea\xd6\xd6\xa1\xf1\x31\x76\xf7\xe8\xb4\x43\xc9\x84\x78\xf6\x59\xf3\xca\xe1\xf1\x54\xe2\x61\x71\x44\x24\x45\x46\xa4\x0e\xf2\x67\x47\xcd\x53\x9d\x56\xd4\x02\xad\xd3\x0d\xaa\x6d\x8e\xa8\xc5\xcd\x11\xfd\x51\xda\xf4\xd2\xe6\x90\x5a\xbc\x2d\x04\xfa\x06\x6d\xa8\x39\xaa\xd3\xaa\xba\xc2\x8b\x65\xba\xae\x16\xd4\xbc\xaa\x82\xae\xd1\x4d\x5a\xa3\x3a\xad\xeb\x7f\x0d\xaa\x0b\x9d\xf8\x26\x35\x68\x8d\x9a\xb4\x2c\xa8\xa1\x0f\x6a\xd0\x0a\xd5\xd4\x0b\xd4\xa0\x86\xa0\x9b\x54\xa3\x65\xaa\xab\x39\x75\x99\x1a\xb4\x4a\x4d\x35\xaf\xae\x08\x35\x2f\xa8\xa9\x23\x6f\xd0\x0a\x35\xe9\x2d\xbe\xd0\xf4\x43\x35\x4f\x1b\x54\x57\x97\x69\x9d\x9a\x8c\xd6\x71\x9e\x5a\xec\x38\x4d\x2d\xf6\x38\xab\x17\x41\xce\x52\x55\x3e\x9e\x96\xa9\xc6\xc7\xdf\x7d\x9a\x9b\xcf\xdc\x82\xec\x0f\x4c\xd3\xd5\x8b\x86\xe8\x0d\xaa\xa9\x2a\xd5\x04\x93\x59\xd6\x95\xae\xa9\x25\x41\xeb\x54\x13\x6a\x5e\x5d\xa2\xba\xaa\xaa\x4b\x1a\xa1\xde\xf9\x6e\x8d\x93\x5a\xef\xea\x6c\x3d\x2d\x35\xc7\x74\xd9\xe3\xb2\x34\xed\x4d\x07\x6e\x50\x93\xe9\x70\xe5\xcc\x96\x56\x99\xd4\xb7\x7a\xb1\xe0\x70\x5d\xe7\xbe\xa1\xe6\x58\x15\xdf\xa3\x0d\x6a\xa8\x79\xaa\x99\x2e\x73\xe3\xbb\xca\xe3\x57\x5e\xa7\xba\x46\xa8\x19\x21\xbd\x46\x4d\x5a\x55\x0b\x5c\xc8\xdd\xd1\xce\xed\x07\xdc\x95\xd9\x5f\xeb\xc6\xec\x89\xb7\xc1\x08\xb4\xca\x14\xd6\xa9\x46\x6f\xea\x9e\x36\x78\x50\x2b\xac\x18\xda\xa0\x15\xa3\xd6\xff\x0b\xa5\xbe\xce\xaa\x52\x73\x74\xdd\x2c\xab\xad\x11\xd6\xe9\x3a\x35\x39\x50\x53\x73\xd4\x50\xcf\x51\xed\xce\xb8\xad\x73\xe7\xdb\x30\x0d\xf5\x7c\x27\x88\xe9\x10\x0b\x6e\x81\xd6\xf8\x9a\xdd\x19\xe8\x1a\x35\xd5\x8b\x5c\x9a\xd6\x5b\x63\x5f\x4b\x3b\x66\x16\x5a\x4d\x5a\xa3\x0c\xf8\x73\x5a\x67\x28\xa3\xa4\x6b\xea\xaa\xd1\x07\x3b\xeb\x56\x28\xa6\x79\x74\x43\xdf\xa8\x3a\x4f\xa1\x4a\x2b\xb4\xaa\xae\xb6\x87\xdb\x33\x6f\xf9\xf6\xac\x96\x9e\x37\xf4\x85\xe0\xab\xba\x4c\x4d\xba\x61\x2f\x46\x0d\xf4\x8a\xa6\x52\xd5\x0e\xbd\x65\x74\xd4\x75\x2d\x98\xd4\x2b\x54\x57\xf3\x6a\x49\xbd\xd0\x76\x2f\xf1\x6c\xd5\x92\xd0\x48\x66\xc0\xba\x78\x7e\x19\xf4\x2a\x35\xd5\x52\xfb\xba\xac\xa8\x6a\x7b\xf9\x7d\x76\x34\xfb\x86\xaa\x9a\x40\x53\x3d\xaf\xaf\xca\x7a\x3b\xc4\xdf\x92\x8e\x97\x5a\xb4\xd4\xa2\xee\x5e\x55\x5d\xe1\xff\xea\x7b\xc7\x53\x5b\xc4\xd3\x5b\xc4\xc7\x7a\xc7\x5b\x32\x58\xe8\xbd\x6f\x67\xba\xd5\x6e\x97\x2a\xb6\xaa\xa1\xe3\x2b\xa3\x7f\x73\xd0\x4b\xfc\x9d\x33\xbf\x3b\x6e\xf9\xe9\x0e\x7f\xac\xed\xd3\x77\xf4\x08\xda\xd3\xfe\x6f\x00\x00\x00\xff\xff\x48\x0a\xc3\x20\x20\x10\x00\x00"

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
	"templates/views/widget.html":      templatesViewsWidgetHtml,
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
