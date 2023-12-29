// Package context 预定义 context.Context 日志字段提取函数.
package context

import (
	"context"
)

func getStringContext(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}
	switch v := ctx.Value(key).(type) {
	case string:
		return v, true
	case *string:
		return *v, true
	}
	return "", false
}
