package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LastCommit(t *testing.T) {
	t.Parallel()

	actual, err := LastCommitHash()
	assert.NoError(t, err)
	assert.Regexp(t, `^[0-9-a-f]{40}$`, actual)
}
