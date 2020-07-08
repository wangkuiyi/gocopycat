package main

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPkgName(t *testing.T) {
	nameOfGoGetPkg := func(t *testing.T, fullPkg string) {
		cmd := exec.Command("go", "get", "-u", fullPkg)
		cmd.CombinedOutput() // Too many cases that intractable to trace errors.

		assert.Greater(t, len(path.Join(getGoPath(), "src")), 0)

		dir := path.Join(os.Getenv("GOPATH"), "src", fullPkg)
		f, e := fullPkgName(dir)
		assert.NoError(t, e)
		assert.Equal(t, fullPkg, f)
	}

	pkgs := []string{
		"robpike.io/ivy",
		"github.com/golang/go/src/pkg/go/ast",
	}

	for _, pkg := range pkgs {
		t.Run(pkg, func(t *testing.T) {
			nameOfGoGetPkg(t, pkg)
		})
	}
}
