package transaction

import (
	"context"
	"testing"
)

func TestContextKey(t *testing.T) {
	type (
		typea string
		typeb string
	)
	var (
		keya interface{} = typea("key")
		keyb interface{} = typeb("key")
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, keya, 123)
	vala := ctx.Value(keya).(int)
	valb, ok := ctx.Value(keyb).(int)
	if ok {
		t.Error("failed")
	}
	if vala == valb {
		t.Fail()
	}
}
