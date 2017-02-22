package sync

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func SHA1(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	hash := sha1.New()
	if _, err = io.Copy(hash, file); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
