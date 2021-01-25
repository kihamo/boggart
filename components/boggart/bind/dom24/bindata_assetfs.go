// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// DO NOT EDIT!

package dom24

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x57\x5d\x6f\x9b\x3c\x14\xbe\xef\xaf\x38\xa2\xd5\x7b\x47\x51\xf2\xee\xa2\xea\x28\x52\xb7\xaa\x52\xa4\x29\xab\xa6\xee\xba\x32\xe0\x24\xd6\x1c\x1b\xd9\x4e\x9a\x08\xf1\xdf\x27\xdb\x90\x00\x4d\x08\x26\x59\x73\x11\x81\x7d\xbe\xfc\x9c\xe7\x1c\x1f\xf2\x1c\x52\x3c\x23\x0c\x83\x97\x70\xa6\x30\x53\x1e\x14\xc5\x15\x00\x40\x9e\x03\x99\xc1\x2d\x4a\x12\xbe\x62\x4a\x56\xcb\xe5\x96\x40\x6c\x8e\xe1\xa6\xdc\x85\xfb\x87\xc3\x92\x61\x4a\xd6\x90\x50\x24\xe5\x83\x27\xf8\xbb\x17\xed\x76\xda\xbb\x09\xa7\xfe\x32\xf5\x47\x63\xd0\x4f\x72\x59\x3d\x6d\xa4\x3f\x1a\xb7\xf4\xda\xba\x9b\xb7\x0c\x31\x4c\x0f\x48\x7d\x94\x54\x44\x51\x7c\x44\xd2\x48\x2f\xc6\x51\x9e\xef\x4e\x76\xfb\xba\xcd\x30\x14\x05\x5c\xd7\x17\x27\x4f\x50\x14\x61\xb0\x18\x77\xd8\x59\xd1\xca\x29\x43\x6b\x60\x68\x1d\x23\xe1\x0b\x32\x5f\x28\x30\xd1\xbe\x29\xce\x69\xcc\x37\x1d\xb1\x18\x3b\x94\x44\x21\xaa\xa1\x44\x51\x26\xb1\x4f\x09\xfb\xe3\x45\x21\xa9\x36\x66\x48\xc2\x0c\xf9\xc9\x02\xaf\x05\x67\xfe\x2a\xf3\xa2\x30\x20\x51\x18\xa0\x28\x0c\x28\xe9\x88\x33\x58\xd1\x8e\xdd\x7a\x86\x28\x46\x62\x46\x36\xda\x72\x4a\xd6\x47\xc0\xee\xd8\x6a\xe4\xa1\xa2\x5b\x87\x6f\x89\x13\x45\x38\xdb\x9f\xdd\x68\x00\x61\x6b\x4e\x92\xae\x1c\xb6\xbd\x7d\x64\xde\x29\x8d\xd2\x87\xbf\xc0\x28\xc5\xa2\x87\x32\x58\xea\x54\xfa\xd9\x8a\x52\x9b\x6c\xaf\xc1\xa6\xc7\x34\x15\x58\xca\x93\xec\xd9\x99\x3c\x0e\xa7\x8b\x48\x13\x8b\x0a\x41\x9f\xb0\x19\x77\x04\xa6\x2c\xce\x2f\x3b\x1b\x09\x3f\x56\x77\xed\x9f\x6e\x28\xa3\x3b\x06\xde\xb3\xe0\x4b\x0f\x6e\xea\x7d\xa2\xd3\x3d\xb2\xa0\xf5\xf3\x62\x34\xa4\x12\x9c\xcd\x1b\xc8\x3f\x4f\x7e\x1a\xd4\xcb\xad\x7e\x9e\x83\xde\xae\x7b\x64\x01\xfe\x01\x94\xaf\xfc\xd3\x81\xfc\xce\x97\x19\x62\xdb\x3a\x98\x61\x2c\x20\xe8\x6f\x75\x17\xfd\x64\x3a\xb5\xe1\xdf\xc3\x01\x0f\x93\xe9\xb4\xf7\xc9\x2e\x9a\xa8\x1e\x22\xf6\x72\xdc\x5f\x12\x02\x31\x89\x4c\xb7\x92\xa7\x62\x3e\xaf\x31\x29\x14\x77\x5e\x5f\x0d\x3d\x23\xdc\xd0\x04\xf3\xef\x4b\x25\x48\x86\x53\x60\xfc\x5d\xa0\xac\xa7\x39\x6b\x52\x77\x44\x17\x79\xd1\x5f\xb8\x74\x10\xed\xf8\xf1\x84\x14\xb6\x04\x09\x03\xb5\x38\xc3\xd0\x0b\x16\x84\xa7\x17\x31\xf5\xb8\xd4\x19\xbf\x88\xa9\x67\x42\x07\x1d\x2f\x0c\x5c\x50\xd5\xb6\x1d\x73\x16\xf3\x74\xeb\x54\xce\xe5\x34\xa8\xf6\x65\xa0\x27\xc2\x41\xe5\xd1\x8c\xc4\x99\x3d\xa9\xe9\x55\xb5\x40\x6e\x35\x89\x2c\xc4\x0e\x18\x1c\x33\x66\x89\x34\xd4\x5c\x55\x8a\x65\xf7\xa8\x1b\x9e\xc8\x17\xb4\x5d\xea\xf1\xa6\x28\x14\xde\x28\x5f\xae\x92\x04\x4b\x99\xe7\x80\xa9\xc4\xd5\x6a\xaa\x81\x16\x7a\x91\xe9\x28\x1c\x2a\xb7\xfa\x9d\xf0\xed\x6c\xcf\x1c\xad\x3d\x86\x4a\x32\x67\x3e\x5f\x29\x1f\x51\xe5\x81\x19\xb9\xed\xa9\x6d\x35\x5a\x6f\x96\xfa\x76\x4e\x1d\x72\x8c\x12\x97\x0b\x46\x4c\xd8\xe1\x80\xbf\x11\x4a\xcf\x8e\xd6\x24\xec\x6a\x88\x6a\x23\x59\xb6\xff\xe8\x2f\x92\xff\xae\xef\xfe\xbf\x1b\x7d\x75\x23\xe1\xa0\x2a\xb8\x04\xc9\x74\xb3\xfb\xfd\xeb\xc7\xe0\x84\x21\x58\x08\x3c\x33\x49\x39\x62\xd7\xab\x52\x1a\x2b\x06\xb1\x62\x66\xbc\xb5\x0f\x09\xb7\x2b\x1b\xe9\x81\x42\x62\x8e\xd5\x83\xf7\x16\x53\xa4\xbf\xa1\x06\x85\x03\x87\x38\x94\xf2\x77\x46\x39\x4a\x3f\x12\xe8\x69\xb7\x73\x06\x89\xc0\x0e\x3a\xe7\xb0\xcf\xcd\x95\x4b\xe6\xdd\xee\x24\xf7\x88\xc2\xc0\xe1\x56\x0a\x03\x33\xe6\x7c\xe2\x30\x78\xe2\x34\x61\x50\x7e\xcb\x3a\x7d\x34\x1f\x58\x6e\x2d\xb5\x5e\x9b\x81\xec\xdf\xf6\x4f\x7f\x03\x00\x00\xff\xff\xea\x6e\xd2\x2d\xf0\x11\x00\x00"

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