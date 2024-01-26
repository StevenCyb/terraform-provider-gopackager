package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrUnableToGetWorkingDirectory is an error returned when the current working directory cannot be retrieved.
var ErrUnableToGetWorkingDirectory = errors.New("unable to get current working directory")

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
// It returns the binary location, the SHA256 hash of the binary and an error if any.
func (c *Compiler) Compile(conf Config) (binaryLocation string, hash string, err error) {
	if err := conf.Verify(); err != nil {
		return binaryLocation, hash, err
	}

	var cmd *exec.Cmd

	workingDir, err := os.Getwd()
	if err != nil {
		return "", "", ErrUnableToGetWorkingDirectory
	}

	binaryLocation = filepath.Join(workingDir, conf.GetDestination())
	workingDir = filepath.Dir(conf.source)

	if strings.HasSuffix(conf.source, ".go") {
		cmd = exec.Command("go", "build", "-mod=mod", "-o", binaryLocation, conf.source)
	} else {
		cmd = exec.Command("go", "build", "-mod=mod", "-o", binaryLocation)
	}

	cmd.Dir = workingDir
	cmd.Env = append(os.Environ(), "GOOS="+conf.goos, "GOARCH="+conf.goarch)

	if combinedOutput, err := cmd.CombinedOutput(); err != nil {
		return binaryLocation, hash, fmt.Errorf(
			"unable to compile binary: %w, \n\tcommand: %s, \n\toutput: %s",
			err, cmd.String(), string(combinedOutput))
	}

	hashSHA256 := sha256.New()
	binaryContent, err := os.ReadFile(binaryLocation)
	if err != nil {
		return binaryLocation, hash, err
	}

	hashSHA256.Write(binaryContent)
	hash = hex.EncodeToString(hashSHA256.Sum(nil))

	return binaryLocation, hash, nil
}
