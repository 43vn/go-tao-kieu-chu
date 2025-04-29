package api

import (
	"encoding/base64"
	"html"
	"os"
)

func workingDir() string {
	dir, _ := os.Getwd()
	return dir
}

// Mã hóa lỗi thành base64
func B64E(errorMsg string) string {
	return base64.URLEncoding.EncodeToString([]byte(errorMsg))
}

// Giải mã base64 thành chuỗi
func B64D(encoded string) string {
	decodedBytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return ""
	}
	// Chuyển byte slice thành chuỗi
	decodedString := string(decodedBytes)

	// Xóa HTML entities để tránh các lỗ hổng XSS
	cleanString := html.UnescapeString(decodedString)

	return cleanString
}
