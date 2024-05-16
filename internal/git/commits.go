package git

import (
	"os/exec"
	"strings"
)

func LastCommitHash() (string, error) {
	var commitID string

	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	commitID = strings.TrimSpace(string(output))

	return commitID, nil
}
