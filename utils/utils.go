package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
)

// FileMD5 returns the file md5 hash hex
func FileMD5(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return ""
	}

	return hex.EncodeToString(md5hash.Sum(nil)[:])
}

// FileContentType returns the file content-type
func FileContentType(filepath string) string {
	mimeType := mime.TypeByExtension(path.Ext(filepath))
	if mimeType != "" {
		return mimeType
	}
	
	fileData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return ""
	}

	return http.DetectContentType(fileData)
}
