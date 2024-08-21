package tool

import "encoding/base64"

// 解码 Base64 编码的密码
func DecodeBase64(encodedPassword string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedPassword)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
