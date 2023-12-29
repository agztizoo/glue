package transaction

import (
	"context"
	"errors"
	"testing"
)

func TestPanicWhenStartTransactionErrorIssue(t *testing.T) {
	tm := NewManager(func(ctx context.Context) interface{} {
		return "ctx_key"
	}, func(ctx context.Context) interface{} {
		return "new_test_db"
	}, func(ctx context.Context, db interface{}, callback func(db interface{}) error) error {
		return errors.New("start transaction error")
	})
	tm.Transaction(context.Background(), func(ctx context.Context) error {
		return nil
	})
}
