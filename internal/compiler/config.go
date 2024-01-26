package compiler

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// Error when the source is not set.
	ErrSourceNotSet = errors.New("source not set")
	// Error when the destination is not set.
	ErrDestinationNotSet = errors.New("destination not set")
	// Error when the GOOS is not set.
	ErrGOOSNoSet = errors.New("GOOS not set")
	// Error when the GOARCH is not set.
	ErrGOARCHNoSet = errors.New("GOARCH not set")
	// Error when the source file does not exist.
	ErrSourceFileNotExists = errors.New("source file does not exist")
)

// Configuration for the compiler.
type Config struct {
	source      string
	destination string
	goos        string
	goarch      string
}

// NewConfig creates a new config.
func NewConfig() *Config {
	return &Config{}
}

// Set the path to source.
func (c *Config) Source(path string) *Config {
	c.source = strings.ReplaceAll(path, `"`, "")

	return c
}

// Set the path to destination.
func (c *Config) Destination(name string) *Config {
	c.destination = strings.ReplaceAll(name, `"`, "")

	return c
}

// Set the GOOS.
func (c *Config) GOOS(good string) *Config {
	c.goos = strings.ReplaceAll(good, `"`, "")

	return c
}

// Set the GOARCH.
func (c *Config) GOARCH(goarch string) *Config {
	c.goarch = strings.ReplaceAll(goarch, `"`, "")

	return c
}

// Verifies the config.
func (c *Config) Verify() error {
	switch {
	case c.source == "":
		return ErrSourceNotSet
	case c.destination == "":
		return ErrDestinationNotSet
	case c.goos == "":
		return ErrGOOSNoSet
	case c.goarch == "":
		return ErrGOARCHNoSet
	}

	if state, err := os.Stat(c.source); os.IsNotExist(err) || state.IsDir() {
		return fmt.Errorf("%w: %s", ErrSourceFileNotExists, c.source)
	}

	return nil
}

// Get the `Source` value.
func (c *Config) GetSource() string {
	return c.source
}

// Get the `Destination` value.
func (c *Config) GetDestination() string {
	return c.destination
}

// Get the `GOOS` value.
func (c *Config) GetGOOS() string {
	return c.goos
}

// Get the `GOARCH` value.
func (c *Config) GetGOARCH() string {
	return c.goarch
}
