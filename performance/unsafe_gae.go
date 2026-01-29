//go:build appengine
// +build appengine

package performance

func UnsafeBytes2String(b []byte) string {
	return string(b)
}

func UnsafeString2Bytes(s string) []byte {
	return []byte(s)
}
