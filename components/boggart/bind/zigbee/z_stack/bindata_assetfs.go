// Code generated by go-bindata.
// sources:
// templates/views/widget.html
// DO NOT EDIT!

package zstack

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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x5a\x7b\x8f\xdb\xb8\x11\xff\x7f\x3f\xc5\x54\xcd\x41\x32\x50\x4b\xde\x6d\xaf\x2d\x7c\xf6\x1e\xd2\x64\x8b\x2c\xb0\xc9\x2d\x92\xed\xa1\xb8\xc3\x61\x41\x8b\x23\x9b\x3e\x99\x54\x48\xca\x8f\x1a\xfe\xee\x05\xf5\xb0\x65\xaf\x24\x5b\xb6\xe3\xe4\xf2\xc7\x86\x22\xe7\xf1\x9b\x07\x87\x23\x5a\xcb\x25\x50\x0c\x18\x47\xb0\x7c\xc1\x35\x72\x6d\xc1\x6a\x75\xd5\xa3\x6c\x0a\x7e\x48\x94\xea\x5b\x52\xcc\xac\xdb\x2b\x00\x80\xe2\xac\x2f\xc2\xf6\x84\xb6\xaf\x6f\xc0\x8c\xd4\x24\x1f\xcd\x55\xfb\xfa\x26\xa3\xdf\xe5\x99\x3f\x47\x84\x63\x58\x58\x7d\x49\xa1\x99\x0e\x71\x87\x22\xa1\x1a\xdd\xdc\x2e\x97\xc0\xae\xff\xc9\xc1\xfa\x85\x0d\xff\x85\x08\x14\xa7\xcc\x47\x0b\x5c\x58\xad\x7a\xde\xe8\xa6\x84\xab\x88\x38\x44\x22\x03\x36\xb7\x6e\x7b\x1e\x65\xd3\x1d\x10\x25\x53\x5b\xb8\x72\xef\x94\xe8\x08\x84\x9c\x80\x14\x21\xf6\x2d\x33\xb4\x80\xf8\x9a\x09\xde\xb7\x7e\xcc\x06\x0a\xb5\x66\x7c\xa8\x2c\x98\xa0\x1e\x09\xda\xb7\x22\xa1\xb4\x05\x8c\xf6\xad\xf5\xda\x4b\xc1\x89\x70\x4d\x06\x21\xe6\x30\xd2\x87\xe4\x6f\x5b\x69\xc9\x22\xa4\x40\x75\x5b\xa2\x8a\x04\x57\x6c\x8a\xc0\xc5\x4c\x92\xc8\x02\xa5\x17\x06\xcf\x8c\x51\x3d\xea\x5e\x77\x3a\xdf\x55\xc8\x4f\x75\x8c\x90\xd0\xba\x75\x59\xbd\x98\x09\xc8\x11\x4e\x68\xdb\x64\xc1\x8d\xb5\x09\xd6\xcf\x44\x32\x83\x38\x8f\x93\x1e\x9d\x26\x2d\x8c\x0f\x12\xd5\xf3\xea\x60\x1b\xde\x3d\x46\x0f\x04\x5d\x9c\xe4\x14\xba\x41\xfd\x88\x72\xc2\x34\x8c\x05\xe3\x6b\xec\x35\xca\x73\xfe\x5a\x82\x84\x88\xf1\x28\xd6\xa0\x17\x11\xf6\x2d\x7f\x84\xfe\xef\x03\x31\xb7\x72\xf7\x8d\x55\x5b\xcd\x98\xf6\x47\x16\x70\x32\xc1\xbe\x15\x25\x30\xda\x29\x0c\x93\x7d\x5b\x13\x53\xe3\xda\xbe\x25\xb8\x05\x06\x78\x00\xf8\x19\x9c\x48\x32\xae\xc1\x4d\x09\x9f\x0d\x61\x0b\x2c\x2d\x4d\x0c\x56\x2b\x48\x54\x22\x5d\x2e\x01\x39\x35\x13\xde\x1e\xa3\x6a\xcd\xde\x13\xb2\x46\x1e\x7f\xb8\x7b\x0b\xc8\x4d\xde\xd1\xb3\x7a\x3c\xf5\x8c\x1b\x22\x7d\x56\x71\x14\x09\xa9\x4d\xbd\xdc\xc7\x05\xcd\x63\x95\x00\x37\x31\x4a\x06\x75\xb1\x31\x58\x32\x53\x4f\x89\x4d\x66\x1d\x86\x0a\x0f\x36\x29\xf3\x78\xb0\x05\x02\x56\xab\x6c\xb4\x11\x47\x99\x5a\xcf\x24\x68\xf6\x87\x23\xc7\x93\x90\x7f\xb1\xac\xca\xe0\x4f\x51\x2a\x26\x78\x9d\xa6\x66\xf9\xf7\x24\x09\x57\x49\x7a\x48\x9c\x32\x23\xbb\x49\x1a\x2e\x97\x6b\x48\xee\x5a\xd2\xc7\x4c\xd0\x5e\x29\xe7\xdc\x47\x8f\x52\xd0\xd8\xd7\xc7\x82\xcf\xd8\x4d\xfa\x39\x5b\x46\x2d\x22\x93\x16\xad\xcb\x19\xf2\x73\xaa\xb9\x89\x21\x53\x17\x8a\x98\xdf\x93\xb1\x90\x1f\x31\x44\x92\xa4\xb4\xbb\xb5\xc6\x78\x71\xed\x72\x66\xbd\x27\x8c\x83\x96\x18\x1e\x1b\x21\x23\xe0\x49\x62\x78\x51\xd4\xef\x88\xa4\x33\x22\xf1\xe4\xcd\x91\x0b\x3a\xd3\xde\xd8\x5f\x6f\xb2\x7a\xc1\x78\x20\xce\x57\x2c\xde\x26\x8d\x6c\x72\x34\x34\x75\x84\x01\xe2\xa6\x8d\xf0\xb3\x4e\x37\xd5\xe5\xc2\x98\xe1\x56\x9a\xe8\x93\x80\x27\x02\xbe\x06\xf2\xfb\xbb\xbb\x3b\x20\x94\x4a\x54\xea\x28\x03\x18\x22\x3e\x67\x02\x2e\x95\x7b\x1c\xf5\x4c\xc8\xdf\xcf\x97\x7e\x1f\x32\x81\x8f\xaf\x3f\xc0\xfd\xdb\xa6\x7e\xc8\xe0\xb8\x11\xe1\xcf\x6c\xff\xd1\x7e\xce\x28\xe6\xc0\x71\xae\x91\x53\xa4\x27\x5a\x90\x8b\x79\xfe\x8a\xa6\xf8\x23\xc2\x79\xf3\x6a\x9e\x9b\x90\xb1\x5f\x2a\x15\x93\x1e\x53\x9d\x2f\x13\xef\x39\xd3\x8c\x84\xec\x7f\xc4\xbc\x38\xa7\x2d\x6c\x43\x57\x6c\x60\xb9\x8c\x33\xd3\x79\x4c\x49\xc8\x0a\xad\x28\x17\x1a\xd6\x53\x07\xf4\xa2\xa7\xba\xaa\xe7\xd5\xbc\x4a\xf6\xbc\xe4\x6d\xbe\xe4\x56\xc1\x0b\x84\x9c\xd4\x5e\x52\x14\x1e\xb3\x61\xfe\xdf\xb7\x7d\x7b\x93\xd6\x5e\x05\x84\x53\xe3\xbb\x48\x30\xae\xd5\x57\xbf\xc5\xd9\x7f\xd1\x42\x34\x49\xe7\x4f\xbb\x72\xa9\xbb\x6e\xa9\xdd\x2d\x3d\x3d\x7a\x59\x2f\x76\x0f\xaf\xba\x1b\x91\x22\x7f\xe9\xc9\x77\x28\xf3\x1b\x12\x91\x01\x0b\x99\x66\xd8\x9c\xb9\xac\xd9\x39\x94\xf7\x3d\xe1\x71\x40\x7c\x1d\x4b\x94\xe0\x0b\xda\x5c\xc2\x03\x51\x1a\x14\x22\xdf\xc7\x59\xbd\xed\x6b\xef\x8e\xea\xee\x8d\x96\x4b\x90\x84\x0f\x11\x5e\xa5\x7d\x0f\x74\xfb\x90\xb5\x40\x95\x35\x74\x4f\x46\x24\x25\x2f\x13\xe7\x66\x39\xf1\xfa\xd0\x76\x64\x9b\xdb\x64\x44\xc6\xfa\x5a\x7d\xd2\x92\xf1\x61\x53\x11\xc5\xbc\x68\xca\x9b\xa6\x85\x79\x2f\x6c\xaa\x3d\xb9\x06\x09\xc0\xea\xcc\xbf\xfb\xaf\xb5\x96\x57\x4c\x95\x37\x82\x1e\xd0\x5e\x6a\x7a\xdb\x53\xbe\x64\x51\x7e\x49\x43\xa2\x28\x64\x7e\x72\x14\x79\x63\x32\x25\xe9\xa2\x75\x4b\x85\x1f\x4f\x90\x6b\x77\x26\x99\x46\x87\x12\x8d\x4f\x22\xc5\xec\xd8\x05\x93\x4c\xae\x7d\x42\xe4\xee\xbf\x85\x9c\x10\x0d\xd6\x4d\xa7\xf3\xf7\x76\xe7\xba\xdd\xb9\x79\xba\xfe\xbe\xdb\xf9\x5b\xb7\xf3\xfd\x2f\x9d\x7f\x74\x3b\x1d\x0b\x56\x2b\xbb\xd5\xea\x79\xa9\x8a\xdb\x6a\xa4\xd5\x69\x59\x7f\x12\x55\x9e\x42\xd5\x27\xd0\x1f\xac\x22\xde\x65\x47\x49\xe3\x92\xf0\x28\x45\xc0\x42\x2c\x36\x8e\x07\x57\x51\x0e\x7e\x18\x2b\x8d\xb2\x79\x1d\xfc\x29\xd6\xc7\x33\x5f\xb2\xfc\x7f\x23\x95\x70\x9b\x2f\x6f\x1c\x0c\x67\xbe\xdf\xf2\x0c\x38\x4b\x43\xba\xd6\xe0\xde\xbf\x6d\xd2\x7f\x6e\xf8\xb2\xbc\x3a\x96\xfd\x9e\xbf\x49\xb3\xe3\x81\x29\x7d\x9c\x88\x9f\x62\x7d\xac\x8c\xe3\x8e\x94\x12\x09\x47\x1d\x2b\x27\xb4\xdc\xe7\x2c\x83\x87\x36\xdc\x1b\x9d\x57\x85\x1f\x54\xcd\xd6\xb0\x72\x20\xcb\x65\x72\x4d\xc2\xfc\x77\x4f\xef\x1f\xc0\x49\xc7\xff\xf9\xf8\x00\x96\x47\x89\x1a\x0d\x04\x91\xd4\x23\x4a\xa1\x56\xde\x14\x39\x15\x52\x79\xeb\xfa\xaa\xcc\xfb\x5d\x7b\xa0\x3c\x5f\xa5\xb3\x4f\xe9\xec\x40\x08\xad\xb4\x24\x91\x3b\x61\xdc\xf5\xcd\x56\x0e\x48\xa8\xb0\x75\x46\xad\x9b\xba\x9e\x03\xd8\xcc\xd4\x03\x28\xf7\xca\x58\x9d\xd1\x27\xde\x58\x79\xe3\xcf\x31\xca\x85\x5b\x70\x8b\xc1\x32\xfe\x12\xbe\x18\x28\xa3\xb0\x32\x00\x5f\x44\xe7\xc6\xdb\x3b\xba\x0b\x61\xb8\x80\xf2\xcc\xf6\xca\xd8\x6f\xab\x4f\xf7\xc8\x41\xcd\xd4\x7a\x77\xbd\x72\xf2\xbe\xaa\xe5\x4a\x24\x74\xe1\x04\x31\x4f\x7e\x40\x07\xa7\x05\xcb\xab\x9d\x4d\xce\x82\xe4\x05\xd2\x45\x29\x85\x54\xe0\x0c\x11\x9c\x10\x79\x3e\xd1\x82\x4e\xab\xac\x08\x14\x4e\x10\x43\x97\x1c\x3c\x99\x88\x8a\x92\xc1\x71\x06\x8f\x1f\x84\x66\xc1\xc2\x59\x56\x16\xa5\xe4\xcd\xb7\x0b\xf6\x9d\x91\x65\xff\xa5\x9a\x0e\xe7\xba\x0b\x49\x97\x98\x02\x58\xad\xea\xa8\x17\x91\x11\x8a\x7b\x84\x8e\x18\xc5\x6e\xea\xfe\x6a\x22\xd3\x8d\x31\x3e\xec\x82\xbd\x8e\xdc\x5f\xed\x52\xea\x55\xeb\x87\x32\xbf\x95\x94\xd5\xc2\x16\x2f\x4e\xaf\x03\xa7\x62\xdf\x47\xa5\xde\x11\x4e\x43\x94\x8e\xdc\x8d\xa3\xf9\xc7\x02\x70\xa4\xc9\xe6\x38\xd4\xd0\xef\xf7\xc1\x0e\x08\x0b\x91\xda\x65\xc4\xf0\x85\x02\x22\xdd\x09\x2a\x45\x86\x35\xfe\xfb\x76\x62\xb1\x4a\xaf\xb2\x52\xc7\x65\xb8\xe1\x4f\xc6\x73\x31\x4f\xcb\xec\xb9\x9c\xf7\x29\x8d\xdf\x19\xdd\xa7\xf6\x4a\xbc\x84\x03\xaf\xb6\x9f\xb6\x1e\x5f\x39\xf6\x9f\xf3\xaf\x72\x20\xf9\xcd\xfe\x57\x46\x7f\xb3\x5b\xc9\xdd\xea\x10\xd7\x55\xe9\x45\x51\x32\xff\xa6\x44\x56\x62\x0e\xa0\xbf\x25\xdc\x6e\x55\xdb\x87\x09\xad\x1e\x31\x55\x43\x64\x8a\x35\xf4\x61\xf9\xb2\x6e\xfd\x70\x55\xba\xcd\xd0\x8d\xa4\x88\x1c\xdb\xc4\xc2\x6e\xa5\x9b\x2d\xff\x12\xa1\x32\x63\x8c\x96\x5f\x73\x4e\x46\xed\xd6\x6f\xd0\x87\xfc\x39\xfb\xc8\xc0\x6e\xc1\x8f\x60\x0b\x6e\x43\x17\x6c\x11\x04\x76\x65\xd2\x36\xd2\x31\x25\xa1\x53\x1a\xbe\x17\x53\xaf\x5c\x32\x26\xf3\x8a\x8c\x4e\x53\x2f\xc8\xc4\xa7\xdf\x60\x55\x39\x3f\x96\xe1\x86\x34\xfd\x7e\xab\x8a\xd4\x80\xee\x26\x7f\xcb\xd7\xb3\x54\xef\xee\x54\xc1\x97\xe6\xec\x98\x58\x7c\xce\xc7\xeb\x9b\x81\x4d\x57\xf5\xff\x00\x00\x00\xff\xff\xe7\x0c\xdc\x89\xbb\x27\x00\x00"

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
