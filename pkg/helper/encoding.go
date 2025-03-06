package helper

import "encoding/base64"

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
