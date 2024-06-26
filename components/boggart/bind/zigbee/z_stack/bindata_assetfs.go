// Code generated by go-bindata. DO NOT EDIT.
// sources:
// templates/views/widget.html
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

var _templatesViewsWidgetHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x5a\x7b\x8f\xdb\xb8\x11\xff\x7f\x3f\xc5\x54\xcd\x41\x32\x50\x4b\xde\x6d\xaf\x2d\x7c\xf6\x1e\xd2\x64\x8b\x2c\xb0\xc9\x2d\x92\xed\xa1\xe8\xe1\xb0\xa0\xc5\x91\x4d\x9f\x4c\xea\x48\xca\x8f\x1a\xfe\xee\x05\xf5\xb0\x65\xaf\x24\x5b\xb6\xd7\x49\xf3\xc7\x86\x22\xe7\xf1\x9b\x07\x87\x23\x5a\xcb\x25\x50\x0c\x18\x47\xb0\x7c\xc1\x35\x72\x6d\xc1\x6a\x75\xd5\xa3\x6c\x0a\x7e\x48\x94\xea\x5b\x52\xcc\xac\xdb\x2b\x00\x80\xe2\xac\x2f\xc2\xf6\x84\xb6\xaf\x6f\xc0\x8c\xd4\x24\x1f\xcd\x55\xfb\xfa\x26\xa3\xdf\xe5\x99\x3f\x47\x84\x63\x58\x58\x7d\x49\xa1\x99\x0e\x71\x87\x22\xa1\x1a\xdd\xdc\x2e\x97\xc0\xae\xff\xce\xc1\xfa\x0f\x1b\xfe\x03\x11\x28\x4e\x99\x8f\x16\xb8\xb0\x5a\xf5\xbc\xd1\x4d\x09\x57\x11\x71\x88\x44\x06\x6c\x6e\xdd\xf6\x3c\xca\xa6\x3b\x20\x4a\xa6\xb6\x70\xe5\xde\x29\xd1\x11\x08\x39\x01\x29\x42\xec\x5b\x66\x68\x01\xf1\x35\x13\xbc\x6f\xfd\x98\x0d\x14\x6a\xcd\xf8\x50\x59\x30\x41\x3d\x12\xb4\x6f\x45\x42\x69\x0b\x18\xed\x5b\xeb\xb5\x97\x82\x13\xe1\x9a\x0c\x42\xcc\x61\xa4\x0f\xc9\xdf\xb6\xd2\x92\x45\x48\x81\xea\xb6\x44\x15\x09\xae\xd8\x14\x81\x8b\x99\x24\x91\x05\x4a\x2f\x0c\x9e\x19\xa3\x7a\xd4\xbd\xee\x74\xbe\xab\x90\x9f\xea\x18\x21\xa1\x75\xeb\xb2\x7a\x31\x13\x90\x23\x9c\xd0\xb6\xc9\x82\x1b\x6b\x13\xac\x9f\x89\x64\x06\x71\x1e\x27\x3d\x3a\x4d\x5a\x18\x1f\x24\xaa\xe7\xd5\xc1\x36\xbc\x7b\x8c\x1e\x08\xba\x38\xc9\x29\x74\x83\xfa\x11\xe5\x84\x69\x18\x0b\xc6\xd7\xd8\x6b\x94\xe7\xfc\xb5\x04\x09\x11\xe3\x51\xac\x41\x2f\x22\xec\x5b\xfe\x08\xfd\xdf\x06\x62\x6e\xe5\xee\x1b\xab\xb6\x9a\x31\xed\x8f\x2c\xe0\x64\x82\x7d\x2b\x4a\x60\xb4\x53\x18\x26\xfb\xb6\x26\xa6\xc6\xb5\x7d\x4b\x70\x0b\x0c\xf0\x00\xf0\x77\x70\x22\xc9\xb8\x06\x37\x25\x7c\x36\x84\x2d\xb0\xb4\x34\x31\x58\xad\x20\x51\x89\x74\xb9\x04\xe4\xd4\x4c\x78\x7b\x8c\xaa\x35\x7b\x4f\xc8\x1a\x79\xfc\xe1\xee\x3d\x20\x37\x79\x47\xcf\xea\xf1\xd4\x33\x6e\x88\xf4\x59\xc5\x51\x24\xa4\x36\xf5\x72\x1f\x17\x34\x8f\x55\x02\xdc\xc4\x28\x19\xd4\xc5\xc6\x60\xc9\x4c\x3d\x25\x36\x99\x75\x18\x2a\x3c\xd8\xa4\xcc\xe3\xc1\x16\x08\x58\xad\xb2\xd1\x46\x1c\x65\x6a\x3d\x93\xa0\xd9\x1f\x8e\x1c\x4f\x42\xfe\x6a\x59\x95\xc1\x9f\xa2\x54\x4c\xf0\x3a\x4d\xcd\xf2\xef\x49\x12\xae\x92\xf4\x90\x38\x65\x46\x76\x93\x34\x5c\x2e\xd7\x90\xdc\xb5\xa4\xcf\x99\xa0\xbd\x52\xce\xb9\x8f\x1e\xa5\xa0\xb1\xaf\x8f\x05\x9f\xb1\x9b\xf4\x73\xb6\x8c\x5a\x44\x26\x2d\x5a\x97\x33\xe4\xe7\x54\x73\x13\x43\xa6\x2e\x14\x31\x7f\x24\x63\x21\x3f\x63\x88\x24\x49\x69\x77\x6b\x8d\xf1\xe2\xda\xe5\xcc\xfa\x48\x18\x07\x2d\x31\x3c\x36\x42\x46\xc0\x93\xc4\xf0\xa2\xa8\x3f\x10\x49\x67\x44\xe2\xc9\x9b\x23\x17\x74\xa6\xbd\xb1\xbf\xde\x64\xf5\x82\xf1\x40\x9c\xaf\x58\xbc\x4f\x1a\xd9\xe4\x68\x68\xea\x08\x03\xc4\x4d\x1b\xe1\x67\x9d\x6e\xaa\xcb\x85\x31\xc3\xad\x34\xd1\x27\x01\x4f\x04\x7c\x0d\xe4\xf7\x77\x77\x77\x40\x28\x95\xa8\xd4\x51\x06\x30\x44\x7c\xce\x04\x5c\x2a\xf7\x38\xea\x99\x90\xbf\x9d\x2f\xfd\x3e\x65\x02\x1f\xdf\x7e\x82\xfb\xf7\x4d\xfd\x90\xc1\x71\x23\xc2\x9f\xd9\xfe\xa3\xfd\x9c\x51\xcc\x81\xe3\x5c\x23\xa7\x48\x4f\xb4\x20\x17\xf3\xfc\x15\x4d\xf1\x47\x84\xf3\xe6\xd5\x3c\x37\x21\x63\xbf\x54\x2a\x26\x3d\xa6\x3a\x5f\x26\xde\x73\xa6\x19\x09\xd9\x7f\x89\x79\x71\x4e\x5b\xd8\x86\xae\xd8\xc0\x72\x19\x67\xa6\xf3\x98\x92\x90\x15\x5a\x51\x2e\x34\xac\xa7\x0e\xe8\x45\x4f\x75\x55\xcf\xab\x79\x95\xec\x79\xc9\xdb\x7c\xc9\xad\x82\x17\x08\x39\xa9\xbd\xa4\x28\x3c\x66\xc3\xfc\xbf\x6f\xfb\xf6\x26\xad\xbd\x0a\x08\xa7\xc6\x77\x91\x60\x5c\xab\xaf\x7e\x8b\xb3\xff\xa2\x85\x68\x92\xce\x9f\x76\xe5\x52\x77\xdd\x52\xbb\x5b\x7a\x7a\xf4\xb2\x5e\xec\x1e\x5e\x75\x37\x22\x45\xfe\xd2\x93\xef\x50\xe6\x77\x24\x22\x03\x16\x32\xcd\xb0\x39\x73\x59\xb3\x73\x28\xef\x47\xc2\xe3\x80\xf8\x3a\x96\x28\xc1\x17\xb4\xb9\x84\x07\xa2\x34\x28\x44\xbe\x8f\xb3\x7a\xdb\xd7\xde\x1d\xd5\xdd\x1b\x2d\x97\x20\x09\x1f\x22\xbc\x49\xfb\x1e\xe8\xf6\x21\x6b\x81\x2a\x6b\xe8\x9e\x8c\x48\x4a\x5e\x26\xce\xcd\x72\xe2\xed\xa1\xed\xc8\x36\xb7\xc9\x88\x8c\xf5\xad\xfa\xa2\x25\xe3\xc3\xa6\x22\x8a\x79\xd1\x94\x37\x4d\x0b\xf3\x5e\xd8\x54\x7b\x72\x0d\x12\x80\xd5\x99\x7f\xf7\x6f\x6b\x2d\xaf\x98\x2a\xef\x04\x3d\xa0\xbd\xd4\xf4\xb6\xa7\x7c\xc9\xa2\xfc\x92\x86\x44\x51\xc8\xfc\xe4\x28\xf2\xc6\x64\x4a\xd2\x45\xeb\x96\x0a\x3f\x9e\x20\xd7\xee\x4c\x32\x8d\x0e\x25\x1a\x9f\x44\x8a\xd9\xb1\x0b\x26\x99\x5c\xfb\x82\xc8\xdd\x7f\x0a\x39\x21\x1a\xac\x9b\x4e\xe7\xaf\xed\xce\x75\xbb\x73\xf3\x74\xfd\x7d\xb7\xf3\x97\x6e\xe7\xfb\x76\xe7\x6f\xdd\x4e\xc7\x82\xd5\xca\x6e\xb5\x7a\x5e\xaa\xe2\xb6\x1a\x69\x75\x5a\xd6\x9f\x44\x95\xa7\x50\xf5\x09\xf4\x7f\x56\x11\xef\xb2\xa3\xa4\x71\x49\x78\x94\x22\x60\x21\x16\x1b\xc7\x83\xab\x28\x07\x3f\x8c\x95\x46\xd9\xbc\x0e\xfe\x14\xeb\xe3\x99\x2f\x59\xfe\xbf\x91\x4a\xb8\xcd\x97\x37\x0e\x86\x33\xdf\x6f\x79\x06\x9c\xa5\x21\x5d\x6b\x70\xef\xdf\x37\xe9\x3f\x37\x7c\x59\x5e\x1d\xcb\x7e\xcf\xdf\xa5\xd9\xf1\xc0\x94\x3e\x4e\xc4\x4f\xb1\x3e\x56\xc6\x71\x47\x4a\x89\x84\xa3\x8e\x95\x13\x5a\xee\x73\x96\xc1\x43\x1b\xee\x8d\xce\xab\xc2\x0f\xaa\x66\x6b\x58\x39\x90\xe5\x32\xb9\x26\x61\xfe\x87\xa7\x8f\x0f\xe0\xa4\xe3\x7f\x7d\x7e\x00\xcb\xa3\x44\x8d\x06\x82\x48\xea\x11\xa5\x50\x2b\x6f\x8a\x9c\x0a\xa9\xbc\x75\x7d\x55\xe6\xfd\xae\x3d\x50\x9e\xaf\xd2\xd9\xa7\x74\x76\x20\x84\x56\x5a\x92\xc8\x9d\x30\xee\xfa\x66\x2b\x07\x24\x54\xd8\x3a\xa3\xd6\x4d\x5d\xcf\x01\x6c\x66\xea\x01\x94\x7b\x65\xac\xce\xe8\x13\x6f\xac\xbc\xf1\xef\x31\xca\x85\x5b\x70\x8b\xc1\x32\x7e\x0d\x5f\x0c\x94\x51\x58\x19\x80\x57\xd1\xb9\xf1\xf6\x8e\xee\x42\x18\x2e\xa0\x3c\xb3\xbd\x32\xf6\xdb\xea\xd3\x3d\x72\x50\x33\xb5\xde\x5d\x6f\x9c\xbc\xaf\x6a\xb9\x12\x09\x5d\x38\x41\xcc\x93\x1f\xd0\xc1\x69\xc1\xf2\x6a\x67\x93\xb3\x20\x79\x81\x74\x51\x4a\x21\x15\x38\x43\x04\x27\x44\x9e\x4f\xb4\xa0\xd3\x2a\x2b\x02\x85\x13\xc4\xd0\x25\x07\x4f\x26\xa2\xa2\x64\x70\x9c\xc1\xe3\x27\xa1\x59\xb0\x70\x96\x95\x45\x29\x79\xf3\xed\x82\x7d\x67\x64\xd9\x7f\xaa\xa6\xc3\xb9\xee\x42\xd2\x25\xa6\x00\x56\xab\x3a\xea\x45\x64\x84\xe2\x1e\xa1\x23\x46\xb1\x9b\xba\xbf\x9a\xc8\x74\x63\x8c\x0f\xbb\x60\xaf\x23\xf7\x67\xbb\x94\x7a\xd5\xfa\xa1\xcc\x6f\x25\x65\xb5\xb0\xc5\x8b\xd3\xeb\xc0\xa9\xd8\xf7\x51\xa9\x0f\x84\xd3\x10\xa5\x23\x77\xe3\x68\xfe\xb1\x00\x1c\x69\xb2\x39\x0e\x35\xf4\xfb\x7d\xb0\x03\xc2\x42\xa4\x76\x19\x31\xbc\x52\x40\xa4\x3b\x41\xa5\xc8\xb0\xc6\x7f\xdf\x4e\x2c\x56\xe9\x55\x56\xea\xb8\x0c\x37\xfc\xc1\x78\x2e\xe6\x69\x99\x3d\x97\xf3\xbe\xa4\xf1\x3b\xa3\xfb\xd4\x5e\x89\x97\x70\xe0\xd5\xf6\xd3\xd6\xe3\x1b\xc7\xfe\x63\xfe\x55\x0e\x24\xbf\xd9\xff\xc2\xe8\xaf\x76\x2b\xb9\x5b\x1d\xe2\xba\x2a\xbd\x28\x4a\xe6\xdf\x94\xc8\x4a\xcc\x01\xf4\xb7\x84\xdb\xad\x6a\xfb\x30\xa1\xd5\x23\xa6\x6a\x88\x4c\xb1\x86\x3e\x2c\x5f\xd6\xad\x1f\xae\x4a\xb7\x19\xba\x91\x14\x91\x63\x9b\x58\xd8\xad\x74\xb3\xe5\x5f\x22\x54\x66\x8c\xd1\xf2\x4b\xce\xc9\xa8\xdd\xfa\x15\xfa\x90\x3f\x67\x1f\x19\xd8\x2d\xf8\x11\x6c\xc1\x6d\xe8\x82\x2d\x82\xc0\xae\x4c\xda\x46\x3a\xa6\x24\x74\x4a\xc3\xf7\x62\xea\x8d\x4b\xc6\x64\x5e\x91\xd1\x69\xea\x05\x99\xf8\xf4\x1b\xac\x2a\xe7\xc7\x32\xdc\x90\xa6\xdf\x6f\x55\x91\x1a\xd0\xdd\xe4\x6f\xf9\x7a\x96\xea\xdd\x9d\x2a\xf8\xd2\x9c\x1d\x13\x8b\xcf\xf9\x78\x7d\x33\xb0\xe9\xaa\xfe\x17\x00\x00\xff\xff\x9c\x8d\xb9\xd9\xbb\x27\x00\x00"

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
