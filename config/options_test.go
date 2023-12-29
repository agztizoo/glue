package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestOptionFiles(t *testing.T) {
	cases := []struct {
		opts  *options
		files []string
	}{
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "", nil
				},
			},
		},
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "config.base.yml", nil
				},
			},
			files: []string{"config.base.yml"},
		},
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "config.base.yml", nil
				},
				fileLoader: func() ([]string, error) {
					return []string{
						"config.base-1.yml",
						"config.base-2.yml",
					}, nil
				},
			},
			files: []string{"config.base.yml"},
		},
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "", nil
				},
				fileLoader: func() ([]string, error) {
					return []string{
						"config.base-1.yml",
						"config.base-2.yml",
					}, nil
				},
			},
			files: []string{
				"config.base-1.yml",
				"config.base-2.yml",
			},
		},
		{
			opts: &options{
				fileLoader: func() ([]string, error) {
					return []string{
						"config.base-1.yml",
						"config.base-2.yml",
					}, nil
				},
			},
			files: []string{
				"config.base-1.yml",
				"config.base-2.yml",
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			files, err := c.opts.Files()
			if err != nil {
				t.Errorf("expect not error, got: %v", err)
			}
			if !reflect.DeepEqual(files, c.files) {
				t.Errorf("expect: %v, got: %v", c.files, files)
			}
		})
	}
}

func TestOptionFiles_Error(t *testing.T) {
	cases := []struct {
		opts *options
	}{
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "", errors.New("load error")
				},
			},
		},
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "", errors.New("load error")
				},
				fileLoader: func() ([]string, error) {
					return []string{
						"config.base-1.yml",
						"config.base-2.yml",
					}, nil
				},
			},
		},
		{
			opts: &options{
				dynamicLoader: func() (string, error) {
					return "", nil
				},
				fileLoader: func() ([]string, error) {
					return nil, errors.New("load error")
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			_, err := c.opts.Files()
			if err == nil {
				t.Error("expect error, got nil")
			}
		})
	}
}

func TestWithFileLoader(t *testing.T) {
	opts := &options{}
	fn := WithFileLoader(func() ([]string, error) {
		return []string{}, nil
	})
	fn(opts)

	if opts.fileLoader == nil {
		t.Error("expect fileLoader not nil")
	}
}
