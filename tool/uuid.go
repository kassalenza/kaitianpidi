package tool

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	// 生成一个随机的UUID
	id := uuid.New()
	// 将UUID转换为字符串并返回
	return id.String()
}
