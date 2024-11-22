package zap

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/app/log"
	"github.com/natefinch/lumberjack"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Logger)

func NewLogger(cfg Config, opts ...Option) *Logger {
	l := &Logger{
		log:    newZapLogger(cfg),
		msgKey: log.DefaultMessageKey,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen&1 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		fmt.Println(keyvals[i].(string))
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

// Config 是日志配置结构体
type Config struct {
	Format     string // "console" 或 "json"
	Output     string // "stdout", "file", "both"
	Level      string // 日志级别: "debug", "info", "warn", "error", "fatal"
	Filename   string // 日志文件路径
	MaxSize    int    // 文件最大大小 (MB)
	MaxBackups int    // 最大备份数量
	MaxAge     int    // 最大保存天数
	Compress   bool   // 是否启用压缩
}

// newZapLogger 创建 Zap 日志实例
func newZapLogger(cfg Config) *zap.Logger {
	var cores []zapcore.Core

	// 设置日志级别
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志格式
	var encoder zapcore.Encoder
	switch cfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		panic("unsupported log format: " + cfg.Format)
	}

	// 日志输出目标
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleWriter := zapcore.Lock(os.Stdout)
		cores = append(cores, zapcore.NewCore(encoder, consoleWriter, level))
	}
	if cfg.Output == "file" || cfg.Output == "both" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, level))
	}

	// 构建 core
	core := zapcore.NewTee(cores...)

	// 返回 zap logger
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
