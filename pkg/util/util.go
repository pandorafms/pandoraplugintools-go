package util

import (
	"crypto/md5"
	"encoding/hex"
	"runtime"
	"time"
)

func GenerateMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func GetOS() string {
	return runtime.GOOS
}

// PandoraTimestampLayout is the timestamp format expected by the Pandora server.
const PandoraTimestampLayout = "2006/01/02 15:04:05"

func Now(timezone ...string) string {
	loc := time.Local
	if len(timezone) > 0 && timezone[0] != "" {
		if tz, err := time.LoadLocation(timezone[0]); err == nil {
			loc = tz
		}
	}
	return time.Now().In(loc).Format(PandoraTimestampLayout)
}
