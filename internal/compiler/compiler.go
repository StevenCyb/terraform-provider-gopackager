package compiler

import (
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
	Compile(conf Config) (binaryLocation string, err error)
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
func (c *Compiler) Compile(conf Config) (binaryLocation string, err error) {
	if err := conf.Verify(); err != nil {
		return binaryLocation, err
	} else if conf.source, err = filepath.Abs(conf.source); err != nil {
		return "", fmt.Errorf("unable to get absolute path of source: %w", err)
	} else if conf.destination, err = filepath.Abs(conf.destination); err != nil {
		return "", fmt.Errorf("unable to get absolute path of destination: %w", err)
	}

	if strings.Contains(conf.destination, "/") {
		destinationDirectory := filepath.Dir(conf.destination)

		if info, err := os.Stat(destinationDirectory); os.IsNotExist(err) || !info.IsDir() {
			if err := os.MkdirAll(filepath.Dir(destinationDirectory), 0755); err != nil {
				return "", fmt.Errorf("unable to create destination directory: %w", err)
			}
		}
	}

	args := []string{"build", "-mod=mod", "-o", conf.destination}
	if strings.HasSuffix(conf.source, ".go") {
		args = append(args, filepath.Base(conf.source))
		conf.source = filepath.Dir(conf.source)
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = conf.source
	cmd.Env = append(os.Environ(), "GOOS="+conf.goos, "GOARCH="+conf.goarch)
	if combinedOutput, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf(
			"unable to compile binary: %w, \n\tcommand: %s, \n\toutput: %s",
			err, cmd.String(), string(combinedOutput))
	}

	return conf.destination, nil
}
