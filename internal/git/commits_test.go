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
