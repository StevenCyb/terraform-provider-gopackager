package packager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccInterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ ZIPI = &ZIP{}
	var _ ZIPI = &MockZIP{}
}

func TestAccZIPZip(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("test.zip")
	})

	zip := ZIP{}
	err := zip.Zip("test.zip", map[string]string{
		"packager.go":      "packager.go",
		"packager_mock.go": "a/packager_mock.go",
		"packager_test.go": "b/packager_test.go",
		"../packager":      "c/packager",
	})

	assert.NoError(t, err)
}
