package log

import (
	"testing"
)

func TestLog(_ *testing.T) {
	logger := DefaultLogger
	logger.Log(LevelInfo, "testlog", "v1", "k1", "v1")
	With(logger, "k1", "v1").Log(LevelInfo)

	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}
