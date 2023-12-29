package db

import (
	"context"
	"testing"

	"gorm.io/gorm"
)

func TestRouteWithKey(t *testing.T) {
	dbs := map[string]*gorm.DB{
		"db": new(gorm.DB),
	}
	router := RouteWithKey(dbs, func(ctx context.Context) string {
		db, ok := ctx.Value("db").(string)
		if !ok {
			return ""
		}
		return db
	})

	t.Run("db", func(t *testing.T) {
		const key = "db"
		ctx := context.WithValue(context.Background(), "db", key)
		db := router(ctx)
		if db == nil {
			t.Error("expect get not nil db")
		}
	})

	t.Run("not exist key", func(t *testing.T) {
		const key = "not exists key"
		ctx := context.WithValue(context.Background(), "db", key)
		db := router(ctx)
		if db != nil {
			t.Errorf("expect get nil db, got: %v", db)
		}
	})
}
