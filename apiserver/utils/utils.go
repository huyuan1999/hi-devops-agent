package utils

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
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

func Client(timeout time.Duration) http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func FatalError(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
	fmt.Println(err.Error())
	os.Exit(1)
}

func PrintError(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	}
	fmt.Println(err.Error())
}
