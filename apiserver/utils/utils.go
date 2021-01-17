package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// 判断路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断路径是否存在, 且为目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断路径是否存在, 且为文件
func IsFile(path string) bool {
	if Exists(path) {
		return !IsDir(path)
	}
	return false
}

func Md5Sum(s string) string {
	ret := md5.Sum([]byte(s))
	return hex.EncodeToString(ret[:])
}
