package db

import (
	"context"
	"testing"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func TestNewMustInjectFromContextScope(t *testing.T) {
	scope := NewMustInjectFromContextScope("tenant", func(ctx context.Context) string {
		if v, ok := ctx.Value("tenantKey").(string); ok {
			return v
		}
		return ""
	})
	p := testdb_newprovider_rw_with_scopes(t, []string{"inject"}, scope)
	db := p.UseDB(context.WithValue(context.Background(), "tenantKey", "username"))

	t.Run("read db", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true}).
			First(&TestDBModel{}).Statement.SQL.String()
		const expect = "SELECT * FROM `test_db_models` WHERE `tenant` = ? ORDER BY `test_db_models`.`id` LIMIT 1"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("write db", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true}).
			Clauses(dbresolver.Write).
			First(&TestDBModel{}).Statement.SQL.String()
		const expect = "SELECT * FROM `test_db_models` WHERE `tenant` = ? ORDER BY `test_db_models`.`id` LIMIT 1"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true}).
			Updates(&TestDBModel{Name: "name"}).Statement.SQL.String()
		const expect = "UPDATE `test_db_models` SET `name`=? WHERE `tenant` = ?"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true}).
			Delete(&TestDBModel{ID: 1}).Statement.SQL.String()
		const expect = "DELETE FROM `test_db_models` WHERE `tenant` = ? AND `test_db_models`.`id` = ?"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("create", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true}).
			Create(&TestDBModel{ID: 1}).Statement.SQL.String()
		const expect = "INSERT INTO `test_db_models` (`tenant`,`name`,`id`) VALUES (?,?,?)"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})
}

func TestNewInjectFromContextScope_EmptyValue(t *testing.T) {
	scope := NewInjectFromContextScope("tenant", func(ctx context.Context) string {
		return ""
	}, true)
	p := testdb_newprovider_rw_with_scopes(t, []string{"inject"}, scope)
	db := p.UseDB(context.Background())

	t.Run("read db", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true, Context: context.Background()}).
			First(&TestDBModel{}).Statement.SQL.String()
		const expect = "SELECT * FROM `test_db_models` ORDER BY `test_db_models`.`id` LIMIT 1"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("write db", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true, Context: context.Background()}).
			Clauses(dbresolver.Write).
			First(&TestDBModel{}).Statement.SQL.String()
		const expect = "SELECT * FROM `test_db_models` ORDER BY `test_db_models`.`id` LIMIT 1"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true, Context: context.Background()}).
			Updates(&TestDBModel{Name: "name"}).Statement.SQL.String()
		const expect = "UPDATE `test_db_models` SET `name`=?"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true, Context: context.Background()}).
			Delete(&TestDBModel{ID: 1}).Statement.SQL.String()
		const expect = "DELETE FROM `test_db_models` WHERE `test_db_models`.`id` = ?"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})

	t.Run("create", func(t *testing.T) {
		t.Parallel()
		sql := db.Session(&gorm.Session{DryRun: true, Context: context.Background()}).
			Create(&TestDBModel{ID: 1}).Statement.SQL.String()
		const expect = "INSERT INTO `test_db_models` (`tenant`,`name`,`id`) VALUES (?,?,?)"
		if sql != expect {
			t.Errorf("expect: %s, got: %s", expect, sql)
		}
	})
}

func TestNewInjectFromContextScope_Required(t *testing.T) {
	scope := NewInjectFromContextScope("tenant", func(ctx context.Context) string {
		return ""
	}, false)
	p := testdb_newprovider_rw_with_scopes(t, []string{"inject"}, scope)

	defer func() {
		if e := recover(); e == nil {
			t.Error("expect panic")
		}
	}()
	p.UseDB(context.Background()).First(&TestDBModel{})
}
