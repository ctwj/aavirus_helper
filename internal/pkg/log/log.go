package log

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// 日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

type Options struct {
	Filename   string //日志保存路径
	LogLevel   Level  //日志记录级别
	MaxSize    int    //日志分割的尺寸 MB
	MaxAge     int    //分割日志保存的时间 day
	MaxBackups int    //分割日志最大保留数量 个 0 代表都保留
	Stacktrace Level  //记录堆栈的级别
	IsStdOut   bool   //是否标准输出console输出
	Compress   bool   //是否压缩
	LogType    string //日志类型,普通 或 json
}

type Logger interface {
	SetLogLevel(Level)
	GetLogLevel() Level
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Writer() io.Writer
	WithContext(ctx context.Context) Logger
}

var (
	once   sync.Once
	logger Logger

	level = InfoLevel

	prefix string
)

func insertPrefix(format string) string {
	if len(prefix) > 0 {
		return prefix + " " + format
	}
	return format
}

func withLevel(fn func(v ...interface{}), l Level, v ...interface{}) {
	if l < level {
		return
	}
	fn(v...)
}

func withLevelf(fn func(format string, v ...interface{}), l Level, format string, v ...interface{}) {
	if l < level {
		return
	}
	fn(format, v...)
}

func Debug(v ...interface{}) {
	withLevel(logger.Debug, DebugLevel, v...)
}

func Debugf(format string, v ...interface{}) {
	withLevelf(logger.Debugf, DebugLevel, format, v...)
}

func Info(v ...interface{}) {
	withLevel(logger.Info, InfoLevel, v...)
}

func Infof(format string, v ...interface{}) {
	withLevelf(logger.Infof, InfoLevel, format, v...)
}

func Warn(v ...interface{}) {
	withLevel(logger.Warn, WarnLevel, v...)
}

func Warnf(format string, v ...interface{}) {
	withLevelf(logger.Warnf, WarnLevel, format, v...)
}

func Error(v ...interface{}) {
	withLevel(logger.Error, ErrorLevel, v...)
}

func Errorf(format string, v ...interface{}) {
	withLevelf(logger.Errorf, ErrorLevel, format, v...)
}

func Panic(v ...interface{}) {
	withLevel(logger.Panic, PanicLevel, v...)
}

func Panicf(format string, v ...interface{}) {
	withLevelf(logger.Panicf, PanicLevel, format, v...)
}

func Fatal(v ...interface{}) {
	withLevel(logger.Fatal, FatalLevel, v...)
}

func Fatalf(format string, v ...interface{}) {
	withLevelf(logger.Fatalf, FatalLevel, format, v...)
}

func WithContext(ctx context.Context) Logger {
	return logger.WithContext(ctx)
}

func SetLevel(l Level) {
	level = l
	logger.SetLogLevel(l)
}

func GetLevel() Level {
	return level
}

func SetPrefix(p string) {
	prefix = p
}

// Set service name
func Name(name string) {
	prefix = fmt.Sprintf("[%s]", name)
}

func InitDefault(o *Options) {
	once.Do(func() {
		logger = NewZapLogger(o)
		level = o.LogLevel
	})
}

func SetDefault(l Logger) {
	logger = l
}

func Default() Logger {
	return logger
}
