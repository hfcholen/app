package logrus

import (
	"testing"

	"github.com/hfcholen/app/log"
)

func TestLogrus(t *testing.T) {
	config := Config{
		Format:     "json",
		Output:     "file",
		Level:      "info",
		Filename:   "./logs/app.log",
		MaxSize:    100,  // 文件最大 100MB
		MaxBackups: 3,    // 最多备份 3 个文件
		MaxAge:     7,    // 保存 7 天
		Compress:   true, // 启用压缩
	}
	// 示例 1：仅输出到控制台，文本格式
	logger1 := NewLogger(config,
		WithMultiOutput(config),
		WithTextFormat(),
	)

	//// 示例 2：仅输出到文件，JSON 格式
	//logger2 := NewLogger(log.LevelDebug,
	//	WithFileOutput(config),
	//	WithJSONFormat(),
	//)
	//
	//// 示例 3：同时输出到控制台和文件，文本格式
	//logger3 := NewLogger(log.LevelDebug,
	//	WithMultiOutput(config),
	//	WithTextFormat(),
	//)
	log.SetLogger(logger1)
	log.Info(111)

	//log.SetLogger(logger2)
	log.Info(222)

	//log.SetLogger(logger3)
	log.Info(333)
}
