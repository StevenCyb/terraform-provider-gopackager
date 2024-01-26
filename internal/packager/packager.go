package packager

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

// ZIPI is an interface for ZIP type.
type ZIPI interface {
	Zip(zipPath string, files map[string]string) (string, error)
}

// Provide ZIP packaging.
type ZIP struct{}

// NewZIP creates a new ZIP instance.
func New() *ZIP {
	return &ZIP{}
}

// Zip given files.
// `files` is a map of file (including path) to the file path inside of the ZIP.
// Returns ZIP file SHA256 hash and an error if any.
func (z ZIP) Zip(zipPath string, files map[string]string) (string, error) {
	if err := os.Remove(zipPath); err != nil && !os.IsNotExist(err) {
		return "", err
	}

	archive, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	defer zipWriter.Close()

	for source, destination := range files {
		err := filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				relativePath, err := filepath.Rel(source, path)
				if err != nil {
					return err
				}

				zipEntryPath := filepath.Join(destination, relativePath)

				f, err := os.Open(path)
				if err != nil {
					return err
				}

				defer f.Close()

				writer, err := zipWriter.Create(zipEntryPath)
				if err != nil {
					return err
				}

				if _, err := io.Copy(writer, f); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return "", err
		}
	}

	hashSHA256 := sha256.New()
	binaryContent, err := os.ReadFile(zipPath)
	if err != nil {
		return "", err
	}

	hashSHA256.Write(binaryContent)
	hash := hex.EncodeToString(hashSHA256.Sum(nil))

	return hash, nil
}
