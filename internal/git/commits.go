package git

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// LastCommitHash returns the last commit hash that modified any files in the repository.
// This provides consistent behavior with git_trigger mode and includes all file types
// (Go files, static resources, configuration files, templates, etc.).
func LastCommitHash() (string, error) {
	var commitID string

	cmd := exec.Command("git", "--no-pager", "log", "-1", "--pretty=format:%H")
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

// IsDirty checks if the working directory has uncommitted changes for any files.
// This provides consistent behavior with git_trigger mode and includes all file types
// (Go files, static resources, configuration files, templates, etc.).
func IsDirty() (bool, error) {
	// Check if there are any staged or unstaged changes for any files
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// If output is not empty, there are uncommitted changes
	return strings.TrimSpace(string(output)) != "", nil
}

// IsDirtyForPath checks if the working directory has uncommitted changes for ANY files in the specified path.
// This includes all file types: Go files, static resources, configuration files, templates, etc.
func IsDirtyForPath(path string) (bool, error) {
	// Check if there are any staged or unstaged changes in the specified path (all files)
	cmd := exec.Command("git", "status", "--porcelain", "--", path)
	output, err := cmd.Output()
	if err != nil {
		// If git command fails, assume path is outside repository or doesn't exist
		//nolint:nilerr
		return false, nil
	}

	// If output is not empty, there are uncommitted changes
	return strings.TrimSpace(string(output)) != "", nil
}

// GetModifiedFilesContent returns a deterministic hash of the content of all modified files in the path.
func GetModifiedFilesContent(path string) (string, error) {
	// Get list of modified files
	cmd := exec.Command("git", "status", "--porcelain", "--", path)
	output, err := cmd.Output()
	if err != nil {
		// If git command fails, assume path is outside repository or doesn't exist
		//nolint:nilerr
		return "", nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		// No modified files, return empty hash
		return "", nil
	}

	// Collect file contents in a deterministic order
	var fileContents []string
	var filePaths []string

	for _, line := range lines {
		if line == "" {
			continue
		}
		// Parse git status output (format: "XY filename")
		if len(line) < 4 {
			continue
		}
		filePath := strings.TrimSpace(line[3:])
		filePaths = append(filePaths, filePath)
	}

	// Sort files for deterministic order
	sort.Strings(filePaths)

	// Read file contents
	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			// File might be deleted, skip it
			continue
		}
		fileContents = append(fileContents, fmt.Sprintf("%s:%s", filePath, string(content)))
	}

	// Create a hash of all file contents combined
	if len(fileContents) == 0 {
		return "", nil
	}

	hasher := sha256.New()
	for _, content := range fileContents {
		hasher.Write([]byte(content))
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// GetStableCommitHashWithFallback returns a stable commit hash, with fallback for dirty working directory.
// If the working directory is clean, it returns the actual commit hash.
// If dirty, it returns a deterministic hash based on HEAD commit + modified files content.
func GetStableCommitHashWithFallback() (string, error) {
	// First check if working directory is dirty
	isDirty, err := IsDirty()
	if err != nil {
		return "", err
	}

	// Get the current HEAD commit
	headCommit, err := LastCommitHash()
	if err != nil {
		return "", err
	}

	if !isDirty {
		// Clean working directory, return actual commit hash
		return headCommit, nil
	}

	// Dirty working directory, create fallback hash
	modifiedContent, err := GetModifiedFilesContent(".")
	if err != nil {
		return "", err
	}

	if modifiedContent == "" {
		// No modified files content but still dirty (maybe staged changes), use HEAD + "dirty" indicator
		hasher := sha256.New()
		hasher.Write([]byte(headCommit + ":dirty"))
		fallbackHash := hex.EncodeToString(hasher.Sum(nil))
		return fallbackHash[:40], nil
	}

	// Create deterministic fallback hash: SHA256(HEAD_COMMIT + MODIFIED_CONTENT)
	hasher := sha256.New()
	hasher.Write([]byte(headCommit + ":" + modifiedContent))
	fallbackHash := hex.EncodeToString(hasher.Sum(nil))

	return fallbackHash[:40], nil // Return first 40 characters to match git commit hash length
}

// GetStableCommitHashWithFallbackForPath returns a stable commit hash for a specific path, with fallback for dirty working directory.
// If the working directory is clean for the path, it returns the actual commit hash for that path.
// If dirty, it returns a deterministic hash based on the path's HEAD commit + modified files content.
func GetStableCommitHashWithFallbackForPath(path string) (string, error) {
	// First check if working directory is dirty for this path
	isDirty, err := IsDirtyForPath(path)
	if err != nil {
		return "", err
	}

	// Get the current HEAD commit for this path
	headCommit, err := LastCommitHashForPath(path)
	if err != nil {
		return "", err
	}

	if !isDirty {
		// Clean working directory for this path, return actual commit hash
		return headCommit, nil
	}

	// Dirty working directory, create fallback hash
	modifiedContent, err := GetModifiedFilesContent(path)
	if err != nil {
		return "", err
	}

	if modifiedContent == "" {
		// No modified files content but still dirty (maybe staged changes), use HEAD + "dirty" indicator
		hasher := sha256.New()
		hasher.Write([]byte(headCommit + ":dirty:" + path))
		fallbackHash := hex.EncodeToString(hasher.Sum(nil))
		return fallbackHash[:40], nil
	}

	// Create deterministic fallback hash: SHA256(HEAD_COMMIT + MODIFIED_CONTENT)
	hasher := sha256.New()
	hasher.Write([]byte(headCommit + ":" + modifiedContent))
	fallbackHash := hex.EncodeToString(hasher.Sum(nil))

	return fallbackHash[:40], nil // Return first 40 characters to match git commit hash length
}
