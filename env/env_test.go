package env

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestInitFromBazelEnv(t *testing.T) {
	cases := []struct {
		src string
		tws string
		wdr string
	}{
		{
			src: "",
			tws: "",
			wdr: "",
		},
		{
			src: "/abc",
			tws: "",
			wdr: "",
		},
		{
			src: "",
			tws: "/def",
			wdr: "",
		},
		{
			src: "/abc",
			tws: "/def",
			wdr: "/abc/def",
		},
		{
			src: "/abc/",
			tws: "/def/",
			wdr: "/abc/def",
		},
		{
			src: "/abc/",
			tws: "def/",
			wdr: "/abc/def",
		},
		{
			src: "/abc/",
			tws: "def",
			wdr: "/abc/def",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			workDir = ""
			os.Setenv("TEST_SRCDIR", c.src)
			os.Setenv("TEST_WORKSPACE", c.tws)
			initFromBazelEnv()
			if WorkDir() != c.wdr {
				t.Errorf("expect: %s, got: %s", c.wdr, WorkDir())
			}
		})
	}
}

func TestInitFromEnv(t *testing.T) {
	cases := []struct {
		env string
		wdr string
	}{
		{
			env: "",
			wdr: "",
		},
		{
			env: "/abc/def",
			wdr: "/abc/def",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			os.Setenv("WORKPATH", c.env)
			initFromEnv()
			if WorkDir() != c.wdr {
				t.Errorf("expect: %s, got: %s", c.wdr, WorkDir())
			}
		})
	}
}

func TestInitFromGolang(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("XPC_SERVICE_NAME", "com.jetbrains.goland.20308"); err != nil {
		t.Fatal(err)
	}

	initFromGoland()

	if wd != WorkDir() {
		t.Errorf("expect: %s, got: %s", wd, WorkDir())
	}
}

func TestInitFromOSArgs(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		pth string
		wdr string
	}{
		{
			pth: "/abc/123/def/cat",
			wdr: "/abc/123",
		},
		{
			pth: "/abc/123/def/cat.exe",
			wdr: "/abc/123",
		},
		{
			pth: "/abc//123//def/./cat",
			wdr: "/abc/123",
		},
		{
			pth: "/abc//123/../def/./cat",
			wdr: "/abc",
		},
		{
			pth: "abc//123/../def/./cat",
			wdr: filepath.Join(wd, "abc"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i+1), func(t *testing.T) {
			os.Args = []string{c.pth}
			initFromOSArgs()
			if WorkDir() != c.wdr {
				t.Errorf("expect: %s, got: %s", c.wdr, WorkDir())
			}
		})
	}
}
