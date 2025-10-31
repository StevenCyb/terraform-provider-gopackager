package hasher

import (
	"archive/zip"
	"os"
	"path/filepath"
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
			SHA256Base64: "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg=",
			SHA512Base64: "7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w==",
		}, combined)
	})

	t.Run("HashDir", func(t *testing.T) {
		t.Parallel()

		t.Run("Success", func(t *testing.T) {
			t.Parallel()

			// Create a temporary directory for testing
			tempDir := t.TempDir()

			// Create test files and directories
			testFiles := map[string]string{
				"file1.txt":        "content of file 1",
				"file2.txt":        "content of file 2",
				"subdir/file3.txt": "content of file 3",
				"subdir/file4.txt": "content of file 4",
			}

			for relPath, content := range testFiles {
				fullPath := filepath.Join(tempDir, relPath)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				assert.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(content), 0644)
				assert.NoError(t, err)
			}

			// Test HashDir
			result, err := hasher.HashDir(tempDir)
			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Verify all hash fields are populated and non-empty
			assert.NotEmpty(t, result.MD5)
			assert.NotEmpty(t, result.SHA1)
			assert.NotEmpty(t, result.SHA256)
			assert.NotEmpty(t, result.SHA512)
			assert.NotEmpty(t, result.SHA256Base64)
			assert.NotEmpty(t, result.SHA512Base64)

			// Test deterministic behavior - hashing the same directory twice should give same result
			result2, err := hasher.HashDir(tempDir)
			assert.NoError(t, err)
			assert.Equal(t, result, result2)
		})

		t.Run("Directory_Order_Independence", func(t *testing.T) {
			t.Parallel()

			// Create two temporary directories
			tempDir1 := t.TempDir()
			tempDir2 := t.TempDir()

			// Create the same files in different order
			testFiles := []struct {
				path    string
				content string
			}{
				{"file1.txt", "content 1"},
				{"file2.txt", "content 2"},
				{"subdir/file3.txt", "content 3"},
			}

			// Create files in tempDir1 in order
			for _, file := range testFiles {
				fullPath := filepath.Join(tempDir1, file.path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				assert.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(file.content), 0644)
				assert.NoError(t, err)
			}

			// Create files in tempDir2 in reverse order
			for i := len(testFiles) - 1; i >= 0; i-- {
				file := testFiles[i]
				fullPath := filepath.Join(tempDir2, file.path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				assert.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(file.content), 0644)
				assert.NoError(t, err)
			}

			// Both directories should produce the same hash
			hash1, err := hasher.HashDir(tempDir1)
			assert.NoError(t, err)

			hash2, err := hasher.HashDir(tempDir2)
			assert.NoError(t, err)

			assert.Equal(t, hash1, hash2)
		})

		t.Run("Empty_Directory", func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			result, err := hasher.HashDir(tempDir)
			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Even empty directory should have valid hashes
			assert.NotEmpty(t, result.MD5)
			assert.NotEmpty(t, result.SHA256)
		})

		t.Run("Nonexistent_Directory", func(t *testing.T) {
			t.Parallel()

			_, err := hasher.HashDir("/nonexistent/directory")
			assert.Error(t, err)
		})

		t.Run("With_Symlinks", func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			// Create a regular file
			regularFile := filepath.Join(tempDir, "regular.txt")
			err := os.WriteFile(regularFile, []byte("regular content"), 0644)
			assert.NoError(t, err)

			// Create a symlink
			symlinkPath := filepath.Join(tempDir, "symlink.txt")
			err = os.Symlink("regular.txt", symlinkPath)
			assert.NoError(t, err)

			result, err := hasher.HashDir(tempDir)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.SHA256)
		})
	})
}
