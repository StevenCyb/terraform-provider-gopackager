package getter

import "github.com/stretchr/testify/mock"

// MockGit is an mock type for the Git type.
type MockGit struct {
	mock.Mock
}

// Mocks the Get method.
func (m *MockGit) Get(destination string, repoURL string, branch string) error {
	ret := m.Called(destination, repoURL, branch)

	return ret.Error(0)
}
