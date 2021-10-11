package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	camelRe = regexp.MustCompile("(_)([a-zA-Z]+)")
	snakeRe = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// 字符串转为驼峰
func CamelCase(str string) string {
	camel := camelRe.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

// 字符串转为驼峰(首字母小写)
func CamelCaseLowerFirst(str string) string {
	camel := CamelCase(str)
	for i, v := range camel {
		return string(unicode.ToLower(v)) + camel[i+1:]
	}
	return camel
}

// 驼峰式写法转为下划线蛇形写法
func SnakeCase(str string) string {
	snake := snakeRe.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// string数组去重复
func RemoveRepeat(arr []string) []string {
	newArr := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			// 写入空结构体不占用空间
			temp[item] = struct{}{}
			newArr = append(newArr, item)
		}
	}
	return newArr
}
