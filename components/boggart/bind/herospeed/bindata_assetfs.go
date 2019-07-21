// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
// DO NOT EDIT!

package herospeed

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x57\x41\x6f\xeb\x36\x0c\xbe\xe7\x57\x10\xda\xdb\x92\x6c\x8d\x9d\x3c\xec\x30\xb4\x49\x76\xd8\x3b\x6c\x40\x07\x3c\x14\x7d\xbb\x16\xb2\x45\x27\xca\x14\xc9\x95\x64\xe7\x05\x46\xfe\xfb\x20\xcb\x76\x9c\x26\x75\x9d\xae\x3d\x34\x34\x45\x7d\x24\x3f\x51\xa4\x5d\x14\xc0\x30\xe1\x12\x81\xc4\x4a\x5a\x94\x96\xc0\xe1\x30\x98\x33\x9e\x43\x2c\xa8\x31\x0b\xa2\xd5\x8e\x2c\x07\x00\x00\x6d\x6d\xac\xc4\x64\xcb\x26\xb3\xcf\xe0\x24\xb3\xad\xa5\xef\x66\x32\xfb\x5c\xd9\xbf\xdc\xf3\xfd\x29\xa5\x12\x45\x6b\xf5\xdc\xa2\x8e\xe2\xd4\xa6\xb1\xd3\x4a\xe0\x82\x58\x1a\x5d\x42\x6a\x2c\x33\x51\x03\x4a\x9a\x83\xa4\xf9\xc4\xd2\xc8\x40\x44\xf5\x93\x13\xc8\x11\x46\x70\x73\xc9\x57\x83\x24\x78\x65\x9b\x6a\x34\x28\x2d\xb5\x5c\x49\x52\x14\xc0\x13\xc0\x67\x08\x68\xec\x14\x40\x1c\x6b\xb5\x4f\xa7\xcb\xd1\x19\xa1\x64\x70\x38\x2c\xe7\x14\xd6\x1a\x93\xc5\x2b\xfb\x36\x34\xa7\x26\xd6\x3c\xb5\xb7\xb9\xe2\x6c\x34\x1d\xdf\xb9\xbd\xc2\x20\x1c\x0e\x45\x01\xc1\x03\x3e\x67\x68\x6c\xf0\xed\xe1\x3e\xf8\x4a\xed\xda\xab\x3d\x38\x59\x3a\xd0\xd9\x6f\x12\xc8\x57\x8d\x39\xc7\x1d\x81\x00\x0e\x87\x79\x48\x97\xf3\x50\xf0\x0f\xc8\x2d\x56\x32\xe1\xab\x4c\xfb\xf5\x77\x26\x7a\x06\xf2\x9e\xac\x7f\xf7\x68\x8b\x13\xb0\x4b\x54\xfc\x71\xea\xed\x6d\x42\xe6\x61\x26\x5e\x59\x69\x95\xa7\xa5\xd1\xe4\xf5\x02\xed\x2a\xd4\x36\x82\xd3\x40\x42\x19\x82\x27\x10\xb8\xec\x40\x73\x7f\x3d\xe8\xec\xdc\x5f\x46\xe5\xaa\x1d\x5b\x71\x08\x84\xf2\xff\xc4\x58\xcd\x53\x64\xc0\xa8\xa5\xa5\x86\x80\xb1\x7b\x17\xfe\x8e\x33\xbb\xbe\x9d\x4d\xa7\x3f\xbe\x11\xe0\xd1\xc9\x1a\x29\xeb\x6b\xab\xfb\x19\x56\xc0\x2f\x3a\xcf\xaf\xed\xca\xa7\x9a\x6e\xd1\xa2\xae\x8f\xda\xae\xaf\x82\x3e\x22\xfd\x43\x45\x86\xd7\xa2\xcc\xc3\xbe\xa9\x38\xcc\x2b\x08\x8a\x14\xdb\xf7\x4f\xa4\x28\x40\x53\xb9\x42\xf8\xa4\x52\x57\x16\x37\xf0\x29\x77\xe9\xc0\xed\x02\x82\x93\x82\xe9\x53\x2f\xc7\x28\xae\x38\x26\xbf\x81\x39\x3e\xab\x20\x3c\x91\x3d\x33\x7e\x09\xe1\xc3\xbf\x16\xa1\xff\x71\x80\x67\xcd\xb7\x8f\xbe\x07\xd8\xef\x50\xe6\x61\x79\x95\xde\xbc\xd7\x55\xc3\x7b\x1b\xaf\xd5\x86\x12\xa5\xb7\x93\x95\x56\x59\x0a\x69\x26\xc4\x44\xf3\xd5\xba\xab\x21\xbd\x86\xc3\x65\x9a\x59\x0f\xd4\x73\x77\x89\x20\x68\x84\xcd\x80\x8d\xac\xbc\x1a\xa1\x44\x89\x32\x6b\x95\x04\xbb\x4f\x71\x41\xfc\x03\x69\x81\x82\x03\xe6\xb1\xaa\x04\x99\xa8\x52\x30\x5b\x02\x4a\xc6\x82\xc7\xff\x96\x73\xcb\x0d\xbc\x07\x4c\x34\x9a\xf5\x68\x7c\x65\x08\x65\x18\xbc\x21\x95\x1a\x48\xe8\xc4\xec\x65\x4c\xc0\x72\xeb\xfa\x5f\xd3\x16\x2a\x17\xbe\x31\x90\xe5\x3c\xec\x98\xab\x17\xdd\x84\x3e\xc3\x2b\x77\xb5\xc6\x69\xd7\x38\xac\x78\xf8\x89\xa9\x9d\x14\x8a\xb2\xc5\xac\x1f\x93\x1f\x40\x57\xed\xf2\x9c\xb2\x2f\xcd\xca\xbb\x39\xa3\xd7\xdc\xf9\xb2\x2a\xfb\x76\x61\xc6\xf3\x3e\x57\xb8\x9f\x19\xdf\xae\xc0\xe8\xb8\xef\x31\x11\x28\xc7\xea\x82\x94\x73\x15\xd6\xe8\x6e\x6f\xfd\xc4\x59\x53\xd6\x04\xc2\xb7\xbb\x47\x77\xe7\xea\xc8\xe0\x95\xa5\x0b\xea\x17\xaa\xd6\x63\x25\x56\x3f\xc7\x70\x06\xad\xaf\x0a\x37\xee\x9a\xf7\x93\xa2\x00\xe3\x5e\x33\xe3\x3f\x1f\xff\xbe\x87\x91\x97\xbf\x3d\xdc\x03\x09\x19\x35\xeb\x48\x51\xcd\x42\x6a\x0c\x5a\x13\xe6\x28\x99\xd2\x26\x6c\xde\x4a\x4c\x20\xd1\x4e\x22\x13\xc6\xc6\x6b\x1f\xbd\x36\x52\xca\x1a\xab\x69\x1a\x6c\xb9\x0c\x62\x63\x08\x24\x54\x18\x1c\x7f\xa0\x57\x8d\x26\x55\xd2\xf0\x1c\xeb\x00\x8e\x9a\xee\x00\x2e\xb3\xb2\x31\x6d\x4e\x7a\xbf\xdb\xfd\xff\x4c\xc2\x8d\x09\x37\xcf\x19\xea\x7d\xd0\xa2\xd0\xc5\xbd\x39\xe3\xed\xe3\x4e\x6c\xd3\x71\x60\xe7\x7e\x5f\x4c\xc5\xb9\xff\x4a\xa8\x26\x05\x4d\x53\xc1\xe3\x92\x9a\xf0\xf8\x09\xd1\x6a\x65\x3b\x2e\x99\xda\x05\xa7\xb3\x01\x16\x90\x64\xb2\xe4\x77\x34\x86\xe2\xa4\xbe\x3f\x8d\x86\x3f\x54\xd6\xc3\x71\x40\xad\xd5\xa3\xa1\xd1\xf1\xf0\x06\x86\x3d\xfb\xee\x93\x5d\x0c\x7f\x19\x49\xdc\xc1\x17\x6a\x71\x34\x1e\x07\x2b\xb4\x8f\x7c\xeb\xc4\xbb\xc6\xd7\xe1\x6e\x30\xb8\xec\x53\x23\x65\xfb\x51\x1d\x20\x9c\x45\x58\xa5\x64\xd0\xfe\x25\x2d\xea\x9c\x8a\xd1\xc5\x2c\x6f\x1c\x75\xb5\xee\x49\x7b\xe5\x13\xaf\xf6\xb8\x0f\xb7\x9f\x61\x36\x9d\x4e\xdb\x41\x55\xf2\x3c\xf4\x44\x2e\x07\xa7\x9d\xa5\x91\xfe\x0b\x00\x00\xff\xff\xe4\x49\x25\x15\x29\x10\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x4e\x3f\x6f\xd3\x40\x1c\x7d\x0d\xe5\x5f\x10\x08\x31\x31\x30\x1c\x03\x15\x0c\x2e\x76\x60\x40\x4e\x9c\x20\x4a\x2b\x55\x34\x52\x14\x99\x8a\xf5\xd4\x5c\x1c\x0b\xe7\x2e\xba\xb3\x5b\x40\x1d\x42\x41\x30\x50\x89\x09\x21\x31\x50\xf1\x0d\x02\x28\x52\x40\x4d\x56\xd6\xdf\xad\x0c\x0c\x8c\x4c\x7c\x04\x94\xc4\x50\xf1\x96\xf7\x9e\xdf\x7b\x3f\xdf\x8f\x0b\x8b\x6f\x00\xe0\x24\x80\x4b\x00\x42\x00\x67\x00\xec\x63\x8e\x11\x80\xb3\x00\xbe\x00\x38\x05\xe0\x1b\x80\xd3\x00\x28\xdf\x7c\xcf\xf9\x27\x80\xe3\x00\x7e\x01\xa8\x2d\x00\xbf\x01\x5c\x04\xe0\x16\x80\x73\x00\xd6\x0b\xc0\x79\x00\x0f\x72\x8e\x73\x7e\x92\xf3\xab\x02\xb0\x80\x23\x1c\x03\x50\xc8\x6f\x4f\x71\x22\xe7\xc5\xfc\x3f\x58\x51\xb2\x1d\x47\x99\xe6\x69\xac\x24\xee\xaa\x1d\x99\x28\xde\x42\x83\x6b\xde\x15\xa9\xd0\x68\x68\xb1\x1d\x8b\x1d\x34\x45\x5b\x0b\xd3\xc1\x26\x4f\x32\x81\xa6\xe8\x29\x9d\x3a\x75\x13\xc5\x2d\xe7\x4e\x16\x19\x27\x54\x3e\x6b\x89\xed\xdb\x0f\xe3\x0e\xef\xaa\x65\x9d\x15\x37\xb8\x49\x9d\x50\x73\x69\x12\x9e\x2a\xed\xb3\x7b\xb3\x88\xd5\x33\xcd\xbb\xaa\xa5\x58\xe5\xbf\x7e\xb5\xb8\xc1\x65\x94\xf1\x48\x38\xa1\xe0\x5d\x9f\xfd\xf3\x3e\x6b\x66\xc6\xc4\x5c\x16\xeb\xeb\xf5\x55\x67\x53\x68\x13\x2b\xe9\x33\x6f\xd9\x2d\xae\x28\x99\x0a\x99\x3a\xe1\xe3\x9e\xf0\x59\x2a\x1e\xa5\xd7\x7b\x09\x8f\x65\x99\x6d\x75\xb8\x36\x22\x0d\xee\x87\x6b\xce\xad\xa3\xde\xf4\x3d\x6d\xa1\x9d\x55\xb9\xa5\x5a\xb1\x8c\x7c\x56\x6c\x24\x99\xe6\x89\xb3\xa6\x74\xd7\xf8\x4c\xf6\x66\xd6\x04\x37\xca\x6c\x2e\x03\x79\xc5\x73\x83\xc0\x63\x4b\x4b\x6c\x2a\xdd\xcb\x81\xe7\xb1\x1a\x73\x99\x3f\xf3\xd5\xa0\xf4\x37\xaa\x04\x37\xa7\xf2\xea\xac\x56\xf1\x5c\xb6\xbb\x3b\x9f\x54\x83\x92\x7b\x8d\xd5\x98\xc7\x7c\x56\x2a\x83\xde\xd1\x84\xc6\xf6\x39\x8d\xe8\xb3\x7d\x66\xfb\x34\xb0\x2f\x68\x64\x5f\x83\x3e\xd0\x57\x1a\xd8\x97\x34\xb0\x7b\x76\x1f\x74\x40\x83\x69\x4a\x87\x34\xb4\x7b\xb6\x0f\x3a\xb0\x7d\x9a\xd8\xa7\x74\x48\x93\xf9\x87\xf7\xf4\x91\xc6\x34\xa1\x4f\x34\x9a\x4f\xde\xd2\x78\x76\x61\x48\x63\x1a\xd1\x10\x7f\x02\x00\x00\xff\xff\x72\x24\x38\x60\x9c\x02\x00\x00"

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
