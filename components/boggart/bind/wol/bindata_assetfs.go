// Code generated for package wol by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
package wol

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x54\x4d\x6f\xdb\x30\x0c\xbd\xf7\x57\x10\x3c\x6d\x03\x5c\x37\x29\xd0\x0d\x83\x1d\xa0\xd8\x65\x87\x74\x87\x7d\x60\xc7\x42\xb6\xe8\x5a\x9d\x3e\x3c\x89\xce\xd2\x19\xf9\xef\x83\x3f\x92\x38\x89\x97\xe5\xb0\x20\x40\x28\xe9\x51\x7c\x79\x8f\x62\xd3\x80\xa4\x42\x59\x02\xcc\x9d\x65\xb2\x8c\xb0\xd9\x5c\x25\x52\xad\x20\xd7\x22\x84\x14\xbd\xfb\x85\x8b\x2b\x00\x80\xf1\x6e\xee\x74\x64\x64\x34\x9b\x43\x1b\x05\xb3\x8d\xd6\x21\x9a\xcd\x07\xfc\x71\xce\xfa\xb1\x12\x96\xf4\xe8\xf4\x14\xc1\x8a\x35\x1d\x21\x3a\x54\x39\x5f\x34\x0d\xa8\xd9\x3b\x0b\xf8\x5d\xfc\x20\x70\x16\x96\xf7\x9f\x10\xae\x61\xb3\x49\xe2\x72\x3e\x91\x33\xe6\xab\x49\xf8\x42\xad\x71\x91\xc4\x52\xad\x8e\x28\x4c\x6c\x1d\xb0\xda\x6a\x33\x51\xa3\x70\xde\x6c\x81\x6d\x1c\x95\xce\xab\xdf\xce\xb2\xd0\xd0\xad\xb5\xc8\x48\x47\x9a\x0a\x46\xf0\x4e\x53\x0f\x43\x30\xc4\xa5\x93\x29\x56\x2e\x30\x82\x92\x29\x06\xb2\x12\xc1\xba\x95\xd0\x4a\x0a\xa6\xd3\x6a\xc7\xcc\x14\x93\xe9\xab\x3c\x79\x57\x57\x13\xfc\x76\x59\x1d\x8d\x16\x9b\x62\xa5\xc5\x0b\x79\xdc\x5b\x69\xd9\x3b\xdd\x13\x85\xc1\xd8\xdb\xad\xaf\xb7\x93\xb6\x1e\x7f\x76\xd6\x3c\xdc\x7f\x00\x21\xa5\xa7\x10\x7a\x6b\x20\x09\x95\xb0\xbb\x66\xa2\x9f\xb5\xf2\x24\x71\xf1\x26\x89\xdb\x83\x33\x8c\xe3\x8e\xd0\x19\xc0\x69\x3b\xde\x6d\x59\xdf\x5d\xc4\xba\xbb\x45\xd9\xaa\x66\xe0\x97\x8a\x52\x64\x5a\x33\x1e\xb8\x39\x88\x83\x60\x85\xa1\x14\x8d\xc8\x7b\xaf\xba\xc0\x88\xb5\x26\xfb\xc4\x65\x8a\xb3\xb7\x08\x2b\xa1\x6b\x4a\xf1\xe6\xe6\xfd\xe1\xf7\x9c\x2d\xa7\x9d\x77\xc9\xd1\xe8\x9f\x6b\xfb\x18\x9c\x56\x72\xb2\xb1\xa7\x12\x2e\x6b\x97\xbf\x69\x6b\x64\xe4\x8a\x22\x10\x47\xb7\xff\x12\x36\xab\x99\x9d\x1d\x75\x76\xaf\x71\xa8\x33\xa3\xf6\x2a\x67\x6c\x21\x63\x1b\x85\x3a\xcf\xdb\xa6\xd9\xbf\xf2\x2f\x5d\x52\xff\xbc\xfb\xbb\xfe\xa7\x8e\x49\xdc\x0a\x71\x76\x0e\x8c\x96\x43\x38\xfc\x34\x0d\x90\x95\xed\xa0\xbc\x1a\x0d\xd0\xe7\xd0\xcd\x4e\xe8\x9f\x43\x60\xc1\x2a\xff\xf8\xf5\x61\x09\xaf\xfa\xf8\xdb\xe7\x25\x60\x2c\x45\x28\x33\x27\xbc\x8c\x45\x08\xc4\x21\x5e\x91\x95\xce\x87\x78\x78\xf5\xce\xc7\xcf\xa3\xc5\xb5\x51\xf6\xba\xbd\xb9\x10\x3a\xd0\xeb\xb6\xc0\xae\xfa\x9f\x00\x00\x00\xff\xff\xf7\xf4\x1d\x65\xbd\x05\x00\x00"

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
