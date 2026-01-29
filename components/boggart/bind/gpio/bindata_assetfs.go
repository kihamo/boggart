// Code generated for package gpio by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
// locales/ru/LC_MESSAGES/widget.mo
package gpio

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x54\x4d\x8f\xdb\x20\x10\xbd\xef\xaf\x18\xd1\x4a\x1b\x1f\x6c\x2b\x59\x69\x55\xa9\x49\xae\x6d\xa5\x4a\xed\x3f\xa8\x58\x18\xaf\xd9\x62\xa0\x80\x93\x6c\xa3\xfc\xf7\x0a\xfc\x91\xc4\xeb\xd8\x5b\x5f\x3c\x66\x9e\x67\x1e\xf3\x1e\x1c\x8f\xc0\xb1\x10\x0a\x81\x30\xad\x3c\x2a\x4f\xe0\x74\xba\x3b\x1e\x41\x14\xa0\x10\xb2\x4a\x73\x84\x65\x58\x5b\x3b\x66\x85\xf1\xe0\x5f\x0d\x6e\x08\x35\x46\x0a\x46\xbd\xd0\x2a\x7f\xa1\x3b\xda\x24\xc9\xf6\x0e\x00\xe0\xe3\x82\x6b\x56\x57\xa8\x7c\x92\x59\xa4\xfc\x75\x51\xd4\x8a\x05\x2c\x2c\x12\x38\x46\x4c\x83\xbb\xff\x20\x71\x87\xf2\x3e\xc9\x58\x49\xd5\x33\xf6\xc0\x2b\x5c\x87\x75\xe8\xef\x93\xcc\xd5\x4f\x95\xf0\x8b\xe4\x73\x9f\x3f\xb5\x71\x78\xaf\xf3\x86\xca\x36\x6c\x02\x15\x8f\xd4\xb9\xd8\x01\x93\xd4\xb9\x0d\xb1\x7a\xdf\xb2\xbc\x5c\x65\x5a\xa6\x15\x4f\x97\x2b\x08\x91\xab\xba\xe8\xe0\xd2\xe5\xaa\xc5\x0f\xff\x39\xfc\x32\x54\xa1\xbc\xc8\xbe\x45\x78\xe1\x25\x0e\x10\x11\x55\xae\xb6\x61\xc8\xcb\x4f\x0a\xc8\x97\x9f\xdf\x7e\x10\xc8\xe0\x74\x5a\xe7\xe5\x6a\x04\x7c\x49\x54\x22\xb5\x85\x38\x90\xed\x3a\xe7\x62\x37\xe8\x3d\xb2\x74\x45\xa7\xd3\x78\xa4\x47\xa1\x6d\xd5\x01\x43\x9c\x96\xda\x8a\xbf\x5a\x79\x2a\x21\x7e\x4b\xfa\x84\x32\x95\x58\x78\x02\x56\x4b\x6c\x60\x04\x2a\xf4\xa5\xe6\x1b\x62\xb4\xf3\x04\x04\xdf\x10\x87\x63\x1d\x86\x6c\x62\xd1\x67\xab\x6b\x73\x03\x1c\x7f\x88\x5d\x03\x81\x0d\x89\x4e\x21\x67\xc5\x94\xb7\x5a\x36\xb4\xa0\xd5\xef\xa1\x93\xef\x61\x54\xbd\xb1\xa7\x57\xe1\x7b\x53\x3e\xc8\x70\x9b\x4e\x1e\xdb\x4d\xf0\x7d\xeb\xa9\xc7\x8e\xd3\xe3\x7f\x71\x1a\x1e\xbf\x29\xfc\x65\xd7\x99\xd2\xd0\x0f\x75\x1e\x17\xb1\x42\x99\xba\x3b\xf4\xac\x44\xf6\xfb\x49\x1f\x7a\x11\x5e\x5c\xea\xf6\xc2\xb3\x92\x80\xa2\x15\xf6\x1a\x05\x17\xb4\x61\xb3\x17\xfc\x03\x0b\x63\x85\xf2\x90\xc5\xf5\x04\x88\xb7\x35\x86\xcb\x06\x62\x55\xe4\xfd\x71\x85\xfc\x1d\x5b\x98\x13\x02\xc6\x8f\xc3\xf0\x09\x4d\xa5\xc3\xb9\x01\xc7\x6a\xe6\xca\xb9\x9d\xff\x9c\xa7\x5e\xb0\x77\x4c\x1d\x7a\x61\x27\x86\xd1\xbb\xf1\xab\x78\x2e\x1b\x33\x9e\x29\x9e\xad\xaa\xf7\xe7\x5c\x73\xc5\xcd\xb2\xcf\xcd\xfc\x24\xa6\x4b\x4d\x8c\xf3\x46\x6a\x9d\x87\x51\x4d\xde\x51\x17\x9f\x6d\xd8\xbe\x7a\x3e\xff\x02\x00\x00\xff\xff\x74\x34\xa2\x29\xa3\x06\x00\x00"

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

var _localesRuLc_messagesWidgetMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x91\xcb\x6b\x13\x51\x14\xc6\xbf\x24\xad\x4a\x14\x17\xbe\x56\x2e\x8e\x4a\x8b\x2e\xa6\x26\x55\x44\xa6\x9d\xc6\x07\x2d\x8a\x09\x48\x8d\x5d\x88\x9b\x21\xb9\x9d\x0c\x4e\x66\xc2\xbd\x33\xad\x42\x17\x8d\x6e\x0a\x16\x04\xc1\x95\xba\x72\xe7\x2a\xad\x94\xa6\x2d\x89\x5b\xdd\xc8\xb9\xe0\xc2\x95\x7f\x89\x0b\x99\x49\xda\xd4\x22\xb8\xf3\x6c\xe6\x3c\xbe\xef\x9c\x1f\x73\x7f\x9e\x18\x7a\x03\x00\x47\x00\x9c\x05\xf0\x18\xc0\x51\x00\xef\xd0\x8b\x2f\x00\x86\x00\x7c\x05\x30\x0c\xe0\x1b\x80\x0c\x80\xef\x00\x2e\x00\xf8\x01\xe0\x0c\x80\x5f\x00\xce\x03\x38\x99\x02\x4e\x03\xb8\x96\x02\x0a\x29\xa0\x98\x02\x8e\x03\x58\x4d\xf7\xbe\x6f\xd3\xc0\x31\x00\x1f\xd3\xbd\x5b\x9b\x69\xe0\x26\x80\x53\x19\xe0\x11\x80\xb9\x0c\x50\x00\xb0\x96\x01\x52\x00\xd2\x7d\xb6\xe1\x3e\xc7\x21\x00\x87\x31\x88\xcc\x6e\x72\xc7\x75\x6a\x28\x8a\x05\xe1\xa1\x18\x2c\xe2\x81\x08\xa9\xe6\x3a\x35\xf2\xe2\x16\xcd\xdb\xae\x27\xaa\xb4\xe8\x86\x35\x12\x52\x06\x92\x46\xd4\x41\x8d\x8a\x2a\x15\xa1\x7a\x6d\x2f\x58\xfc\x97\x73\x20\xd9\x35\xce\x8a\x46\x20\x43\xa3\xa4\x1c\xb7\x6a\xdc\x8a\x1c\x65\x94\x03\x93\xaa\x62\xe1\xc6\x13\xb7\x66\xd7\x83\x31\x19\x65\x8b\xb6\x0a\x8d\xb2\xb4\x7d\xe5\xd9\x61\x20\x4d\xba\x97\x8c\xa8\x14\x49\xbb\x1e\x54\x03\x9a\xfc\x43\x3f\x95\x2d\xda\xbe\x13\xd9\x8e\x30\xca\xc2\xae\x9b\xb4\x57\x9b\x34\x1b\x29\xe5\xda\x7e\xb6\x74\xb7\x34\x6d\xcc\x09\xa9\xdc\xc0\x37\x29\x3f\x96\xcb\xde\x0e\xfc\x50\xf8\xa1\x51\x7e\xd6\x10\x26\x85\xe2\x69\x78\xb9\xe1\xd9\xae\x3f\x41\x95\x9a\x2d\x95\x08\xad\x87\xe5\x19\xe3\xfa\x40\x17\xf3\xcc\x0b\x69\x4c\xfb\x95\xa0\xea\xfa\x8e\x49\xd9\xfb\x5e\x24\x6d\xcf\x98\x09\x64\x5d\x99\xe4\x37\x92\x52\x59\x57\x26\xa8\x97\x5a\xfe\x48\x3e\x67\x59\x79\x1a\x1d\xa5\x38\xcd\x9d\xb3\xf2\x79\x2a\x50\x8e\xcc\xa4\x9e\xb2\xc6\x77\x47\x93\xd6\xd5\x38\xbd\x98\xc8\x26\xf3\x39\x5a\x5a\xea\x59\xa6\xac\xf1\xdc\x25\x2a\x50\x9e\x4c\x1a\x9f\x00\xbf\xd6\x2f\x75\x93\xbb\xbc\xcd\x6d\xde\x02\x7f\xd0\xcb\xdc\xe5\x75\xde\xe0\x8e\x5e\x05\xbf\xe7\x36\x6f\x0e\x66\x4d\xfd\x9c\x5b\xdc\x49\x14\xdb\xdc\x22\x5e\xdf\x73\x77\xf9\x13\x77\x49\xbf\xe8\xfb\x3b\xfa\x15\xf1\x26\xb7\xe2\x55\x7a\x59\xaf\x70\x9b\x77\xb8\xa5\x9b\x7a\x95\x74\x93\xb8\x9b\x74\xd6\x12\xdf\x56\xfc\xba\xf1\x72\xfe\xcc\x1b\x7a\x25\x5e\x1f\xef\xd9\x7f\x6b\x27\xe6\xd9\x7f\xad\xcd\x5b\x83\x5b\x7d\xd6\xbf\xe0\x75\xfa\xf8\xff\x03\xae\x33\xf8\x55\x07\xd1\x7e\x07\x00\x00\xff\xff\x66\x8e\xf0\x7e\xf1\x03\x00\x00"

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
