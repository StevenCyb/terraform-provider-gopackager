package compiler

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccInterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ CompilerI = &Compiler{}
	var _ CompilerI = &MockCompiler{}
}

func TestAccCompiler(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("binary")
		os.Remove("binary2")
	})

	t.Run("RootPath", func(t *testing.T) {
		t.Parallel()

		conf := NewConfig().
			Source("../../").
			Destination("binary").
			GOOS("linux").
			GOARCH("amd64")
		assert.NotNil(t, conf)

		compiler := New()
		binaryPath, err := compiler.Compile(*conf)
		assert.NoError(t, err)
		assert.NotEmpty(t, binaryPath)
		assert.True(t, strings.HasSuffix(binaryPath, "binary"))
	})

	t.Run("Main", func(t *testing.T) {
		t.Parallel()

		conf := NewConfig().
			Source("../../main.go").
			Destination("binary2").
			GOOS("linux").
			GOARCH("amd64")
		assert.NotNil(t, conf)

		compiler := New()
		binaryPath, err := compiler.Compile(*conf)
		assert.NoError(t, err)
		assert.NotEmpty(t, binaryPath)
		assert.True(t, strings.HasSuffix(binaryPath, "binary2"))
	})
}
