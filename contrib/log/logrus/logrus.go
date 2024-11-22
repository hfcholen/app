package logrus

import (
	"io"
	"os"

	"github.com/app/log"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var _ log.Logger = (*Logger)(nil)

// Logger 实现 log.Logger 接口
type Logger struct {
	log *logrus.Logger
}

// Option 配置选项
type Option func(*Logger)

// WithJSONFormat 配置为 JSON 格式
func WithJSONFormat() Option {
	return func(l *Logger) {
		l.log.SetFormatter(&logrus.JSONFormatter{})
	}
}

// WithTextFormat 配置为纯文本格式
func WithTextFormat() Option {
	return func(l *Logger) {
		l.log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// WithFileOutput 配置文件输出
func WithFileOutput(cfg Config) Option {
	return func(l *Logger) {
		l.log.SetOutput(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,    // 每个日志文件最大 100MB
			MaxBackups: cfg.MaxBackups, // 最多保留 5 个备份文件
			MaxAge:     cfg.MaxAge,     // 日志保留 30 天
			Compress:   cfg.Compress,   // 启用压缩
		})
	}
}

// WithConsoleOutput 配置为控制台输出
func WithConsoleOutput() Option {
	return func(l *Logger) {
		l.log.SetOutput(os.Stdout)
	}
}

// WithMultiOutput 同时支持控制台和文件输出
func WithMultiOutput(cfg Config) Option {
	return func(l *Logger) {
		l.log.SetOutput(&multiOutput{
			consoleWriter: os.Stdout,
			fileWriter: &lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,    // 每个日志文件最大 100MB
				MaxBackups: cfg.MaxBackups, // 最多保留 5 个备份文件
				MaxAge:     cfg.MaxAge,     // 日志保留 30 天
				Compress:   cfg.Compress,   // 启用压缩
			},
		})
	}
}

// NewLogger 创建 Logger 实例
func NewLogger(cfg Config, opts ...Option) *Logger {
	l := &Logger{
		log: newLogrusLogger(cfg),
	}
	// 应用所有配置项
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Log 实现 log.Logger 接口
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	logrusLevel := convertLogLevel(level)
	if logrusLevel > l.log.GetLevel() {
		return nil
	}

	var (
		fields logrus.Fields = make(map[string]interface{})
		msg    string
	)

	// 解析 keyvals
	if len(keyvals) > 0 {
		if len(keyvals)&1 != 0 {
			keyvals = append(keyvals, "") // 补齐 keyvals
		}
		for i := 0; i < len(keyvals); i += 2 {
			key, ok := keyvals[i].(string)
			if !ok {
				continue
			}
			if key == log.DefaultMessageKey { // 默认消息键
				msg, _ = keyvals[i+1].(string)
				continue
			}
			fields[key] = keyvals[i+1]
		}
	}

	// 写入日志
	if len(fields) > 0 {
		l.log.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.log.Log(logrusLevel, msg)
	}

	return nil
}

// convertLogLevel 转换 log.Level 为 logrus.Level
func convertLogLevel(level log.Level) logrus.Level {
	switch level {
	case log.LevelDebug:
		return logrus.DebugLevel
	case log.LevelInfo:
		return logrus.InfoLevel
	case log.LevelWarn:
		return logrus.WarnLevel
	case log.LevelError:
		return logrus.ErrorLevel
	case log.LevelFatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

// multiOutput 支持多路输出
type multiOutput struct {
	consoleWriter io.Writer
	fileWriter    io.Writer
}

// Write 实现 io.Writer 接口，支持控制台和文件同时输出
func (m *multiOutput) Write(p []byte) (n int, err error) {
	if n, err = m.consoleWriter.Write(p); err != nil {
		return
	}
	if _, err = m.fileWriter.Write(p); err != nil {
		return
	}
	return
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

// newLogrusLogger 根据 Config 创建 logrus.Logger
func newLogrusLogger(cfg Config) *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	switch cfg.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// 定义日志输出
	var writers []io.Writer
	if cfg.Output == "stdout" || cfg.Output == "both" {
		writers = append(writers, os.Stdout)
	}
	if cfg.Output == "file" || cfg.Output == "both" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, fileWriter)
	}

	// 设置输出目标
	if len(writers) > 1 {
		logger.SetOutput(io.MultiWriter(writers...))
	} else if len(writers) == 1 {
		logger.SetOutput(writers[0])
	}

	return logger
}
