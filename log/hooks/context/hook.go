package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

// New 返回从 context.Context 注入日志字段 Hook.
func New(opts ...Option) logrus.Hook {
	h := &hook{
		valuers: make(map[string]Valuer),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

type hook struct {
	valuers map[string]Valuer
}

func (h *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *hook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}
	for key, valuer := range h.valuers {
		value, ok := valuer(ctx)
		if ok {
			entry.Data[key] = value
		}
	}
	return nil
}

// Valuer 实现从 context.Context 提取值.
//
// 返回 value, false 代表值不存在.
type Valuer func(context.Context) (interface{}, bool)

// StringValuer Valuer 实现从 context.Context 提取值字符串值.
//
// 返回 value, false 代表值不存在.
type StringValuer func(context.Context) (string, bool)

// Option 代表 Hook 构造选项.
type Option func(*hook)

// WithValuer 构造字段注入选项.
//
// name 为注入字段名, valuer 返回 value, false 时不注入.
func WithValuer(name string, valuer Valuer) Option {
	return func(h *hook) {
		h.valuers[name] = valuer
	}
}

// WithMustValuer 构造字段注入选项.
//
// name 为注入字段名, value 为 valuer 返回值.
func WithMustValuer(name string, valuer func(context.Context) interface{}) Option {
	return func(h *hook) {
		h.valuers[name] = func(ctx context.Context) (interface{}, bool) {
			return valuer(ctx), true
		}
	}
}

// WithStringValuer 构造字段注入选项.
//
// name 为注入字段名, valuer 返回 value, false 时不注入.
func WithStringValuer(name string, valuer StringValuer) Option {
	return func(h *hook) {
		h.valuers[name] = func(ctx context.Context) (interface{}, bool) {
			return valuer(ctx)
		}
	}
}

// WithMustStringValuer 构造字段注入选项.
//
// name 为注入字段名, value 为 valuer 返回值.
func WithMustStringValuer(name string, valuer func(context.Context) string) Option {
	return func(h *hook) {
		h.valuers[name] = func(ctx context.Context) (interface{}, bool) {
			return valuer(ctx), true
		}
	}
}
