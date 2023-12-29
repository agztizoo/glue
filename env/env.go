package env

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	workDir string
)

func init() {
	for _, f := range []func(){
		initFromBazelEnv,
		initFromEnv,
		initFromGoland,
		initFromOSArgs,
	} {
		f()
		if workDir != "" {
			return
		}
	}
}

func initFromBazelEnv() {
	var (
		src = os.Getenv("TEST_SRCDIR")
		ws  = os.Getenv("TEST_WORKSPACE")
	)
	if src != "" && ws != "" {
		workDir = filepath.Join(src, ws)
	}
}

func initFromEnv() {
	workDir = os.Getenv("WORKPATH")
}

func initFromGoland() {
	x := os.Getenv("XPC_SERVICE_NAME")
	if strings.HasPrefix(x, "com.jetbrains.goland.") {
		ws, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		workDir = ws
	}
}

// 支持如下形式的目录结构.
//
// .
// ├── bin
// │   └── program
// └── conf
//     ├── config.yml
func initFromOSArgs() {
	workDir, _ = filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
}

// WorkDir 返回当前工作目录.
//
// 通过函数的方式返回，预防运行时修改.
func WorkDir() string {
	return workDir
}
