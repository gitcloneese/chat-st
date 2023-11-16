package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
)

func HmacSha256(key string, data []byte) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(data)
	return mac.Sum(nil)
}

func Base64Std(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Url(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Escape(str string) string {
	return url.QueryEscape(str)
}
