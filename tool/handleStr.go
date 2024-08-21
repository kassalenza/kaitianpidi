package tool

import "strings"

// 去掉字符串后的\n
func TrimSuffixNewLine(data string) string {
	if strings.HasSuffix(data, "\n") {
		return strings.TrimSuffix(data, "\n")
	}
	return data
}
