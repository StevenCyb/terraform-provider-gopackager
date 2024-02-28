package getter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit(t *testing.T) {
	t.Parallel()

	repo := "https://github.com/StevenCyb/terraform-provider-gopackager.git"
	branch := "main"
	key := "94c5a5bb8e8bbc2549fa6cfebd600e364ca4f82a"
	git := &Git{}

	t.Run("KeyGeneration", func(t *testing.T) {
		assert.Equal(t, key, git.generateKey(repo, branch))
	})

	t.Run("Get_Clone", func(t *testing.T) {
		destination, err := git.Get(".", repo, branch)
		assert.NoError(t, err)
		assert.Equal(t, "94c5a5bb8e8bbc2549fa6cfebd600e364ca4f82a", destination)
	})

	t.Run("Get_Pull", func(t *testing.T) {
		destination, err := git.Get(".", repo, branch)
		assert.NoError(t, err)
		assert.Equal(t, "94c5a5bb8e8bbc2549fa6cfebd600e364ca4f82a", destination)

		os.RemoveAll(destination)
	})
}
