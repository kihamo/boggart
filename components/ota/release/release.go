package release

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/kihamo/boggart/components/ota"
)

var separator = []byte("|")

func GenerateReleaseID(rl ota.Release) string {
	hasher := md5.New()
	hasher.Write(rl.Checksum())

	if releaseFile, ok := rl.(*LocalFileRelease); ok {
		hasher.Write(separator)
		hasher.Write([]byte(releaseFile.Path()))
	}

	return hex.EncodeToString(hasher.Sum(nil))
}
