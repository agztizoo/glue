package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// WithInitializeHook 实现 gorm.Dialector Initialize 方法 Hook 能力.
func WithInitializeHook(dial Dialector, hooks ...func(*gorm.DB) error) Dialector {
	return func(opts *Options) (gorm.Dialector, error) {
		dl, err := dial(opts)
		if err != nil {
			return nil, err
		}
		return &initializeHook{
			Dialector: dl,
			hooks:     hooks,
		}, nil
	}
}

type initializeHook struct {
	gorm.Dialector
	hooks []func(*gorm.DB) error
}

func (h *initializeHook) Initialize(db *gorm.DB) error {
	if err := h.Dialector.Initialize(db); err != nil {
		return err
	}
	for _, hook := range h.hooks {
		if hook == nil {
			continue
		}
		if err := hook(db); err != nil {
			return err
		}
	}
	return nil
}

// 连接池默认配置.
var (
	DefaultMaxIdleConns    = 100
	DefaultMaxOpenConns    = 200
	DefaultConnMaxLifeTime = 300 * time.Second
)

// ConnPoolHook 实现数据库初始化时, 连接池初始化.
//
// https://github.com/go-sql-driver/mysql#important-settings
func ConnPoolHook(maxIdleConns, maxOpenConns int, connMaxLifeTime time.Duration) func(*gorm.DB) error {
	return func(db *gorm.DB) error {
		d, err := db.DB()
		if err != nil {
			return err
		}

		if maxIdleConns <= 0 {
			maxIdleConns = DefaultMaxIdleConns
		}
		if maxOpenConns <= 0 {
			maxOpenConns = DefaultMaxOpenConns
		}
		if connMaxLifeTime <= 0 {
			connMaxLifeTime = DefaultConnMaxLifeTime
		}
		d.SetMaxIdleConns(maxIdleConns)
		d.SetMaxOpenConns(maxOpenConns)
		// 兼容 go1.4
		if s, ok := (interface{}(d)).(interface {
			SetConnMaxIdleTime(time.Duration)
		}); ok {
			s.SetConnMaxIdleTime(connMaxLifeTime)
		}
		return nil
	}
}

// NewTagProcessor 创建 gorm tag 处理器.
func NewTagProcessor(tagName string, marshaller, unmarshaller TagHandler) *TagProcessor {
	return &TagProcessor{
		tagName:      tagName,
		marshaller:   marshaller,
		unmarshaller: unmarshaller,
	}
}

// TagProcessor 代表数据库 Tag 处理器.
type TagProcessor struct {
	// 需要处理的 Tag 名称.
	//
	// e.g. encrypt
	tagName      string
	marshaller   TagHandler
	unmarshaller TagHandler
}

// Marshal 实现数据编码.
func (p *TagProcessor) Marshal(db *gorm.DB) {
	if p.marshaller == nil {
		return
	}
	p.loopFields(db, p.marshaller)
}

// Unmarshal 实现数据解码.
func (p *TagProcessor) Unmarshal(db *gorm.DB) {
	if p.unmarshaller == nil {
		return
	}
	p.loopFields(db, p.unmarshaller)
}

func (p *TagProcessor) loopFields(db *gorm.DB, handler TagHandler) {
	if db.Error != nil {
		return
	}
	if db.Statement.Schema == nil {
		return
	}
	if p.tagName == "" {
		return
	}

	for _, field := range db.Statement.Schema.Fields {
		tagValue, ok := field.Tag.Lookup(p.tagName)
		if !ok {
			continue
		}
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue.Index(i))
				fieldValue, ok, err := handler(db.Statement.Context, tagValue, fieldValue, isZero)
				if err != nil {
					db.AddError(err)
					return
				}
				if !ok {
					continue
				}
				if err := field.Set(db.Statement.ReflectValue.Index(i), fieldValue); err != nil {
					db.AddError(err)
					return
				}
			}
		case reflect.Struct:
			fieldValue, isZero := field.ValueOf(db.Statement.ReflectValue)
			fieldValue, ok, err := handler(db.Statement.Context, tagValue, fieldValue, isZero)
			if err != nil {
				db.AddError(err)
				return
			}
			if !ok {
				continue
			}
			if err := field.Set(db.Statement.ReflectValue, fieldValue); err != nil {
				db.AddError(err)
				return
			}
		}
	}
}

// TagHandler 代表标记 Tag 字段值处理程序.
//
// tagValue： tag 标记对应的值.
// fieldValue: 字段值.
// isZero: fieldValue 是否对应类型的零值.
// retVal: 处理后的值.
// changed: 是否需要更新字段值.
type TagHandler func(ctx context.Context, tagValue string, fieldValue interface{}, isZero bool) (retVal interface{}, changed bool, err error)

// StringTagHandler 字符串类型 Tag 字段处理程序.
type StringTagHandler func(ctx context.Context, tagValue string, fieldValue string) (retVal string, err error)

// WrapStringTagHandler 将 StringTagHandler 转换为 TagHandler.
//
// handler: fieldValue，retVal 都为 string.
//
// 转换后:
//   - string
//   - *string
//   - []byte
func WrapStringTagHandler(handler StringTagHandler) TagHandler {
	return func(ctx context.Context, tagValue string, fieldValue interface{}, isZero bool) (interface{}, bool, error) {
		switch v := fieldValue.(type) {
		case string:
			rt, err := handler(ctx, tagValue, v)
			if err != nil {
				return v, false, err
			}
			return rt, rt != v, err
		case *string:
			if v == nil {
				return nil, false, nil
			}
			rt, err := handler(ctx, tagValue, *v)
			if err != nil {
				return nil, false, err
			}
			return &rt, rt != *v, nil
		case []byte:
			if v == nil {
				return nil, false, nil
			}
			rt, err := handler(ctx, tagValue, string(v))
			if err != nil {
				return nil, false, err
			}
			return []byte(rt), rt != string(v), nil
		default:
			if fieldValue == nil {
				return fieldValue, false, nil
			}
			return fieldValue, false, fmt.Errorf("non-stringable type: %T", fieldValue)
		}
	}
}

// NonEmptyStringTagHandler 过滤空字符串 Tag 处理程序.
func NonEmptyStringTagHandler(handler StringTagHandler) StringTagHandler {
	return func(ctx context.Context, tagValue string, fieldValue string) (string, error) {
		if fieldValue == "" {
			return fieldValue, nil
		}
		return handler(ctx, tagValue, fieldValue)
	}
}
