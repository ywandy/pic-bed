package storage

import (
	"net/http"
	"net/url"
	"strings"
)

var imgPerfix = []string{"png", "jpg", "jpeg", "gif", "bmp", "tif", "webp", "exif"}

func GetFileContentType(out []byte) (string, error) {
	contentType := http.DetectContentType(out)
	return contentType, nil
}

func IsWebUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return strings.HasPrefix(u.Scheme, "http")
}

func IsKnownContentType(t string) string {
	for _, iPerfix := range imgPerfix {
		if strings.Contains(t, iPerfix) {
			return iPerfix
		}
	}
	return ""
}
