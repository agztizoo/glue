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

// envReset 用于恢复测试插桩.
func envReset() func() {
	tempWorkDirFunc := workDirFunc
	tempShardProfileFunc := shardProfileFunc
	tempProfileFunc := profileFunc
	tempBase := base
	tempEnvFunc := envFunc
	tempFileStateFunc := fileStateFunc
	return func() {
		workDirFunc = tempWorkDirFunc
		shardProfileFunc = tempShardProfileFunc
		profileFunc = tempProfileFunc
		base = tempBase
		envFunc = tempEnvFunc
		fileStateFunc = tempFileStateFunc
	}
}

func TestDynamicFile(t *testing.T) {
	defer envReset()()
	temp := t.TempDir()
	os.Mkdir(filepath.Join(temp, "conf"), os.ModePerm)

	cases := []struct {
		init    func(t *testing.T)
		shard   string
		profile string
		base    string

		file string
	}{
		{},
		{
			shard: "ci-shard",
		},
		{
			profile: "ci",
		},
		{
			base: "not_exists.yml",
		},
		{
			shard:   "ci-shard-1",
			profile: "ci-1",
			base:    "conf/config.base-1.yml",
			init: func(t *testing.T) {
				file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "ci-shard-1"))
				err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			},
			file: filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "ci-shard-1")),
		},
		{
			shard:   "ci-shard-2",
			profile: "ci-2",
			base:    "conf/config.base-2.yml",
			init: func(t *testing.T) {
				file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "ci-2"))
				err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			},
			file: filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "ci-2")),
		},
		{
			shard:   "ci-shard-3",
			profile: "ci-3",
			base:    "conf/config.base-3.yml",
			init: func(t *testing.T) {
				file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "base-3"))
				err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
				if err != nil {
					t.Fatal(err)
				}
			},
			file: filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "base-3")),
		},
		{
			shard:   "",
			profile: "ci-4",
			base:    "conf/config.base-4.yml",
			init: func(t *testing.T) {
				for _, profile := range []string{"ci-shard-4", "ci-4", "base-4"} {
					file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, profile))
					err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
					if err != nil {
						t.Fatal(err)
					}
				}
			},
			file: filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "ci-4")),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			shardProfileFunc = func() string {
				return c.shard
			}
			profileFunc = func() string {
				return c.profile
			}
			workDirFunc = func() string {
				return temp
			}
			base = c.base
			if c.init != nil {
				c.init(t)
			}
			file, err := dynamicFile()
			if err != nil {
				t.Errorf("expect not err, got: %v", err)
			}
			if file != c.file {
				t.Errorf("expect file: %s, got: %s", c.file, file)
			}
		})
	}
}

func TestDynamicFile_StatsError(t *testing.T) {
	defer envReset()()
	temp := t.TempDir()
	os.Mkdir(filepath.Join(temp, "conf"), os.ModePerm)

	base = "conf/config.base-stats-error.yml"
	file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "base-stats-error"))
	err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	workDirFunc = func() string {
		return temp
	}
	fileStateFunc = func(file string) (os.FileInfo, error) {
		return nil, errors.New("stat file error")
	}
	if _, err := dynamicFile(); err == nil {
		t.Errorf("expect error, got nil")
	}
}

func TestEnvAwareFile(t *testing.T) {
	defer envReset()()
	temp := t.TempDir()
	os.Mkdir(filepath.Join(temp, "conf"), os.ModePerm)
	file := filepath.Join(temp, fmt.Sprintf(DynamicFilePattern, "test"))
	err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		patterns []string
		files    []string
		env      string
		tce      bool
	}{
		{
			files: []string{},
		},
		{
			patterns: []string{"conf/config.aware.%s.yml"},
			files:    []string{},
		},
		{
			patterns: []string{"conf/config.not_exists.%s.yml"},
			files:    []string{},
		},
		{
			env:      "test",
			tce:      true,
			patterns: []string{DynamicFilePattern},
			files: []string{
				file,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			workDirFunc = func() string {
				return temp
			}
			envFunc = func() string {
				return c.env
			}
			files, err := EnvAwareFile(c.patterns...)()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(files, c.files) {
				t.Errorf("expect: %v, got: %v", c.files, files)
			}
		})
	}
}

func TestEnvAwareFile_StatsError(t *testing.T) {
	defer envReset()()
	temp := t.TempDir()
	os.Mkdir(filepath.Join(temp, "conf"), os.ModePerm)

	base = "conf/config.base-stats-error.yml"
	file := filepath.Join(temp, fmt.Sprintf(EnvAwareFilePattern, "base-stats-error"))
	err := ioutil.WriteFile(file, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	envFunc = func() string {
		return "base-stats-error"
	}
	workDirFunc = func() string {
		return temp
	}
	fileStateFunc = func(file string) (os.FileInfo, error) {
		return nil, errors.New("stat file error")
	}
	if _, err := EnvAwareFile(EnvAwareFilePattern)(); err == nil {
		t.Errorf("expect error, got nil")
	}
}

func TestEnvAwareFile_EnvError(t *testing.T) {
	defer envReset()()

	envFunc = func() string {
		return "dev"
	}
	if _, err := EnvAwareFile(EnvAwareFilePattern)(); err == nil {
		t.Errorf("expect error, got nil")
	}
}
