// Code generated for package neptun by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
package neptun

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x98\xdf\x8f\x1a\x37\x10\xc7\xdf\xf3\x57\x8c\xdc\x8b\xd4\x4a\x05\x8e\xe4\xa5\xba\xc2\xe5\xa1\x55\x5f\x5a\xb5\x55\x92\xb6\x8f\xab\x61\x3d\x80\x1b\xaf\xbd\xf1\x98\xdd\x20\xc4\xff\x5e\x79\x7f\x70\x7b\xc0\xc1\x1e\x6c\xa5\x68\xef\x01\x61\xec\xf1\x67\x3c\x5f\x7b\x06\xcb\x9b\x0d\x48\x9a\x2b\x43\x20\x62\x6b\x3c\x19\x2f\x60\xbb\x7d\x35\x91\x2a\x83\x58\x23\xf3\x54\x38\x9b\x8b\xfb\x57\x00\x00\xcd\xde\xd8\xea\x41\x22\x07\xe3\x37\x10\x5a\x9c\xd4\xad\x2f\x3c\x18\xbf\xa9\xec\xf7\xe7\x7c\x89\x52\x34\xa4\x1b\xa3\x87\x16\xf5\x2a\x1e\xdb\xec\xec\x9c\xd5\x34\x15\x1e\x67\xc7\x48\x3b\xcb\x95\xae\x81\x06\x33\x30\x98\x0d\x3c\xce\x18\x66\xe8\xa2\xd0\x10\x0f\x18\xad\xf8\x98\xaf\x1d\x49\xab\xca\x36\x75\xc4\x64\x3c\x7a\x65\x8d\xd8\x6c\x40\xcd\x81\x3e\xc3\x10\xe3\xd0\x01\x22\xa8\x56\xfb\x0c\x7d\x19\x05\x23\x32\x12\xb6\xdb\xfb\x09\xc2\xd2\xd1\x7c\xfa\xc4\xbc\x7f\x31\x43\x8e\x9d\x4a\xfd\x5d\x66\x95\xfc\xf6\xf6\xbb\x1f\xc3\x5c\xcd\x04\xdb\xed\x66\x03\xc3\xf7\xf4\x79\x45\xec\x87\x7f\xbd\xff\x6d\xf8\x27\xfa\x65\xd9\x5d\xc2\xc5\x7d\x80\x8e\x7f\x30\x20\x7e\xa6\x4c\xc5\x24\x60\x08\xdb\xed\x64\x84\xf7\x93\x91\x56\x1d\x84\x16\xdb\x95\xf1\xe4\xf8\xd2\x10\x9b\xf3\x2f\x09\xf5\x5d\x09\x9a\xd6\x9c\x63\xa1\xff\xb4\xf3\x71\x3e\xf8\xc9\x68\xa5\x9f\x18\x69\x9c\x44\x8f\xb3\xc1\xd3\x67\xf1\xd4\x99\x6c\x12\x42\x0f\xcc\x51\x12\x94\x8a\x81\x32\x27\x68\x05\x71\x6e\x5d\x52\x21\x43\x53\x40\x15\xbe\x78\x27\x20\x21\xbf\xb4\x72\x2a\x52\xcb\x5e\x80\x92\x53\xc1\xe4\xbd\x32\x0b\x16\x20\xd1\xe3\xc0\xdb\xc5\x22\xcc\xcc\x50\x2b\x89\xde\xba\x33\xde\x4e\xef\xd6\xc9\xa9\xc5\x62\x43\x06\x51\x23\x60\x4d\x20\xfd\xc0\x11\xa7\xd6\x70\x88\xd7\xd8\xdc\x61\x2a\x80\xfd\x3a\xac\x2b\x57\xd2\x2f\xef\xc6\xb7\xb7\xaf\xcf\x2c\xec\xc1\xc3\x92\x50\xb6\xb5\x75\xed\x0c\x2b\x70\xbd\xee\x44\x0e\x42\xe5\x1a\x1f\x9e\x27\xf8\xa6\x3e\x50\x7e\xd9\x15\xfa\x83\xb6\xfe\x62\xee\x03\xe6\x6f\xd4\x2b\x7a\x2e\x65\x32\x6a\x2b\x51\x60\x3e\x43\xf8\x99\x95\xeb\x76\xb6\x9b\x0d\x38\x34\x0b\x82\x1b\xf5\x3d\xdc\x64\x21\x0a\xb8\x9b\xc2\xb0\x68\x71\x9b\x43\x07\xcf\xdf\x6a\x19\x74\x2b\x9d\x0d\x7f\x5f\x25\x33\x72\xa5\x6c\x2d\xe3\x3b\x84\x14\x9b\x78\x21\xa2\xdc\xbf\xd7\xc3\xb7\x73\x48\xde\x0a\xb8\x01\xa3\xf4\xee\x53\xf1\x8b\xdd\x7d\x96\x83\xf6\x5b\xbb\xab\x9f\x6d\xb9\xed\x36\x77\x32\x2a\xf2\xff\x6c\xbd\xa9\x8a\xfd\xcb\xad\x2d\x7f\xa4\xc5\x1f\xed\x4b\x2d\x00\x97\xa4\x6e\x19\xf1\x2f\xda\x5a\x07\x39\xf2\x52\x99\xc5\x2e\xf2\x0b\x12\x70\x0e\xc3\x79\x60\x45\x15\xab\xbc\x7b\x94\x4e\xd6\x54\x5d\x22\x1e\xdd\x4b\xca\x31\x63\x1f\x86\x8a\x0c\xfa\x3f\xd2\xf3\x0a\x7d\x94\x63\x0f\x0b\x67\x57\x29\xa0\x26\xe7\xaf\xd5\x28\xf0\xa2\x82\x17\x15\xbc\x9e\xe8\xf4\x81\x62\x6b\x64\x87\x42\x71\x01\xec\xa1\x52\xff\x28\x47\x9a\x98\x81\xc9\xb0\x75\xa0\x6d\x0e\x33\xf4\x9e\xdc\xfa\x4a\xc9\xf2\x8a\x1c\x95\xe4\x48\xdb\x3c\xaa\xc8\xbd\xd5\x8e\xb9\x73\xd1\x98\x7b\xa2\x56\xb3\x76\x79\x4c\x21\xd6\x96\x3b\xa8\xf2\x8d\x0a\xe6\x31\x8d\x2a\x6a\x4f\x34\x7b\x54\xc7\xba\x13\xed\x51\x35\xeb\x9f\x6a\xfb\x79\x99\xa2\x72\x21\xba\xc4\x4a\xea\x38\x3f\x2b\x74\x14\xd0\x3d\x51\x6f\x3f\x4f\xd9\xa3\xbf\x56\xb6\xfd\x2c\x2d\x98\x3d\xd1\xeb\x20\x47\xbb\x10\xec\x20\x43\xfb\xa4\xd8\xc7\xdc\x96\x72\x71\x17\x29\xe9\x73\x5b\xea\xc4\x7d\xca\xc2\x8f\x98\x72\x5d\xed\xc1\x9a\x0e\xaf\x18\x1e\x53\xae\x0b\x7e\x64\x4d\x0f\x6f\x1a\xbf\xd2\x7a\x66\xd1\x49\xd0\x36\xfe\x74\xa5\x5a\x9f\x2a\x56\x14\x58\x5f\x97\x3e\x9d\x3f\xdb\x9c\x7f\x32\x9a\x8c\xe6\xd6\x25\x27\x1e\xcb\x47\x52\x65\x4f\xbd\xc9\x1f\x1d\x3a\xd2\xbd\xd7\xd5\xf8\x59\x35\xab\xaf\xdd\x92\xff\x0b\x00\x00\xff\xff\xdb\x36\x0b\x34\xe2\x1a\x00\x00"

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
