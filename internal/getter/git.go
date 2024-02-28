package getter

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"os/exec"
	"path"
)

// GitI is an interface for Git puller.
type GitI interface {
	Get(destination, repoURL, branch string) (string, error)
}

// Git is a struct for Git puller.
type Git struct{}

// New creates a new Git puller instance.
func New() GitI {
	return &Git{}
}

// Clone clones the repository to the destination.
func (g *Git) Get(destination, repoURL, branch string) (string, error) {
	key := g.generateKey(repoURL, branch)
	destination = path.Join(destination, key)

	if fileInfo, err := os.Stat(destination); err == nil && fileInfo.IsDir() {
		return destination, g.pull(destination)
	}

	return destination, g.clone(destination, repoURL, branch)
}

// pull pulls the repository to the destination.
func (g *Git) pull(destination string) error {
	cmd := exec.Command("git", "-C", destination, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// clone clones the repository to the destination.
func (g *Git) clone(destination, repoURL, branch string) error {
	var cmd *exec.Cmd

	if branch == "" {
		cmd = exec.Command("git", "clone", repoURL, destination)
	} else {
		cmd = exec.Command("git", "clone", "-b", branch, repoURL, destination)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// generateKey generates a key for the repository.
func (g *Git) generateKey(repoURL, branch string) string {
	hashSHA1 := sha1.New()
	hashSHA1.Write([]byte(repoURL + "," + branch))
	key := hex.EncodeToString(hashSHA1.Sum(nil))

	return key
}
