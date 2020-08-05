// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// DO NOT EDIT!

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\x5d\x8b\xe3\x36\x14\x7d\xdf\x5f\x71\x31\xec\xa3\x93\xd8\x81\x52\x4a\x66\x60\xdb\xed\xbc\x15\x16\x0a\xd3\xc7\xa0\x44\xd7\x89\xa8\x2d\xbb\x92\xf2\x31\x84\xfc\xf7\xe2\xaf\xf8\x33\x4e\xbc\x2b\x79\x1c\x76\x16\x76\x26\x89\x75\xee\xb9\x3e\xe7\xea\x4a\xa3\xf8\x74\x02\x8a\x1e\xe3\x08\xd6\x3a\xe4\x0a\xb9\xb2\xe0\x7c\xfe\xb4\xa0\x6c\x0f\x6b\x9f\x48\xf9\x64\x89\xf0\x60\x3d\x7f\x02\x00\x28\x7f\xba\x0e\x7d\x3b\xa0\xb6\xe3\x42\xfc\x4a\x06\xf9\xab\xa3\xb4\x1d\x37\x1b\x5f\xc7\x1c\x97\x11\xe1\xe8\x97\xae\x36\x47\x28\xa6\x7c\xac\x8d\x48\x46\x6d\xdd\xe7\xd3\x09\x98\xf3\x2b\x07\xeb\x8f\x9d\x10\xc8\x15\x04\x48\xe4\x4e\x60\x80\x5c\x49\x0b\x26\x70\x3e\x2f\xa6\x5b\xb7\x05\x5c\x4e\xdc\x47\x22\x3c\x76\xb4\x9e\x17\x53\xca\xf6\xb5\x5c\x5a\x3e\xaa\xa4\x97\x8b\xd4\xc2\xa1\xc8\xca\xc7\x7c\x64\xfa\x26\xf9\x69\x4b\x25\x58\x84\x14\x28\x51\x24\xfd\x9c\x2a\x5b\xa0\x8c\x42\x2e\xd9\x1e\x81\x87\x07\x41\x22\x0b\xa4\x7a\xf3\xf1\xc9\x3a\x30\xaa\xb6\xbf\x39\xb3\xd9\xe7\x16\x96\x94\x69\x8b\x84\x5e\xbb\x26\xda\x2f\x64\xc0\x3c\xbf\x80\xda\xb1\x5b\xae\x55\x88\xfa\x57\x21\x66\xae\xa5\xda\xde\x08\x16\xfa\x32\x22\xfc\xc9\x9a\x97\xe2\xbc\x12\x7f\x87\x77\x46\x28\x50\x5f\x51\xae\x05\x8b\x14\x0b\xf9\x2d\xec\x62\x7a\xed\x26\x63\x4c\x87\x34\xab\x90\xbe\x75\xa5\xd3\x21\x5d\x3a\x80\x16\xf9\xfe\x8d\x82\x11\x1f\xf8\x2e\x58\xa1\xb8\x64\x7c\x85\x3a\xff\x17\xa3\x3d\x98\xc8\x04\xbb\x4c\xb1\x13\x14\x22\x14\xf1\xa4\xeb\x42\x66\xfc\x97\xf2\xc2\xa3\xb2\x29\xe1\x9b\x98\xbb\xe6\x42\x5b\xf8\xc9\x9f\x19\xc9\x5d\x39\xa2\x2f\xf1\xee\x84\x3a\xc9\xf7\x71\x2d\xdc\x4d\xcb\xe9\x2d\xd6\x8a\x05\x5f\x71\xcf\xd6\x08\xb2\xaf\x13\xd7\xeb\x07\x7a\x97\xc1\xef\x3b\xe6\x27\x33\x1b\x7b\xd6\xc0\x2a\x06\x2e\x63\xa0\xfe\x02\xa8\xc7\x1e\xce\xfd\x12\x73\x62\xfd\xe4\x25\x14\x01\x51\x60\xb9\xb3\xd9\x2f\xf6\xcc\xb1\x67\xae\xa5\xb9\x1c\x86\xf3\xfa\x85\x89\xe0\x40\x04\xc2\x1e\x85\x2c\xf7\xa9\xfb\x1c\xf7\x32\xf8\x32\x83\xeb\xf7\xbd\x9d\x61\x38\xf7\x1b\xfc\x26\xa6\xff\x70\x7e\x7f\xa1\x54\xa0\x94\x1f\x36\xbf\x87\xcd\xf5\x2e\xcf\x51\x1d\x42\xf1\x2f\x90\xfb\x4d\xd1\x59\x0b\xdf\xc2\x43\xef\x65\x3e\x8a\x31\xfa\xed\x2f\x85\x1d\xcc\xf3\x54\x84\xcf\x13\xd7\x83\x7f\x62\x15\x38\xf3\x2f\xff\xb3\x84\x8c\x16\x41\x22\x3f\xc8\x5d\x10\x10\xf1\xf6\x2e\xde\x83\x47\xd6\x2a\x14\x7d\xbb\x41\xa2\xcd\x32\xc3\x1a\xaa\x85\x6a\xf8\x11\xd5\xc4\x25\xb1\x01\x6a\x23\xa3\x7a\xaf\x1a\x79\x11\xf8\xdf\x0e\xf9\xfa\x2e\x62\x28\xad\x16\x39\xce\xc0\x32\x51\x0d\x3d\x8e\xba\x28\x92\x7a\xec\xbd\xc1\x87\xdf\x3f\x97\xdf\x69\x93\x51\x44\x30\xcf\x03\xe7\xbb\x16\x01\xe5\x98\xea\xff\x79\xe4\xc1\x2d\xa7\x57\x3b\xbf\x72\x1e\xd1\xf0\xed\x73\xf7\x89\x15\xd4\x4f\xad\xbe\x6d\x89\xc4\x52\x3d\xf4\xc7\xba\x3f\x80\x9d\xf7\xc1\xde\x38\x8b\xd3\x38\x57\x5e\x43\x5f\x91\x4d\xdf\x03\x91\x7d\x8a\x5a\x46\xf1\xad\xfd\xf8\x54\x49\xa6\x47\x5b\x4c\x63\x93\xa4\xd6\x0b\x5f\x1b\x33\xa3\x96\x8d\xee\xf9\xd1\xa6\xa3\x6b\x40\x47\x77\x54\x3a\xba\x83\xe8\x38\x37\xa0\xe3\x7c\x54\x3a\xce\x1f\xb1\x5f\x97\x0f\x6b\x82\x08\x45\xff\xae\x43\x32\x98\xd6\xb6\xd3\x1a\x74\x28\x9f\xbf\x34\x7c\xae\xa7\x63\x68\xc2\x54\x69\x34\x75\x9e\xd6\xa0\x63\x91\xd2\x58\xef\xa9\xd2\x68\x6a\x3e\xad\x41\xc7\x22\xe5\xa3\xb7\x9f\xef\x3e\x1f\xd4\xda\x78\x9a\x11\x87\xf2\xf7\xda\xdf\x01\x66\x5b\x4e\x89\x43\x53\xbf\x69\x46\x1c\x85\x82\xc6\x3a\x4d\x89\x43\x53\x9b\x69\x46\x1c\x85\x82\x0f\xdf\x60\xe2\x9b\x90\x40\xf8\xc6\xef\xbd\xc7\x89\x31\x7a\x37\x38\x8d\x88\x43\x79\x4c\x71\xd3\x5c\x49\xca\xd9\x98\x5a\x91\x0b\x0e\x5d\x3b\x9b\x46\xc4\x91\x68\x68\x6e\x57\x53\x70\xe8\xda\xd2\x34\x22\x8e\x44\xc3\x87\xef\x36\xba\xbe\xf2\x32\xb0\xbd\x69\x8b\x3c\x90\xed\xdd\x5f\x73\x0d\xb2\xd9\xa9\x70\x69\xdd\xf4\xb4\x45\x1e\x91\xae\x86\xb7\x40\x15\x2e\xad\x5b\xa1\xb6\xc8\x23\xd2\x75\x5c\xad\x6a\x31\xbd\xf2\x94\xe6\x62\x9a\x3c\x3a\xdb\xf9\xb8\x6e\xe9\x6d\xf6\x32\xfb\x55\xe4\xfd\x7f\x00\x00\x00\xff\xff\xa4\x67\x7e\x9b\xf0\x2c\x00\x00"

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
