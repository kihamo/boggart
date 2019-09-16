// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
// DO NOT EDIT!

package esphome

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x93\xc1\x6e\xc2\x30\x0c\x86\xef\x3c\x85\x15\xb1\xdb\x4a\x29\xa7\x09\x05\x5e\x00\x69\x97\x3d\xc0\x14\xb0\x37\x22\x95\xa4\x4a\x2c\xa0\x8a\xfa\xee\x53\x9b\xb2\x75\xac\x19\x1b\x97\xd6\x89\xff\xdf\xae\xbf\x26\x21\x00\xd2\x9b\x36\x04\x62\x67\x0d\x93\x61\x01\x4d\x33\x91\xa8\x8f\xb0\x2b\x95\xf7\x2b\xe1\xec\x49\xac\x27\x00\x00\xc3\xdd\x9d\x2d\xb3\x03\x66\xc5\x02\xda\xc8\x1f\x2e\xd1\xd9\x67\xc5\xa2\xd7\x5f\x7b\xce\xaf\x95\x32\x54\x0e\xb2\x3f\x15\x97\xaf\xf8\xae\xe9\x74\xac\xb6\x25\x5d\x94\x71\xd1\x3d\x33\xcf\x4e\x57\x84\x80\x8a\x55\xdc\x47\xce\x1c\xf9\xca\x1a\xaf\x8f\x04\xc6\x9e\x9c\xaa\x04\x78\xae\x4b\x5a\x89\x93\x46\xde\x2f\x8b\xf9\xfc\x61\xa4\x4b\xec\xb4\x27\x85\xe3\xb9\x98\x77\xe9\x64\x5f\x60\x1d\x02\xe8\xe2\xc9\x80\xd8\x50\x2d\x60\x06\x4d\x23\x73\xde\xff\xc3\xf7\xac\x0e\x74\x97\xf1\x85\x15\xff\xc9\x29\xf3\xd4\x20\xad\x2f\x89\x40\xf2\xd6\x62\x9d\x2e\x1b\x02\x38\x65\xde\x09\xa6\xfa\x11\xa6\x64\x58\x73\x0d\xcb\x15\xcc\xba\x50\x93\x6f\x4f\x58\x7a\x92\x9b\x6c\xb1\x1d\xb5\xaf\x3b\xdb\x50\x1d\xe7\xfc\xe5\x7f\x8d\xd8\x5a\xb8\xf7\xf8\x3a\xb6\x37\x8d\x69\xb0\x10\xf9\x90\xc1\x14\x04\x99\x27\xf0\xca\xbc\x3b\xdb\x57\x97\x27\x47\x7d\x1c\xdc\xb6\xaf\x65\x1f\xf6\xaf\xcf\x9e\x1f\x01\x00\x00\xff\xff\xe5\x94\xaa\x1c\xf1\x03\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x8e\x3f\x4f\x14\x5d\x14\xc6\x7f\x0c\xbc\xaf\xba\x36\x4a\x67\x62\x71\x34\x91\x48\x31\xb8\x83\x16\x66\x96\x01\xa3\x40\x42\x60\xd1\xe0\x6a\x7f\x65\x2f\xcb\xe8\xee\xbd\x9b\xf9\x63\x24\xa1\x91\xc6\x44\x8d\x54\x26\xc6\xce\x6f\xa0\x24\x24\x88\xb2\x85\xb5\xc5\x1d\x0b\x13\x1b\x3f\x8b\xd9\xdd\x41\xa3\xa7\x39\xcf\xef\x3e\xcf\xb9\xe7\xfc\x1c\x1f\x7b\x0d\x70\x02\x38\x0f\x34\x80\xd3\xc0\x4b\x86\x75\x00\x8c\x02\x1f\x81\x31\xe0\x08\x38\x05\x7c\x01\xfe\x03\xbe\x96\xfe\x37\xc0\x03\xbe\x03\x73\x23\xf0\x03\x38\x09\x9c\xf3\xe0\x7f\xe0\xa2\x07\x67\x81\x49\x0f\xc6\x81\x5a\xd9\x57\x3d\x38\x03\x28\x0f\x46\xca\x3b\x8e\xcb\x2b\x77\x8c\x96\x3c\x56\xf6\xfe\x7f\x2c\xeb\x2d\x56\x55\x47\x73\xfb\xc1\x43\xbd\x9e\xc9\xd2\x3c\x77\x33\x95\x69\xec\xc6\x06\xd6\xb0\xa6\xbb\x36\xc9\xfc\x7a\xda\x8a\x9b\xfe\xcd\xbc\x95\xfa\x0d\x1b\x4a\x53\x3f\xbe\xf1\x28\xde\x54\x1d\x3b\x95\xe4\x95\x15\x95\x66\x7e\x23\x51\x26\x6d\xab\xcc\x26\xa1\x2c\x0f\x2c\xa9\xe7\x89\xea\xd8\xa6\x95\x99\xbf\xf2\xb3\x95\x15\x65\x5a\xb9\x6a\x69\xbf\xa1\x55\x27\x94\xdf\x1c\xca\x5a\x9e\xa6\xb1\x32\x95\xfa\x52\x7d\xc1\xbf\xaf\x93\x34\xb6\x26\x94\x60\xaa\x5a\xb9\x65\x4d\xa6\x4d\xe6\x37\xb6\xba\x3a\x94\x4c\x3f\xc9\xae\x74\xdb\x2a\x36\x35\x59\xdf\x54\x49\xaa\xb3\xe8\x5e\x63\xd1\xbf\xfe\x27\xd7\xbf\x67\x43\x27\xfe\x82\x59\xb7\xcd\xd8\xb4\x42\xa9\xdc\x69\xe7\x89\x6a\xfb\x8b\x36\xe9\xa4\xa1\x98\xee\x00\xd3\xe8\x6a\x4d\x86\x32\x32\x97\x82\x6a\x14\x05\x32\x31\x21\x7d\x59\xbd\x10\x05\x81\xcc\x49\x55\xc2\x01\xcf\x46\xd3\xc7\xd6\x4c\x74\xad\x2f\x2f\x0f\x62\x33\x41\x55\xb6\xb7\x87\x23\xb3\xd1\x74\x75\x52\xe6\x24\x90\x50\xa6\x6b\xb8\xb7\xee\x53\xf1\xaa\x78\x86\x7b\xe3\x3e\x17\xbb\x2c\xcd\x8b\xeb\xb9\x0f\xc5\x73\xb7\xef\x0e\x8b\x1d\xf7\x1e\xf7\xce\xf5\x8a\xa7\xc5\x8e\xeb\x15\xbb\xee\xc8\x1d\xb8\x7d\xdc\x5e\xf1\xc2\x1d\x0e\x07\xdd\xbe\x3b\x72\x3d\xdc\xde\x3f\x0f\xbf\x02\x00\x00\xff\xff\xec\x45\x18\x2a\x72\x02\x00\x00"

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
