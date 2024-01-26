package packager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ ZIPI = &ZIP{}
	var _ ZIPI = &MockZIP{}
}

func TestZIPZip(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("test.zip")
	})

	zip := ZIP{}
	hash, err := zip.Zip("test.zip", map[string]string{
		"packager.go":      "packager.go",
		"packager_mock.go": "a/packager_mock.go",
		"packager_test.go": "b/packager_test.go",
		"../packager":      "c/packager",
	})

	assert.NoError(t, err)
	assert.Regexp(t, `^[0-9a-zA-Z]+$`, hash)
}
