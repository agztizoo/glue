package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glue/env"
)

var (
	// 动态规则配置文件格式.
	DynamicFilePattern = "conf/config.%s.yml"
	// 环境感知配置文件格式.
	EnvAwareFilePattern = "conf/config.%s.yml"

	// 基础配置文件.
	//
	// 链接阶段设置基础配置文件, 用于支持分片测试.
	//
	// 设置方式:
	//  1. go build / IDE:
	//   -ldflags="-X github.com/.../config.base=conf/config.base.yml"
	//   -ldflags="-X github.com/.../config.base=/User/username/xxx/conf/config.base.yml"
	//  2. bazel 设置:
	//   - 构建目标 go_test 或 go_binary 添加
	//     x_defs = {"github.com/.../config.base": "conf/config.base.yml"}
	//
	// Tips: base 链接阶段使用的字段, 不要修改命名.
	base string
)

// 用于测试插桩.
var (
	fileStateFunc    = os.Stat
	shardProfileFunc = env.ShardProfile
	profileFunc      = env.Profile
	workDirFunc      = env.WorkDir

	envFunc = func() string {
		if e := strings.ToLower(os.Getenv("ENV")); e != "" {
			return e
		}
		return "dev"
	}
)

// dynamicFile 基于动态规则返回配置文件.
func dynamicFile() (string, error) {
	files := []string{
		fmt.Sprintf(DynamicFilePattern, shardProfileFunc()),
		fmt.Sprintf(DynamicFilePattern, profileFunc()),
		base,
	}
	for _, file := range files {
		if file == "" {
			continue
		}
		if !filepath.IsAbs(file) {
			file = filepath.Join(workDirFunc(), file)
		}
		if _, err := fileStateFunc(file); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		return file, nil
	}
	return "", nil
}

// EnvAwareFile 返回环境感知配置文件加载函数.
//
// 例:
//   当 pattern 为 conf/config.yml 时, 加载 conf/config.+env.GetEnv()+".yml"
func EnvAwareFile(patterns ...string) func() ([]string, error) {
	return func() ([]string, error) {
		e := envFunc()
		files := make([]string, 0, len(patterns))
		for _, p := range patterns {
			file := fmt.Sprintf(p, e)
			if !filepath.IsAbs(file) {
				file = filepath.Join(workDirFunc(), file)
			}
			if _, err := fileStateFunc(file); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, err
			}
			files = append(files, file)
		}
		return files, nil
	}
}
