// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
// DO NOT EDIT!

package homie

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5c\xef\x8e\xdb\x36\x12\xff\x9e\xa7\x98\xaa\x5b\xd8\x6e\xd7\x96\x77\xd3\xe6\xee\x1c\xdb\x41\xae\xb9\xa2\x3d\xa4\x49\xd0\x6c\x5b\xe0\x72\xb9\x05\x2d\xd2\x36\x37\x14\xa9\x92\x94\xbd\xce\xc2\xef\x7e\x20\xa9\x7f\x96\x25\x5b\x76\xbc\x69\x03\x44\x40\xb7\x12\x35\x9c\x21\x7f\x33\x9c\x19\x8e\xe8\xdc\xdd\x01\x26\x53\xca\x09\x78\x81\xe0\x9a\x70\xed\xc1\x7a\xfd\x60\x88\xe9\x02\x02\x86\x94\x1a\x79\x11\x9a\x91\xae\xa6\x9a\x11\x6f\xfc\x00\x00\xa0\xf8\xd2\xb6\x5f\x33\x32\xd5\xc9\x4b\x4b\x30\x7f\x98\x3f\x98\xeb\xee\x0e\xe8\xc5\xdf\x39\x78\xcf\xc8\x82\x06\x04\xbe\x52\x1e\x9c\x01\xa7\x2c\xfb\xaf\xc7\x51\x48\x60\xbd\x1e\x4e\x24\xf8\xe3\x07\x5b\xbd\xa7\xd0\x63\x48\xe9\xeb\x38\xc2\x48\x1b\xc2\x0d\x0a\x2b\x54\x85\x88\xb1\xf1\x50\x05\x92\x46\x1a\xf4\x2a\x22\x23\x0f\x45\x11\xa3\x01\xd2\x54\x70\xff\x06\x2d\x90\x7b\xe9\x8d\xb1\x08\xe2\x90\x70\xdd\x5b\x4a\xaa\x49\xdb\xf0\xbc\x12\xaf\xb5\xa4\x7c\xd6\x6e\xdd\xdd\x6d\x08\xeb\xfd\x20\x64\x88\x34\x78\x97\xfd\xfe\xa3\x6e\xff\xa2\xdb\xbf\xbc\xba\xf8\x6e\xd0\xff\x76\xd0\xff\xee\x3f\xfd\xbf\x0d\xfa\x7d\x03\x59\xab\xd3\x19\xfa\x8e\xfd\x78\xe8\xbb\xb1\x94\x67\x41\x38\x2e\x8e\x7c\xe8\xa7\x38\x0d\x7d\x4c\x17\xc9\xac\xb7\xd1\x95\x74\x36\xdf\x80\xb7\x40\x31\x15\x32\xec\xce\xa4\x88\x23\x88\x62\xc6\xba\x65\xda\x32\x3d\xe5\x51\xac\x5d\x87\x12\x95\xa5\x64\x68\x42\x58\x4a\x3b\xd1\xbc\x96\xb2\xa0\x17\xc1\x99\xb1\x9f\x0a\x95\x58\x96\x93\x58\x6b\xc1\x13\x7d\xb8\x07\xaf\x20\x01\x8c\x14\x1a\x08\x77\xb3\x44\x92\x53\x3e\xb3\xf7\x2a\xf4\x2a\x39\xa6\x17\x46\x1a\x75\xb5\x98\xcd\x18\x19\x79\xa1\xc0\x88\x79\x49\x1b\x92\x33\xa2\x47\xde\x97\xae\x71\x3f\x13\x4b\xe7\x4c\x7c\xe4\x65\xb6\xfa\xbd\xe0\x53\x2a\x43\x90\x44\x69\x24\x35\x60\x6b\xbb\x1e\xf4\x60\xbd\x6e\xcc\x74\x22\xf0\xaa\xc0\xf3\x77\xca\x18\x4c\x48\x89\x67\x0f\x9e\x4a\x02\x2b\x11\x83\x8a\x25\x79\x72\xa0\x88\x00\x31\x36\x41\xc1\xbb\x91\x87\x02\x63\xe8\xed\x56\xc2\xbd\xd5\x79\x5c\xa3\x39\xab\x18\x9a\xd9\x10\x52\x30\x45\xdd\x48\x2c\x89\xec\x8a\xe9\xd4\x83\x32\x14\xbf\x38\x86\xc9\xc0\xc6\x43\x9f\x56\xf3\x1d\xfa\x4e\xc1\xb5\xf6\x52\x5a\x01\x1b\x7d\x0f\x32\x14\x8c\xf8\x8c\xc8\xbf\x9a\x9d\x90\xd3\x5a\xc9\x92\x46\x04\x02\xc3\x7e\x76\x7f\x96\x42\x0e\xb5\x13\x22\x91\x22\x72\xdb\x48\x7e\xa7\x11\xf9\x00\x0b\x19\xfa\xd6\xf9\x94\x5c\x97\x73\x8c\x15\x8f\xc9\x6d\xea\x39\x8b\x3e\x2e\x60\x04\xc9\x29\xbd\x35\xc3\xd8\x7e\x2b\xc5\xb2\x22\x8e\x05\x82\x75\x43\xdc\xbd\xb8\x04\x73\xa7\xc2\xf4\xee\x56\x75\x2f\x2e\x6b\x7c\xef\xed\x75\x84\x38\x61\x3b\xbc\xed\xed\x75\x31\x6e\x6e\x50\xcd\x2f\xc7\xe5\xa8\x88\xb4\x96\x74\x12\x6b\xa2\x1c\x8c\x43\x7f\x7e\x59\xd1\x33\xce\x1c\x34\x47\x0b\xe0\x68\x31\x41\xd2\xf9\x7d\xb0\xe3\xb9\xd6\x42\xb0\x89\xb8\xad\xd1\xe9\x90\xd1\xf1\x10\x15\xe6\xcd\x50\xa4\x48\x97\x51\xfe\xce\x1b\x17\x95\x6d\x74\x1d\xcc\xc9\x42\x0a\xde\x35\x31\xc0\xe8\x74\xe8\xa3\xf1\xd0\x67\x15\xca\x1d\xfa\x31\xab\x68\xdd\xa5\x96\x1d\x8a\xde\x06\x32\xcd\x4e\x2a\x64\x68\x34\x61\x24\x0b\x99\xf6\xc1\xfe\xed\x2a\x2d\x69\x44\xb0\x5d\x02\xae\x1d\xeb\xae\x24\x2a\x12\x5c\xd1\x05\x01\x2e\x96\x12\x45\x1e\x28\xbd\x32\x96\xbc\xa4\x58\xcf\x07\x17\xfd\xfe\x57\x75\xd0\xe9\x39\x41\xb8\xee\x9d\xdc\xb1\x86\xf4\x3c\x1d\x5f\x88\xbb\xc6\xae\x2e\xbd\x5c\xfd\x4f\x53\xbd\xa7\x6a\xd7\xf3\x9d\xac\xf2\x9e\xbf\x21\x16\xef\xed\x35\xf4\xeb\x86\x66\xfa\xec\x98\x90\xf1\x4a\xf5\x03\xb9\xbb\x03\x69\x9c\x30\x9c\x99\xd4\xed\x1c\xce\x16\x66\x30\x30\x18\x41\xcf\xf9\x2a\x75\x9d\x1b\x74\x9d\xd3\x87\x7d\xc8\x39\x02\x6c\xa6\x7c\x96\xa6\x88\xbe\xae\x19\x72\xb1\xc3\x4e\x02\xc8\xb2\x17\xf2\x47\xc2\xd7\xa3\x61\xc4\x88\xc9\x09\x6d\xb6\xd8\x73\x6e\xd7\xdb\x35\xf0\x0d\x91\x91\x24\xa9\x1d\x45\x08\x63\xca\x67\x83\xbe\x37\x1e\x06\x02\x67\xa6\x79\xa3\x04\xf7\x8e\x4b\x50\xff\xfd\xfa\xe5\x8b\x9e\xb2\xf9\x29\x9d\xae\xdc\x63\x84\xa4\x22\x6d\xab\x86\x65\x8a\xfe\x7a\xdd\x39\x07\x1e\x33\x76\x0e\x97\x1b\x09\xa9\x19\xc7\x78\xe8\x47\x92\x34\x42\x86\x30\x55\x9b\xd1\x55\x90\x67\xd2\x1b\xf1\xae\xcf\x01\x32\x34\x77\xaa\xb8\xde\xa0\xf7\x4b\x18\xfa\x35\x66\x3d\xf4\xad\x87\x38\x32\xfc\x24\xa9\x30\x17\xd8\xd9\xfa\x5f\x3c\xe0\xbc\x30\xe3\xfc\x1c\x64\x3e\xb9\x20\x03\xcd\xdc\x65\x21\x42\xfc\xf4\xac\x49\x50\xd9\xea\xf6\x02\x85\x8d\xa2\xd1\x56\xc7\xab\x55\x74\x5c\xc7\xa7\x52\xa2\xd5\x01\x3d\x53\x3c\x31\x55\x11\x43\xab\x01\x70\xc1\x89\xd1\xfa\xce\xe0\x79\xfa\x50\x58\x08\x83\xe7\x70\x66\x1c\x80\x0d\x81\x99\x27\xa8\x9f\x42\xd3\xa8\x27\x30\xe9\xfd\xf4\xac\xe7\xaa\x13\x8d\xe3\x5f\xd6\xd5\xa8\xf2\xe8\xce\x46\x9d\x47\x77\xb6\x2a\x3d\xb4\x77\x93\x08\x72\x16\x49\x11\x11\xa9\x29\x51\x06\x6c\x27\xed\x55\xde\xd6\x2c\x0e\xd1\xe9\x06\xa3\x06\x9d\xf6\x3b\x81\x64\xb1\xef\x9f\x05\x34\x58\xe9\xdb\xf4\x7b\x4c\xa6\x42\x40\xbe\xc0\x12\x7c\x56\x1e\x9c\x35\x5a\x63\x3b\x99\x25\x89\xe7\xe1\x9c\x76\xc7\xee\x12\x65\x73\x78\xf6\x25\xab\xe5\x2b\x5f\xb5\xd1\x79\x66\x05\x2b\x6b\x4c\x87\x99\x44\x3e\x80\x83\x75\xe3\x56\x4a\x2a\xfb\xe0\x75\xba\x97\xa1\xd5\x51\xce\x11\x36\x5e\xfe\xca\xa9\x3e\x56\x5a\x73\x1d\x42\xe3\x8c\xaf\xc0\xbb\x99\x26\x2b\xf3\xb5\xe3\xe5\x1f\x9b\x71\xfe\x19\xd9\x66\x2e\x33\xcd\x3b\x15\xd1\x9a\xf2\xd9\xfd\xa7\x9e\x90\xa6\x72\xf8\xc3\x93\xd0\xef\x93\x1d\xd6\xe7\x2c\xb4\x79\x16\x3a\x15\x32\x04\x29\x4c\xe2\x63\x6e\x3d\x70\x85\xbc\x91\xf7\x24\xb9\x49\x2d\xc1\x83\x90\xe8\xb9\xc0\x23\x2f\x12\x4a\x7b\x40\xf1\xc8\xcb\xdf\x6d\x94\x44\x17\x88\x51\x8c\xb4\x90\xb5\x09\xe9\xc7\x4a\x7d\xe1\x44\xe9\x6f\xa9\xd6\x72\x51\xa8\xb5\xbc\x8c\x0c\x4c\x47\x65\xa8\x8d\x8a\x2d\xb0\xd7\x3f\xee\x8d\x6b\xcd\x0b\x2f\xf4\x1c\xce\x84\x9d\x90\xcd\x39\x8b\x5e\xa0\x9e\x79\xc3\xb4\xd3\xf1\xb5\x51\xe9\x34\xb9\x5b\x5e\x6f\x49\x58\x9b\xdc\x12\xbc\x89\x10\xac\x51\x95\x65\x68\x3f\x5d\x25\x65\x93\x60\x4e\x82\x77\x66\xe1\xe7\xb5\x95\xae\x5a\x52\x1d\xcc\x3d\xe0\x28\x74\x65\xe9\xd2\x1c\xdc\x22\xa8\x6a\xcf\x86\xd6\x8e\x24\xe5\x3a\x23\xb0\x1a\xef\x80\xa7\xa5\x51\xfc\x7a\x0d\x56\x2c\xc1\x99\x03\x06\x7f\xef\x9c\x6d\x25\x85\x4e\x41\x48\x68\x6f\x4d\x9e\x72\xed\x75\xaa\xdb\x1f\x7d\xeb\x75\x0e\x86\x85\xc7\xe1\x84\x48\x6f\xe3\xe3\xa0\xf1\x26\xd2\x60\x5c\x8f\x8b\xad\xde\x6c\xbc\xfa\x2d\xa9\xe7\xd4\x63\x16\x21\xad\x89\xe4\x23\xef\x7f\x6f\xba\xdf\xbc\x7d\xf2\xa6\xdf\xfd\xc7\xdb\xaf\xcf\xbc\x0f\x04\x24\xae\x43\x24\xfe\xb4\x20\x39\x0a\x91\xad\x49\x4f\x99\x40\x66\xd6\x9f\xc8\xa4\x73\x3b\x68\xff\xb7\xe7\x6e\x3a\x4f\x0e\x01\x00\x71\x5c\xa1\x7a\x57\xf6\x34\x56\x31\x23\xd0\x66\x84\x97\xd7\xe7\x45\xbf\xdf\xcc\x2e\x34\xb9\xd5\x48\x12\x54\x85\x0b\x48\xa2\xe8\x7b\x13\xc4\xae\x53\x32\x0f\xa4\x58\xaa\x91\xf7\xed\x11\x3e\x65\x5c\x01\xe2\xd0\x4f\x39\x37\x43\xe4\x50\xad\x1b\xee\x1f\x43\xe7\x0d\x14\xfa\x67\x96\x78\x61\xdf\x36\xa2\x76\xeb\x30\xf4\x0d\x6a\xa7\x4a\xca\x85\x46\xd7\x84\x1b\x49\xf8\x93\xca\xcb\x5f\x5e\x3d\xfd\x9c\x94\x9f\xb4\x34\x7c\x78\x09\x78\x77\x06\x78\x58\x09\xf8\xb5\x46\x3a\x56\xcd\x53\x5e\x9b\xff\xa5\x06\x2c\x63\xce\x5d\x8d\x20\xe3\x97\x34\x39\x86\xb9\xaf\xca\xdf\x13\x84\x57\xf9\x5b\xbb\x28\x3e\x64\xb5\x1f\x36\xdb\x57\x52\xcc\x24\x51\x07\xcd\x77\xbf\x9b\x2d\x1e\x87\x4b\x05\x34\xac\x4c\x55\x74\xed\x4e\x90\x84\xe2\x43\x66\x29\x95\xc8\x6f\x50\x52\x3e\x15\x76\xd7\xb7\x20\x39\xf6\x1b\x14\xee\x68\x8c\x83\xde\x78\x71\xb7\x5d\x4c\x49\x26\x48\x7a\x80\x24\x45\x5d\xeb\xf9\xb9\x58\x5a\x1f\x6f\x45\x2e\x25\xd5\x9a\x70\xeb\xfb\x73\x92\x90\xf2\x91\xd7\xdf\x68\x41\xb7\x79\x27\x2d\x34\x62\xb6\xcb\x86\x79\xa7\xaf\x53\xb9\xb0\x5e\xef\xda\xf6\x6d\xc1\xa6\x22\xc4\x53\xdc\x94\xec\x0a\xce\x56\x36\xa8\x96\x06\x0a\x93\x95\x26\x0a\xc4\x14\xca\xe3\x19\xfa\x86\x45\xd3\x5a\xe3\xd6\xea\x3f\x92\xec\xe3\x19\xfa\xf7\x66\x2f\xa2\xe2\xf0\xc0\x85\x6d\x31\x0a\x92\xbe\x1f\xb0\x32\xf7\x97\xb6\x6a\xaa\x17\xc5\xfc\x64\x2e\x24\x7d\x2f\xb8\x51\x98\x7d\xb6\x27\x80\xba\x8c\x4c\x35\x60\x29\xa2\xf7\x82\x13\x6f\xa3\xdc\xb1\x5d\xd4\x10\x1a\x6d\x57\x41\x6c\x23\x17\x49\x65\xa3\xa6\x48\xb8\x71\x2c\x52\x93\x10\xf2\xb3\x94\xbb\xea\x13\xee\x88\xe4\x54\xc8\x91\xc7\xc8\x82\x30\x2f\x0f\x65\x36\xd9\xea\x26\x67\x28\x5d\x30\x7f\x98\xc6\xf2\x87\x95\xa1\xbc\x7c\xe5\xdf\xd1\x68\x48\x44\x9c\x9c\xb8\xdb\xa1\x9e\x8a\x33\x53\x75\x73\x4c\x46\xf4\x28\x1d\xd1\xa3\x46\x23\x82\x03\xb2\x4c\xa3\x0f\x9d\x0e\xdc\xa5\x9c\xd9\x63\x9e\x67\xba\x55\xea\xda\xad\xdf\x90\xe4\x8f\x98\x4a\x82\x47\x5e\x7a\xb7\x0b\xff\xfa\x35\x58\xf3\xea\x44\x59\xdd\x27\x93\xbf\xfd\x53\x0a\x84\x03\xa4\x34\x84\x44\x29\x34\x23\x9f\xb3\xb9\x83\x4b\xac\xcd\x9c\xd4\x3e\xdf\x34\x49\x35\xb1\xed\xa1\x0a\xaf\x3e\x69\x3f\xf5\xdc\xb1\x36\x06\xb6\x19\xb3\xf3\xa5\xfc\xf5\xbe\x48\xfc\x17\xf4\x62\x09\x62\xce\x87\x25\x0f\x11\x43\x01\x99\x0b\x86\x89\x1c\x79\x88\x11\x99\x7b\xb5\x7b\xf7\x61\x65\x14\x8e\xb3\x83\xcc\x1f\xdc\x87\x25\xfc\x5c\x74\x36\x7f\xba\xae\x3f\xb4\xd4\x93\x41\x65\xac\x21\x7d\x18\x37\xa9\xe0\x7c\xb8\x72\x19\xbf\x56\x82\x51\x5c\xe9\xda\xaa\x3a\x34\x33\x84\x3a\x60\x43\xdc\x15\xd3\xa9\x22\xba\xfb\x70\x1f\xaa\x1b\x47\xf4\x55\x3c\x09\xa9\xde\x3a\xa2\xaf\xe2\x20\xb0\xdb\xa3\x7c\x03\x4a\x38\x4e\x63\xd0\xae\xdf\x07\x1c\x87\xde\xe9\x8a\x36\xe9\x0f\xa1\xe6\x04\xe1\xac\xe6\x7a\x77\x07\x4a\x23\x4d\x83\x1f\xaf\x7e\x7e\x0e\x6d\x77\xff\xeb\x2f\xcf\xc1\xf3\x31\x52\xf3\x89\x40\x12\xfb\x48\x29\xa2\x95\xbf\x20\x1c\x0b\xa9\xfc\xec\x7b\x98\xea\x71\xa2\xbb\x13\xe5\x07\xca\xb5\x5e\xb9\xd6\x89\x10\x5a\x69\x89\xa2\x5e\x48\x79\x2f\x30\xbb\xd5\x29\x62\x8a\x74\x4e\x28\x35\xff\x0e\x97\x0e\x20\x6f\xb9\x9f\x01\xcc\xe9\x6c\xce\x4c\xe6\x70\xe3\xe4\x69\x11\x0a\x29\xc5\xf2\xa4\x93\x4c\xb6\x05\x89\x88\xf4\xb1\x4a\x84\xd3\xb7\xdd\x94\x16\x9c\xbe\xe9\x56\x30\xf3\x2f\xb3\x8f\x66\x76\x3b\x03\x36\x46\x9c\x6f\x35\x67\xee\xe4\x6e\xc3\xce\xdc\x56\x17\x2e\xfa\xfd\xaf\x1e\x67\x2f\xd6\x89\xa1\x59\xc9\x75\x36\x76\xa3\x4e\x68\x61\xfe\x8d\xf2\x6f\xfe\x88\x89\x5c\xf5\x0a\x46\x66\x20\xb9\xb9\x0f\xcb\x9a\x28\x23\xb0\xd6\x9c\xef\x45\x66\x6e\xbb\x25\xd9\x05\xa3\xfe\x08\xc2\x93\xb9\xd7\xae\xa4\xd3\x88\x2f\x2e\xa4\x9b\xc2\x63\x2f\x42\xc1\xbb\x13\xca\xc9\xbe\xfc\x1b\x29\xd9\xc3\x29\x71\xcc\x57\xeb\x4d\x69\xb1\x6e\x0b\x68\x76\x30\x3e\x5b\x66\x67\xed\xf4\x8c\x7c\xa7\x67\x4b\x8e\xed\x69\xcc\x6d\x72\x0d\xed\x4e\x69\x9d\x3e\x4b\x25\xbb\x2f\x19\xca\x6c\x41\x61\x54\x22\x32\x57\x84\x24\x62\x8c\xb0\x5f\x23\x26\x10\x56\x03\xb8\x38\xaf\xa4\x09\x5f\xa0\x90\x0c\xc0\x9b\x52\x19\x2e\x91\x24\xde\x36\x59\x20\x09\xd2\xe4\xa7\x10\xcd\xc8\xd5\x3c\x0e\x27\x1c\x51\xa6\x06\x6e\xce\xdb\xd4\x28\x08\x48\xa4\x09\xfe\x81\x32\xa2\x06\xb0\x31\xfd\x10\x05\x13\xca\x91\x5c\x9d\xf7\x26\x94\x6f\xfe\xda\x6b\xfd\x78\xb3\xce\x32\x67\x37\xaa\x47\x39\xd5\x3f\xa6\x46\x43\xf9\xec\x25\x7f\x2e\x10\x6e\x77\x4a\xb4\x19\x5e\x49\xc8\xfe\x11\x71\xcc\x88\x6c\xcb\x32\x7c\xe6\xa2\x53\x68\x4b\xb3\xda\x62\xa6\x61\x34\x1a\x41\x6b\x8a\x28\x23\xb8\x55\x45\x6c\x2e\x4e\x96\xf0\xea\x85\xd0\x74\xba\x6a\x57\x53\x98\xcb\xee\x68\x07\xd0\xfa\x97\x94\x42\xb6\xb6\x71\xc9\xe8\xc8\xad\x1e\x80\xec\x25\xa9\xd8\x0e\xc2\x55\x64\xf8\x91\x3d\xfc\xe6\x14\x93\x5a\x65\xa4\x97\xf1\xe0\x94\xcf\x06\xd0\xca\x56\xf9\xc3\x56\x25\xf5\xba\xf3\x78\xab\x7d\x9d\x7d\xd7\x6c\x67\xe3\x86\x2f\x0c\x72\x31\x77\x61\xe0\x54\xe0\xbd\x76\xfa\x3b\x21\x7c\x6a\x2f\xc7\x8f\x01\xe0\xa6\xa5\x3f\x28\x85\x5e\x8e\xc5\xb2\xe7\xb6\xd3\x30\xca\x8c\xb9\x1d\x84\xb8\x0a\xd5\xb3\x1e\xba\x41\xb7\x35\x60\x26\xb3\x7e\xf5\xf2\xf5\x55\xcd\x94\x63\xc9\x06\xd0\x4a\x77\xef\x2d\xf8\x06\x82\x10\x57\x93\x26\xd8\x0d\x4a\xcb\x6a\x7b\x7a\xa5\x29\x97\x97\xf2\x59\xbb\xf5\x65\x56\x26\x68\x75\x7a\x2e\xdf\xce\xbc\x5c\x9b\x2c\x8c\xef\xdb\x9e\x90\x6d\xef\x45\xd2\xfe\xff\x19\x99\xa2\x98\xe9\x76\x05\xbc\x0b\x24\x81\xc0\x08\xce\xda\x7a\x4e\x55\xd9\x39\x9c\x0a\x33\xd2\x8b\xa4\x88\xda\x2d\x87\x5c\xab\x53\x4d\x6a\x02\xae\xa1\x55\x44\x52\xc4\xe8\x7b\xd2\xae\x21\x3c\x1a\xdc\xf2\xfc\x2a\x4e\x4e\x6e\x4c\xbd\xdd\xca\x73\x40\x9b\x14\xbe\xa1\xf8\xed\x46\x62\x98\xa4\x84\xa6\xbd\xd5\xe9\x05\x73\xc4\x67\x24\xd7\x4e\x95\x0d\x2e\xd0\xf6\x40\xcd\x55\xd0\x42\x3d\x3c\x26\x5e\xad\xab\x5f\xbf\x83\x51\x86\x33\xc5\xad\xce\x16\x51\x85\x6e\x8d\x5b\x4a\xfb\x18\x55\xb6\x3a\xce\xab\xa7\x27\xab\x6a\x5d\x93\x19\xcb\x9b\x77\x6f\x0b\x22\x93\x53\x51\xad\x0e\x7c\x31\x82\x56\xab\xd6\x13\xee\xe7\xb7\x40\xac\xbd\x3d\xfa\xf5\xfd\x58\x66\x51\xc5\xad\xce\x01\x66\x6a\xfe\x9e\xde\x3a\x8b\xcf\x85\x3d\x43\xa9\xd9\xd8\xac\x8d\x6d\x55\x45\x96\xfd\x51\xa3\x49\xb8\x75\xb1\xc2\xfe\x2b\x18\xa9\xa4\x3a\xca\xfd\xb1\x76\x6f\x98\x68\x16\x22\xea\x01\x2a\x53\x64\xbf\x79\xcc\x37\x5e\xff\x0f\x00\x00\xff\xff\xf1\x5b\xc0\x63\xe4\x44\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x53\xcf\x6f\x13\xd7\x16\xfe\xb0\xb1\xcd\x33\xef\xf1\x1e\x3c\x4a\x7f\x4b\x17\x55\x40\x69\x35\xa9\x9d\x94\xfe\x70\x62\x02\x24\x20\x50\x09\xa4\xe0\xd2\x4d\xa5\xea\x12\x5f\x9c\x11\xf6\x8c\x75\x67\x06\x9a\x8a\x45\xe2\x54\xa5\x08\xaa\x42\xa5\x0a\x54\x51\xaa\xc0\xa2\x5d\x74\xe1\x24\x98\x18\x9a\x4c\xa4\x76\x53\x75\x75\xee\xae\xdd\xf4\x2f\x40\xed\xa6\x9b\x2e\x2a\x55\x77\xe6\xc6\x04\x48\xbc\xb9\xf7\x9e\x39\xdf\x77\xbe\xf3\x9d\xe3\xdf\x36\xad\xfd\x02\x00\x9e\x07\xf0\x1c\x80\xb5\x6b\x80\x9d\x00\xee\xaf\x41\xf4\xbb\x9c\x00\x52\x00\xae\x24\x80\x7f\x01\xb8\x9a\x00\x36\x02\x98\x4a\x00\xeb\x00\xcc\x26\x80\x34\x80\xf9\x04\xb0\x19\xc0\x8f\x09\x60\x0b\x80\x5f\x4c\xfe\x1f\x26\xff\xaf\x04\x90\x00\xb0\x3e\x19\xf3\x6d\x48\x02\x19\x00\x9b\x93\xc0\x5a\x5d\xdb\xc4\xb7\x26\x63\xbe\x1d\xc9\x98\xbf\xcb\x9c\xaf\x9b\xfc\x3d\x26\xff\xa0\xc9\x1b\x32\xf1\x92\x89\xbf\x67\x78\x78\x12\xd8\x06\xa0\x92\x8c\xfb\xf9\xd4\x7c\x9f\x33\xdf\xe7\x0d\xee\x87\x24\x30\xba\x06\xf8\x39\x09\xfc\x1b\x40\x4f\x0a\xf8\x9f\xae\x93\x02\x76\x01\x78\xdb\xbc\xcf\xa7\x80\xa7\x00\x5c\x4b\x01\x7d\x00\xa6\x53\xc0\x41\x00\xbf\xa7\x80\x27\xb5\xbe\x34\xb0\x03\xc0\xd1\x74\xdc\xe7\xd9\x34\xb0\x01\xc0\x87\x69\x60\x13\x80\x0b\xe9\x58\xef\xf5\x74\xdc\xcf\xad\x34\x90\x05\xf0\x5d\x3a\xe6\x9f\x33\xe7\x4f\x69\xe0\x19\x00\xbf\x1a\xdc\x9f\xe9\x58\x57\x26\x13\xf3\x6d\xce\xc4\x3c\x2f\x64\xe2\xfc\x97\x32\xc0\x30\x80\xde\x0c\xf0\x3e\x80\x4b\x26\x7e\x3f\x13\xe3\xfe\x36\xef\xec\x3a\x40\x8f\x54\x6b\xdb\x10\x8f\x16\x4f\x9b\x73\x3d\x1e\xfe\x3d\x61\x76\x21\x09\xe0\x59\xc4\xbe\x69\x4d\x9a\x47\xd7\xd6\xfe\x69\xed\xba\xef\x2d\x06\xf3\xff\x65\x78\xdd\xdf\x7f\x97\xbd\xf5\x1e\x68\xef\x32\xe6\xbd\xd1\xec\x9b\xd6\xa7\x77\xe6\x3f\x88\x7b\x8d\x7e\x7b\xa5\xe4\x63\xd8\xeb\xfb\xd2\x3e\x19\xf8\x02\xfb\xa4\xcb\xcb\x23\xdc\xf3\x59\x4d\x78\x1e\xaf\x08\x0c\x8c\x8a\x91\xd3\x5e\x50\xc3\x80\xeb\x9c\xb2\x2b\xf1\x21\x6b\x4c\x0a\x4f\xf8\xac\x2c\xce\xd8\x23\x62\x79\xd0\xe7\xb2\x13\x1e\x8c\x0e\xb6\xcd\x5b\xba\xf1\xa5\x42\x1e\x0e\x0d\xe2\xb0\x38\x23\xaa\x18\x32\x85\x8e\xf0\x9a\xc0\x11\xb7\x2c\x3c\x1c\xad\xfb\xb6\xeb\x60\x58\xba\x15\x29\x3c\x4f\x5f\xea\x42\xfa\x63\x38\x16\xf3\xe3\xb8\x70\xca\x38\xee\x73\x3f\xf0\x50\xb2\x6b\xc2\x0d\x7c\x94\xc6\xea\x02\x27\x78\x35\x10\x78\xd7\xae\x56\xd9\x49\xf1\x88\x9c\x2e\xb6\x57\x0a\x36\xe6\x06\xcc\x0b\xa4\xe8\xef\x64\x9d\xb5\xeb\x82\x8d\x44\xdd\xad\x9a\x59\x17\x90\x82\x97\xc7\x20\x03\xc7\xb1\x9d\x0a\x8e\x89\xba\x2b\x7d\x6b\xc8\xab\xd8\x65\x6b\x5f\x50\xf1\xac\x92\x5b\xd0\xf0\x3d\xa7\xed\x51\x5e\x73\xbb\x64\x90\x1d\x3e\x5a\xb2\x06\xa4\xe0\xba\x19\x6b\x90\xfb\xa2\xc0\xba\x73\xf9\x37\xad\x5c\x8f\x95\x7f\x8d\x75\xf7\x14\x76\xed\x7a\x39\xd7\x93\xcb\x65\x0f\x73\xcf\xb7\x4a\x92\x3b\x5e\x95\xfb\xae\x2c\xb0\xb7\x22\x0e\x36\x14\x48\x5e\x73\xcb\x2e\xeb\x7b\x88\x78\x77\xf6\x30\x77\x2a\x01\xaf\x08\xab\x24\x78\xad\xc0\x3a\xef\x02\x3b\x16\x78\x9e\xcd\x9d\xec\xd0\xa1\xa1\xfd\xd6\x09\x21\x3d\xdb\x75\x0a\x2c\xdf\x95\xcb\x0e\xb8\x8e\x2f\x1c\xdf\xd2\x3e\x15\x98\x2f\x3e\xf0\x5f\xa9\x57\xb9\xed\xf4\xb2\x91\x51\x2e\x3d\xe1\x17\xdf\x29\x1d\xb0\xde\x78\x90\xa7\xf5\x9c\x12\xd2\xda\xef\x8c\xb8\x65\xdb\xa9\x14\x58\x76\xb8\x1a\x48\x5e\xb5\x0e\xb8\xb2\xe6\x15\x98\x53\x8f\x9e\x5e\xb1\xa7\x97\xc5\xd7\xa2\xb3\x2d\x9f\x2b\x16\xf3\x6c\xfb\x76\xa6\xaf\xb9\xad\xc5\x7c\x9e\xf5\xb3\x1c\x2b\x44\xef\xdd\xc5\xee\xa5\x4f\x7d\xc5\x57\xf5\xf5\xc5\x28\xad\x2f\x9f\x63\xe7\xce\xc5\x90\xdd\xc5\xee\xdc\x4e\xd6\xcf\xf2\xac\xc0\xba\x7b\x41\xd7\xa9\xa9\x26\xd4\x04\xb5\x69\x06\x74\x59\x35\x54\x43\x8d\x53\x9b\xa6\xd5\xa4\x6a\x80\xbe\xa5\xb6\x1a\xa7\x90\xee\x51\x48\x33\xd4\x52\x17\xa8\xa9\x1a\xd4\xa2\xef\xd5\x25\x5a\xa0\x90\x5a\x4c\x4d\x50\x48\x21\x4d\xab\x0b\xd4\xa2\x05\x6a\x53\x0b\xf4\x0d\xb5\xe8\x9e\x9a\x50\x93\x34\x4f\xf3\xd4\x04\x7d\x49\x21\x2d\xa8\x8f\xa8\x4d\xb3\x6a\x52\x8d\x53\x53\x7d\x4c\x6d\xf5\x19\xe8\x6b\x0a\xe9\xb6\x6a\x44\xdc\xe3\x74\x87\x6e\x2f\x91\x68\xde\x69\x5d\x5a\x4d\x50\x93\xa9\x49\x35\x11\x09\x0b\xe9\xae\xbe\xd1\x8c\x66\x5d\x1d\x4c\x8b\x51\xa8\x45\x73\xd4\xa4\x59\x35\xae\x26\x69\x8e\xee\x51\x7b\x15\xa2\x5b\x8f\x05\x43\xfd\xaf\x7a\xd4\x0e\x75\x71\x65\xfc\xa1\x41\x4d\x31\x1e\x5b\x44\x0b\xea\x12\x68\xea\x71\x53\xae\xd1\xbc\xee\xf8\xab\x48\xf4\x45\xd0\x0d\x5a\xec\xb8\x10\x81\x67\xb5\x62\x3d\x0b\x0d\x9f\x59\x2e\x46\xb7\xfa\x78\x3f\x6d\xd5\xd0\xa5\x6e\xa8\x06\x2d\x6a\x4f\x69\x66\x29\x34\xa5\x1a\x7a\x4e\x5a\x2c\xe8\x26\x35\xe9\xae\x1e\x43\x3c\xd1\x9b\xd4\xa6\x45\xd0\x55\x5a\xa0\xa6\x3a\xff\x40\xde\x4a\x1e\xe8\xb6\xb5\xab\xaa\xb1\x92\xa3\x77\x22\x70\xd8\xc5\xe8\xf3\xc8\x19\x33\x07\xed\xc0\xc5\x7e\xd0\x95\x0e\xb6\x33\xca\x4f\x22\x44\x93\x45\xeb\xb4\xc2\x3e\xac\x68\xef\x6a\xfc\x53\x9d\xfd\x30\x5d\xcf\x52\xa8\x1a\xd1\x10\x42\x44\x32\x17\xd5\xa4\x19\x40\x88\x7f\x02\x00\x00\xff\xff\xc9\xa5\xfa\x87\x1b\x08\x00\x00"

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
