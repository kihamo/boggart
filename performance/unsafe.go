// +build !appengine

package performance

import (
	"reflect"
	"unsafe"
)

func UnsafeBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func UnsafeString2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
