package db

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type TestReplicaModel struct {
	ID   int64
	Name string
}

func TestRWOptions_OpenDB_Replicas(t *testing.T) {
	const writeName = "anonymous"

	p := testdb_newprovider_rw(t, "rwoptions")
	ctx := context.Background()

	w := &TestDBModel{ID: testDefaultRecord.ID, Name: writeName}
	if err := p.UseDB(ctx).Updates(w).Error; err != nil {
		t.Fatal(err)
	}

	r := &TestDBModel{}
	p.UseWriteDB(ctx).Where("id = ?", testDefaultRecord.ID).Find(r)
	if r.Name != writeName {
		t.Errorf("expect name: %v, got: %s", writeName, r.Name)
	}

	r = &TestDBModel{}
	p.UseDB(ctx).Clauses(dbresolver.Write).Where("id = ?", testDefaultRecord.ID).Find(r)
	if r.Name != writeName {
		t.Errorf("expect name: %v, got: %s", writeName, r.Name)
	}

	r = &TestDBModel{}
	p.UseDB(ctx).Where("id = ?", testDefaultRecord.ID).Find(r)
	if r.Name != testDefaultRecord.Name {
		t.Errorf("expect name: %v, got: %s", testDefaultRecord.Name, r.Name)
	}
}

func TestRWOptions_OpenDB(t *testing.T) {
	dial := func(opts *Options) (gorm.Dialector, error) {
		if opts.DBName == "" {
			return nil, errors.New("database name not exits")
		}
		return sqlite.Open(filepath.Join(t.TempDir(), opts.DBName)), nil
	}

	cases := []struct {
		name string

		opts  *RWOptions
		conf  *gorm.Config
		error bool
	}{
		{
			name:  "no write and read config and with config",
			opts:  &RWOptions{},
			conf:  &gorm.Config{},
			error: true,
		},
		{
			name:  "no write and read config and without config",
			opts:  &RWOptions{},
			conf:  nil,
			error: true,
		},
		{
			name: "no write config and with config",
			opts: &RWOptions{
				Read: &Options{DBName: "rwoptions_opendb_r_1"},
			},
			conf:  &gorm.Config{},
			error: true,
		},
		{
			name: "no write config and without config",
			opts: &RWOptions{
				Read: &Options{DBName: "rwoptions_opendb_r_2"},
			},
			conf:  nil,
			error: true,
		},
		{
			name: "write error",
			opts: &RWOptions{
				Write: &Options{DBName: ""},
			},
			conf:  &gorm.Config{},
			error: true,
		},
		{
			name: "read error",
			opts: &RWOptions{
				Write: &Options{DBName: "rwoptions_opendb_w_1"},
				Read:  &Options{DBName: ""},
			},
			conf:  &gorm.Config{},
			error: true,
		},
		{
			name: "only write with config",
			opts: &RWOptions{
				Write: &Options{DBName: "rwoptions_opendb_w_2"},
			},
			conf:  &gorm.Config{},
			error: false,
		},
		{
			name: "only write without config",
			opts: &RWOptions{
				Write: &Options{DBName: "rwoptions_opendb_w_3"},
			},
			conf:  nil,
			error: false,
		},
		{
			name: "write and read with config",
			opts: &RWOptions{
				Write: &Options{DBName: "rwoptions_opendb_w_4"},
				Read:  &Options{DBName: "rwoptions_opendb_r_3"},
			},
			conf:  &gorm.Config{},
			error: false,
		},
		{
			name: "write and read without config",
			opts: &RWOptions{
				Write: &Options{DBName: "rwoptions_opendb_w_5"},
				Read:  &Options{DBName: "rwoptions_opendb_r_4"},
			},
			conf:  nil,
			error: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := c.opts.OpenDB(dial, c.conf)
			if c.error && err == nil {
				t.Error("expect error, got nil")
			}
			if !c.error && err != nil {
				t.Errorf("expect not error, got: %v", err)
			}
			if c.error && d != nil {
				t.Errorf("expect database not created, got: %v", d)
			}
			if !c.error && d == nil {
				t.Error("expect database created")
			}
		})
	}
}

func TestRWOptions_OpenDB_Initialize_Error(t *testing.T) {
	dial := func(opts *Options) (gorm.Dialector, error) {
		dl := &sqlite.Dialector{DriverName: "not_exists_driver_name"}
		return dl, nil
	}
	opts := &RWOptions{
		Write: &Options{DBName: "rwoptions_opendb_initialize_error_w"},
		Read:  &Options{DBName: "rwoptions_opendb_initialize_error_r"},
	}
	_, err := opts.OpenDB(dial, &gorm.Config{})
	if err == nil {
		t.Error("expect error, got nil")
	}
}

func TestOptions_OpenDB(t *testing.T) {
	dial := func(opts *Options) (gorm.Dialector, error) {
		if opts.DBName == "" {
			return nil, errors.New("database name not exits")
		}
		return sqlite.Open(filepath.Join(t.TempDir(), opts.DBName)), nil
	}

	cases := []struct {
		name string

		opts  *Options
		conf  *gorm.Config
		error bool
	}{
		{
			name:  "dialect error with config",
			opts:  &Options{DBName: ""},
			conf:  &gorm.Config{},
			error: true,
		},
		{
			name:  "dialect error without config",
			opts:  &Options{DBName: ""},
			conf:  nil,
			error: true,
		},
		{
			name:  "with config",
			opts:  &Options{DBName: "options_opendb_1"},
			conf:  &gorm.Config{},
			error: false,
		},
		{
			name:  "without config",
			opts:  &Options{DBName: "options_opendb_2"},
			conf:  nil,
			error: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d, err := c.opts.OpenDB(dial, c.conf)
			if c.error && err == nil {
				t.Error("expect error, got nil")
			}
			if !c.error && err != nil {
				t.Errorf("expect not error, got: %v", err)
			}
			if c.error && d != nil {
				t.Errorf("expect database not created, got: %v", d)
			}
			if !c.error && d == nil {
				t.Error("expect database created")
			}
		})
	}
}

func TestMultiOptions(t *testing.T) {
	opts := MultiOptions{
		"group1": &Options{DBName: "group1_db"},
		"group2": &Options{DBName: "group2_db"},
	}
	source, err := opts.ToSource(testdb_dial(t), nil, testdb_router)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProvider(source)

	cases := []struct {
		ctx  context.Context
		name string
	}{
		{
			ctx:  testdb_new_context_with_dbname("group1"),
			name: "group1_db",
		},
		{
			ctx:  testdb_new_context_with_dbname("group2"),
			name: "group2_db",
		},
	}

	for _, c := range cases {
		m := &TestDBModel{}
		p.UseDB(c.ctx).Where("id = ?", testDBNameRecordID).Find(m)
		if m.Name != c.name {
			t.Errorf("expect db name: %s, got: %s", c.name, m.Name)
		}
	}
}

func TestMultiRWOptions(t *testing.T) {
	opts := MultiRWOptions{
		"group1": &RWOptions{
			Write: &Options{DBName: "group1_write_db"},
			Read:  &Options{DBName: "group1_read_db"},
		},
		"group2": &RWOptions{
			Write: &Options{DBName: "group2_write_db"},
			Read:  &Options{DBName: "group2_read_db"},
		},
	}
	source, err := opts.ToSource(testdb_dial(t), nil, testdb_router)
	if err != nil {
		t.Fatal(err)
	}
	p := NewProvider(source)

	cases := []struct {
		ctx   context.Context
		name  string
		write bool
	}{
		{
			ctx:  testdb_new_context_with_dbname("group1"),
			name: "group1_read_db",
		},
		{
			ctx:  testdb_new_context_with_dbname("group2"),
			name: "group2_read_db",
		},
		{
			ctx:   testdb_new_context_with_dbname("group1"),
			name:  "group1_write_db",
			write: true,
		},
		{
			ctx:   testdb_new_context_with_dbname("group2"),
			name:  "group2_write_db",
			write: true,
		},
	}

	for _, c := range cases {
		m := &TestDBModel{}
		if c.write {
			p.UseWriteDB(c.ctx).Where("id = ?", testDBNameRecordID).Find(m)
			if m.Name != c.name {
				t.Errorf("expect db name: %s, got: %s", c.name, m.Name)
			}
			p.UseDB(c.ctx).Clauses(dbresolver.Write).Where("id = ?", testDBNameRecordID).Find(m)
		} else {
			p.UseDB(c.ctx).Where("id = ?", testDBNameRecordID).Find(m)
		}
		if m.Name != c.name {
			t.Errorf("expect db name: %s, got: %s", c.name, m.Name)
		}
	}
}
