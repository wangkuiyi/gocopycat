package main

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullPkgName(t *testing.T) {
	const fullPkg = "robpike.io/ivy"
	cmd := exec.Command("go", "get", "-u", fullPkg)
	b, e := cmd.CombinedOutput()
	assert.NoError(t, e, "Failed to run go get -u due to %s", b)

	assert.Greater(t, len(path.Join(getGoPath(), "src")), 0)

	dir := path.Join(os.Getenv("GOPATH"), "src", fullPkg)
	f, e := fullPkgName(dir)
	assert.NoError(t, e)
	assert.Equal(t, fullPkg, f)
}
