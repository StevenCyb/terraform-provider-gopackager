package compiler

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"os/exec"
)

// CompilerI is an interface for the Compiler type.
type CompilerI interface {
	Compile(conf Config) (binaryLocation string, hash string, err error)
}

// Compiler is a type that implements the CompilerI interface.
// It is used to compile the source code into a binary.
type Compiler struct{}

// Create a new Compiler instance.
func New() *Compiler {
	return &Compiler{}
}

// Compile compiles the source code into a binary.
// It takes a Config instance as parameter and Verify it beforehand.
// It returns the binary location, the hash of the binary and an error if any.
func (c *Compiler) Compile(conf Config) (binaryLocation string, hash string, err error) {
	if err := conf.Verify(); err != nil {
		return binaryLocation, hash, err
	}

	binaryLocation = conf.GetDestination()

	cmd := exec.Command("go", "build", "-o", binaryLocation, conf.source)
	cmd.Env = append(os.Environ(), "GOOS="+conf.goos, "GOARCH="+conf.goarch)

	if err := cmd.Run(); err != nil {
		return binaryLocation, hash, err
	}

	hashSHA256 := sha1.New()
	binaryContent, err := os.ReadFile(binaryLocation)
	if err != nil {
		return binaryLocation, hash, err
	}

	hashSHA256.Write(binaryContent)
	hash = hex.EncodeToString(hashSHA256.Sum(nil))

	return binaryLocation, hash, nil
}
