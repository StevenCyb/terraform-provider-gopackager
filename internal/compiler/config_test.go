package compiler

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	expected := Config{
		source:      "main.go",
		destination: "binary",
		goos:        "linux",
		goarch:      "amd64",
	}

	actual := NewConfig()
	assert.NotNil(t, actual)

	actual = actual.Source(expected.source)
	assert.NotNil(t, actual)

	actual = actual.Destination(expected.destination)
	assert.NotNil(t, actual)

	actual = actual.GOOS(expected.goos)
	assert.NotNil(t, actual)

	actual = actual.GOARCH(expected.goarch)
	assert.NotNil(t, actual)

	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)

	assert.Equal(t, expected.source, actual.GetSource())
	assert.Equal(t, expected.destination, actual.GetDestination())
	assert.Equal(t, expected.goos, actual.GetGOOS())
	assert.Equal(t, expected.goarch, actual.GetGOARCH())
}

func TestConfigVerify(t *testing.T) {
	t.Parallel()

	_, mainFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		c := NewConfig().
			Source(strings.Replace(mainFile, "pkg/compiler/config_test.go", "main.go", 1)).
			Destination("binary").
			GOOS("linux").
			GOARCH("amd64")

		err := c.Verify()
		assert.Nil(t, err)
	})

	t.Run("SourceNotSet", func(t *testing.T) {
		t.Parallel()

		c := NewConfig()
		assert.NotNil(t, c)

		err := c.Verify()
		assert.NotNil(t, err)
		assert.Equal(t, ErrSourceNotSet, err)
	})

	t.Run("DestinationNotSet", func(t *testing.T) {
		t.Parallel()

		c := NewConfig()
		assert.NotNil(t, c)

		c = c.Source(mainFile)
		assert.NotNil(t, c)

		err := c.Verify()
		assert.NotNil(t, err)
		assert.Equal(t, ErrDestinationNotSet, err)
	})

	t.Run("GOOSNoSet", func(t *testing.T) {
		t.Parallel()

		c := NewConfig()
		assert.NotNil(t, c)

		c = c.Source(mainFile)
		assert.NotNil(t, c)

		c = c.Destination("binary")
		assert.NotNil(t, c)

		err := c.Verify()
		assert.NotNil(t, err)
		assert.Equal(t, ErrGOOSNoSet, err)
	})

	t.Run("GOARCHNoSet", func(t *testing.T) {
		t.Parallel()

		c := NewConfig()
		assert.NotNil(t, c)

		c = c.Source(mainFile)
		assert.NotNil(t, c)

		c = c.Destination("binary")
		assert.NotNil(t, c)

		c = c.GOOS("linux")
		assert.NotNil(t, c)

		err := c.Verify()
		assert.NotNil(t, err)
		assert.Equal(t, ErrGOARCHNoSet, err)
	})

	t.Run("SourceNotExists", func(t *testing.T) {
		t.Parallel()

		c := NewConfig()
		assert.NotNil(t, c)

		c = c.Source("not_exists")
		assert.NotNil(t, c)

		c = c.Destination("binary")
		assert.NotNil(t, c)

		c = c.GOOS("linux")
		assert.NotNil(t, c)

		c = c.GOARCH("amd64")
		assert.NotNil(t, c)

		err := c.Verify()
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrSourceFileNotExists)
	})
}
