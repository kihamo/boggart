// Code generated for package aqicn by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
package aqicn

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x98\xff\x6e\xe2\x46\x10\xc7\xff\xbf\xa7\x18\x59\x54\x02\xe9\xb0\xc1\x49\xa4\x2a\x72\x2c\x55\x44\xa7\xb4\xba\x04\x54\x72\xcd\x9f\xd1\xe2\x1d\xcc\xaa\xf6\x2e\xb7\x5e\x0e\x2c\x8b\x77\xaa\x54\xf5\x05\xf2\x64\x95\x7f\x10\x1b\x63\xd6\x49\x48\x2a\xf5\x2f\xcc\xee\x8c\xfd\x99\xef\xcc\xac\x07\x92\x04\x28\xce\x19\x47\x30\x3c\xc1\x15\x72\x65\xc0\x76\xfb\x29\x49\x80\xcd\xc1\x9c\x23\xd2\xf4\xab\x43\xd9\x0f\xf0\x02\x12\x45\x57\x86\x14\x6b\xc3\xfd\x04\x00\x50\x5d\xf5\x44\xd0\x0f\x69\x7f\x68\x43\x7a\x15\x85\xbb\xab\x4d\xd4\x1f\xda\x85\x7d\xdd\x67\xf3\xb8\x24\x1c\x83\xca\xee\xa1\x85\x62\x2a\xc0\x9a\x45\x66\xb5\xb0\x5d\x87\xc0\x42\xe2\xfc\xca\x48\x92\x1c\xd5\x1c\x31\x15\x9b\xdf\x7e\xff\x0a\xdb\xad\x01\x8a\x48\x1f\xd5\x95\xf1\x38\x0b\x08\xff\xd3\x70\xf7\xad\xee\x48\x88\xb0\xdd\x3a\x16\x71\x1d\x6b\x61\x37\x3c\xa1\x1a\x5d\x80\x44\xce\xd9\xc6\x70\x1d\x8b\xb2\x1f\x35\xe0\x86\xa5\xbd\x18\x76\xc2\x1e\x3e\x23\x97\x99\x70\x0a\x5d\x8e\xd0\x0d\x90\x17\x88\xbf\x28\x25\xd9\x6c\xa5\x98\xe0\x51\x0f\x06\x3d\xe8\x32\x4e\x71\xd3\xb0\x0b\x83\x9e\xf9\x55\xf8\x22\x4d\x93\x2e\x84\x32\x6d\xda\x40\x8b\xe4\x1d\x31\xad\x9b\xaf\x91\xa8\x05\xca\x3e\xf3\x04\xd7\xb8\x64\x6e\x2c\xf4\x21\x92\x5e\x96\xad\x2e\xf9\xce\x3c\xfe\x98\xba\xb5\x47\xd6\x33\xa7\x4a\x32\xee\xa7\x49\xb5\x34\x5c\x87\x69\x78\xc9\x56\x43\xf0\xaf\x8d\x5e\xe1\xa6\x29\xb9\x7b\x6e\xb5\x72\xd5\x06\xad\x29\x60\xad\x5f\x6b\x49\xbf\x5d\xaa\x23\xcb\x49\x02\xc8\xe9\xfb\x94\x5e\x55\x4d\x58\xae\x82\xa0\x2f\x99\xbf\xd0\x09\x5b\xf5\x56\x18\x2e\x51\x12\xb5\x92\x68\xb8\xce\x2c\xd5\x8a\x0d\x7f\xe6\xd0\xcd\xa5\xba\x67\x21\x9a\xd3\xfc\xe3\x8b\x90\x21\x51\x60\xdc\x0a\x4e\x49\x6c\xf4\xc0\xcc\x44\x9b\xb9\x9f\xe1\xf9\x88\x68\xb2\x1f\x5e\x5c\x0e\xce\x8d\xcc\xf6\x9d\xe4\x7b\x73\x7f\x9e\xbf\xa1\x42\x73\x4d\x03\x9c\xb7\xd5\x6a\xe5\xec\x37\x7f\x25\xdf\x99\x79\x9f\xc6\xcc\x5c\xd8\x69\x6a\xdc\x57\xc4\x86\x0e\x6c\xb7\x97\xe0\x44\x4a\x0a\xee\x97\x87\x6c\xee\x69\xfe\x91\xe9\x55\x6c\x3a\x16\x73\xe1\xe9\xaf\x91\x33\x93\x60\xb9\xc7\x8b\x47\x4b\x33\xa9\xd3\x4c\x24\x46\x51\x2b\xca\xa4\x01\x65\x31\x21\x27\xa1\xdc\xd4\x51\x6e\x56\x21\xa3\x4c\xc5\x7a\x94\x9b\x06\x94\x9f\x4e\x02\xb9\xc6\x75\x1d\x25\x5d\x5a\x0a\x96\xbe\xcf\x75\x2c\xd7\xb8\xfe\x80\x1c\x3d\xd4\x69\x1e\x18\xa7\x10\x2d\x11\xa9\x1e\xe7\xa1\x01\x26\xb4\xa2\x97\xc2\xbc\xdb\x3b\xe0\xbf\x6d\xb0\x49\x68\x5f\x1c\x54\xf5\xad\x7d\xd1\x52\xd1\xa1\x7d\xd1\x94\xbb\x7f\x7c\x2b\x7c\xfa\xfb\xb4\x1e\x0b\x87\x83\x43\xa0\xe1\xa0\x0d\x68\x38\x38\x1d\xe8\x7f\x9a\xc3\xf1\x59\x5d\xb0\xf1\x99\x5e\xae\xf1\x59\x83\x58\xcb\x65\x78\x52\xe6\xee\x84\x5d\xe7\xb8\x1b\xdb\x7a\x90\x3b\x61\x7f\x00\xc9\x48\xd4\x41\x46\x63\x3d\xc7\x48\x7c\x00\xc6\xf4\x50\x90\x69\x9b\x20\xd3\x53\x05\xf9\x90\x09\x41\xfb\x23\xa4\x6e\x2c\xc5\x1a\x76\x85\x4d\x49\x1c\x1d\x29\xe7\x24\x01\x49\xb8\x8f\xd0\x61\x9f\xa1\x43\x49\x0c\x97\x57\x85\x12\x5f\x84\x44\x8f\x44\xca\xbc\x26\x2c\x88\x77\x67\xd4\xd1\x90\x73\xe1\x03\x84\x8e\xe9\xad\xa4\x44\xae\x1e\x29\x89\xcd\x6f\x9c\x6d\xb2\x1b\x9b\xd7\x24\xce\xe7\xa9\x6c\x49\x27\xde\x61\x2b\xdb\x6d\xe3\x75\xc5\x85\xa6\xb8\xfd\x22\x76\x48\x9b\xba\xef\x21\x57\x28\x5b\xee\x01\xf9\x98\x5e\xde\x26\xce\x26\xee\x7d\xf6\xdd\x2c\x38\xb0\xd3\x41\xf0\xb9\xa4\xba\xcd\x56\xbf\x11\xbe\x22\x32\x1d\x31\x3b\x59\x35\xe9\x66\xf2\x12\xe1\xfc\x19\x01\x7d\x89\x18\x19\x6e\xe3\x4b\x61\xc7\x76\xcb\x78\x8a\xd2\x2f\x17\x48\x2a\x6f\x56\xae\xd6\xe2\xfc\xe4\x27\x96\xa7\x7e\xf9\xfb\xa3\x73\xa4\x44\x86\x03\xe8\xb0\xde\x1e\xd2\x0b\x3d\x5e\xc1\xac\x69\xae\x17\x6c\xeb\x9b\x57\xbf\xfb\xba\x6e\x6c\xfd\xa7\xa0\xf2\xb5\xb8\x2c\x3e\x4a\x8a\xf2\xea\xdf\x00\x00\x00\xff\xff\x24\x28\x44\xb1\xa9\x11\x00\x00"

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
