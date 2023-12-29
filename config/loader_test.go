package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type testConfigs struct {
	Key1 string `yaml:"key1"`
	Key2 string `yaml:"key2"`
	Key3 string `yaml:"key3"`
}

// 仅测试使用，更改动态加载.
func withDynamicLoader(dl func() (string, error)) Option {
	return func(opt *options) {
		opt.dynamicLoader = dl
	}
}

func TestLoad(t *testing.T) {
	temp := t.TempDir()
	const (
		context1 = `key1: value1`
		context2 = `key2: value2`
		context3 = `key3: value3`
	)
	newFile := func(name, context string) string {
		file := filepath.Join(temp, name)
		err := ioutil.WriteFile(file, []byte(context), os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		return file
	}

	cases := []struct {
		opts []Option
		exp  *testConfigs
	}{
		{
			opts: []Option{
				withDynamicLoader(func() (string, error) {
					return newFile("config.load-1-1.yml", context1), nil
				}),
				WithFileLoader(func() ([]string, error) {
					file1 := newFile("config.load-1-2.yml", context2)
					file2 := newFile("config.load-1-3.yml", context3)
					return []string{file1, file2}, nil
				}),
			},
			exp: &testConfigs{
				Key1: "value1",
			},
		},
		{
			opts: []Option{
				withDynamicLoader(func() (string, error) {
					return "", nil
				}),
				WithFileLoader(func() ([]string, error) {
					file1 := newFile("config.load-2-2.yml", context2)
					file2 := newFile("config.load-2-3.yml", context3)
					return []string{file1, file2}, nil
				}),
			},
			exp: &testConfigs{
				Key2: "value2",
				Key3: "value3",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			out := &testConfigs{}
			err := Load(out, c.opts...)
			if err != nil {
				t.Errorf("expect not error, got: %v", err)
			}
			if !reflect.DeepEqual(out, c.exp) {
				t.Errorf("expect: %v, got: %v", c.exp, out)
			}

			defer func() {
				if e := recover(); e != nil {
					t.Errorf("expect not panic, got: %v", e)
				}
			}()
			MustLoad(&testConfigs{}, c.opts...)
		})
	}
}

func TestLoad_Error(t *testing.T) {
	cases := []struct {
		opts []Option
	}{
		{
			opts: []Option{
				withDynamicLoader(func() (string, error) {
					return "", errors.New("load error")
				}),
				WithFileLoader(func() ([]string, error) {
					return nil, errors.New("load error")
				}),
			},
		},
		{
			opts: []Option{
				withDynamicLoader(func() (string, error) {
					return "", errors.New("load error")
				}),
			},
		},
		{
			opts: []Option{
				WithFileLoader(func() ([]string, error) {
					return nil, errors.New("load error")
				}),
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			out := &testConfigs{}
			err := Load(out, c.opts...)
			if err == nil {
				t.Error("expect error, got nil")
			}

			defer func() {
				if e := recover(); e == nil {
					t.Error("expect panic")
				}
			}()
			MustLoad(out, c.opts...)
		})
	}
}
