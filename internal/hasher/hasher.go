package hasher

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// HasherI is the interface for Hasher.
type HasherI interface {
	ReadFile(path string) ([]byte, error)
	MD5(binaryContent []byte) string
	SHA1(binaryContent []byte) string
	SHA256(binaryContent []byte) string
	SHA512(binaryContent []byte) string
	SHA256Base64(binaryContent []byte) string
	SHA512Base64(binaryContent []byte) string
	CombinedHash(binaryContent []byte) CombinedHash
}

// CombinedHash is a struct for the combined hash.
type CombinedHash struct {
	MD5          string
	SHA1         string
	SHA256       string
	SHA512       string
	SHA256Base64 string
	SHA512Base64 string
}

// Hasher is a type for hashing files.
type Hasher struct{}

// New creates a new Hasher instance.
func New() *Hasher {
	return &Hasher{}
}

// ReadFile reads a file from the filesystem.
// If file is ZIP file, it reads the content of all files in the ZIP file.
func (h *Hasher) ReadFile(path string) ([]byte, error) {
	if strings.HasSuffix(path, ".zip") {
		return h.readZipFile(path)
	}

	return h.readPlaneFile(path)
}

// reads the file content from a zip file.
func (h *Hasher) readZipFile(path string) ([]byte, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	defer archive.Close()

	var binaryContent []byte

	for _, file := range archive.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, err
		}

		defer fileReader.Close()

		fileContent, err := io.ReadAll(fileReader)
		if err != nil {
			return nil, err
		}

		fileContent = append([]byte("__"+file.Name+"__"), fileContent...)
		binaryContent = append(binaryContent, fileContent...)
	}

	return binaryContent, nil
}

// reads a file content.
func (h *Hasher) readPlaneFile(path string) ([]byte, error) {
	binaryContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return binaryContent, nil
}

// MD5 hashes the binary content with MD5.
func (h *Hasher) MD5(binaryContent []byte) string {
	hashMD5 := md5.New()
	hashMD5.Write(binaryContent)
	hash := hex.EncodeToString(hashMD5.Sum(nil))

	return hash
}

// SHA1 hashes the binary content with SHA1.
func (h *Hasher) SHA1(binaryContent []byte) string {
	hashSHA1 := sha1.New()
	hashSHA1.Write(binaryContent)
	hash := hex.EncodeToString(hashSHA1.Sum(nil))

	return hash
}

// SHA256 hashes the binary content with SHA256.
func (h *Hasher) SHA256(binaryContent []byte) string {
	hashSHA256 := sha256.New()
	hashSHA256.Write(binaryContent)
	hash := hex.EncodeToString(hashSHA256.Sum(nil))

	return hash
}

// SHA512 hashes the binary content with SHA512.
func (h *Hasher) SHA512(binaryContent []byte) string {
	hashSHA512 := sha512.New()
	hashSHA512.Write(binaryContent)
	hash := hex.EncodeToString(hashSHA512.Sum(nil))

	return hash
}

// SHA256Base64 hashes the binary content with SHA256 and encodes it with base64.
func (h *Hasher) SHA256Base64(binaryContent []byte) string {
	hash := h.SHA256(binaryContent)
	hash = base64.StdEncoding.EncodeToString([]byte(hash))

	return hash
}

// SHA512Base64 hashes the binary content with SHA512 and encodes it with base64.
func (h *Hasher) SHA512Base64(binaryContent []byte) string {
	hash := h.SHA512(binaryContent)
	hash = base64.StdEncoding.EncodeToString([]byte(hash))

	return hash
}

// CombinedHash hashes the binary content with all available algorithms.
func (h *Hasher) CombinedHash(binaryContent []byte) CombinedHash {
	return CombinedHash{
		MD5:          h.MD5(binaryContent),
		SHA1:         h.SHA1(binaryContent),
		SHA256:       h.SHA256(binaryContent),
		SHA512:       h.SHA512(binaryContent),
		SHA256Base64: h.SHA256Base64(binaryContent),
		SHA512Base64: h.SHA512Base64(binaryContent),
	}
}
