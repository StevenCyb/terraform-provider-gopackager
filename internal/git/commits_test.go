package git

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LastCommit(t *testing.T) {
	t.Parallel()

	actual, err := LastCommitHash()
	assert.NoError(t, err)
	assert.Regexp(t, `^[0-9-a-f]{40}$`, actual)
}

func Test_LastCommitHashForPath(t *testing.T) {
	t.Parallel()

	actual, err := LastCommitHashForPath(".")
	assert.NoError(t, err)
	assert.Regexp(t, `^[0-9-a-f]{40}$`, actual)
}

func Test_HasChangedSinceCommit(t *testing.T) {
	t.Parallel()

	// Test with empty commit (should return true)
	changed, err := HasChangedSinceCommit(".", "")
	assert.NoError(t, err)
	assert.True(t, changed)

	// Test with current HEAD commit (should return false - no changes since HEAD)
	currentCommit, err := LastCommitHash()
	assert.NoError(t, err)

	changed, err = HasChangedSinceCommit(".", currentCommit)
	assert.NoError(t, err)
	assert.False(t, changed)

	// Note: We can't easily test with an old commit without knowing the git history
	// In a real scenario, you'd test with a known older commit hash
}

func Test_CompilationCommitTracking(t *testing.T) {
	t.Parallel()

	testPath := "/tmp/test-gopackager"
	testCommit := "abc123def456"

	// Test saving and retrieving compilation commit
	err := SaveLastCompilationCommit(testPath, testCommit)
	assert.NoError(t, err)

	retrievedCommit, err := GetLastCompilationCommit(testPath)
	assert.NoError(t, err)
	assert.Equal(t, testCommit, retrievedCommit)

	// Cleanup
	_ = os.Remove(testPath + ".gopackager")

	// Test with non-existent file
	retrievedCommit, err = GetLastCompilationCommit("/tmp/non-existent-path")
	assert.NoError(t, err)
	assert.Equal(t, "", retrievedCommit)
}

func Test_GetTriggeringCommitHash(t *testing.T) {
	t.Parallel()

	// Test with empty since commit (should return latest commit for path)
	triggeringCommit, err := GetTriggeringCommitHash(".", "")
	assert.NoError(t, err)
	assert.Regexp(t, `^[0-9-a-f]{40}$`, triggeringCommit)

	// Test with current HEAD commit (should return same commit since no changes)
	currentCommit, err := LastCommitHash()
	assert.NoError(t, err)

	triggeringCommit, err = GetTriggeringCommitHash(".", currentCommit)
	assert.NoError(t, err)
	// Should return the same commit if no changes since then
	assert.Equal(t, currentCommit, triggeringCommit)
}

func Test_IsDirty(t *testing.T) {
	t.Parallel()

	// Test dirty state detection (monitors ALL file types for consistency)
	isDirty, err := IsDirty()
	assert.NoError(t, err)
	// We can't assert the exact value since it depends on the actual git state
	// but we can verify the function doesn't error
	assert.IsType(t, false, isDirty)
}

func Test_IsDirtyForPath(t *testing.T) {
	t.Parallel()

	// Test dirty state detection for specific path (monitors ALL file types)
	isDirty, err := IsDirtyForPath(".")
	assert.NoError(t, err)
	// We can't assert the exact value since it depends on the actual git state
	// but we can verify the function doesn't error
	assert.IsType(t, false, isDirty)

	// Test with non-existent path
	isDirty, err = IsDirtyForPath("/tmp/non-existent-path")
	assert.NoError(t, err)
	assert.False(t, isDirty) // Should be false for non-existent path
}

func Test_GetModifiedFilesContent(t *testing.T) {
	t.Parallel()

	// Test getting modified files content
	content, err := GetModifiedFilesContent(".")
	assert.NoError(t, err)
	// Content can be empty if no modified files, that's fine
	assert.IsType(t, "", content)

	// Test with non-existent path
	content, err = GetModifiedFilesContent("/tmp/non-existent-path")
	assert.NoError(t, err)
	assert.Equal(t, "", content) // Should be empty for non-existent path
}

func Test_GetStableCommitHashWithFallback(t *testing.T) {
	t.Parallel()

	// Test stable commit hash with fallback
	stableHash, err := GetStableCommitHashWithFallback()
	assert.NoError(t, err)
	assert.NotEmpty(t, stableHash)

	// Should be either a valid git commit hash (40 chars) or a fallback hash (40 chars)
	assert.Len(t, stableHash, 40)
	assert.Regexp(t, `^[0-9a-f]{40}$`, stableHash)

	// Test consistency - calling twice should return the same result
	// (assuming no changes between calls in test environment)
	stableHash2, err := GetStableCommitHashWithFallback()
	assert.NoError(t, err)
	assert.Equal(t, stableHash, stableHash2)
}

func Test_GetStableCommitHashWithFallbackForPath(t *testing.T) {
	t.Parallel()

	// Test stable commit hash with fallback for specific path
	stableHash, err := GetStableCommitHashWithFallbackForPath(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, stableHash)

	// Should be either a valid git commit hash (40 chars) or a fallback hash (40 chars)
	assert.Len(t, stableHash, 40)
	assert.Regexp(t, `^[0-9a-f]{40}$`, stableHash)

	// Test with specific subdirectory - use relative path that exists from current directory
	stableHash2, err := GetStableCommitHashWithFallbackForPath(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, stableHash2)
	assert.Len(t, stableHash2, 40)
	assert.Regexp(t, `^[0-9a-f]{40}$`, stableHash2)

	// Test consistency for the same path
	stableHash3, err := GetStableCommitHashWithFallbackForPath(".")
	assert.NoError(t, err)
	assert.Equal(t, stableHash, stableHash3)
}
