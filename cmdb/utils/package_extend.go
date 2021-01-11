package utils

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

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

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	if Exists(path) {
		return !IsDir(path)
	}
	return false
}

func RemoveDuplicate(arr []string) []string {
	resArr := make([]string, 0)
	tmpMap := make(map[string]interface{})
	for _, val := range arr {
		if _, ok := tmpMap[val]; !ok {
			resArr = append(resArr, val)
			tmpMap[val] = nil
		}
	}
	return resArr
}

func LoopMatchString(s string, matchArray []string) (string, error) {
	for _, match := range matchArray {
		compile, err := regexp.Compile(match)
		if err != nil {
			return "", err
		}
		s = compile.FindString(s)
	}
	return s, nil
}

func Trim(old string) string {
	return strings.Trim(strings.Trim(strings.Trim(old, "\n"), "\t"), " ")
}

func DeleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "  ", " ", -1)
	regstr := "\\s{2,}"
	reg, _ := regexp.Compile(regstr)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spc_index := reg.FindStringIndex(string(s2))
	for len(spc_index) > 0 {
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...)
		spc_index = reg.FindStringIndex(string(s2))
	}
	return string(s2)
}

func Call(i interface{}) {
	value := reflect.ValueOf(i)
	for index := 0; index < value.NumMethod(); index++ {
		value.Method(index).Call([]reflect.Value{})
	}
}
