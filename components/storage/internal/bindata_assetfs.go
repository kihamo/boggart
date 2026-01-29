// Code generated for package internal by go-bindata DO NOT EDIT. (@generated)
// sources:
// locales/ru/LC_MESSAGES/config.mo
package internal

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

var _localesRuLc_messagesConfigMo = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8e\xb1\xaa\x13\x41\x18\x85\x4f\x34\xd5\x96\xa9\x2d\x7e\x0b\x83\x82\x13\x77\xa3\x85\x4c\x76\x12\x51\x12\x0c\x66\x41\xc2\x6a\x65\x33\x66\x27\x9b\xc5\xdd\x99\x65\x66\x56\x14\x52\xd8\xfa\x00\x76\xb6\xf6\x36\x42\x40\x7c\x87\xcd\x0b\xf8\x2c\x92\x4d\xbc\x97\xfb\x57\xe7\xe3\x7c\x07\xfe\xbf\x83\xfe\x37\x00\xb8\x05\xe0\x0e\x80\x87\x00\x6e\x03\x88\x71\xbe\x97\x00\x06\x00\x96\x00\x66\x3d\xe0\xdd\x85\xbf\xf7\x80\xde\x65\xd7\xdd\xc6\xe8\x6d\x91\x33\x2f\xdf\xf7\x9d\x37\x56\xe6\x0a\x6b\x55\x1b\xeb\x59\xe2\xf2\x22\x63\xcf\x9b\xdc\xb1\xd4\x70\xca\xd4\xc7\x67\x1f\x8a\x9d\xac\xcc\xc8\x36\xc1\x4a\x3a\xcf\x52\x2b\xb5\x2b\xa5\x37\x96\xd3\xab\xae\xa2\xa4\xb1\xb2\x32\x99\xa1\xf8\x86\x3f\x0d\x56\x52\xe7\x8d\xcc\x15\x4b\x95\xac\x38\x5d\x31\xa7\x75\xe3\x5c\x21\x75\x90\x2c\x93\x39\x7b\xab\xac\x2b\x8c\xe6\x14\x8d\xc2\xe0\x85\xd1\x5e\x69\xcf\xd2\xcf\xb5\xe2\xe4\xd5\x27\xff\xa8\x2e\x65\xa1\x27\xb4\xd9\x49\xeb\x94\x17\x6f\xd2\x05\x7b\x7a\xed\x9d\xfe\xd9\x2a\xcb\xe6\x7a\x63\xb2\x42\xe7\x9c\x82\xd7\x65\x63\x65\xc9\x16\xc6\x56\x8e\x93\xae\x3b\x74\xe2\xf1\x84\xce\x51\xe8\x7b\x51\x28\x44\x44\xc3\x21\x9d\x62\x78\x57\x44\x11\xcd\x28\x24\xde\xf1\x54\x8c\xff\x57\xb1\x78\x72\x8a\xf7\x3b\x2d\x8e\x42\xda\xef\xcf\x93\xa9\x18\x87\x0f\x68\x46\x11\x71\x1a\x4f\xd0\xfe\x38\x7e\x69\x7f\xb6\x7f\xda\x43\xfb\xbb\x3d\x1c\xbf\xb6\xbf\xf0\x2f\x00\x00\xff\xff\xcd\xb7\x21\x91\xaf\x01\x00\x00"

func localesRuLc_messagesConfigMoBytes() ([]byte, error) {
	return bindataRead(
		_localesRuLc_messagesConfigMo,
		"locales/ru/LC_MESSAGES/config.mo",
	)
}

func localesRuLc_messagesConfigMo() (*asset, error) {
	bytes, err := localesRuLc_messagesConfigMoBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "locales/ru/LC_MESSAGES/config.mo", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"locales/ru/LC_MESSAGES/config.mo": localesRuLc_messagesConfigMo,
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
	"locales": &bintree{nil, map[string]*bintree{
		"ru": &bintree{nil, map[string]*bintree{
			"LC_MESSAGES": &bintree{nil, map[string]*bintree{
				"config.mo": &bintree{localesRuLc_messagesConfigMo, map[string]*bintree{}},
			}},
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
