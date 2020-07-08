package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pkgs = []string{
		"robpike.io/ivy",
		"github.com/golang/go/src/pkg/go/ast",
	}
)

func TestFullPkgName(t *testing.T) {
	nameOfGoGetPkg := func(t *testing.T, fullPkg string) {
		cmd := exec.Command("go", "get", "-u", fullPkg)
		// Too many cases that intractable to trace errors.
		cmd.CombinedOutput()

		assert.Greater(t, len(path.Join(getGoPath(), "src")), 0)

		dir := path.Join(os.Getenv("GOPATH"), "src", fullPkg)
		f, e := fullPkgName(dir, filepath.Base(fullPkg))
		assert.NoError(t, e)
		assert.Equal(t, fullPkg, f)
	}

	for _, pkg := range pkgs {
		t.Run(pkg, func(t *testing.T) {
			nameOfGoGetPkg(t, pkg)
		})
	}
}

func TestShortPkgName(t *testing.T) {
	assert.Equal(t, "ivy", shortPkgName(pkgs[0]))
	assert.Equal(t, "ast", shortPkgName(pkgs[1]))
}
