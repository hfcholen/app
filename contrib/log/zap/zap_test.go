package zap

import (
	"testing"

	"github.com/hfcholen/app/log"
)

func TestZap(t *testing.T) {
	config := Config{
		Format:     "console",
		Output:     "both",
		Level:      "info", // debug,info,warn,error,fatal
		Filename:   "./logs/app.log",
		MaxSize:    100,  // 文件最大 100MB
		MaxBackups: 3,    // 最多备份 3 个文件
		MaxAge:     7,    // 保存 7 天
		Compress:   true, // 启用压缩
	}
	logger := NewLogger(config)

	defer logger.Sync()

	log.SetLogger(logger)
	log.Info("info")
	log.Debug("debug")
	log.Debug("warn")
	log.Error("error")
	log.Fatal("fatal")
}
