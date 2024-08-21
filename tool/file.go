package tool

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// create dir

// create file

// read file

// append to file
// file: abs path
func Write2fileAppend(data string, file string) error {
	fileIns, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		msg := fmt.Sprintf("Error opening file %s: %v\n", file, err)
		err = errors.New(msg)
		return err
	}
	defer fileIns.Close()

	_, err = fileIns.WriteString(data)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file %s: %v\n", file, err)
		err = errors.New(msg)
		return err
	}

	return nil
}

// owerwrite to file
// file: abs path
func Write2fileOverwrite(data string, file string) error {
	fileIns, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		msg := fmt.Sprintf("Error opening file %s: %v\n", file, err)
		err = errors.New(msg)
		return err
	}
	defer fileIns.Close()

	_, err = fileIns.WriteString(data)
	if err != nil {
		msg := fmt.Sprintf("Error writing to file %s: %v\n", file, err)
		err = errors.New(msg)
		return err
	}

	return nil
}

// 判断文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 获取文件的最后修改时间
func GetFileModTime(filename string) (time.Time, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err // 返回零值和错误信息
	}
	return info.ModTime(), nil
}
