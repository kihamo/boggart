// Code generated for package v3 by go-bindata DO NOT EDIT. (@generated)
// sources:
// templates/views/widget.html
package v3

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x5a\x5b\x6f\xf2\x38\x10\x7d\xef\xaf\xb0\x22\xf5\x31\x40\x82\xb4\x5a\xad\x28\x52\x77\xbb\x7d\x5b\xa9\x52\xa5\xee\x23\x32\x78\x02\xd6\x26\x4e\xd6\x36\x97\x0a\xf1\xdf\x3f\xe5\x46\xae\x04\xd2\xda\x6e\xa8\xd4\x96\x8b\xcf\x9c\xe1\x9c\x99\xb1\x09\x1c\x8f\x88\x80\x47\x19\x20\x6b\x15\x32\x09\x4c\x5a\xe8\x74\x7a\x98\x11\xba\x43\x2b\x1f\x0b\xf1\x64\xf1\x70\x6f\xcd\x1f\x10\x42\xa8\xfc\xe8\x2a\xf4\xed\x80\xd8\x8e\x8b\xe2\x5b\x22\xc8\x6f\x1d\x84\xed\xb8\xd9\xfa\x3a\xe6\xb0\x88\x30\x03\xbf\xf4\x6c\x73\x85\xa4\xd2\x87\xda\x8a\x64\xd5\xc6\x9d\x1f\x8f\x88\x3a\xbf\x33\x64\xfd\xb5\xe5\x1c\x98\x44\x01\x60\xb1\xe5\x10\x00\x93\xc2\x42\x23\x74\x3a\xcd\xc6\x1b\xb7\x05\x5c\x4e\xdc\x07\xcc\x3d\x7a\xb0\xe6\xb3\x31\xa1\xbb\x5a\x2e\x2d\x0f\x55\xd2\xcb\x45\x6a\xe1\x90\x78\xe9\x43\xbe\x32\xbd\x93\xfc\xb5\x85\xe4\x34\x02\x82\x08\x96\x38\x7d\x9c\x48\x9b\x83\x88\x42\x26\xe8\x0e\x10\x0b\xf7\x1c\x47\x16\x12\xf2\xd3\x87\x27\x6b\x4f\x89\xdc\xfc\xe1\x4c\x26\x8f\x2d\x2c\x29\xd3\x06\x30\xb9\xf4\x1c\x6f\x7f\x22\x03\xe6\xf9\x05\xc4\x8e\xdd\x72\xad\x42\xd4\x7f\x0a\x31\x73\x2d\xe5\xe6\x4a\xb0\xd0\x17\x11\x66\x4f\xd6\xb4\x14\xe7\x03\xfb\x5b\xb8\x31\x42\x81\x7a\x01\xb1\xe2\x34\x92\x34\x64\xd7\xb0\xb3\xf1\xa5\x17\x19\x63\x3a\xa4\x59\x86\xe4\xb3\x2b\x9d\x0e\xe9\xd2\x05\xa4\xc8\xf7\x1d\x38\xc5\x3e\x62\xdb\x60\x09\xfc\x9c\xf1\x05\xea\xfc\x27\x46\x7b\x68\x24\x12\xec\x22\xc5\x8e\x80\xf3\x90\xc7\x4d\xd7\x85\xcc\xf8\xcf\xe5\x05\x07\x69\x13\xcc\xd6\x31\x77\xcd\x85\xb6\xf0\xa3\xbf\x33\x92\x9b\x72\x04\x5f\xc0\xcd\x09\x75\x92\xef\xe2\x5a\xb8\x99\x96\x91\x6b\xac\x15\x0b\x5e\x60\x47\x57\x80\x44\x5f\x27\x2e\xd7\x0f\xea\x5d\x06\x7f\x6e\xa9\x9f\x74\x36\xf4\xac\x81\x65\x0c\x5c\xc4\x40\xf5\x05\x50\x8f\x6d\xce\xfd\x12\x73\x62\xfd\xe8\x35\xe4\x01\x96\xc8\x72\x27\x93\xdf\xec\x89\x63\x4f\x5c\x4b\x71\x39\x98\xf3\xfa\x95\xf2\x60\x8f\x39\xa0\x1d\x70\x51\x9e\x53\xb7\x39\xee\x65\xf0\x45\x06\x57\xef\x7b\x3b\x83\x39\xf7\x1b\xfc\x3a\xda\xdf\x9c\xdf\xcf\x84\x70\x10\xa2\xa7\xcd\x38\x45\xa9\x77\xb7\x12\xd8\x9c\xa9\x39\xad\x89\x51\xce\x40\xee\x43\xfe\x1f\xc2\xb7\x2b\xaf\xd2\xf0\xb7\x70\xdf\x7b\x2f\x8f\x62\x8c\x7a\xb3\x4b\x61\x8d\x59\x9d\x8a\xf0\x38\x72\x3d\xf4\x6f\xac\x02\xa3\xfe\xf9\x37\x4b\x48\x6b\x11\x24\xf2\x23\xb1\x0d\x02\xcc\x3f\x7f\xc4\x7b\xe4\xe1\x95\x0c\x79\xdf\x96\x4f\xb4\x59\x64\x58\x4d\xb5\x50\x0d\x3f\xa0\x9a\x38\x27\x66\xa0\x36\x32\xaa\x9f\xaa\x91\x57\x0e\xff\x6f\x81\xad\x6e\x22\x46\xa5\x9d\x3f\xc7\x69\xd8\xf2\xab\xa1\x87\x51\x17\x45\x52\xf7\x7d\x00\x48\x8b\x4e\x62\x4e\x3d\x0f\x39\x5f\x1a\x0a\xd2\xd1\x35\x0f\xf2\xc8\xc6\x2d\x27\x17\x27\x81\x74\xf4\x0e\x81\x77\x83\x6d\xbf\x99\x77\x5f\xbc\x40\xf5\x0b\x18\x6f\x1b\x2c\xa0\x54\x25\xfd\xb1\xee\x37\xb0\xd3\x3e\xd8\x2b\x97\x65\x14\x76\xd0\x47\xe8\x4b\xbc\xee\xfb\xde\x78\x97\xa2\x16\x51\xfc\xd2\xbe\xdf\x40\x49\xd3\xb4\xc5\xd4\xd6\x3a\xb5\x09\xf9\xd1\xe8\x97\x5a\x36\xaa\xbb\xa6\x4d\x47\x57\x83\x8e\xee\xa0\x74\x74\x8d\xe8\x38\xd5\xa0\xe3\x74\x50\x3a\x4e\xef\x7c\xdb\x7e\x0e\x22\xe0\xfd\xa7\x0e\xce\x60\x4a\xc7\x4e\x6b\x50\x53\x3e\x3f\x37\x7c\xae\xa7\xa3\xa9\x61\xaa\x34\x8a\x26\x4f\x6b\xd0\xa1\x48\xa9\x6d\xf6\x54\x69\x14\x0d\x9f\xd6\xa0\x43\x91\xf2\xde\xc7\xcf\x97\xaf\x22\x29\x1d\x3c\xcd\x88\xa6\xfc\xbd\xf4\xee\x40\xef\xc8\x29\x71\x28\x9a\x37\xcd\x88\x83\x50\x50\xdb\xa4\x29\x71\x28\x1a\x33\xcd\x88\x83\x50\xf0\xee\x07\x4c\xfc\x22\x04\xc2\x6c\xed\xf7\x3e\xe3\xc4\x18\xb5\x07\x9c\x46\x44\x53\x1e\x13\x58\x37\x77\x92\x72\x36\xba\x76\xe4\x82\x43\xd5\xc9\xa6\x11\x71\x20\x1a\xea\x3b\xd5\x14\x1c\xaa\x8e\x34\x8d\x88\x03\xd1\xf0\xee\xa7\x8d\xaa\x0f\x46\x34\x1c\x6f\xda\x22\x1b\xb2\xbd\xfb\xc3\x10\x23\x87\x9d\x0a\x97\xd2\x43\x4f\x5b\xe4\x01\xe9\xaa\xf9\x08\x54\xe1\x52\x7a\x14\x6a\x8b\x3c\x20\x5d\x87\x35\xaa\x66\xe3\x0b\x5f\xd8\x9b\x8d\x93\x6f\x51\x76\x7e\x73\xb3\x74\x37\xbb\x99\xfd\x2b\xf2\xfe\x15\x00\x00\xff\xff\x32\xfd\x54\xfb\xfb\x2a\x00\x00"

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
