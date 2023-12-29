package db

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestDBModel struct {
	ID     int64
	Tenant string
	Name   string
}

// 测试数据库预置记录.
var testDefaultRecord = &TestDBModel{ID: 3721, Tenant: "default tenant", Name: "default record name"}
var testDBNameRecordID = 3722

// DB 测试路由 Key.
const testdb_routekey = "test_db_route_key"

// 设置路由数据库名.
func testdb_context_with_dbname(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, testdb_routekey, name)
}

// 创建 context 并设置路由数据库名.
func testdb_new_context_with_dbname(name string) context.Context {
	return testdb_context_with_dbname(context.Background(), name)
}

// 测试公用 DB 路由函数.
func testdb_router(ctx context.Context) string {
	name, ok := ctx.Value(testdb_routekey).(string)
	if !ok {
		return ""
	}
	return name
}

// 前缀匹配路由.
func testdb_router_prefix(ctx context.Context) string {
	name := testdb_router(ctx)
	return strings.Split(name, "_")[0]
}

// 空值默认路由.
func testdb_router_with_default(router func(context.Context) string, fallback string) func(context.Context) string {
	return func(ctx context.Context) string {
		name := router(ctx)
		if name == "" {
			return fallback
		}
		return name
	}
}

func testdb_dial(t *testing.T) func(*Options) (gorm.Dialector, error) {
	return func(opts *Options) (gorm.Dialector, error) {
		dsn := filepath.Join(t.TempDir(), fmt.Sprintf("%s.db", opts.DBName))
		dl := sqlite.Open(dsn)
		db, err := gorm.Open(dl, &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&TestDBModel{})
		db.Create(testDefaultRecord)
		db.Create(&TestDBModel{ID: int64(testDBNameRecordID), Tenant: opts.DBName, Name: opts.DBName})
		return dl, nil
	}
}

// 创建指定命名的测试数据库.
func testdb_newdbs(t *testing.T, names ...string) map[string]*gorm.DB {
	dbs := make(map[string]*gorm.DB)
	for _, name := range names {
		dsn := filepath.Join(t.TempDir(), fmt.Sprintf("%s.db", name))
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatal(err)
		}
		db.AutoMigrate(&TestDBModel{})
		db.Create(testDefaultRecord)
		db.Create(&TestDBModel{ID: int64(testDBNameRecordID), Tenant: name, Name: name})
		dbs[name] = db
	}
	return dbs
}

// 创建测试数据源.
func testdb_newsource(t *testing.T, names ...string) Source {
	dbs := testdb_newdbs(t, names...)
	if len(dbs) == 1 {
		name := names[0]
		db := dbs[name]
		return NewSource(name, db)
	}
	return NewSourceWithFunc(testdb_router, RouteWithKey(dbs, testdb_router))
}

// 创建测试 Provider.
func testdb_newprovider(t *testing.T, dbnames ...string) *TransProvider {
	source := testdb_newsource(t, dbnames...)
	return NewProvider(source)
}

func testdb_newprovider_with_dial(t *testing.T, dial Dialector, dbnames ...string) *TransProvider {
	return testdb_newprovider_with_scopes_dial(t, dial, dbnames)
}

func testdb_newprovider_with_scopes(t *testing.T, dbnames []string, scopes ...func(*gorm.DB) *gorm.DB) *TransProvider {
	return testdb_newprovider_with_scopes_dial(t, testdb_dial(t), dbnames, scopes...)
}

func testdb_newprovider_with_scopes_dial(t *testing.T, dial Dialector, dbnames []string, scopes ...func(*gorm.DB) *gorm.DB) *TransProvider {
	if len(dbnames) == 1 {
		opts := &Options{DBName: dbnames[0]}
		source, err := opts.ToSource(dial, nil)
		if err != nil {
			t.Fatal(err)
		}
		return NewProvider(source, scopes...)
	}

	opts := make(MultiOptions)
	for _, name := range dbnames {
		opts[name] = &Options{DBName: name}
	}
	source, err := opts.ToSource(testdb_dial(t), nil, testdb_router)
	if err != nil {
		t.Fatal(err)
	}
	return NewProvider(source, scopes...)
}

func testdb_newprovider_rw(t *testing.T, dbnames ...string) *TransProvider {
	return testdb_newprovider_rw_with_scopes(t, dbnames)
}

func testdb_newprovider_rw_with_scopes(t *testing.T, dbnames []string, scopes ...func(*gorm.DB) *gorm.DB) *TransProvider {
	return testdb_newprovider_rw_with_scopes_dial(t, testdb_dial(t), dbnames, scopes...)
}

// 创建测试读写 Provider.
func testdb_newprovider_rw_with_scopes_dial(t *testing.T, dial Dialector, dbnames []string, scopes ...func(*gorm.DB) *gorm.DB) *TransProvider {
	if len(dbnames) == 1 {
		opts := &RWOptions{
			Write: &Options{DBName: fmt.Sprintf("%s_write", dbnames[0])},
			Read:  &Options{DBName: fmt.Sprintf("%s_read", dbnames[0])},
		}
		source, err := opts.ToSource(dial, nil)
		if err != nil {
			t.Fatal(err)
		}
		return NewProvider(source, scopes...)
	}

	opts := make(MultiRWOptions)
	for _, name := range dbnames {
		opts[name] = &RWOptions{
			Write: &Options{DBName: fmt.Sprintf("%s_write", name)},
			Read:  &Options{DBName: fmt.Sprintf("%s_read", name)},
		}
	}
	source, err := opts.ToSource(testdb_dial(t), nil, testdb_router)
	if err != nil {
		t.Fatal(err)
	}
	return NewProvider(source, scopes...)
}
