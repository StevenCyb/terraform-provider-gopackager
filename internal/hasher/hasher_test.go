package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccInterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var _ HasherI = &Hasher{}
	var _ HasherI = &MockHasher{}
}

func TestAccHasher(t *testing.T) {
	t.Parallel()

	hasher := New()

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			content, err := hasher.ReadFile("hasher_test.go")
			assert.NoError(t, err)
			assert.NotNil(t, content)
		})

		t.Run("Failure", func(t *testing.T) {
			t.Parallel()

			_, err := hasher.ReadFile("does_not_exist.txt")
			assert.Error(t, err)
		})
	})

	t.Run("CombinedHash", func(t *testing.T) {
		t.Parallel()

		content := []byte("test")
		combined := hasher.CombinedHash(content)
		assert.Equal(t, CombinedHash{
			MD5:          "098f6bcd4621d373cade4e832627b4f6",
			SHA1:         "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
			SHA256:       "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			SHA512:       "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			SHA256Base64: "OWY4NmQwODE4ODRjN2Q2NTlhMmZlYWEwYzU1YWQwMTVhM2JmNGYxYjJiMGI4MjJjZDE1ZDZjMTViMGYwMGEwOA==",
			SHA512Base64: "OWY4NmQwODE4ODRjN2Q2NTlhMmZlYWEwYzU1YWQwMTVhM2JmNGYxYjJiMGI4MjJjZDE1ZDZjMTViMGYwMGEwOA==",
		}, combined)
	})
}
