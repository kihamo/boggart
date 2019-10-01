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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\x5b\x6f\xdb\x38\x16\x7e\xef\xaf\x38\xab\xc9\xae\x6d\x20\x92\x9c\x74\x66\x76\xe1\xb1\x33\x68\x9b\x0c\xba\x40\x2f\x41\x92\x76\xd1\x7d\x29\x68\x89\xb2\x99\x4a\xa4\x86\xa4\x7c\xd9\xc0\xff\x7d\x41\x52\x96\x65\x47\xb2\x25\x5f\x1a\x67\xb7\x7e\x48\x04\x89\xfc\x78\x2e\xdf\x39\xa4\x0e\xa9\x87\x07\xf0\x71\x40\x28\x06\xcb\x63\x54\x62\x2a\x2d\x98\xcd\x5e\x74\x7d\x32\x02\x2f\x44\x42\xf4\xac\x18\x0d\xb0\x2d\x89\x0c\xb1\x75\xf1\x02\x00\x20\xff\x30\x60\x3c\xb2\x07\x9c\x25\x31\xc4\x49\x18\xda\x9c\x0c\x86\x32\x6d\xb7\xda\x96\xd0\x38\x91\xa6\x71\xae\x85\x6e\x15\xa2\x3e\x0e\xe7\xed\xfa\x92\x16\xb6\xd2\x2d\x11\x0c\x39\x0e\x7a\xd6\xef\xc8\x93\x84\xd1\x9e\xc7\x68\x40\x06\xc2\xc6\x93\x98\x71\x69\xe5\x30\x40\xe1\x10\x8f\x99\x0b\x91\x78\x1e\x16\xc2\x5c\x47\x05\xc8\x1a\x9d\x64\x7a\x21\x01\x01\xb2\x3d\x36\xb0\x40\xeb\xde\xb3\x1e\x1e\x80\x9c\xfd\x83\x82\x75\xa5\x87\x02\x33\x72\xc2\x91\x92\xc3\x02\x07\x66\x33\xeb\xa2\xeb\x92\x02\xa1\x5d\x54\x41\x93\x90\x1d\x4c\x8d\x21\x46\xbe\xc0\xb2\x54\x15\x35\x74\x3d\x0d\xba\xae\x76\x59\xce\xd1\xae\x4f\x46\x29\x3f\xcc\x65\xfa\x6f\x89\x4b\x5e\x88\x11\x0f\xc8\x44\x0d\xf3\xf8\x29\x67\xe3\x02\x8a\x79\x2c\xb4\x23\xdf\x3e\x3b\x07\x75\x25\xa2\xf9\xd5\x44\xd8\x67\xe7\x25\x54\x9b\x7c\x8d\x11\xc5\xe1\x2a\xcd\x96\x5a\xcc\xf9\x5e\xa0\xad\x6a\xc7\x99\xb2\x94\x44\xfd\x22\xa4\xac\x65\x92\xb1\x96\xa2\x11\x50\x34\xb2\x25\xea\x0b\xe8\x23\xfe\x55\x5d\x58\x0b\x98\x90\x88\xa2\xb1\x32\xa4\x90\xa4\x6d\x63\x8e\x05\xa6\xd2\xd0\x4a\x79\x2a\x00\xfc\x27\x38\x86\x25\x60\xa9\xf8\x9c\x8f\xa9\xee\x8d\xb0\x6a\x84\xa9\x0f\xb3\xd9\x45\x46\xaa\xe2\x7e\xf7\x68\x84\x84\xc7\x49\x2c\x3b\x23\x46\xfc\x66\xbb\xf5\x9b\xea\x1b\x0a\x0c\xb3\xd9\xc3\x03\x38\x37\xf8\xcf\x04\x0b\xe9\x7c\xba\x79\xe7\x5c\x23\x39\x34\xb7\x0d\xb8\x75\x91\xd1\xe6\x76\x2a\x24\x8e\x0c\x63\x14\x35\xba\x6e\x58\x40\x9a\xda\xaa\x19\x1a\x6e\xa7\xde\xbc\xef\x36\x2a\xe6\x23\xb0\x48\xdd\x77\x59\x78\xec\x51\xd9\x80\x84\x78\x6b\x6d\xb3\xce\xbb\xa8\xab\x41\x8a\xf4\xfd\xc3\xa0\x6f\x56\xb8\xeb\x26\x61\xc9\x93\x5c\xa8\x49\xd4\xb7\xcb\x83\x6d\xa9\xc7\x4a\xd0\xe5\x11\xd4\x1d\x08\x90\x8f\xc1\x18\x0a\x08\x5d\x83\x56\x4e\x90\x72\x01\x54\x90\xe2\xdc\x90\x21\x06\xfd\xd7\x16\x92\x93\x18\xfb\xe0\x23\x89\xcc\x7d\x5f\xda\x1c\x8b\x98\x51\xa1\x24\xa1\x6c\xcc\x51\x6c\x81\x90\x53\x25\xfe\x98\xf8\x72\xd8\x39\x6b\xb7\xff\xba\x46\x40\x33\xa2\x4a\xcc\xeb\xdb\x98\x76\x7c\x73\xa3\x14\x70\xe1\xc5\x9f\xe6\x1e\x94\xc3\x2d\x7a\xdf\x91\x08\xef\x04\xf0\x1e\x0b\x81\x06\xbb\x61\xdc\x4d\xe3\xdd\x00\x3e\x09\xcc\xeb\x00\x74\xdd\x4d\xa6\x56\x38\x1b\x9d\xd6\x95\x7d\xe6\x4f\x37\x0f\xf7\xf0\x00\x1c\xd1\x01\x86\x13\x72\x0a\x27\x21\x1b\x40\xa7\x07\x8e\x62\xea\x3a\xa2\x2e\x46\xa9\x4c\x0b\x5f\x99\x44\x0d\xe0\x5c\x33\x41\x74\x40\x68\x93\x54\x20\xdf\x1c\xa0\x6b\xd2\x0c\xc8\x69\x8c\x7b\x16\x8a\xe3\x90\x78\x3a\xb7\xb9\x8b\x1c\x64\x5d\xf8\xcc\x4b\x22\x4c\xa5\x33\xe6\x44\xe2\xa6\x8f\x24\xbe\x63\xb7\x92\x13\x3a\x68\x36\xe6\x32\x28\x6e\x39\x7f\x30\x1e\x21\x09\xd6\x79\xbb\xfd\xab\xdd\x3e\xb3\xdb\xe7\x77\x67\xbf\x74\xda\x3f\x77\xda\xbf\xfc\xbb\xfd\xf7\x4e\xbb\xad\x82\xb5\xd1\x6a\x75\x5d\x83\x7d\x51\x4f\xdc\xf9\x58\x97\x48\xa2\xda\xba\x66\x82\x4e\x63\xbc\x75\x67\xc5\xbd\xca\x9d\x37\x13\x0f\x0c\x5f\x4c\xb2\xde\x44\xd1\xf5\xf4\xeb\xba\x3a\x8b\xad\xcd\x9e\x7a\x0a\x29\x9b\x75\xfe\x6f\x72\xe8\x07\xb4\x63\x16\x7c\x8d\x07\x84\x82\xdc\x35\x99\x5e\x51\x7f\x77\x90\x5b\xf2\x9f\xdd\x00\x5e\x69\x1e\x88\x67\x93\x50\x15\x5f\x75\x46\xd5\xc4\x3d\x48\x4a\x55\xc8\x8e\x5a\x2f\x29\xa6\x3c\x55\x4e\xd5\x42\x68\xa6\x7d\x8f\xcc\xba\x4f\xa1\xaf\xa8\xff\x9d\x26\x83\x61\x12\x21\xfa\xb5\x3f\x95\x58\xe4\xbc\xf6\x0e\xd3\x81\x5e\x13\xd7\xc2\xcb\xf2\x1b\x9e\x48\xdb\xc3\x54\x62\xbe\x21\x4b\x2d\x01\xe4\x96\xc6\x59\xa9\x03\xb2\x2b\x7b\x22\x6a\x80\x41\xbe\x92\xa0\xcc\xba\x6e\xc1\xef\xb3\x31\x0d\x19\xf2\xff\x46\x51\x84\x7b\x45\xf4\x7d\x54\x78\xc8\xd7\x1b\x88\xc7\xa8\x75\xf1\xa8\xb8\xa0\x30\xec\x39\xf4\xe3\x12\xc3\x65\xf6\xe4\x24\xab\x2f\x14\x16\x44\x4a\xd5\x5b\x54\x16\x36\x37\x7d\xbe\x73\xed\xba\x59\x35\x60\x3c\x4a\xdf\x8d\xd4\xa5\x05\xa9\x43\xad\xdf\x2d\x88\xb0\x1c\x32\xbf\x67\xc5\x4c\x48\x0b\x88\xdf\xb3\x04\x96\x92\x50\xf5\xb6\xa3\x26\x5b\x5b\xb2\xc1\x40\xf5\x1c\xa1\x90\xf8\x48\xb2\x4d\x5c\x2d\x9a\xc3\xf7\x33\x57\xc3\x81\xe6\xeb\xb9\xac\x91\x6f\x7b\x2c\xb4\xcf\x72\xaf\xb2\x1f\xe3\x45\x75\x6e\xcb\x59\xef\x33\x0a\x93\x5a\xf3\x66\x35\x7e\x55\x9a\xf7\xa0\xc6\xdc\x07\xd9\x6b\xaf\x23\x74\x7d\xe6\x2b\xa1\x01\xab\x32\xeb\x41\x1d\x7b\xc3\x22\xa7\xa6\xcb\x0a\xcc\x09\x0a\x81\x26\x51\x3f\xf7\xb2\x55\x31\x9b\xe6\xd0\xf2\x52\x3b\x06\xf4\x03\xab\x97\x9a\x2b\x19\x1e\x76\xd2\xf6\x53\xbc\xb4\x0e\xdb\x51\xcd\x4b\x3c\x22\x1e\xbe\x49\xf4\xc4\x0d\xb3\x19\x64\xe3\x44\x84\x26\x72\x51\x85\x39\x26\x13\xdc\xb2\x40\x8e\x11\xc7\x30\xc2\x5c\xe4\xa3\x6b\x57\x9f\xb3\x40\xfe\x0b\x71\xfc\xd9\xc0\x1e\x9d\xde\x6f\x11\xf7\x95\xde\xfb\xd2\x57\xe1\x29\x7d\x8f\x56\xd1\x7d\x3b\x78\xae\xf0\xb1\x3a\xf8\x8a\x7a\x7c\x1a\xcb\x7d\xab\x9d\xc2\x1e\xab\xd6\xaf\x13\x12\xae\xbc\x5d\xd6\x54\x78\x3f\x6f\x04\x4b\x36\xd3\x42\x1d\xfc\xd5\xe0\xbb\xcc\x17\x4a\xd1\x63\x33\xaf\x96\x0a\xfd\x2f\xd8\xf7\x55\x88\x78\x04\x84\x82\x37\x44\x54\xef\x19\xec\x27\x68\x35\xee\x3f\xe9\x1b\x83\x7a\x74\x41\x6b\xd4\x66\x89\x3c\x88\xde\x1f\x13\x79\xac\x8a\xdf\xa1\xf0\xdb\x01\xdc\xad\x60\x8f\xd7\xdb\x5a\xe9\x03\x38\x5b\xe1\x1e\xb1\xaf\x3f\x13\x1f\xb3\x03\x38\x5b\xe3\x1e\xaf\xb7\x8d\xda\x07\x70\xb7\x06\x3e\x62\x7f\xbf\x4a\x7c\x72\x08\x7f\x6b\xdc\xe3\xf5\xf7\x25\x19\xec\x5b\xe3\x4b\x32\x38\x56\x6d\xaf\x26\x92\xa3\x7d\xeb\xab\x41\x0f\xab\xf1\xa2\x1e\x58\xa7\x12\xa3\xd6\x7d\x5f\xbd\x84\x73\x4c\xe5\xa1\x4a\x31\xc0\xd9\x58\xc4\x88\xf6\xac\x73\xab\x64\x03\xbf\xa6\x79\xeb\x95\x9f\x57\xce\x6a\x4d\x84\xfd\x73\xcd\x0a\xb6\x86\xd1\xc7\x03\xd3\x45\xae\xc4\x93\xc5\x21\x38\x7d\xc6\xd0\x63\x54\x72\x16\x5a\xa0\x6b\xd7\x96\xb2\xab\x9d\xda\xd5\x82\x11\x0a\x13\x53\x73\x5e\x32\xf8\xd2\xda\xd6\x69\x9f\x39\xed\x73\x48\xd7\xb6\xbf\x5a\xba\xda\x4d\xfc\x55\x28\x8e\x91\xcf\x68\x38\x05\xb7\xa6\x15\xaa\x57\xa9\xa1\x72\xa5\x1a\x0e\x1b\x93\x4f\xe2\xe6\x4a\x07\x42\x2b\xa3\xd5\x26\xcd\x1e\xc8\x62\xd5\xe5\x46\x26\xad\x8a\xd2\x02\xe5\xed\xbe\x5c\x77\x96\x69\x23\x6c\x3f\x91\x92\xd1\x47\x67\x46\x69\xc0\xac\xd4\x34\xa6\x45\x8e\xef\x94\x8d\x73\xb9\xe2\x03\x1b\xcf\x53\x85\x69\xb9\xa5\x7e\xae\x52\x70\x0b\x4e\xd4\x8b\x9d\x2d\xba\x1c\x6a\x42\xd8\xb0\x71\xb3\x04\xbc\xfd\xea\xff\x09\x33\xf9\xcb\x1f\x21\x5e\x4f\xda\x1f\x21\x5e\xd2\xef\xb9\x86\xf8\xe6\x3d\x60\xa8\xb2\x0f\x0c\x55\xf6\x82\xbb\xae\xe2\xf4\x45\xf9\x0a\x73\xbe\xbf\x27\x19\x47\x03\xbc\x71\x83\x6f\xf9\xc0\x7b\xfe\x1b\x8e\xd2\x1e\xc3\xf3\xdc\x46\x8f\x19\x66\xce\x9b\xe1\xf9\x86\xbe\xeb\x0e\xfa\x97\xab\xbc\xfe\xf1\xc1\x36\xa3\xbf\xf3\x46\xf4\x35\xe2\x72\x75\x8f\x74\x9b\xdd\xe8\x1b\x1c\x21\x42\x41\xc4\xc8\xdb\xf1\x6c\x2d\x93\x28\xac\x8f\xf3\x64\x07\xba\xc4\x29\x9c\xa4\xbc\xd7\x67\xba\xaa\xc6\x40\x31\x5a\x7c\x0a\x27\x31\xe2\xd2\x9c\x85\xed\xf4\x32\x6c\xe7\x3a\xbb\x7b\xc0\x89\x7d\x69\xb4\x0f\x9a\x15\xdb\xbe\x07\x2f\x1d\x6c\xca\x54\x72\x0c\x4f\x6e\x95\x7b\x9d\x4f\x84\xca\x97\xe7\x7b\x1e\x40\xf3\x67\x6b\xfc\x43\xe4\xe0\xef\x79\x62\x67\x43\x25\x20\xcd\xd3\x69\x95\x63\xed\xd1\xc3\x1d\x73\x74\x5a\xf4\x10\x3f\x92\xf4\x1e\x92\xf4\x9b\x95\xaa\xd4\x76\x1f\x4e\x28\x0f\xee\x74\x66\x98\x48\x8e\xe4\x6e\x18\x37\xd8\x63\xdc\x7f\x16\x89\xdd\x3b\x85\x93\x34\x50\x74\x62\xaf\x12\x34\x8b\xd1\x6a\x9e\xd7\x4d\xc1\x9d\x2d\x72\xee\x2a\x84\x76\xf4\x4e\x08\xa9\xa3\x97\xce\xfd\xf4\x89\x14\xae\xc0\x1e\xa3\xe9\x89\xc9\xda\xf0\x24\x58\x8c\x60\x68\x60\x3e\x05\x33\xf8\x5c\xdf\x21\x74\x60\xd0\xb3\x5c\xf6\x5c\x8f\x4e\x6e\x4e\xc4\x1b\x44\x59\x93\xdb\x4a\x1e\x15\xdc\x5e\xb9\x55\xfa\x31\x6c\x4e\xe2\xdc\x97\xd7\x2a\xaa\xb2\xcf\x2a\x8a\xbf\xdc\x7c\x91\x53\x48\x48\x24\x89\xf7\xf6\xee\xfd\x3b\x68\x9a\xeb\x4f\x37\xef\xc0\x72\x7d\x24\x86\x7d\x86\xb8\xef\x22\x21\xb0\x14\xee\x08\x53\x9f\x71\xe1\xf6\x19\x93\x42\x72\x14\xdb\x3e\x92\x58\x87\x5d\x4c\xbc\x6f\x98\xbb\x9e\x10\xee\xca\x3d\x27\x22\xd4\xf1\x84\xb0\x20\x40\xa1\xc0\xad\x9c\x5c\xa9\xe8\xc5\x4a\xdc\x8b\x43\xaa\x10\xb1\x08\x53\xe9\xde\xcf\xaf\xb4\x94\xf7\x8f\x84\xdc\xbf\x85\xee\x8b\x0d\xb4\x3c\xf4\xc2\xf1\x95\x8e\x6a\x2c\x51\x27\x48\xa8\x31\x92\xc0\xf2\x8e\x44\xb8\x29\x5b\xf0\xf0\x88\x73\x27\x0e\xba\x47\x93\xe6\xe3\x07\xea\xa7\x46\xeb\x40\xe3\xfa\xe3\xed\x5d\xe3\xb4\xb0\x45\xc2\xc3\x0e\x9c\x34\x1b\x3f\xcd\xcf\x06\x37\x5a\x4e\xcc\x59\xdc\x6c\x18\x0f\x35\x5a\xc5\xfd\x7c\x24\x51\xa7\x40\x9c\x6c\x64\x12\xe1\x0e\xc8\xc2\xe7\xb3\x62\xc8\xf4\x44\x79\x67\xa1\x78\x93\x17\x69\x3c\xff\x91\x00\x9a\xdc\xe1\x58\x24\xa1\x84\x5e\xaf\x07\x8d\x00\x91\x10\xfb\x8d\x75\x9d\xd4\x8f\xe2\x31\x5c\x7f\x60\x92\x04\xd3\x12\xbb\x2d\x6b\x22\x43\x65\xc4\x2b\xce\x19\x2f\xb1\xe2\x52\x7b\x3c\x91\x1d\xe0\x4e\x64\x3e\x73\xac\xd0\xc1\x38\x09\x57\xc4\x1f\x12\x1f\x77\x0c\xc5\x36\x37\x56\xcb\x2a\x42\x07\x1d\x68\x64\x44\x7e\xd9\x58\xdb\x6b\xd6\xfa\xad\xf4\xf9\x2c\xfb\xf8\xab\x99\xe9\x07\x7f\x51\x96\x4f\xa8\x89\xf6\x43\x19\xff\xd6\x70\xe3\x80\xe6\x17\x95\x47\x78\x4a\x07\x14\xc7\xd3\xa3\xbb\xab\x18\x2b\xf3\xe0\x49\x73\x7e\x26\xac\xe5\x70\x8c\xfc\x69\x73\x11\x72\x85\x39\xa6\xd9\xf8\x49\x05\x74\xa3\xe5\xac\xa4\xbc\x12\x17\x0a\x42\x07\x21\xbe\x44\x12\x5f\xeb\x66\x1d\x90\x3c\x29\xb1\x96\xca\x85\xd3\xd7\xba\xf2\xf8\x46\xad\xbf\xb1\x50\xc6\x5a\x7c\x61\x52\xe2\x91\x90\x79\x48\x51\xa3\x9c\x43\x81\xae\xe9\x76\xa0\xf1\xe5\xcb\x97\x2f\xce\xfb\xf7\xce\xe5\x25\xbc\x7d\xdb\x89\xa2\x8e\x10\xc5\x2e\x28\x30\xe4\x69\x96\x8d\xd4\x7c\xc1\x0b\x53\xb0\xd6\x38\xcd\xd1\x66\x12\x4a\xdb\x3a\x46\x84\x66\xa1\x04\xad\x56\xa1\xdb\x4a\x8d\x6f\x53\x36\x6e\xb4\x1c\x2f\x24\xde\xb7\x0d\xfe\xca\x0b\xd4\xd0\xdd\x0a\x86\x7a\xb1\x3a\xf2\xfc\x3a\x3b\xcf\x57\x3a\xbb\xff\x37\x00\x00\xff\xff\x26\xa9\x5c\x17\x1c\x46\x00\x00"

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
