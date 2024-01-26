package compiler

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ CompilerI = &Compiler{}
	var _ CompilerI = &MockCompiler{}
}

func TestCompiler(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("binary")
	})

	_, currentFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	conf := NewConfig().
		Source(strings.Replace(currentFile, "internal/compiler/compiler_test.go", "main.go", 1)).
		Destination("binary").
		GOOS("linux").
		GOARCH("amd64")
	assert.NotNil(t, conf)

	compiler := New()
	binaryPath, hash, err := compiler.Compile(*conf)
	assert.NoError(t, err)
	assert.NotEmpty(t, binaryPath)
	assert.NotEmpty(t, hash)
	assert.True(t, strings.HasSuffix(binaryPath, "binary"))
	assert.Regexp(t, `^[0-9a-zA-Z]+$`, hash)
}
