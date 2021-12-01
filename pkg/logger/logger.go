package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-module/carbon"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	sourceDir = ""
)

func init() {
	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	sourceDir = strings.TrimSuffix(file, fmt.Sprintf("%slogger%slogger.go", string(os.PathSeparator), string(os.PathSeparator)))
}

// Interface logger interface
type Interface interface {
	LogMode(logger.LogLevel) logger.Interface
	LogLevel(Level) Interface
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

// zap logger for gorm
type Logger struct {
	Config
	log                                            *zap.Logger
	normalStr, traceStr, traceErrStr, traceWarnStr string
}

type Level zapcore.Level

type Config struct {
	ops           Options
	gorm          logger.Config
	lineNumPrefix string
	lineNumLevel  int
	keepSourceDir bool
}

// New logger like gorm2
func New(options ...func(*Options)) *Logger {
	ops := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	return newWithOption(ops)
}

func newWithOption(ops *Options) *Logger {
	writer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	if ops.lumber {
		now := time.Now()
		filename := ops.lumberOps.Filename
		if filename == "" {
			filename = fmt.Sprintf("%s/%04d-%02d-%02d%s", ops.lumberOps.LogPath, now.Year(), now.Month(), now.Day(), ops.lumberOps.LogSuffix)
		}
		ops.lumberOps.Filename = filename
		hook := &ops.lumberOps
		defer hook.Close()
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	}

	enConfig := zap.NewProductionEncoderConfig()
	enConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	enConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(carbon.Time2Carbon(t).ToRfc3339String())
	}
	if ops.colorful {
		enConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		enConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enConfig),
		writer,
		zapcore.Level(ops.level),
	)
	l := zap.New(core)
	return NewWithZap(
		l,
		Config{
			ops: *ops,
			gorm: logger.Config{
				Colorful: ops.colorful,
			},
			lineNumPrefix: ops.lineNumPrefix,
			lineNumLevel:  ops.lineNumLevel,
			keepSourceDir: ops.keepSourceDir,
		},
	)
}

// New with zap
func NewWithZap(zapLogger *zap.Logger, config Config) *Logger {
	var (
		normalStr    = "%v%s "
		traceStr     = "%v%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%v%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%v%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.gorm.Colorful {
		normalStr = logger.Cyan + "%v" + logger.Blue + "%s " + logger.Reset
		traceStr = logger.Cyan + "%v" + logger.Blue + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Cyan + "%v" + logger.Blue + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.Cyan + "%v" + logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	l := Logger{
		log:          zapLogger,
		Config:       config,
		normalStr:    normalStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	return &l
}

// get default logger
func DefaultLogger() *Logger {
	return New(WithLumber(false))
}

// LogMode gorm log mode
// LogMode log mode
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	zapLevel := zapcore.InfoLevel
	switch level {
	case logger.Warn:
		zapLevel = zapcore.WarnLevel
	case logger.Error:
		zapLevel = zapcore.ErrorLevel
	case logger.Silent:
		zapLevel = zapcore.PanicLevel
	}
	newLogger.ops.level = Level(zapLevel)
	return newWithOption(&newLogger.ops)
}

func (l *Logger) LogLevel(level Level) Interface {
	newLogger := *l
	newLogger.ops.level = level
	return newWithOption(&newLogger.ops)
}

// Debug print info
func (l Logger) Debug(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.DebugLevel) {
		requestId := l.getRequestId(ctx)
		l.log.Sugar().Debugf(l.normalStr+format, append([]interface{}{requestId, l.removePrefix(utils.FileWithLineNum(), fileWithLineNum())}, args...)...)
	}
}

// Info print info
func (l Logger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.InfoLevel) {
		requestId := l.getRequestId(ctx)
		l.log.Sugar().Infof(l.normalStr+format, append([]interface{}{requestId, l.removePrefix(utils.FileWithLineNum(), fileWithLineNum())}, args...)...)
	}
}

// Warn print warn messages
func (l Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.WarnLevel) {
		requestId := l.getRequestId(ctx)
		l.log.Sugar().Warnf(l.normalStr+format, append([]interface{}{requestId, l.removePrefix(utils.FileWithLineNum(), fileWithLineNum())}, args...)...)
	}
}

// Error print error messages
func (l Logger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.log.Core().Enabled(zapcore.ErrorLevel) {
		requestId := l.getRequestId(ctx)
		l.log.Sugar().Errorf(l.normalStr+format, append([]interface{}{requestId, l.removePrefix(utils.FileWithLineNum(), fileWithLineNum())}, args...)...)
	}
}

// Trace print sql message
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if !l.log.Core().Enabled(zapcore.DPanicLevel) {
		return
	}
	lineNum := l.removePrefix(utils.FileWithLineNum(), fileWithLineNum())
	elapsed := time.Since(begin)
	elapsedF := float64(elapsed.Nanoseconds()) / 1e6
	sql, rows := fc()
	row := "-"
	if rows > -1 {
		row = fmt.Sprintf("%d", rows)
	}
	requestId := l.getRequestId(ctx)
	switch {
	case l.log.Core().Enabled(zapcore.ErrorLevel) && err != nil && (!l.gorm.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		l.log.Error(fmt.Sprintf(l.traceErrStr, requestId, lineNum, err, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.WarnLevel) && elapsed > l.gorm.SlowThreshold && l.gorm.SlowThreshold != 0:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.gorm.SlowThreshold)
		l.log.Warn(fmt.Sprintf(l.traceWarnStr, requestId, lineNum, slowLog, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.DebugLevel):
		l.log.Debug(fmt.Sprintf(l.traceStr, requestId, lineNum, elapsedF, row, sql))
	case l.log.Core().Enabled(zapcore.InfoLevel):
		l.log.Info(fmt.Sprintf(l.traceStr, requestId, lineNum, elapsedF, row, sql))
	}
}

func (l Logger) GetZapLog() *zap.Logger {
	return l.log
}

func (l Logger) getRequestId(ctx context.Context) string {
	var v interface{}
	vi := reflect.ValueOf(ctx)
	if vi.Kind() == reflect.Ptr {
		if !vi.IsNil() {
			v = ctx.Value(l.ops.requestIdCtxKey)
		}
	}
	requestId := ""
	if v != nil {
		requestId = fmt.Sprintf("%v ", v)
	}
	return requestId
}

func fileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, sourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

func (l Logger) removePrefix(s1 string, s2 string) string {
	res1 := l.removeBaseDir(s1)
	res2 := l.removeBaseDir(s2)
	if len(res1) < len(res2) && !strings.HasPrefix(s1, sourceDir) {
		return res1
	}
	return res2
}

func (l Logger) removeBaseDir(s string) string {
	if !l.keepSourceDir && strings.HasPrefix(s, sourceDir) {
		s = strings.TrimPrefix(s, path.Dir(sourceDir)+"/")
	}
	if strings.HasPrefix(s, l.lineNumPrefix) {
		s = strings.TrimPrefix(s, l.lineNumPrefix)
	}
	arr := strings.Split(s, "@")
	if len(arr) == 2 {
		if l.lineNumLevel > 0 {
			s = fmt.Sprintf("%s@%s", l.getParentDir(arr[0], l.lineNumLevel), arr[1])
		}
	}
	return s
}

func (l Logger) getParentDir(dir string, index int) string {
	d, b := filepath.Split(filepath.Clean(dir))
	parentDir := ""
	if index > 0 {
		parentDir = l.getParentDir(d, index-1)
	}
	if parentDir != "" {
		return fmt.Sprintf("%s%s%s", parentDir, string(os.PathSeparator), b)
	}
	return b
}
