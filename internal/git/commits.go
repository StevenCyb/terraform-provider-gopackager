package git

import (
	"os"
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

// LastCommitHashForPath returns the last commit hash that modified files in the specified path.
func LastCommitHashForPath(path string) (string, error) {
	var commitID string

	cmd := exec.Command("git", "--no-pager", "log", "-1", "--pretty=format:%H", "--", path)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	commitID = strings.TrimSpace(string(output))

	return commitID, nil
}

// GetTriggeringCommitHash returns the commit hash that would trigger a rebuild.
// This is the last commit that modified files in the specified path since the last compilation.
func GetTriggeringCommitHash(basePath, sinceCommit string) (string, error) {
	if sinceCommit == "" {
		// If no previous commit, return the latest commit for the path
		return LastCommitHashForPath(basePath)
	}

	// Get all commits that modified the path since the given commit
	cmd := exec.Command("git", "--no-pager", "log", "--pretty=format:%H", sinceCommit+"..HEAD", "--", basePath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	commits := strings.TrimSpace(string(output))
	if commits == "" {
		// No commits since the last compilation, return the since commit
		return sinceCommit, nil
	}

	// Return the oldest commit in the range (first one that triggered the change)
	commitLines := strings.Split(commits, "\n")
	if len(commitLines) > 0 {
		// Return the last commit in the list (oldest in time, first that made changes)
		return commitLines[len(commitLines)-1], nil
	}

	return sinceCommit, nil
}

// HasChangedSinceCommit checks if any files in the specified path have changed since the given commit.
func HasChangedSinceCommit(basePath, sinceCommit string) (bool, error) {
	if sinceCommit == "" {
		// If no previous commit, consider it changed
		return true, nil
	}

	// Use git diff to check if there are any changes in the path between the given commit and HEAD
	cmd := exec.Command("git", "diff", "--quiet", sinceCommit, "HEAD", "--", basePath)
	err := cmd.Run()

	if err != nil {
		// git diff --quiet returns non-zero exit code when there are differences
		// Check if it's because there are actual differences or an error
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means there are differences
			if exitError.ExitCode() == 1 {
				return true, nil
			}
		}
		// Other errors (like invalid commit hash)
		return false, err
	}

	// Exit code 0 means no differences
	return false, nil
}

// GetLastCompilationCommit retrieves the commit hash when the target was last compiled.
// This would typically be stored in a metadata file or state.
func GetLastCompilationCommit(targetPath string) (string, error) {
	// For now, we'll use a simple approach - check if there's a .gopackager metadata file
	metadataFile := targetPath + ".gopackager"

	content, err := os.ReadFile(metadataFile)
	if err != nil {
		// No metadata file exists, return empty string
		//nolint:nilerr
		return "", nil
	}

	return strings.TrimSpace(string(content)), nil
}

// SaveLastCompilationCommit saves the current commit hash as the last compilation commit.
func SaveLastCompilationCommit(targetPath, commitHash string) error {
	metadataFile := targetPath + ".gopackager"
	return os.WriteFile(metadataFile, []byte(commitHash), 0644)
}
