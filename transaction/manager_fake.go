package transaction

import "context"

func NewFakeManager() Manager {
	f := &fake{}
	return NewManager(f.ctxKeyF, f.lookupDB, f.transaction)
}

type fakeContextKey string
type fake struct{}

// 事务上下文在 context 中存储的 key.
func (f *fake) ctxKeyF(context.Context) interface{} {
	return fakeContextKey("fake_transaction_context")
}

// 实现通过 context 查找 DB, 非事务上下文中 DB.
func (f *fake) lookupDB(context.Context) interface{} {
	return "new_fake_db"
}

// 实现事务执行并通过回调返回新 DB.
func (f *fake) transaction(_ context.Context, db interface{}, callback func(db interface{}) error) error {
	// fake 实现不开启事务.
	return callback(db)
}
