package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/agztizoo/glue/transaction"
)

// Provider 定义 *gorm.DB 提供者.
//
// 例：
//   特定租户访问单独的数据库集群,其他租户访问默认集群.
type Provider interface {
	// UseDB 实现通过 context 选择数据库.
	//
	// 如果在事务上下文内，返回写库.
	//
	// 不在事务上下文内时, 依据执行语句动态选择读库或写库.
	//
	// 无匹配 DB 时 panic.
	UseDB(context.Context) *gorm.DB

	// UseWriteDB 实现通过 context 选择写库.
	//
	// 无匹配 DB 时 panic.
	UseWriteDB(context.Context) *gorm.DB
}

// NewProvider 创建支持事务管理的 db.Provider
//
// scopes 在新会话创建后通过 db.Scopes(scopes...) 应用.
func NewProvider(source Source, scopes ...func(*gorm.DB) *gorm.DB) *TransProvider {
	p := &TransProvider{Source: source, scopes: scopes}
	lookupDB := func(ctx context.Context) interface{} {
		return p.lookupDB(ctx, true)
	}
	p.Manager = transaction.NewManager(p.getCtxKey, lookupDB, p.transaction)
	return p
}

// ToProvider 转换 *TransProvider 为 Provider.
//
// 用于依赖注入的工厂函数.
func ToProvider(tp *TransProvider) Provider {
	return tp
}

// ToTransactionManager 转换 *TransProvider 为 transaction.Manager.
//
// 用于依赖注入的工厂函数.
func ToTransactionManager(tp *TransProvider) transaction.Manager {
	return tp
}

// TransProvider 实现支持事务上下文的 DB Provider.
type TransProvider struct {
	Source
	transaction.Manager

	scopes []func(*gorm.DB) *gorm.DB
}

var _ transaction.Manager = new(TransProvider)

type transCtxKey string

// getCtxKey 返回事务上下文存储到 context 的 Key.
//
// 事务上下文实现了 transaction.TransContext
//
// 返回的 key 需要转换为私有类型, 防止内容污染.
func (p *TransProvider) getCtxKey(ctx context.Context) interface{} {
	name := p.getWriteDBName(ctx)
	return transCtxKey(name)
}

// lookupDB 查找非事务上下文 DB.
func (p *TransProvider) lookupDB(ctx context.Context, write bool) *gorm.DB {
	if write {
		return p.getWriteDB(ctx).Clauses(dbresolver.Write)
	}
	return p.getReadDB(ctx)
}

// findTransDB 查找事务上下文 DB.
func (p *TransProvider) findTransDB(ctx context.Context) *gorm.DB {
	tc, ok := ctx.Value(p.getCtxKey(ctx)).(transaction.TransContext)
	if !ok {
		return nil
	}
	return tc.GetTransDB().(*gorm.DB)
}

// transaction 执行数据库事务.
func (p *TransProvider) transaction(ctx context.Context, db interface{}, callback func(db interface{}) error) error {
	return db.(*gorm.DB).Transaction(func(db *gorm.DB) error {
		return callback(db)
	})
}

func (p *TransProvider) useDB(ctx context.Context, write bool) *gorm.DB {
	db := p.findTransDB(ctx)
	if db == nil {
		db = p.lookupDB(ctx, write)
	}
	if db == nil {
		panic("matching database not found")
	}
	sess := &gorm.Session{Context: ctx}
	// 保护逻辑
	return db.Session(sess).Scopes(p.scopes...)
}

// UseDB 实现通过 context 选择数据库.
//
// 如果在事务上下文内，返回写库.
//
// 不在事务上下文内时, 依据执行语句动态选择读库或写库.
//
// 无匹配 DB 时 panic.
func (p *TransProvider) UseDB(ctx context.Context) *gorm.DB {
	return p.useDB(ctx, false)
}

// UseWriteDB 实现通过 context 选择写库.
//
// 无匹配 DB 时 panic.
func (p *TransProvider) UseWriteDB(ctx context.Context) *gorm.DB {
	return p.useDB(ctx, true)
}
