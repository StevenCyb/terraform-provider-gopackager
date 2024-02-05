package hasher

import (
	"archive/zip"
	"os"
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

		t.Run("Plane_Success", func(t *testing.T) {
			t.Parallel()

			content, err := hasher.ReadFile("hasher_test.go")
			assert.NoError(t, err)
			assert.NotNil(t, content)
		})

		t.Run("ZIP_Success", func(t *testing.T) {
			t.Parallel()

			testFileName := "unit_test.zip"

			file, err := os.Create(testFileName)
			assert.NoError(t, err)

			zipWriter := zip.NewWriter(file)

			zipFile, err := zipWriter.Create("test.txt")
			assert.NoError(t, err)

			_, err = zipFile.Write([]byte("test"))
			assert.NoError(t, err)

			zipFile2, err := zipWriter.Create("dir/test2.txt")
			assert.NoError(t, err)

			_, err = zipFile2.Write([]byte("test2"))
			assert.NoError(t, err)

			zipWriter.Close()
			file.Close()

			t.Cleanup(func() {
				os.Remove(testFileName)
			})

			content, err := hasher.ReadFile(testFileName)
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
			SHA512:       "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff",
			SHA256Base64: "OWY4NmQwODE4ODRjN2Q2NTlhMmZlYWEwYzU1YWQwMTVhM2JmNGYxYjJiMGI4MjJjZDE1ZDZjMTViMGYwMGEwOA==",
			SHA512Base64: "ZWUyNmIwZGQ0YWY3ZTc0OWFhMWE4ZWUzYzEwYWU5OTIzZjYxODk4MDc3MmU0NzNmODgxOWE1ZDQ5NDBlMGRiMjdhYzE4NWY4YTBlMWQ1Zjg0Zjg4YmM4ODdmZDY3YjE0MzczMmMzMDRjYzVmYTlhZDhlNmY1N2Y1MDAyOGE4ZmY=",
		}, combined)
	})
}
