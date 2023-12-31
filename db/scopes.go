package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// NewInjectFromContextScope 实现 context 进行 SQL 注入.
//
// 从 context 读取值并注入 SQL 语句.
//
// optional 为 true, 当 keyf 返回 value 为空不进行注入.
// optional 为 false, 当 keyf 返回 value 为空 panic.
func NewInjectFromContextScope(field string, keyFrom func(context.Context) string, optional bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		value := keyFrom(db.Statement.Context)
		if value == "" {
			if !optional {
				panic(fmt.Errorf("failed to inject required query field `%s`", field))
			}
			return db
		}
		return db.Where(fmt.Sprintf("`%s` = ?", field), value)
	}
}

// NewMustInjectFromContextScope 实现 context 进行 SQL 注入.
//
// 从 context 读取值并注入 SQL 语句.
//
// 当 keyf 返回 value 为空 panic.
func NewMustInjectFromContextScope(field string, keyf func(context.Context) string) func(*gorm.DB) *gorm.DB {
	return NewInjectFromContextScope(field, keyf, false)
}
