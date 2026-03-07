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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x99\xdf\x6f\xdb\x36\x10\xc7\xdf\xfb\x57\x1c\x38\x17\x68\x81\x59\x8e\xdb\x97\x21\xb3\xd3\x87\xfd\x78\xd9\xb0\x0d\x6d\xb7\x3d\x0a\xb4\x78\xb2\xb9\xd2\xa4\xca\xa3\xe5\x18\x86\xff\xf7\x81\x12\xa5\x28\x8e\x63\xcb\xb6\x36\x14\xea\x43\x10\x86\x3c\x7e\xc8\xfb\xea\xee\x1b\xc1\xde\x6e\x41\x60\x2a\x35\x02\x4b\x8c\x76\xa8\x1d\x83\xdd\xee\xc5\x44\xc8\x1c\x12\xc5\x89\xa6\xcc\x9a\x35\xbb\x7b\x01\x00\xd0\x9c\x4d\x8c\x1a\x2e\xc5\x70\xfc\x06\xfc\x88\x96\xd5\xe8\x9e\x86\xe3\x37\x21\x7e\x7f\xcf\x7d\x9c\x71\x8d\xaa\xb1\xfa\x34\xa2\xba\xc5\xe3\x98\x3a\xce\x1a\x85\x53\xe6\xf8\xec\x10\xa9\x8e\x5c\xa9\x0a\xa8\x79\x0e\x9a\xe7\x43\xc7\x67\x04\x33\x6e\x63\x3f\x60\x0f\x18\x25\xe9\xd0\x59\x35\x49\xc9\x10\x9b\x59\x24\xd4\x8e\x3b\x69\x34\xdb\x6e\x41\xa6\x80\x9f\x21\xe2\x89\x9f\x00\xe6\x55\xab\xce\xf4\x73\x39\xfa\x20\xd4\x02\x76\xbb\xbb\x09\x87\x85\xc5\x74\xfa\xcc\xbe\x7f\x78\xce\x29\xb1\x32\x73\xb7\xb9\x91\xe2\xd5\xcd\xeb\xef\xfd\x5e\x45\x08\xbb\xdd\x76\x0b\xd1\x7b\xfc\xbc\x42\x72\xd1\x9f\xef\x7f\x8d\xfe\xe0\x6e\x51\x4e\x97\x70\x76\xe7\xa1\xe3\xef\x34\xb0\x1f\x31\x97\x09\x32\x88\x60\xb7\x9b\x8c\xf8\xdd\x64\xa4\x64\x07\xa9\x25\x66\xa5\x1d\x5a\xba\x34\xc5\xe6\xfe\x4b\x52\x7d\x57\x82\xa6\x15\xe7\x50\xea\x3f\xd4\x67\x9c\x4e\x7e\x32\x5a\xa9\x67\x56\x1a\x95\xe8\xf8\x6c\xf8\x7c\x2d\x1e\xab\xc9\x26\xc1\xcf\x40\xca\x05\x42\xa9\x18\x48\x7d\x84\x56\x10\x53\x63\x97\x01\xe9\x87\x0c\x42\xfa\xec\x1d\x83\x25\xba\x85\x11\x53\x96\x19\x72\x0c\xa4\x98\x32\x42\xe7\xa4\x9e\x13\x03\xc1\x1d\x1f\x3a\x33\x9f\xfb\x9d\x39\x57\x52\x70\x67\xec\x89\xd3\x8e\x3f\xad\xa3\x5b\x8b\xcb\xfa\x0e\xc2\x46\xc2\x0a\x41\xb8\xa1\x45\xca\x8c\x26\x9f\xaf\x36\x6b\xcb\x33\x06\xe4\x36\xfe\x5e\x6b\x29\xdc\xe2\x76\x7c\x73\xf3\xf2\xc4\xc5\x1e\x4e\x58\x20\x17\x6d\x63\x6d\xbb\xc0\x00\xae\xee\xbd\x14\x43\xef\x5c\xe3\xa7\xf5\x04\xdf\x54\x05\xe5\x16\x5d\xa1\x3f\x28\xe3\xfe\x0b\xee\x4f\xda\xcb\x2f\x2e\x04\x3f\x70\xfe\xe2\x6a\x85\xe7\x52\x26\xa3\xb6\xda\x7b\xe6\x19\x4f\x74\x66\xc4\xa6\x5d\xec\x76\x0b\x96\xeb\x39\xc2\x40\x7e\x0b\x83\xdc\x67\x01\xb7\x53\x88\x42\x41\xc7\xc5\x0c\xb5\xa9\xea\x06\x71\x90\xa4\x73\x4f\x79\x25\xb5\xc0\x7b\x18\x54\x34\xf2\xff\xa7\x52\x39\x27\x18\xc8\xd7\xe7\x30\xcf\xaa\xd1\x72\x83\xf0\xcf\xa6\x4c\x28\xfa\x6d\xb5\x9c\xa1\x2d\x1f\x4d\x4b\x0d\x0f\x83\x8a\x2a\xbc\x02\x23\xd3\x42\x9b\x28\x54\x5d\xe9\xde\x65\x01\x6d\x30\xd8\xf0\x23\x67\x2f\xd7\xb4\x79\x58\x2a\x3c\xfc\xf2\x0b\x14\xbc\x97\xd1\xdb\x14\x96\x6f\x19\x0c\x40\x4b\x55\xff\x84\x24\x8b\x52\x3e\xfb\x90\xf6\xb5\x5c\x67\xd1\xb6\xf2\xdb\x55\xf3\x64\x54\x38\xe9\x49\xe7\x0e\xe2\x7e\xbd\x2e\xfd\x7b\x56\xbc\xb2\x7c\xad\x8e\x77\xa6\x92\x8d\xb6\xf9\x59\x19\x63\x61\xcd\x69\x21\xf5\xbc\xce\xfc\x9c\x1e\xa9\x5d\x20\x4a\x3d\x2b\x0e\xac\xff\xc1\x07\xce\x10\xfe\x72\x7d\xa4\x25\x07\x73\x6b\x56\x19\x70\x85\xd6\x5d\xab\x91\xe7\xc5\x05\x2f\x2e\x78\x3d\xd1\xe9\x03\x26\x46\x8b\x0e\x85\xa2\x02\xd8\x43\xa5\xfe\x96\x16\x15\x12\x01\xa1\x26\x63\x41\x99\x35\xcc\xb8\x73\x68\x37\x57\x4a\xb6\x0e\xe4\xb8\x24\xc7\xca\xac\xe3\x40\xee\xad\x76\x44\x9d\x8b\x46\xd4\x13\xb5\x9a\xde\xe5\x78\x06\x89\x32\xd4\x81\xcb\x37\x1c\xcc\xf1\x2c\x0e\xd4\x9e\x68\xf6\xc8\xc7\xba\x13\xed\x91\x9b\xf5\x4f\xb5\xfd\xbe\xcc\xb8\xb4\x3e\xbb\xa5\x11\xd8\x71\x7f\x06\x74\xec\xd1\x3d\x51\x6f\xbf\x4f\xc9\x71\x77\xad\x6c\xfb\x5d\x5a\x30\x7b\xa2\xd7\x93\x1e\xed\x42\xb0\x27\x1d\xda\x27\xc5\x3e\xae\x4d\x29\x17\x75\xd1\x92\x6e\x6d\x4a\x9d\xa8\x4f\x5d\xf8\x91\x67\x54\xb9\x3d\x18\xdd\xe1\x2b\x86\xe3\x19\x55\x86\x1f\x1b\xdd\xc3\x37\x8d\x5f\x70\x33\x33\xdc\x0a\x50\x26\xf9\x74\xa5\x5a\x9f\x02\x2b\xf6\xac\x2f\x4b\x9f\xce\x3f\xb6\x39\xfd\x91\xd1\x64\x94\x1a\xbb\x3c\xf2\xb5\xc3\x48\xc8\xfc\xb9\x6f\x37\x0e\x2e\x1d\x98\xde\x9b\x6a\xfc\x19\x86\xe1\x57\x7d\xe5\x7f\x03\x00\x00\xff\xff\x39\x9e\x1c\x53\x2c\x1c\x00\x00"

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
