// Code generated by go-bindata.
// sources:
// templates/views/message.html
// DO NOT EDIT!

package internal

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
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

var _templatesViewsMessageHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x58\xdf\x6f\xdb\x36\x10\x7e\xcf\x5f\x71\xe0\xc3\x90\x0e\x90\x15\x3b\x43\x31\x60\x92\x81\x3d\x0c\x18\x86\x0e\x03\xd6\x6e\x0f\x7b\x09\xce\xe6\x29\x66\xca\x1f\x2a\x49\xb9\x6e\x0c\xff\xef\x03\x45\x2a\x56\x13\x59\xb5\x5b\x07\x68\x5e\x12\x92\x77\x3c\x7e\xfc\xee\xf3\x91\xd4\x76\x0b\x9c\x2a\xa1\x09\xd8\xd2\x68\x4f\xda\x33\xd8\xed\x2e\x0a\x2e\xd6\xb0\x94\xe8\x5c\xc9\xac\xf9\xc8\xe6\x17\x00\x00\xfd\xd1\xa5\x91\x99\xe2\xd9\x74\x06\xa1\xe5\x54\xd7\xda\xb8\x6c\x3a\x4b\xfe\x8f\xe7\x6c\x6e\x6a\xd4\x24\x7b\xd6\xa7\x1e\x5e\x78\x49\x8f\x3c\x5a\xaf\xd5\x6c\xbe\xdd\x82\x98\xfe\xac\x81\xbd\x25\xcd\x41\x91\x73\x78\x4b\x0c\x26\xb0\xdb\x15\xf9\x6a\x36\x30\xa9\x0f\x58\x12\xda\x4a\x6c\xd8\xbc\xc8\xb9\x58\x3f\xc2\x30\x30\xf4\x19\xac\x8e\x9c\xa7\x6b\x04\x50\x15\x4c\xc8\x5a\x63\x03\x75\x63\x18\x50\x92\xf5\xd0\xfe\xcd\x38\xea\x5b\xb2\x5d\x47\x38\x25\x9c\xc3\xc5\xe0\xde\xdb\x30\x8b\xc6\x7b\xa3\xc1\x7f\xaa\xa9\x64\xb1\xc3\xf6\x7b\x33\x8e\x18\x70\xf4\xd8\x85\x4a\x8b\x31\x40\x2b\x30\x5b\x09\xce\x49\x97\xcc\xdb\x86\xd8\xfc\x07\x2f\x14\xb9\x5f\x8a\x3c\x86\x19\x5e\x70\xbb\x1d\xdb\xd4\x53\xbe\xd2\x9c\x90\x99\xdd\xee\xe2\x10\x4d\x29\x69\xa7\x10\x25\x74\x65\xbe\x6f\x9a\xc6\x36\x75\x32\x51\x45\x65\xac\xea\x10\x87\x76\xb6\x32\x56\xdc\x1b\xed\x51\x42\xdb\x97\xb8\x20\x99\x49\xaa\x3c\x03\x6b\x24\x45\x37\x06\x8a\xfc\xca\xf0\x92\xd5\xc6\x79\x06\x82\x97\xcc\x91\xe6\x99\x42\x21\x19\xe0\xd2\x0b\xa3\x4b\x96\xaf\x8d\x58\x52\x9e\x20\xe7\x0c\xb4\x59\xa3\x14\x1c\x3d\x1d\x60\xb4\x97\x16\xe1\x49\x45\x0c\xb7\xd6\x34\xf5\x81\x1c\xb4\xb3\x5a\x90\xc1\xb7\x64\x6b\x23\x1b\x45\xfb\x2c\x18\xed\xad\x91\x71\x1b\x90\xea\xc8\x75\x57\x46\xae\x07\xab\xc8\x01\xe6\x63\x35\xf8\x37\xc5\x0f\x75\x00\x0a\x57\xa3\x7e\x28\x5d\xf4\xa1\x11\x96\x38\x9b\xff\x58\xe4\xc1\x30\x02\x38\x6f\xf1\x8c\x38\x3c\x2d\x7e\xaf\x3b\xd0\xaf\x8f\x06\x5d\x08\x5d\x37\x3e\xc9\xd3\x86\x12\xc0\x40\x09\x5d\xb2\x29\x03\x85\x9b\x92\x4d\xaf\xae\x18\x38\x4f\x75\x3b\xd4\x97\x41\xe2\x2d\x66\xb6\xe3\x54\xa3\xa2\x7d\x6f\x8d\xb2\xa1\x92\x05\x4d\xc6\x21\xd8\xed\x18\x74\x2c\xf4\xf9\x18\xe1\x61\x50\xb0\x5f\x32\x7d\xb3\x48\x5c\x4d\xc4\x9f\x4f\x23\x6f\x63\xf8\x17\x2d\x91\xab\x24\x91\xeb\x4e\x20\x57\x93\x31\x89\x24\x46\xa3\x42\x52\xa7\x27\x90\x76\xe4\x45\xe9\x03\xdf\x93\x3d\xab\x42\xfa\xea\x88\xc1\x5f\x84\x3e\x1c\x49\x5a\xfa\xa1\xc4\x43\x34\xcd\x7a\x59\x8f\xfb\x4a\x7a\x88\x9d\x93\xf2\xfd\xb0\xaa\xa9\xc3\x01\xd2\x29\xe8\x0e\x35\xb1\x78\xa6\xd3\x87\x56\x4c\x21\x36\xc4\xf1\x40\x62\x44\x42\xed\x21\x14\x5b\xec\xe1\xd0\xdb\x5f\xe4\xfe\x40\x4d\x70\x59\x91\x42\x49\xaf\xba\x9b\x5c\x5c\xe9\x64\x48\xe6\xbd\x43\x8d\x43\xa0\x92\xe5\x78\x58\x7f\xb5\x13\xce\x05\x0c\xe5\x27\xe7\x86\x70\x45\xc3\xf1\xb0\x7e\x0d\xfe\x67\xa3\x4b\xe1\xfd\x6a\x90\xad\xd6\x70\x02\x59\xc1\xff\x5c\xa8\xee\x71\x85\x76\x08\x55\x34\x1c\x8f\xea\xbf\xe0\x0f\x97\xe7\xc0\x44\x56\x09\x39\x84\x29\x1a\x8e\xc7\xf4\x5b\xf0\xff\x0a\x4c\x45\x1e\xa3\x7e\x7f\x55\xd9\xd3\xc6\x3f\xdf\xa1\xfd\xae\x8d\xfe\x22\x6a\x72\x20\x02\x2d\xe1\x60\x55\xb6\xe4\xc4\x7d\x78\xb4\xdc\x74\x6e\xe1\xe2\xfe\xd1\x95\xec\xa7\xae\x52\x47\x22\x43\x99\x8e\xad\xa1\x1a\x1d\x4e\xed\x60\x6d\x75\xd3\x45\x7a\x4e\x4d\x48\x7d\xe3\x8c\x14\x7c\xf0\xc9\x3c\x34\xe1\x38\xfd\x1c\x62\x5c\xf1\xcc\x54\x95\x23\x9f\x5d\x7f\x89\xee\xf4\xca\xeb\xde\x38\x2c\xdd\x96\x5c\xb3\x50\x62\x2f\xc8\x85\xd7\xb0\xf0\x3a\x73\xcd\x72\x49\xce\xb1\xcf\x3f\x20\x74\xbf\xc0\xb1\xd7\xdd\xd7\xf1\x58\xe4\x81\x88\xd1\x2f\x0c\xbd\x6e\x6a\xa6\x7f\xbd\x97\x61\xef\xdb\xcc\x5d\x7b\x4c\x84\x11\xe7\xd1\x8b\xe5\xef\xef\xfe\x7c\x03\x97\xb1\xfd\xcf\xdf\x6f\x80\xe5\x1c\xdd\x6a\x61\xd0\xf2\x1c\x9d\x23\xef\xf2\x35\x69\x6e\xac\xcb\xd3\xeb\xce\xd8\xfc\xae\xd7\x99\x28\xa1\x27\x21\x6a\x85\xd2\xd1\xab\x14\x3c\xae\xfc\x7f\x00\x00\x00\xff\xff\x7f\xf0\x0d\x30\x14\x12\x00\x00")

func templatesViewsMessageHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesViewsMessageHtml,
		"templates/views/message.html",
	)
}

func templatesViewsMessageHtml() (*asset, error) {
	bytes, err := templatesViewsMessageHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/views/message.html", size: 4628, mode: os.FileMode(420), modTime: time.Unix(1537231244, 0)}
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
	"templates/views/message.html": templatesViewsMessageHtml,
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
			"message.html": &bintree{templatesViewsMessageHtml, map[string]*bintree{}},
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