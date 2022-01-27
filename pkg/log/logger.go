package log

import (
	"context"
	"fmt"
	"github.com/piupuer/go-helper/pkg/constant"
	"gorm.io/gorm/logger"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	logDir    = ""
	helperDir = ""
)

func init() {
	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	logDir = regexp.MustCompile(`logger\.go`).ReplaceAllString(file, "")
	helperDir = regexp.MustCompile(`go-helper.pkg.log.logger\.go`).ReplaceAllString(file, "")
}

// Interface logger interface
type Interface interface {
	Options() Options
	WithFields(fields map[string]interface{}) Interface
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})
}

type Config struct {
	ops  Options
	gorm logger.Config
}

func New(options ...func(*Options)) (l Interface) {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	switch ops.category {
	case constant.LogCategoryZap:
		l = newZap(ops)
	case constant.LogCategoryLogrus:
		l = newLogrus(ops)
	default:
		l = newLogrus(ops)
	}
	return l
}

func getRequestId(ctx context.Context) (id string) {
	if interfaceIsNil(ctx) {
		return
	}
	// get value from context
	requestIdValue := ctx.Value(constant.MiddlewareRequestIdCtxKey)
	if item, ok := requestIdValue.(string); ok && item != "" {
		id = item
	}
	return
}

func fileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, logDir) || strings.HasSuffix(file, "_test.go")) && !strings.Contains(file, "src/runtime") {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

func removePrefix(s1 string, s2 string, ops Options) string {
	res1 := removeBaseDir(s1, ops)
	res2 := removeBaseDir(s2, ops)
	if strings.HasPrefix(s1, logDir) {
		return res2
	}
	f1 := len(res1) <= len(res2)
	f2 := strings.HasPrefix(s1, logDir)
	// src/runtime may be in go routine
	if strings.Contains(res2, "src/runtime") || (f1 || !f1 && f2) {
		return res1
	}
	return res2
}

func removeBaseDir(s string, ops Options) string {
	sep := string(os.PathSeparator)
	if !ops.keepSourceDir && strings.HasPrefix(s, helperDir) {
		s = strings.TrimPrefix(s, path.Dir(helperDir)+"/")
	}
	if strings.HasPrefix(s, ops.lineNumPrefix) {
		s = strings.TrimPrefix(s, ops.lineNumPrefix)
	}
	arr := strings.Split(s, "@")
	if len(arr) == 2 {
		arr1 := strings.Split(arr[0], sep)
		arr2 := strings.Split(arr[1], sep)
		if ops.lineNumLevel > 0 {
			if ops.lineNumLevel < len(arr1) {
				arr1 = arr1[len(arr1)-ops.lineNumLevel:]
			}
		}
		if !ops.keepVersion {
			arr2 = arr2[1:]
		}
		s1 := strings.Join(arr1, sep)
		s2 := strings.Join(arr2, sep)
		if !ops.keepVersion {
			s = fmt.Sprintf("%s%s%s", s1, sep, s2)
		} else {
			s = fmt.Sprintf("%s@%s", s1, s2)
		}
	}
	return s
}

func interfaceIsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return i == nil
}
