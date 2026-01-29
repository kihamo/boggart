// Code generated for package webos by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
package webos

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x54\x4f\x6f\xdb\x3e\x0c\xbd\xf7\x53\x10\x3c\xfd\x7e\x03\x1c\x23\xc9\x50\xec\x60\xe7\xbc\x43\x7b\xd9\xba\x73\x20\x5b\x74\xa3\x42\x96\x32\x91\x49\xb3\x1a\xf9\xee\x83\xfc\x27\x75\x13\x2f\xcb\x61\xb9\x84\xb2\x9e\xc4\x47\xbe\x27\x36\x0d\x68\xaa\x8c\x23\xc0\xd2\x3b\x21\x27\x08\xc7\xe3\x5d\xa6\xcd\x1e\x4a\xab\x98\x73\x0c\xfe\x15\x57\x77\x00\x00\xe3\xaf\xa5\xb7\x49\xad\x93\xf9\x02\x62\xc4\xf5\x10\x1d\x38\x99\x2f\x7a\xfc\xf9\x99\xc3\x7a\xab\x1c\xd9\xd1\xee\x25\x42\x8c\x58\x3a\x43\xb4\xa8\xcd\x62\xd5\x34\x60\xe6\x5f\x1c\xe0\x93\x57\x2c\x08\x33\x38\x1e\xb3\x74\xb3\x98\x40\x8f\x99\x5a\x52\xa1\x32\x07\x5c\x65\xa9\x36\xfb\xb3\xe4\x13\x9f\x3e\xf0\x19\xba\x32\x91\xa3\xf2\xa1\x1e\x80\x31\x4e\x36\x3e\x98\x37\xef\x44\x59\x68\xd7\x56\x15\x64\x13\x4b\x95\x20\x04\x6f\xa9\x83\x21\xd4\x24\x1b\xaf\x73\xdc\xfa\x58\x85\xd1\x39\x32\x39\x8d\xe0\xfc\x5e\x59\xa3\x95\xd0\x65\xb6\x73\x66\x46\xa8\xee\xb2\x3c\x07\xbf\xdb\x4e\xf0\x3b\x9d\x6a\x69\x44\x6c\x8e\x35\x31\xab\x67\xc2\x77\x15\x9d\x04\x6f\x3b\xa6\xd0\x6b\xba\x1c\x24\x5d\x4e\x2a\x3a\xf5\x3b\x29\xf3\x38\x24\x88\xda\x40\xc6\x5b\xe5\x4e\x3e\xa2\x9f\x3b\x13\x48\xe3\xea\x53\x96\xc6\x8d\x2b\x94\xd3\x96\xd0\x15\xc0\xa5\x13\xef\x07\xd6\xf7\x37\xb3\xce\x84\x0e\xa2\x02\xa9\x0f\x2a\xf6\x3d\x81\x40\x6c\xde\x54\x61\x69\x3d\xc0\xa2\x8a\xaf\x9c\xe3\x67\x04\xa7\x6a\xca\x51\x3a\x1f\x46\x05\xfb\x70\x28\x72\x5c\x6e\xd3\xc0\xac\xdd\x6e\xed\x3a\x5c\x76\xad\xfa\x0b\x4f\xde\xb2\x35\x6a\x89\x75\x6b\xf6\xd6\xe8\x49\xcb\x4f\x1d\xb8\xcd\x48\x7f\x6a\x7a\xad\x13\x5f\x55\x4c\x92\x2c\xff\xd6\xf1\x62\x27\xe2\xdd\xc8\xf3\xf2\x6b\x4b\x39\xf2\xae\xa8\x8d\x9c\x6c\x59\x88\x83\x42\x5c\xc2\xbb\xb2\x24\x66\x7c\x7f\xf9\xdf\xdb\x43\xdd\xc3\xef\xee\xfa\x97\x7d\xcc\xd2\xd8\x88\xab\x13\x62\xb4\xec\xc3\xfe\xaf\x69\x80\x9c\x8e\xc3\xf3\x6e\x34\x54\x5f\xb8\x9d\xa7\xd0\xbd\x11\x16\x25\xa6\xfc\xfa\xf4\xf8\x00\xff\x75\xf1\x8f\x6f\x0f\x80\xa9\x56\xbc\x29\xbc\x0a\x3a\x55\xcc\x24\x9c\xee\xc9\x69\x1f\x38\xed\xe7\x81\x0f\xe9\xcb\x68\x31\xab\x8d\x9b\xc5\x9b\x2b\x65\x99\xfe\x8f\x09\x4e\xd9\x7f\x07\x00\x00\xff\xff\xe4\xce\x09\xe9\xd1\x05\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x4e\xcf\x6b\x13\x41\x14\xfe\x1a\xab\x42\xbc\x68\xf1\x28\xf8\x44\x2c\x7a\xd8\xba\x1b\x45\x64\x93\x4d\x44\x6d\x41\x4c\x40\x62\xf4\x3e\x64\x27\xdb\xc5\x64\x36\xec\xec\x16\x0b\x39\xd4\x58\x41\xc1\xab\x57\x41\xbc\xe8\xad\x16\x02\x2d\x31\x9b\x83\xff\xc0\x5b\xf0\xec\xdd\x8b\x7f\x83\xec\x26\x4d\x15\xfa\x5d\xde\xf7\x6b\xde\xbc\x5f\x2b\xcb\x1f\x00\xe0\x0c\x80\x4b\x00\xea\x00\xce\x01\x18\x60\x86\x2f\x00\x2e\x00\xf8\x0a\xe0\x2c\x80\x31\x80\xf3\x00\x7e\x00\xb8\x08\xe0\x27\x80\x65\x00\xbf\x01\xd4\x96\x80\x3f\x00\xae\x02\x58\x2b\x00\x2b\x00\x9a\x05\xe0\x32\x00\x59\x00\xee\x00\x78\x33\xf7\xbf\x17\x80\x25\x1c\xa3\x30\xdf\x73\x84\xd3\xf3\x99\xdd\x75\x2a\x23\x0f\xe5\x96\xdf\x96\xe4\x6b\x0a\x3a\x9d\xae\xaf\x24\x1a\x52\x6b\xe1\x2d\x66\x16\xc9\x5e\x3f\xda\x5e\x18\x5a\x2a\x97\x74\xdc\x6e\x4b\xad\xf1\x54\x2a\x17\x4d\xd9\x0f\xc2\xc8\x68\x68\xcf\x77\x8d\xfb\xb1\xa7\x8d\x56\x60\x93\x2b\xb7\xee\xbd\xf0\x37\x45\x2f\x58\x0b\xe3\x62\x5d\xe8\xc8\x68\x85\x42\xe9\xae\x88\x82\xd0\xa6\xc7\x79\x44\x8d\x38\x14\xbd\xc0\x0d\xa8\xf2\x5f\xbf\x5a\xac\x0b\xe5\xc5\xc2\x93\x46\x4b\x8a\x9e\x4d\x0b\x6d\x53\x33\xd6\xda\x17\xaa\xd8\x78\xd4\x58\x37\x9e\xcb\x50\xfb\x81\xb2\xc9\x5a\x33\x8b\x0f\x02\x15\x49\x15\x19\xad\xed\xbe\xb4\x29\x92\x2f\xa3\x9b\xfd\xae\xf0\x55\x99\xda\x9b\x22\xd4\x32\x72\x9e\xb5\x36\x8c\xbb\xc7\xbd\xec\x9e\x8e\x0c\x8d\x75\xd5\x0e\x5c\x5f\x79\x36\x15\x9f\x74\xe3\x50\x74\x8d\x8d\x20\xec\x69\x9b\x54\x3f\x97\xda\xb9\x55\xa6\x19\x75\xd4\x35\xcb\x74\x1c\x8b\x56\x57\x29\xa3\xe6\x15\xc7\xb2\xa8\x46\x26\xd9\xb9\xae\x3a\xa5\xa3\xa8\xe2\xdc\xce\xe8\xf5\xbc\x56\xb1\x4c\x1a\x0c\x66\x4f\xaa\x4e\xc9\xbc\x41\x35\xb2\xc8\xa6\x52\x19\xfc\x39\x7d\x95\x0e\xd3\x1d\x4e\xf8\x30\x63\xbc\xcf\x09\x71\x92\xee\xa6\xbb\x3c\xe6\x3d\x3e\xe4\x09\xf8\x13\x27\x9c\xf0\xb7\xf4\x1d\x8f\x78\xc2\x07\x3c\x3a\xc1\x22\x9e\xa6\xaf\xf3\x0d\xc9\xc9\x79\x96\xf2\x94\x47\xe9\x5b\x9e\xcc\xfe\x18\xf2\x34\xdd\xe1\x3d\xde\xe7\x71\xde\x4a\xc0\x1f\xff\x31\x0f\xd2\x61\xfa\x1e\x7f\x03\x00\x00\xff\xff\xa9\xd1\x1a\xe4\xcf\x02\x00\x00"

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
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
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
