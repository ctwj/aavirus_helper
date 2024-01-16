package log

import (
	"context"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Define key
const (
	TraceIDKey = "trace_id"
)

type (
	traceIDKey struct{}
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type ZapLogger struct {
	logger    *zap.SugaredLogger
	zapLogger *zap.Logger
	atom      zap.AtomicLevel
	options   *Options
	writer    io.Writer
}

func (l *ZapLogger) Sync() {
	l.logger.Sync()
}

func (l *ZapLogger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}
func (l *ZapLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}
func (l *ZapLogger) Info(v ...interface{}) {
	l.logger.Info(v...)
}
func (l *ZapLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}
func (l *ZapLogger) Warn(v ...interface{}) {
	l.logger.Warn(v...)
}
func (l *ZapLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}
func (l *ZapLogger) Error(v ...interface{}) {
	l.logger.Error(v...)
}
func (l *ZapLogger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}
func (l *ZapLogger) Panic(v ...interface{}) {
	l.logger.Panic(v...)
}
func (l *ZapLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}
func (l *ZapLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}
func (l *ZapLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func (l *ZapLogger) Println(v ...interface{}) {
	l.Info(v...)
}

func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	var fields []interface{}

	if v := FromTraceIDContext(ctx); v != "" {
		fields = append(fields, zap.String(TraceIDKey, v))
	}
	return &ZapLogger{
		logger: l.logger.With(fields...),
	}
}

func (l *ZapLogger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func (l *ZapLogger) SetLogLevel(level Level) {
	l.atom.SetLevel(parseLevel(level))
}

func (l *ZapLogger) GetLogLevel() Level {
	return parseZapLevel(l.atom.Level())
}

func (l *ZapLogger) Zap() *zap.Logger {
	return l.zapLogger
}

func (l *ZapLogger) Writer() io.Writer {
	return l.writer
}

func NewZapLogger(o *Options) *ZapLogger {

	writers := []zapcore.WriteSyncer{}
	osFileout := zapcore.AddSync(&lumberjack.Logger{
		Filename:   o.Filename,
		MaxSize:    o.MaxSize, // megabytes
		MaxAge:     o.MaxAge,  // days
		MaxBackups: o.MaxBackups,
		LocalTime:  true,
		Compress:   o.Compress,
	})
	if o.IsStdOut {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	writers = append(writers, osFileout)
	w := zapcore.NewMultiWriteSyncer(writers...)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(parseLevel(o.LogLevel)) //改变日志级别

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	var enc zapcore.Encoder
	if o.LogType != "json" {
		enc = zapcore.NewConsoleEncoder(cfg)
	} else {
		enc = zapcore.NewJSONEncoder(cfg)
	}

	core := zapcore.NewCore(
		//这里控制json 或者不是json 类型
		enc,
		w,
		atom,
	)
	logger := zap.New(
		core,
		zap.AddStacktrace(parseLevel(o.Stacktrace)),
		// zap.AddCaller(),
		// zap.AddCallerSkip(2),
	)

	loggerSugar := logger.Sugar()
	return &ZapLogger{
		logger:    loggerSugar,
		zapLogger: logger,
		atom:      atom,
		options:   o,
		writer:    w,
	}

}

func parseLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	}

	return zapcore.DebugLevel
}

func parseZapLevel(level zapcore.Level) Level {
	switch level {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.PanicLevel:
		return PanicLevel
	case zapcore.FatalLevel:
		return FatalLevel
	}
	return DebugLevel
}
