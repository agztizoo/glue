package db

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

type TestHookModel struct {
	ID   int64
	Name string
}

func TestWithInitializeHook(t *testing.T) {
	cases := []struct {
		opts  *Options
		hook  func(*gorm.DB) error
		error bool
	}{
		{
			opts:  &Options{DBName: "hook_1"},
			hook:  func(*gorm.DB) error { return errors.New("hook error") },
			error: true,
		},
		{
			opts:  &Options{DBName: "hook_2"},
			hook:  nil,
			error: false,
		},
		{
			opts:  &Options{DBName: "hook_3"},
			hook:  func(*gorm.DB) error { return nil },
			error: false,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			dl := WithInitializeHook(testdb_dial(t), c.hook)
			d, err := c.opts.OpenDB(dl, nil)
			if c.error && err == nil {
				t.Error("expect error, got nil")
			}
			if !c.error && err != nil {
				t.Errorf("expect not error, got: %+v", err)
			}
			if !c.error && d == nil {
				t.Error("expect database created")
			}
			if c.error {
				return
			}

			d.AutoMigrate(&TestHookModel{})
			const name = "test_name"
			if err := d.Create(&TestHookModel{ID: 1, Name: name}).Error; err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWithInitializeHook_Hooked(t *testing.T) {
	var hooked bool
	dial := WithInitializeHook(testdb_dial(t), func(*gorm.DB) error {
		hooked = true
		return nil
	})

	opts := &Options{DBName: "hook_hooked"}
	d, err := opts.OpenDB(dial, nil)
	if err != nil {
		t.Errorf("expect not error, got: %+v", err)
	}
	if !hooked {
		t.Error("expect hooked")
	}
	d.AutoMigrate(&TestHookModel{})
}

func TestConnPoolHook(t *testing.T) {
	cases := []struct {
		opts *Options

		maxIdleConns    int
		maxOpenConns    int
		connMaxLifeTime time.Duration
	}{
		{
			opts: &Options{DBName: "conn_pool_1"},
		},
		{
			opts:            &Options{DBName: "conn_pool_2"},
			maxIdleConns:    10,
			maxOpenConns:    100,
			connMaxLifeTime: 100 * time.Second,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			hook := ConnPoolHook(c.maxIdleConns, c.maxOpenConns, c.connMaxLifeTime)
			d, err := c.opts.OpenDB(testdb_dial(t), nil)
			if err != nil {
				t.Fatal(err)
			}
			if err := hook(d); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConnPoolHook_Error(t *testing.T) {
	hook := ConnPoolHook(10, 100, 100*time.Second)
	if err := hook(&gorm.DB{Config: &gorm.Config{}}); err == nil {
		t.Error("expect not error")
	}
}
