package log

import (
	"context"
	"log"
)

// 功能：核心日志记录器的定义，提供标准化的日志接口和基本实现。
// 应用场景：Kratos 日志组件的主入口，定义核心功能，如记录日志、设置输出目标、指定日志格式等。

// DefaultLogger is default logger.
var DefaultLogger = NewStdLogger(log.Writer())

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type logger struct {
	logger    Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (c *logger) Log(level Level, keyvals ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, keyvals...)
	return c.logger.Log(level, kvs...)
}

// With with logger fields.
func With(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, kv...)
	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	switch v := l.(type) {
	default:
		return &logger{logger: l, ctx: ctx}
	case *logger:
		lv := *v
		lv.ctx = ctx
		return &lv
	case *Filter:
		fv := *v
		fv.logger = WithContext(ctx, fv.logger)
		return &fv
	}
}

// filter.go

// 功能：定义日志过滤器的实现，用于过滤掉不符合特定条件的日志。例如，可以根据日志级别过滤掉低优先级日志。
// 应用场景：适用于希望仅记录特定级别（如 Error 或 Warning）的日志，或根据动态条件过滤日志。
// global.go

// 功能：提供全局日志对象的管理，类似单例模式，使日志可以在全局范围内被使用。
// 应用场景：当需要统一管理日志实例，方便在多个模块或服务之间共享相同的日志配置时。
// helper.go

// 功能：提供日志的辅助工具类，封装常用的日志记录操作，比如快速记录 Info、Error 等级别的日志。
// 应用场景：在代码中简化日志调用逻辑，减少重复代码，提供一致性的日志记录体验。
// helper_writer.go

// 功能：实现 io.Writer 接口，使日志组件能够与标准库或第三方库兼容。例如，可以将日志直接写入到 os.Stdout 或文件中。
// 应用场景：适用于需要将日志与其他依赖 io.Writer 的系统（如 HTTP 响应或流式输出）集成的场景。
// level.go

// 功能：定义日志级别（如 Debug、Info、Warn、Error 等）以及与之相关的操作。
// 应用场景：为日志分类提供基础，方便开发者根据优先级控制日志输出策略。
// log.go

// 功能：核心日志记录器的定义，提供标准化的日志接口和基本实现。
// 应用场景：Kratos 日志组件的主入口，定义核心功能，如记录日志、设置输出目标、指定日志格式等。
// std.go

// 功能：实现标准输出日志器（如 stdout 和 stderr），用于将日志直接输出到终端。
// 应用场景：用于调试和开发环境中直接查看日志信息。
// value.go

// 功能：定义日志的结构化字段操作，例如支持以 key-value 形式记录日志附加信息。
// 应用场景：在分布式系统或需要结构化日志的场景中使用（如 ELK、Loki），便于后续分析和检索。
