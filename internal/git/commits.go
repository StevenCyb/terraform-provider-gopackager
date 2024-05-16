package git

import (
	"os/exec"
	"strings"
)

func LastCommitHash() (string, error) {
	var commitID string

	cmd := exec.Command("git", "--no-pager", "log", "-1", "--pretty=format:%H", "--", "*.go", "--", "go.mod", "--", "go.sum")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	commitID = strings.TrimSpace(string(output))

	return commitID, nil
}
